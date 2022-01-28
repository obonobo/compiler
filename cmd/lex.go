package cmd

import (
	"flag"
	"fmt"
	"path"
	"strings"
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

`

type LexParams struct {
	output     string
	inputFiles []string
}

func lexCmd(config *Config) (usage func(), action func(args []string) int) {
	lexerCmd := flag.NewFlagSet(LEX, flag.ExitOnError)
	lexerCmdOutput := lexerCmd.String("output", "", "")
	lexerCmdO := lexerCmd.String("o", "", "")
	lexerCmd.Usage = func() {
		fmt.Printf(
			LEXUSAGE,
			path.Base(config.Command),
			LEX,
			strings.ToUpper(string(LEX[0]))+LEX[1:])
	}

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

		return Lex(params)
	}
}

// LEX subcommand
func Lex(params LexParams) (exit int) {
	return exit
}
