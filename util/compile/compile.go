//
// Contains utilities for compiling full moon programs
//
package compile

import (
	"bytes"
	"fmt"
	"io"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	tdp "github.com/obonobo/esac/core/tabledrivenparser"
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	tds "github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/core/token/codegen"
	"github.com/obonobo/esac/core/token/visitors"
	"github.com/obonobo/esac/util"
)

// Compiles source code into moon assembly
func TagsBased(src io.Reader) (string, error) {
	return Compile(src, codegen.NewTagsBasedCodeGenVisitor)
}

// Compiles the given source using the provided codegen visitor factory to
// create the codegen visitor.
func Compile[V token.Visitor](
	src io.Reader,
	codeGeneratorFactory func(out, funcOut, dataOut func(string)) V,
) (string, error) {
	chrs := chuggingcharsource.MustChuggingReader(src)
	errs := make([]error, 0, 1024)
	assembly, assemblyData := new(bytes.Buffer), new(bytes.Buffer)
	scnr := scanner.NewObservableScanner(tds.NewScanner(chrs, scannertable.TABLE()))
	prsr := tdp.NewParserNoComments(scnr,
		parsertable.TABLE(),
		func(e *tdp.ParserError) { errs = append(errs, e) }, nil,
		token.Comments()...)

	if !prsr.Parse() {
		return "", fmt.Errorf("parse failed")
	}

	// Callbacks
	logErr := func(e *visitors.VisitorError) { errs = append(errs, e) }
	logAsm := func(line string) { fmt.Fprintln(assembly, line) }
	logData := func(line string) { fmt.Fprintln(assemblyData, line) }

	// Apply visitors
	prsr.AST().Root.Accept(visitors.NewSymTabVisitor(logErr))
	prsr.AST().Root.Accept(visitors.NewSemCheckVisitor(logErr))
	prsr.AST().Root.Accept(codegen.NewMemSizeVisitor())
	prsr.AST().Root.AcceptOnce(codeGeneratorFactory(logAsm, nil, logData))

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("%v", util.Join(errs, "\n"))
	}
	return fmt.Sprintf(
		"%v Main:\n%v\n%v Data:\n%v",
		token.MOON_COMMENT,
		assembly.String(),
		token.MOON_COMMENT,
		assemblyData.String()), err
}
