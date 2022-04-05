package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	ccs "github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/parser"
	"github.com/obonobo/esac/core/scanner"
	tdp "github.com/obonobo/esac/core/tabledrivenparser"
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	tds "github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/core/token/codegen"
	"github.com/obonobo/esac/core/token/visitors"
)

const (
	outderivation   = "outderivation"
	outsyntaxerrors = "outsyntaxerrors"
	scratchDotMoon  = "scratch.moon"
)

// Libraries that should be linked with our moon program
var libs = []string{
	"stdlib/lib.m",
}

func main() {
	// chrs := ccs.MustChugging("../../resources/src/bubblesort.src")
	// chrs := ccs.MustChugging("../../resources/src/polynomial.src")
	// chrs := ccs.MustChuggingReader(bytes.NewBufferString(TYPECHECK_FAIL_1))
	// chrs := ccs.MustChuggingReader(bytes.NewBufferString(TYPECHECK_FAIL_2))
	chrs := ccs.MustChuggingReader(bytes.NewBufferString(CODEGEN))

	errs := make([]error, 0, 1024)
	assembly, assemblyData := new(bytes.Buffer), new(bytes.Buffer)
	scnr := scanner.NewObservableScanner(tds.NewScanner(chrs, scannertable.TABLE()))
	var prsr parser.Parser = tdp.NewParserNoComments(scnr,
		parsertable.TABLE(),
		func(e *tdp.ParserError) { collect(&errs, e) }, nil,
		token.Comments()...)

	if prsr.Parse() {
		// Output callbacks
		logErr := func(e *visitors.VisitorError) { collect(&errs, e) }
		logAsm := logLine(assembly)
		logData := logLine(assemblyData)

		// Print the ast for debugging purposes
		fmt.Printf("\n%v\n", prsr.AST().TreeString())

		// Apply visitors
		prsr.AST().Root.Accept(visitors.NewSymTabVisitor(logErr))
		prsr.AST().Root.Accept(visitors.NewSemCheckVisitor(logErr))
		prsr.AST().Root.Accept(codegen.NewMemSizeVisitor())
		prsr.AST().Root.AcceptOnce(codegen.NewTagsBasedCodeGenVisitor(logAsm, logData))

		// Write output
		token.WritePrettySymbolTable(os.Stdout, prsr.AST().Root.Meta.SymbolTable)
		writeOutSymTabFile(prsr.AST().Root.Meta.SymbolTable)
		writeOutAst(prsr.AST())
		fmt.Printf("\n%v Main:\n%v\n", token.MOON_COMMENT, assembly.String())
		fmt.Printf("%v Data:\n%v\n", token.MOON_COMMENT, assemblyData.String())

		// Write and run the moon program
		writeMoonProgram(assembly.String(), assemblyData.String())
		runMoonProgram()
	} else {
		fmt.Println("Parse failed...")
	}

	for i, e := range errs {
		if ee := new(visitors.Warning); errors.As(e, &ee) {
			errs[i] = fmt.Errorf("WARNING: %w", ee)
		}
	}

	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
	writeOutSemTabErrs(errs)
}

func runMoonProgram() {
	fmt.Fprintf(os.Stderr, "Running %v:\n--------------------\n", scratchDotMoon)
	out, err := exec.Command("./moon", append([]string{scratchDotMoon}, libs...)...).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	fmt.Fprintf(os.Stderr, string(out))
}

func writeMoonProgram(assembly, assemblyData string) {
	fh, err := os.Create(scratchDotMoon)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	fmt.Fprintf(fh, "%v Main:\n%v\n", token.MOON_COMMENT, assembly)
	fmt.Fprintf(fh, "%v Data:\n%v\n", token.MOON_COMMENT, assemblyData)
}

func logLine(w io.Writer) func(string) {
	return func(s string) { fmt.Fprintln(w, s) }
}

func writeOutAst(ast token.AST) {
	file := "outast"
	if fh, err := os.Create(file); err == nil {
		ast.Print(fh)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to write '%v': %v\n", file, err)
	}
}

func writeOutSemTabErrs(errs []error) {
	file := "outsemerrs"
	if fh, err := os.Create(file); err == nil {
		for _, e := range errs {
			fmt.Fprintln(fh, e)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Failed to write '%v': %v\n", file, err)
	}
}

func writeOutSymTabFile(table token.SymbolTable) {
	if fh, err := os.Create("outsymtab"); err == nil {
		token.WritePrettySymbolTable(fh, table)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to write 'outsymbtab': %v\n", err)
	}
}

func collect(errs *[]error, e error) {
	*errs = append(*errs, e)
}

const SHADOWING = `
struct Parent2 {
	public func do_something() -> void;
};

impl Parent2 {
	func do_something() -> void {}
}

struct Parent inherits Parent2 {
	private let a: float;
	private let b: float;
};

struct Child inherits Parent {
	private let a: float;
	private let b: integer;
	public func do_something(x: integer) -> void;
};

impl Child {
	func do_something(x: integer) -> void {}
}
`

const TYPECHECK_FAIL_1 = `
struct MyGuy {
	public func new() -> MyGuy;
	public func say_hello() -> void;
	public func do_something(val: integer[2][1]) -> integer;
};

impl MyGuy {
	func new() -> MyGuy {
		let guy: MyGuy;
		return (guy);
	}

	func say_hello() -> void {
		write(10);
	}

	func do_something(val: integer[2][1]) -> integer {
		let first: integer[1];
		first = val[0];
		return (first[0]);
	}
}

func double(val: float) -> float {
	return (val * val);
}

func main() -> integer {
	let guy: MyGuy;
	let x: integer;
	let y: integer;
	let arg1: integer[2][1];

	guy.say_hello();
	double(guy.do_something(arg1));

	if (1 == 1) then {
		write(double(10.0));
	} else;

	y = 10;
	x = y;
	return (x);
}
`

const TYPECHECK_FAIL_2 = `
// struct Parent2 {
// 	public func do_something() -> void;
// };

// impl Parent2 {
// 	func do_something() -> void {}
// }

// struct Parent inherits Parent2 {
// 	private let a: float;
// 	private let b: float;
// };

// struct Child inherits Parent {
// 	private let a: float;
// 	private let b: integer;
// 	public func do_something(x: integer) -> void;
// };

// impl Child {
// 	func do_something(x: integer) -> void {}
// }

// impl MyImplementation {
// 	func do_something(x: integer[2]) -> void {
// 		let result: float;
// 		let result2: integer[2][4][5];
// 		write(x);
// 	}

// 	func do_something(y: integer) -> void {}

// 	func and_another_one() -> float {
// 		return (2.9);
// 	}
// }

// struct MyImplementation {
// 	public func do_something(x: integer[2]) -> void;
// 	public func do_something(y: integer) -> void;
// 	public func do_something(y: integer) -> void;
// 	public func and_another_one() -> float;
// };

// func top_level() -> void {}
// func top_level(y: integer) -> void {}
// func top_level(x: integer, y: float) -> void {}

struct MyGuy {
	public func new() -> MyGuy;
	public func say_hello() -> void;
	public func do_something(val: integer[2][1]) -> integer;
};

impl MyGuy {
	func new() -> MyGuy {
		let guy: MyGuy;
		return (guy);
	}

	func say_hello() -> void {
		write(10);
	}

	func do_something(val: integer[2][1]) -> integer {
		let first: integer[1];
		first = val[0];
		return (first[0]);
	}
}

// func double(val: integer) -> integer {
// 	return (val * val);
// }

func double(val: float) -> float {
	return (val * val);
}

func main() -> integer {
	let guy: MyGuy;
	let x: integer;
	let y: integer;
	let arg1: integer[2][1];

	guy.say_hello();
	double(guy.do_something(arg1));

	if (1 == 1) then {
		write(double(10.0));
	} else;

	y = 10;
	x = y;
	return (x);
}
`

const CODEGEN = `
func main() -> void {
	let x: integer;
	let y: integer;
	write(1 + 5);
}
`
