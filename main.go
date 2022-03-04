package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/obonobo/esac/cmd"
	ccs "github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/parser"
	tdp "github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	tds "github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/reporting"
)

const (
	outderivation   = "outderivation"
	outsyntaxerrors = "outsyntaxerrors"
)

func main2() {
	cmd.RunAndExit()
}

func main() {
	outder, err := os.Create(outderivation)
	if err != nil {
		panic(err)
	}
	defer outder.Close()
	outsyn, err := os.Create(outsyntaxerrors)
	if err != nil {
		panic(err)
	}
	defer outsyn.Close()

	outderivationLogger := log.New(outder, "", 0)
	outsyntaxerrorsLogger := log.New(outsyn, "", 0)

	// chrs := ccs.MustChugging("resources/a2/src/polynomial.src")
	// chrs := ccs.MustChugging("resources/a2/src/bubblesort.src")
	// chrs := ccs.MustChugging("resources/a2/src/bubblesort-with-errors.src")
	// chrs := ccs.MustChugging("resources/a2/src/polynomial-with-errors-2.src")
	// chrs := ccs.MustChugging("resources/a2/src/polynomial-with-errors.src")
	// chrs := ccs.MustChugging("resources/a2/src/something-else.src")
	chrs := ccs.MustChuggingReader(bytes.NewBufferString(`
	func other(x: integer) -> void {
		write(1);
	}

	func main() -> void {
		// write(id1 + id2 * id3);
	}
	`))

	scnr := tds.NewTableDrivenScanner(chrs, scannertable.TABLE())
	var prsr parser.Parser = tdp.NewTableDrivenParserIgnoringComments(scnr,
		compositetable.TABLE(),
		reporting.ErrSpool(outsyntaxerrorsLogger),
		reporting.RuleSpool(outderivationLogger),
		token.Comments()...)
	valid := prsr.Parse()

	fmt.Println(valid)
	fmt.Print(prsr.AST().TreeString())
}
