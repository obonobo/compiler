package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/core/token/visitors"
	"github.com/obonobo/esac/util"
)

const BUILD = "build"

var BUILD_USAGE = strings.TrimLeft(`
usage: %v %v [-o output] [-d/--outdir <outdir>] [-D|--debug] [input files]

%v builds source code. Runs all compilation steps.

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

	--debug
		Also creates the .outderivation, .outsyntaxerrors, .outlextokens,
		.outsymboltables, .outsemanticerrors, and .outlexerrors files.

`, "\n")

const (
	OUT_SYMBOL_TABLES   = "outsymboltables"
	OUT_SEMANTIC_ERRORS = "outsemanticerrors"
	BUFSIZE             = 1024
)

type BuildParams struct{ ParseParams }

/*
TODO:

We have to support 2 new files: my_file.outsymboltables and
my_file.outsemanticerrors
*/

func buildCmd(config *Config) (usage func(), action func(args []string) int) {
	buildCmd := flag.NewFlagSet(BUILD, flag.ExitOnError)
	buildCmd.Usage = func() {
		fmt.Printf(
			BUILD_USAGE,
			path.Base(config.Command),
			BUILD, titleCase(BUILD))
	}

	var params BuildParams
	buildCmd.StringVar(&params.output, "o", "", "")
	buildCmd.StringVar(&params.output, "output", "", "")
	buildCmd.StringVar(&params.outdir, "d", "", "")
	buildCmd.StringVar(&params.outdir, "outdir", "", "")
	buildCmd.BoolVar(&params.debug, "debug", false, "")

	return buildCmd.Usage, func(args []string) int {
		params.parseFrom(buildCmd, args)
		fmt.Printf(`The '%v build' subcommand has not yet been implemented...`+"\n", config.Command)
		fmt.Printf(`Run '%v help' for usage.`+"\n", config.Command)
		return Build(params)
	}
}

type buildOutputLocations struct {
	outputLocations
	outsymtab    *os.File
	outsemerrors *os.File
}

func Build(params BuildParams) (exit int) {
	// Open all files

	return 1
}

// Builds a single file
func build(params BuildParams, outLocations buildOutputLocations, src io.Reader) error {
	chrs, err := chuggingcharsource.ChuggingReader(src)
	if err != nil {
		return err
	}

	// Output
	var (
		parserErrors  = make([]*tabledrivenparser.ParserError, 0, BUFSIZE)
		visitorErrors = make([]*visitors.VisitorError, 0, BUFSIZE)
		rules         = make([]token.Rule, 0, BUFSIZE)
	)

	scnr := scanner.NewObservableScanner(
		tabledrivenscanner.NewScanner(chrs, scannertable.TABLE()))

	prsr := tabledrivenparser.NewParserNoComments(scnr,
		parsertable.TABLE(),
		util.Appendback(&parserErrors),
		util.Appendback(&rules),
		token.Comments()...)

	// Apply the visitors
	if prsr.Parse() {
		var (
			logErr   = util.Appendback(&visitorErrors)
			symtab   = visitors.NewSymTabVisitor(logErr)
			semcheck = visitors.NewSemCheckVisitor(logErr)
		)
		ast := prsr.AST().Root
		ast.Accept(symtab)
		ast.Accept(semcheck)
	} else {

	}

	return nil
}

func (p *BuildParams) parseFrom(buildCmd *flag.FlagSet, args []string) {
	buildCmd.Parse(args)
	p.inputFiles = buildCmd.Args()
	p.outputMode = outputMode(p.output)
	p.outdir = outdir(p.outdir)
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0]) + s[1:])
}
