package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/parser"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/reporting"
)

const PARSE = "parse"
const OUTAST = "outast"

var PARSE_USAGE = strings.TrimLeft(`
usage: %v %v [-o output] [-d/--outdir <outdir>] [-D|--debug] [input files]

%v converts the input files to tokens and then consumes the token stream to
convert it to an AST. This command produces a file for every input file:
'myfile.outast'

If no input files are specified, input is read from STDIN.

Flags:

	-o, --output [outfile|-]
		An alternative output location. If this flag is specified,
		the file described above will not be created, only the specified
		file will be created. Specify '-' to print the token stream to
		STDOUT.

	-d, --outdir [outdir]
		An alternative output location for the files. The default output
		location is the current directory.

	-D, --debug
		Also creates the .outderivation, .outsyntaxerrors, .outlextokens,
		and .outlexerrors files.

`, "\n")

const (
	OUT_DERIVATION    = "outderivation"
	OUT_SYNTAX_ERRORS = "outsyntaxerrors"
)

type ParseParams struct {
	LexParams
	debug bool
	input *os.File
}

func parseCmd(config *Config) (usage func(), action func(args []string) (exit int)) {
	parseCmd := flag.NewFlagSet(PARSE, flag.ExitOnError)
	parseCmd.Usage = func() {
		fmt.Printf(
			PARSE_USAGE,
			path.Base(config.Command),
			PARSE, titleCase(PARSE))
	}

	params := ParseParams{}
	parseCmd.StringVar(&params.output, "o", "", "")
	parseCmd.StringVar(&params.output, "output", "", "")
	parseCmd.StringVar(&params.outdir, "d", "", "")
	parseCmd.StringVar(&params.outdir, "outdir", "", "")
	parseCmd.BoolVar(&params.debug, "debug", false, "")

	return parseCmd.Usage, func(args []string) (exit int) {
		parseCmd.Parse(args)
		params.inputFiles = parseCmd.Args()
		params.outputMode = outputMode(params.output)
		params.outdir = outdir(params.outdir)
		if exit := checkInputFiles(config, params.LexParams, PARSE); exit != 0 {
			return exit
		}
		if len(params.LexParams.inputFiles) == 0 {
			params.input = os.Stdin
		}
		return Parse(params)
	}
}

// PARSE subcommand
func Parse(params ParseParams) (exit int) {
	if exit := makeOutputDirIfNotExists(params.outdir); exit != EXIT_CODE_OKAY {
		return exit
	}

	output, close, exit := openOutputLocations(params)
	if exit != EXIT_CODE_OKAY {
		return exit
	}
	defer close()

	var i int
	moreThanOne := len(params.inputFiles) > 1
	for _, out := range sortedOutputLocations(output) {
		file := out.file
		if moreThanOne {
			to := os.Stdout
			if i > 0 {
				fmt.Fprintln(to)
			}
			fmt.Fprintf(to, "%v:\n", file)
		}
		parse(out.locations, params)
		i++
	}

	return EXIT_CODE_OKAY
}

// Parses a single input file
func parse(out *outputLocations, params ParseParams) {
	scnr := createScanner(out.source)
	outlextokens, outlexerrors := reporting.StreamTokensSplitErrors(scnr.Subscribe())
	outsyntaxerrors := make(chan tabledrivenparser.ParserError, 1024)

	var outderivation chan token.Rule
	if params.debug {
		outderivation = make(chan token.Rule, 1024)
	}

	prsr := createParser(scnr, outsyntaxerrors, outderivation)

	var wait sync.WaitGroup
	if params.debug {
		wait.Add(4)
		goWriteTo(&wait, outlextokens, out.outlextokens)
		goWriteTo(&wait, outlexerrors, out.outlexerrors, os.Stderr)
		goWriteTo(&wait, outderivation, out.outderivation)
		goWriteTo(&wait, outsyntaxerrors, out.outsyntaxerrors, os.Stderr)
	} else {
		wait.Add(3)

		// Eat outlextokens
		go func() {
			for range outlextokens {
			}
			wait.Done()
		}()

		// Log the errors
		goWriteTo(&wait, outlexerrors, os.Stderr)
		goWriteTo(&wait, outsyntaxerrors, os.Stderr)
	}

	// Write the AST
	if prsr.Parse() {
		prsr.AST().Print(out.outast)
	}
	wait.Wait()
}

// Asynchronously writes from channel to writer(s), calls wait.Done() upon
// completion
func goWriteTo[T any](
	wait *sync.WaitGroup,
	from <-chan T,
	to ...io.Writer,
) {
	go func() {
		for item := range from {
			for _, writer := range to {
				fmt.Fprintln(writer, item)
			}
		}
		wait.Done()
	}()
}

// This function opens all files that need to be opened per params
//
// We need to open 4-5 output files:
// 1. myfile.outlextokens
// 2. myfile.outlexerrors
// 3. myfile.outderivation
// 4. myfile.outsyntaxerrors
// 5. myfile.outast iff params.outputMode is OUT_MODE_TOFILE
func openOutputLocations(params ParseParams) (
	out map[string]*outputLocations,
	close func(),
	exit int,
) {
	chugged, exit := openAndChugFilesHandleStdin(params)
	if exit != EXIT_CODE_OKAY {
		return nil, func() {}, exit
	}

	outputs := make(map[string]*outputLocations, len(chugged))
	close = func() {
		for _, v := range outputs {
			v.close()
		}
	}

	for file, source := range chugged {
		out, exit := openAllFiles(file, source.src, params)
		if exit != EXIT_CODE_OKAY {
			close() // Close all previously opened files
			return nil, func() {}, exit
		}
		out.index = source.i
		outputs[file] = out
	}
	return outputs, close, EXIT_CODE_OKAY
}

type indexedSource struct {
	i   int
	src scanner.CharSource
}

func openAndChugFilesHandleStdin(params ParseParams) (map[string]indexedSource, int) {
	if params.input != nil {
		// Then we have a single file, read from stdin
		in, err := chuggingcharsource.ChuggingReader(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, EXIT_CODE_NOT_OKAY
		}
		return map[string]indexedSource{"a.in": {0, in}}, EXIT_CODE_OKAY
	}
	return openAndChugFiles(params.inputFiles)
}

// A big record holding all information output information
type outputLocations struct {
	source          scanner.CharSource
	outlextokens    *os.File
	outlexerrors    *os.File
	outsyntaxerrors *os.File
	outderivation   *os.File
	outast          *os.File
	close           func()
	index           int // sort
}

func openAllFiles(
	inputFile string,
	chrs scanner.CharSource,
	params ParseParams,
) (*outputLocations, int) {
	out := &outputLocations{source: chrs, close: func() {}}
	if params.debug {
		if exit := setupFiles(out, params, inputFile); exit != EXIT_CODE_OKAY {
			return nil, exit
		}
	}
	if exit := setupOutputFile(out, params, inputFile); exit != EXIT_CODE_OKAY {
		return nil, exit
	}
	return out, EXIT_CODE_OKAY
}

func setupOutputFile(output *outputLocations, params ParseParams, file string) (exit int) {
	switch params.outputMode {
	case OUT_MODE_NORMAL, OUT_MODE_TOFILE:
		outast, err := os.Create(path.Join(params.outdir, filename(file, OUTAST)))
		if err != nil {
			output.close()
			fmt.Fprintln(os.Stderr, failedToOpenFileError(err))
			return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}
		output.outast = outast
	case OUT_MODE_STDOUT:
		output.outast = os.Stdout
	default:
		output.close()
		fmt.Fprintf(os.Stderr,
			"invalid output mode (%v), should be %v, %v, or %v\n",
			params.outputMode, OUT_MODE_NORMAL, OUT_MODE_STDOUT, OUT_MODE_TOFILE)
		return EXIT_CODE_NOT_OKAY
	}

	oldClose := output.close
	output.close = func() {
		oldClose()
		if output.outast != os.Stdout {
			output.outast.Close()
		}
	}
	return EXIT_CODE_OKAY
}

func setupFiles(output *outputLocations, params ParseParams, file string) (exit int) {
	files, close, exit := openN(
		path.Join(params.outdir, filename(file, OUT_LEX_ERRORS)),
		path.Join(params.outdir, filename(file, OUT_LEX_TOKENS)),
		path.Join(params.outdir, filename(file, OUT_SYNTAX_ERRORS)),
		path.Join(params.outdir, filename(file, OUT_DERIVATION)))

	if exit != EXIT_CODE_OKAY {
		return exit
	}

	output.outlexerrors = files[0]
	output.outlextokens = files[1]
	output.outsyntaxerrors = files[2]
	output.outderivation = files[3]
	output.close = close

	return EXIT_CODE_OKAY
}

// Opens N files, logs all errors, merges file closing functions
func openN(paths ...string) (files []*os.File, close func(), exit int) {
	files = make([]*os.File, 0, len(paths))
	close = func() {}
	for _, p := range paths {
		fh, err := os.Create(p)
		if err != nil {
			close()
			fmt.Fprintln(os.Stderr, failedToOpenFileError(err))
			return nil, nil, EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}
		files = append(files, fh)
		oldClose := close
		close = func() {
			oldClose()
			fh.Close()
		}
	}
	return files, close, EXIT_CODE_OKAY
}

func failedToOpenFileError(wrap error) error {
	return fmt.Errorf("failed to open file: %w", wrap)
}

func outdir(dir string) string {
	if dir == "-" || dir == "" {
		return "."
	}
	return dir
}

func createParser(
	scnr scanner.Scanner,
	errc chan<- tabledrivenparser.ParserError,
	rulec chan<- token.Rule,
) parser.Parser {
	// return tabledrivenparser.NewParser(scnr, parsertable.TABLE(), errc, rulec)
	return tabledrivenparser.NewParser(scnr, parsertable.TABLE(),
		func(e *tabledrivenparser.ParserError) { errc <- *e },
		func(r token.Rule) { rulec <- r })
}

func createScanner(chrs scanner.CharSource) *scanner.ObservableScanner {
	return scanner.NewObservableScanner(
		scanner.IgnoringComments(
			tabledrivenscanner.NewScanner(chrs, scannertable.TABLE()),
			token.Comments()...))
}

// Replaces file extension
func filename(file, extension string) string {
	if len(file) == 0 {
		return "." + extension
	}
	file = path.Base(file)
	if strings.Contains(file, ".") {
		file = strings.TrimRight(
			strings.TrimRightFunc(file, func(r rune) bool { return r != '.' }), ".")
	}
	return fmt.Sprintf("%v.%v", file, extension)
}

type ele struct {
	file      string
	locations *outputLocations
}

func sortedOutputLocations(m map[string]*outputLocations) []ele {
	out := make([]ele, 0, len(m))
	for k, v := range m {
		out = append(out, ele{k, v})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].locations.index < out[j].locations.index
	})
	return out
}
