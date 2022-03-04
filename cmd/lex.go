package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/reporting"
)

const LEX = "lex"

const LEXUSAGE = `usage: %v %v [-o output] [input files]

%v converts the input files to tokens and save them to files. Two files will be
generated for every input file: 'myfile.outlextokens', and 'myfile.outlexerrors'.

If no input files are specified, input is read from STDIN.

Flags:

	-o, --output [outfile|-]
		An alternative output location. If this flag is specified,
		the 2 files described above will not be created, only the
		specified file will be created and it will contain both valid
		and error tokens. Specify '-' to print the token stream to STDOUT.

	-d, --outdir [outdir]
		An alternative output location for the two files ('.outlextokens' and
		'.outlexerrors'). The default output location is the current directory.

`

const (
	OUT_MODE_STDOUT = iota // Output tokens to STDOUT
	OUT_MODE_NORMAL        // Output tokens to the two files: `.outlextokens`, `.outlexerrors`
	OUT_MODE_TOFILE        // Output tokens to a single file
)

const (
	EXIT_CODE_OKAY = iota
	EXIT_CODE_NOT_OKAY
	EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
	EXIT_CODE_CANNOT_CREATE_OUTPUT_DIR
)

const (
	OUTLEXTOKENS = "outlextokens"
	OUTLEXERRORS = "outlexerrors"
)

type LexParams struct {
	output     string
	outdir     string
	outputMode int
	inputFiles []string
}

func lexCmd(config *Config) (usage func(), action func(args []string) int) {
	lexerCmd := flag.NewFlagSet(LEX, flag.ExitOnError)
	lexerCmd.Usage = func() {
		fmt.Printf(
			LEXUSAGE,
			path.Base(config.Command),
			LEX,
			strings.ToUpper(string(LEX[0]))+LEX[1:])
	}

	lexerCmdOutput := lexerCmd.String("output", "", "")
	lexerCmdO := lexerCmd.String("o", "", "")

	lexerCmdOutdir := lexerCmd.String("outdir", "", "")
	lexerCmdD := lexerCmd.String("d", "", "")

	return lexerCmd.Usage, func(args []string) int {
		var params LexParams
		lexerCmd.Parse(args)

		// Input files
		params.inputFiles = lexerCmd.Args()

		// Output location
		if lexerCmdO != nil && *lexerCmdO != "" {
			params.output = *lexerCmdO
		}
		if lexerCmdOutput != nil && *lexerCmdOutput != "" {
			params.output = *lexerCmdOutput
		}

		// Output dir
		if lexerCmdD != nil && *lexerCmdD != "" {
			params.outdir = *lexerCmdD
		}
		if lexerCmdOutdir != nil && *lexerCmdOutdir != "" {
			params.outdir = *lexerCmdOutdir
		}

		params.outputMode = outputMode(params.output)

		if exit, msg := checkParams(config, params); exit != 0 {
			fmt.Println(msg)
			return exit
		}

		return Lex(params)
	}
}

// LEX subcommand
func Lex(params LexParams) (exit int) {
	switch params.outputMode {
	case OUT_MODE_NORMAL:
		return lexNormal(params)
	case OUT_MODE_TOFILE:
		return lexToFile(params)
	case OUT_MODE_STDOUT:
		return lexToStdout(params)
	}

	fmt.Println("Something went wrong lexing...")
	return EXIT_CODE_NOT_OKAY
}

// In normal mode, the program writes errors and tokens to different files
func lexNormal(params LexParams) int {
	outdir := params.outdir
	if outdir == "-" || outdir == "" {
		outdir = "."
	}

	// Create output directory if it doesn't already exist
	if code := makeOutputDirIfNotExists(outdir); code != EXIT_CODE_OKAY {
		return code
	}

	chugged, exit := openAndChugFiles(params.inputFiles)
	if exit != EXIT_CODE_OKAY {
		return exit
	}

	// Lex and write output files
	var wait sync.WaitGroup
	for file, source := range chugged {
		outTokens, err := os.Create(
			path.Join(outdir, inputFileNameToOutputFileName(file, OUTLEXTOKENS)))

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			wait.Wait()
			return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}

		defer outTokens.Close()

		outErrors, err := os.Create(
			path.Join(outdir, inputFileNameToOutputFileName(file, OUTLEXERRORS)))

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			wait.Wait()
			return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}

		defer outErrors.Close()

		out, errs := reporting.StreamLinesSplitErrors(
			compositetable.NewTableDrivenScanner(source), -1)

		wait.Add(1)
		go func() {
			for line := range out {
				fmt.Fprintln(outTokens, line)
			}
			wait.Done()
		}()

		wait.Add(1)
		go func() {
			for line := range errs {
				fmt.Fprintln(outErrors, line)
			}
			wait.Done()
		}()
	}

	wait.Wait()
	return EXIT_CODE_OKAY
}

func inputFileNameToOutputFileName(name string, extension string) string {
	base := path.Base(name)
	trimExtension := strings.TrimRightFunc(base, func(r rune) bool { return r != '.' })
	return trimExtension + extension
}

func lexToFile(params LexParams) int {
	out, err := os.Open(params.output)
	if err != nil {
		return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
	}
	defer out.Close()
	return lexTo(params, out)
}

func lexToStdout(params LexParams) int {
	return lexTo(params, os.Stdout)
}

func reportFileErrors(errs []error, out io.Writer) int {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(out, capitalizeFirstLetter(err))
		}
		return EXIT_CODE_NOT_OKAY
	}
	return EXIT_CODE_OKAY
}

func lexTo(params LexParams, to io.Writer) int {
	chugged, exit := openAndChugFiles(params.inputFiles)
	if exit != EXIT_CODE_OKAY {
		return exit
	}

	lines := streamCharSources(chugged)

	var i int
	for fileName, out := range lines {
		fmt.Fprintf(to, "Scanner output for file '%v':\n", fileName)
		fmt.Fprintln(to, "----------------------------------------------")
		for s := range out {
			fmt.Fprintln(to, s)
		}
		if i < len(lines)-1 {
			fmt.Fprintln(to)
		}
		i++
	}

	return EXIT_CODE_OKAY
}

// Capitalizes the first letter of an error message
func capitalizeFirstLetter(err error) string {
	s := err.Error()
	return fmt.Sprint(strings.ToUpper(s[:1]) + s[1:])
}

func streamCharSources(
	sources map[string]scanner.CharSource,
) (lines map[string]<-chan string) {
	lines = make(map[string]<-chan string, len(sources))
	for fileName, chugger := range sources {
		s := compositetable.NewTableDrivenScanner(chugger)
		lines[fileName] = reporting.StreamLines(s, -1)
	}
	return lines
}

func openAndChugFiles(inputFiles []string) (map[string]scanner.CharSource, int) {
	files, errs := openInputFiles(inputFiles)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	if e := reportFileErrors(errs, os.Stdout); e != EXIT_CODE_OKAY {
		return nil, e
	}

	chuggers, errs := chugFiles(files)
	if e := reportFileErrors(errs, os.Stdout); e != EXIT_CODE_OKAY {
		return nil, e
	}

	ret := make(map[string]scanner.CharSource, len(files))
	for i, f := range files {
		ret[f.Name()] = chuggers[i]
	}

	return ret, EXIT_CODE_OKAY
}

func chugFiles(files []*os.File) ([]*chuggingcharsource.ChuggingCharSource, []error) {
	out := make([]*chuggingcharsource.ChuggingCharSource, 0, len(files))
	errs := make([]error, 0, len(files))
	for _, f := range files {
		var c chuggingcharsource.ChuggingCharSource
		err := c.ChugReader(f)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read input file '%v'", f.Name()))
		} else {
			out = append(out, &c)
		}
	}
	return out, errs
}

// Open a set of files and report errors
func openInputFiles(files []string) ([]*os.File, []error) {
	out := make([]*os.File, 0, len(files))
	errs := make([]error, 0, len(files))
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open input file '%v'", file))
		} else {
			out = append(out, f)
		}
	}
	return out, errs
}

// Validate params. If the exit code returned by this function is non-zero, then
// the program should fail early and report the reason for failure
func checkParams(config *Config, params LexParams) (exit int, msg string) {
	if params.outdir != "" && params.output != "" {
		fmt.Fprintln(os.Stderr, ""+
			"WARNING: flag '-d'/'--outdir' has no "+
			"effect when used alongside flag '-o'/'--output'")
	}

	if len(params.inputFiles) < 1 {
		return 1, fmt.Sprintf(""+
			"Please provide input files to the lex command.\n"+
			"usage: %v %v [-o output] [input files]", path.Base(config.Command), LEX)
	}
	return exit, ""
}

// Determines which mode of output we need to be using
func outputMode(outputParam string) int {
	switch outputParam {
	case "-":
		return OUT_MODE_STDOUT
	case "":
		return OUT_MODE_NORMAL
	default:
		return OUT_MODE_TOFILE
	}
}

func makeOutputDirIfNotExists(outdir string) (exit int) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		if err := os.MkdirAll(outdir, 0o775); err != nil {
			return EXIT_CODE_CANNOT_CREATE_OUTPUT_DIR
		}
	}
	return EXIT_CODE_OKAY
}
