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
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
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
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}

	// struct Hey inherits Yo {
	// 	private let a: float;
	// 	public let b: integer;
	// 	public func doIt(A: float, B: float) -> Yo;
	// 	public func hey(x: float, y: integer) -> Hey;
	// };

	// struct Hey2 {
	// 	private let a: float;
	// 	public func doIt(A: float, B: float) -> Yo;
	// };

	// impl MyImplementation {
	// 	func do_something(x: integer[2]) -> void {
	// 		write(x);
	// 	}

	// 	func and_another_one() -> float {
	// 		return (2.9);
	// 	}
	// }

	// func other() -> void {
	// 	id3 = 12;
	// }

	// func other() -> void {
	// 	id1[1][2][3].id2(1).id3 = 12;
	// }

	// func other() -> void {
	// 	// read(id1[1][2][3].id2(1).id3);
	// 	if (1 < 2) then {
	// 		write(1);
	// 	} else {
	// 		read(id1);
	// 		while (1 == 1) {
	// 			// noop
	// 		};
	// 	};
	// }

	// func other(x: integer) -> void {
	// 	write(1);
	// }

	// func other(x: integer[2][3][4]) -> integer {
	// 	write(1 + 1 - 1 * 1 < 5 / 5 + 3 - 2);
	// }

	// func other(x: integer[2][3][4], z: float, y: integer[2][3][4]) -> integer {
	// 	return(1 + 1 - 1 * 1 < 5 / 5 + 3 - 2);
	// }

	// func main() -> void {
	// 	// write(id1 + id2 * id3);
	// }
	`))


	scnr := tds.NewTableDrivenScanner(chrs, scannertable.TABLE())
	var prsr parser.Parser = tdp.NewParserNoComments(scnr,
		parsertable.TABLE(),
		reporting.ErrSpool(outsyntaxerrorsLogger),
		reporting.RuleSpool(outderivationLogger),
		token.Comments()...)
	valid := prsr.Parse()

	fmt.Println(valid)
	fmt.Print(prsr.AST().TreeString())
}
