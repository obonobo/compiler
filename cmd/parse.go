package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

// TODO: this file is incomplete
// TODO: this file is incomplete
// TODO: this file is incomplete
// TODO: this file is incomplete
// TODO: this file is incomplete

const PARSE = "parse"

var PARSE_USAGE = strings.TrimLeft(`
usage: %v %v [-o output] [input files]

%v converts the input files to tokens and then consumes the token stream to
convert it to an AST. This command produces two files for every input file:
'myfile.outderivation', and 'myfile.outsyntaxerrors'.

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

`, "\n")

const (
	OUT_DERIVATION    = "outderivation"
	OUT_SYNTAX_ERRORS = "outsyntaxerrors"
)

type ParseParams struct{ LexParams }

func parseCmd(config *Config) (usage func(), action func(args []string) (exit int)) {
	parseCmd := flag.NewFlagSet(PARSE, flag.ExitOnError)
	parseCmd.Usage = func() {
		fmt.Printf(
			PARSE_USAGE,
			path.Base(config.Command),
			PARSE, strings.ToUpper(string(PARSE[0]))+PARSE[1:])
	}

	params := struct{ LexParams }{}
	parseCmd.StringVar(&params.LexParams.output, "o", "", "")
	parseCmd.StringVar(&params.LexParams.output, "output", "", "")
	parseCmd.StringVar(&params.LexParams.outdir, "d", "", "")
	parseCmd.StringVar(&params.LexParams.outdir, "outdir", "", "")

	return parseCmd.Usage, func(args []string) (exit int) {
		// Using some functions from lex.go to aid our parsing
		parseCmd.Parse(args)
		params.LexParams.inputFiles = parseCmd.Args()
		params.LexParams.outputMode = outputMode(params.output)
		if exit, msg := checkParams(config, params.LexParams); exit != 0 {
			fmt.Println(msg)
			return exit
		}
		if params.outdir == "-" || params.outdir == "" {
			params.outdir = "."
		}
		return Parse(params)
	}
}

// PARSE subcommand
func Parse(params ParseParams) (exit int) {
	switch params.outputMode {
	case OUT_MODE_NORMAL:
		return parseNormal(params)
	case OUT_MODE_TOFILE:
		return parseToFile(params)
	case OUT_MODE_STDOUT:
		return parseToStdout(params)
	}

	fmt.Println("Something went wrong, not sure where to print results...")
	return EXIT_CODE_NOT_OKAY
}

// TODO: complete this function
func parseNormal(params ParseParams) (exit int) {
	if exit := makeOutputDirIfNotExists(params.outdir); exit != EXIT_CODE_OKAY {
		return exit
	}

	chugged, exit := openAndChugFiles(params.inputFiles)
	if exit != EXIT_CODE_OKAY {
		return exit
	}

	var wait sync.WaitGroup
	for file, source := range chugged {
		outDerivation, err := os.Create(
			path.Join(params.outdir, inputFileNameToOutputFileName(file, OUT_DERIVATION)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			wait.Wait()
			return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}
		defer outDerivation.Close()

		outErrors, err := os.Create(
			path.Join(params.outdir, inputFileNameToOutputFileName(file, OUT_SYNTAX_ERRORS)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			wait.Wait()
			return EXIT_CODE_CANNOT_OPEN_OUTPUT_FILE
		}
		defer outErrors.Close()

		// TODO: remove this line when you finish the CLI part
		fmt.Println(source)
	}

	wait.Wait()
	return EXIT_CODE_OKAY
}

func parseToFile(params ParseParams) (exit int) {
	return EXIT_CODE_NOT_OKAY
}

func parseToStdout(params ParseParams) (exit int) {
	return EXIT_CODE_NOT_OKAY
}
