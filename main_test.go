package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	parsertable "github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/core/token/visitors"
	"github.com/obonobo/esac/internal/testutils"
)

// This is the transition table that we are using for our scanner
var table = scannertable.TABLE

func TestSymTabVisitor_SimpleImplBeforeStruct(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}

	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation   |
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                MyImplementation
				+-------------------------------------------------------------------------------+
				| Name            | Kind    | Type           | Link                             |
				+-------------------------------------------------------------------------------+
				| do_something    | FuncDef | (public) void  | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDef | (public) float | ⊙---> and_another_one()          |
				+-------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
	`)
}

func TestSymTabVisitor_SimpleStructBeforeImpl(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};

	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                MyImplementation
				+-------------------------------------------------------------------------------+
				| Name            | Kind    | Type           | Link                             |
				+-------------------------------------------------------------------------------+
				| do_something    | FuncDef | (public) void  | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDef | (public) float | ⊙---> and_another_one()          |
				+-------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
	`)
}

func TestSymTabVisitor_DuplicateStructDefined(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}

	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};

	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation   |
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                MyImplementation
				+-------------------------------------------------------------------------------+
				| Name            | Kind    | Type           | Link                             |
				+-------------------------------------------------------------------------------+
				| do_something    | FuncDef | (public) void  | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDef | (public) float | ⊙---> and_another_one()          |
				+-------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
			duplicate definition for 'MyImplementation' (defined on line 14, and again on line 19)
	`)
}

func TestSymTabVisitor_ImplMissingStructMethod(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		// func and_another_one() -> float {
		// 	return (2.9);
		// }
	}

	struct MyImplementation {
		public let x: integer;
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation   |
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                MyImplementation
				+------------------------------------------------------------------------------+
				| Name         | Kind    | Type             | Link                             |
				+------------------------------------------------------------------------------+
				| x            | VarDecl | (public) integer | ____                             |
				| do_something | FuncDef | (public) void    | ⊙---> do_something(integer[2])   |
				+------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+
			impl 'MyImplementation' is missing method 'public func and_another_one() -> float' defined in struct (line 17)
	`)
}

func TestSymTabVisitor_StructMissingImplMethod(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}

	struct MyImplementation {
		public let x: integer;
		// public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation   |
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                 MyImplementation
				+---------------------------------------------------------------------------------+
				| Name            | Kind    | Type             | Link                             |
				+---------------------------------------------------------------------------------+
				| x               | VarDecl | (public) integer | ____                             |
				| do_something    | FuncDef | void             | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDef | (public) float   | ⊙---> and_another_one()          |
				+---------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
			struct 'MyImplementation' is missing method 'func do_something(integer[2]) -> void' defined in impl (line 3)
	`)
}

func TestSymTabVisitor_NoStructFoundForImpl(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	// struct MyImplementation {
	// 	public func do_something(x: integer[2]) -> void;
	// 	public func and_another_one() -> float;
	// };

	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}
	`, `
			                            Global
			+--------------------------------------------------------------+
			| Name             | Kind    | Type | Link                     |
			+--------------------------------------------------------------+
			| MyImplementation | ImplDef | ____ | ⊙---> MyImplementation   |
			+--------------------------------------------------------------+

				                            MyImplementation
				+----------------------------------------------------------------------+
				| Name            | Kind    | Type  | Link                             |
				+----------------------------------------------------------------------+
				| do_something    | FuncDef | void  | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDef | float | ⊙---> and_another_one()          |
				+----------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
			malformed type: no struct found for impl 'MyImplementation', impl methods must first be declared in a struct
	`)
}

func TestSymTabVisitor_NoImplFoundForStruct(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func and_another_one() -> float;
	};

	// impl MyImplementation {
	// 	func do_something(x: integer[2]) -> void {
	// 		let result: float;
	// 		let result2: integer[2][4][5];
	// 		write(x);
	// 	}
	`, `
			                              Global
			+-----------------------------------------------------------------+
			| Name             | Kind       | Type | Link                     |
			+-----------------------------------------------------------------+
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation   |
			+-----------------------------------------------------------------+

				                                 MyImplementation
				+--------------------------------------------------------------------------------+
				| Name            | Kind     | Type           | Link                             |
				+--------------------------------------------------------------------------------+
				| do_something    | FuncDecl | (public) void  | ⊙---> do_something(integer[2])   |
				| and_another_one | FuncDecl | (public) float | ⊙---> and_another_one()          |
				+--------------------------------------------------------------------------------+

					      do_something(integer[2])
					+----------------------------------+
					| Name | Kind  | Type       | Link |
					+----------------------------------+
					| x    | Param | integer[2] | ____ |
					+----------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+
			malformed type: no impl found for struct 'MyImplementation', struct methods declared but not defined
	`)
}

func TestSymTabVisitor_DuplicateMethodDefinitions(t *testing.T) {
	t.Parallel()
	assertSymbolTableOutput(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			let result: float;
			let result2: integer[2][4][5];
			write(x);
		}

		func do_something(x: integer) -> void {}
		func do_something(y: integer) -> void {}

		func and_another_one() -> float {
			return (2.9);
		}
	}

	struct MyImplementation {
		public func do_something(x: integer[2]) -> void;
		public func do_something(x: integer) -> void;
		public func do_something(y: integer) -> void;
		public func and_another_one() -> float;
	};

	func top_level() -> void {}
	func top_level(x: integer) -> void {}
	func top_level(y: integer) -> void {}
	func top_level(x: integer, y: float) -> void {}
	`, `
			                                  Global
			+--------------------------------------------------------------------------+
			| Name             | Kind       | Type | Link                              |
			+--------------------------------------------------------------------------+
			| MyImplementation | ImplDef    | ____ | ⊙---> MyImplementation            |
			| MyImplementation | StructDecl | ____ | ⊙---> MyImplementation            |
			| top_level        | FuncDef    | void | ⊙---> top_level()                 |
			| top_level        | FuncDef    | void | ⊙---> top_level(integer)          |
			| top_level        | FuncDef    | void | ⊙---> top_level(integer, float)   |
			+--------------------------------------------------------------------------+

				                                MyImplementation
				+-------------------------------------------------------------------------------+
				| Name            | Kind    | Type           | Link                             |
				+-------------------------------------------------------------------------------+
				| do_something    | FuncDef | (public) void  | ⊙---> do_something(integer[2])   |
				| do_something    | FuncDef | (public) void  | ⊙---> do_something(integer)      |
				| and_another_one | FuncDef | (public) float | ⊙---> and_another_one()          |
				+-------------------------------------------------------------------------------+

					            do_something(integer[2])
					+---------------------------------------------+
					| Name    | Kind    | Type             | Link |
					+---------------------------------------------+
					| x       | Param   | integer[2]       | ____ |
					| result  | VarDecl | float            | ____ |
					| result2 | VarDecl | integer[2][4][5] | ____ |
					+---------------------------------------------+

					       do_something(integer)
					+-------------------------------+
					| Name | Kind  | Type    | Link |
					+-------------------------------+
					| x    | Param | integer | ____ |
					+-------------------------------+

					       and_another_one()
					+---------------------------+
					| Name | Kind | Type | Link |
					+---------------------------+
					+---------------------------+

				         top_level()
				+---------------------------+
				| Name | Kind | Type | Link |
				+---------------------------+
				+---------------------------+

				       top_level(integer)
				+-------------------------------+
				| Name | Kind  | Type    | Link |
				+-------------------------------+
				| x    | Param | integer | ____ |
				+-------------------------------+

				    top_level(integer, float)
				+-------------------------------+
				| Name | Kind  | Type    | Link |
				+-------------------------------+
				| x    | Param | integer | ____ |
				| y    | Param | float   | ____ |
				+-------------------------------+
			duplicate definition for 'do_something' (defined on line 9, and again on line 10)
			duplicate definition for 'do_something' (defined on line 19, and again on line 20)
			'MyImplementation::do_something' has been overloaded 2 times: do_something(integer[2]), do_something(integer)
			duplicate definition for 'top_level' (defined on line 25, and again on line 26)
			'Global::top_level' has been overloaded 3 times: top_level(), top_level(integer), top_level(integer, float)
	`)
}

func assertSymbolTableOutput(t *testing.T, input, output string) {
	prefix := "			"
	output = clean(output, prefix)
	prsr, errs := createErrorLoggingParser(input)
	if !prsr.Parse() {
		t.Fatalf("Parse failed...")
	}

	// Apply the visitor
	visitorOut, visitorErr := new(bytes.Buffer), new(bytes.Buffer)
	out := func() string { return visitorOut.String() + errs() + visitorErr.String() }
	ast := prsr.AST()
	ast.Root.Accept(visitors.NewSymTabVisitor(func(e *visitors.VisitorError) {
		fmt.Fprintln(visitorErr, e)
	}))
	token.WritePrettySymbolTable(visitorOut, ast.Root.Meta.SymbolTable)

	// Assert output
	if expected, actual := output, out(); expected != actual {
		t.Fatalf("\nExpected output:\n%v\n\nActual output:\n%v", expected, actual)
	}
}

// Tests parsing some source that has multiple "missing token" errors meant to
// be caught by the statement closer
func TestParseWithStatementCloserErrors(t *testing.T) {
	t.Parallel()

	var (
		expectedValid  = false
		expectedErrors = strings.TrimLeft(strings.ReplaceAll(`
		Syntax error on line 4, column 8: unexpected token 'opencubr', should be 'id'
		Syntax error on line 9, column 16: unexpected token 'id', should be 'inherits', or 'opencubr'
		Syntax error on line 12, column 2: unexpected token 'let', should be 'closecubr', 'private', or 'public'
		Syntax error on line 18, column 2: unexpected token 'func', should be 'closecubr', 'private', or 'public'
		Syntax error on line 41, column 3: unexpected token 'opencubr', should be 'arrow'
		`, "\t", ""), "\n")
	)

	src := strings.TrimLeft(testutils.POLYNOMIAL_WITH_ERRORS_2_SRC, "\n")
	errc, close, out := errSpool()
	par := tabledrivenparser.NewParserNoComments(
		tabledrivenscanner.NewScanner(
			chuggingcharsource.MustChuggingReader(bytes.NewBufferString(src)),
			scannertable.TABLE()),
		parsertable.TABLE(),
		func(e *tabledrivenparser.ParserError) { errc <- *e },
		nil, token.Comments()...)

	valid := par.Parse()
	close()

	// Assert that the parse is not valid
	if expected, actual := expectedValid, valid; expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	// Assert errors
	if expected, actual := expectedErrors, out(); expected != actual {
		t.Errorf("\nExpected error output:\n%v\nBut got:\n%v", expected, actual)
	}
}

// Tests parsing the `polynomial-with-errors-2.src` file
func TestParsePolynomialWithErrors2Src(t *testing.T) {
	t.Parallel()
	assertParse(t, testutils.POLYNOMIAL_WITH_ERRORS_2_SRC, false)
}

// Tests parsing the `polynomial-with-errors.src` file
func TestParsePolynomialWithErrorsSrc(t *testing.T) {
	t.Parallel()
	assertParse(t, testutils.POLYNOMIAL_WITH_ERRORS_SRC, false)
}

// Tests parsing the `polynomial.src` file
func TestParsePolynomialSrc(t *testing.T) {
	t.Parallel()
	assertParse(t, testutils.POLYNOMIAL_SRC, true)
}

// Tests parsing the `bubblesort.src` file
func TestParseBubbleSortSrc(t *testing.T) {
	t.Parallel()
	assertParse(t, testutils.BUBBLESORT_SRC, true)
}

// Tests the ability to lex nested comments
func TestImbricatedComments(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		input  string
		tokens func(string) []token.Token
	}{
		{
			name: "simple imbricated comment",
			input: `/* this is an imbricated
			/* block comment
			*/
			*/`,
			tokens: func(input string) []token.Token {
				return []token.Token{{
					Id:     token.BLOCKCMT,
					Lexeme: token.Lexeme(input),
					Line:   1,
					Column: 1,
				}}
			},
		},
		{
			name: "imbricated minefield",
			input: `
			/* this is an imbricated
			/* block comment
			*/
			*/
			struct abc {
				field1: integer;
				public func do_something(x: float, y: integer) -> float;
			};

			/*
			/*
			/*
					TRIPLE IMBRICATED!!!!
			*/
			*/
			*/

			// ====== struct implementations ====== //
			impl POLYNOMIAL {
				func evaluate(x: float) -> float {
					return (0);
				}
			}
			/*
			/*
														/*
				/*
				/*
						/*
					OMG OVERKILL!!!!
									*/
							*/
					*/
					*/
							*/
			*/
			`,
			tokens: func(s string) []token.Token {
				return []token.Token{
					{
						Id:     token.BLOCKCMT,
						Lexeme: "/* this is an imbricated\n\t\t\t/* block comment\n\t\t\t*/\n\t\t\t*/",
						Line:   2,
						Column: 4,
					},
					{
						Id:     token.STRUCT,
						Lexeme: "struct",
						Line:   6,
						Column: 4,
					},
					{
						Id:     token.ID,
						Lexeme: "abc",
						Line:   6,
						Column: 11,
					},
					{
						Id:     token.OPENCUBR,
						Lexeme: "{",
						Line:   6,
						Column: 15,
					},
					{
						Id:     token.ID,
						Lexeme: "field1",
						Line:   7,
						Column: 5,
					},
					{
						Id:     token.COLON,
						Lexeme: ":",
						Line:   7,
						Column: 11,
					},
					{
						Id:     token.INTEGER,
						Lexeme: "integer",
						Line:   7,
						Column: 13,
					},
					{
						Id:     token.SEMI,
						Lexeme: ";",
						Line:   7,
						Column: 20,
					},
					{
						Id:     token.PUBLIC,
						Lexeme: "public",
						Line:   8,
						Column: 5,
					},
					{
						Id:     token.FUNC,
						Lexeme: "func",
						Line:   8,
						Column: 12,
					},
					{
						Id:     token.ID,
						Lexeme: "do_something",
						Line:   8,
						Column: 17,
					},
					{
						Id:     token.OPENPAR,
						Lexeme: "(",
						Line:   8,
						Column: 29,
					},
					{
						Id:     token.ID,
						Lexeme: "x",
						Line:   8,
						Column: 30,
					},
					{
						Id:     token.COLON,
						Lexeme: ":",
						Line:   8,
						Column: 31,
					},
					{
						Id:     token.FLOAT,
						Lexeme: "float",
						Line:   8,
						Column: 33,
					},
					{
						Id:     token.COMMA,
						Lexeme: ",",
						Line:   8,
						Column: 38,
					},
					{
						Id:     token.ID,
						Lexeme: "y",
						Line:   8,
						Column: 40,
					},
					{
						Id:     token.COLON,
						Lexeme: ":",
						Line:   8,
						Column: 41,
					},
					{
						Id:     token.INTEGER,
						Lexeme: "integer",
						Line:   8,
						Column: 43,
					},
					{
						Id:     token.CLOSEPAR,
						Lexeme: ")",
						Line:   8,
						Column: 50,
					},
					{
						Id:     token.ARROW,
						Lexeme: "->",
						Line:   8,
						Column: 52,
					},
					{
						Id:     token.FLOAT,
						Lexeme: "float",
						Line:   8,
						Column: 55,
					},
					{
						Id:     token.SEMI,
						Lexeme: ";",
						Line:   8,
						Column: 60,
					},
					{
						Id:     token.CLOSECUBR,
						Lexeme: "}",
						Line:   9,
						Column: 4,
					},
					{
						Id:     token.SEMI,
						Lexeme: ";",
						Line:   9,
						Column: 5,
					},
					{
						Id:     token.BLOCKCMT,
						Lexeme: "/*\n\t\t\t/*\n\t\t\t/*\n\t\t\t\t\tTRIPLE IMBRICATED!!!!\n\t\t\t*/\n\t\t\t*/\n\t\t\t*/",
						Line:   11,
						Column: 4,
					},
					{
						Id:     token.INLINECMT,
						Lexeme: "// ====== struct implementations ====== //\n",
						Line:   19,
						Column: 4,
					},
					{
						Id:     token.IMPL,
						Lexeme: "impl",
						Line:   20,
						Column: 4,
					},
					{
						Id:     token.ID,
						Lexeme: "POLYNOMIAL",
						Line:   20,
						Column: 9,
					},
					{
						Id:     token.OPENCUBR,
						Lexeme: "{",
						Line:   20,
						Column: 20,
					},
					{
						Id:     token.FUNC,
						Lexeme: "func",
						Line:   21,
						Column: 5,
					},
					{
						Id:     token.ID,
						Lexeme: "evaluate",
						Line:   21,
						Column: 10,
					},
					{
						Id:     token.OPENPAR,
						Lexeme: "(",
						Line:   21,
						Column: 18,
					},
					{
						Id:     token.ID,
						Lexeme: "x",
						Line:   21,
						Column: 19,
					},
					{
						Id:     token.COLON,
						Lexeme: ":",
						Line:   21,
						Column: 20,
					},
					{
						Id:     token.FLOAT,
						Lexeme: "float",
						Line:   21,
						Column: 22,
					},
					{
						Id:     token.CLOSEPAR,
						Lexeme: ")",
						Line:   21,
						Column: 27,
					},
					{
						Id:     token.ARROW,
						Lexeme: "->",
						Line:   21,
						Column: 29,
					},
					{
						Id:     token.FLOAT,
						Lexeme: "float",
						Line:   21,
						Column: 32,
					},
					{
						Id:     token.OPENCUBR,
						Lexeme: "{",
						Line:   21,
						Column: 38,
					},
					{
						Id:     token.RETURN,
						Lexeme: "return",
						Line:   22,
						Column: 6,
					},
					{
						Id:     token.OPENPAR,
						Lexeme: "(",
						Line:   22,
						Column: 13,
					},
					{
						Id:     token.INTNUM,
						Lexeme: "0",
						Line:   22,
						Column: 14,
					},
					{
						Id:     token.CLOSEPAR,
						Lexeme: ")",
						Line:   22,
						Column: 15,
					},
					{
						Id:     token.SEMI,
						Lexeme: ";",
						Line:   22,
						Column: 16,
					},
					{
						Id:     token.CLOSECUBR,
						Lexeme: "}",
						Line:   23,
						Column: 5,
					},
					{
						Id:     token.CLOSECUBR,
						Lexeme: "}",
						Line:   24,
						Column: 4,
					},
					{
						Id:     token.BLOCKCMT,
						Lexeme: "/*\n\t\t\t/*\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t/*\n\t\t\t\t/*\n\t\t\t\t/*\n\t\t\t\t\t\t/*\n\t\t\t\t\tOMG OVERKILL!!!!\n\t\t\t\t\t\t\t\t\t*/\n\t\t\t\t\t\t\t*/\n\t\t\t\t\t*/\n\t\t\t\t\t*/\n\t\t\t\t\t\t\t*/\n\t\t\t*/",
						Line:   25,
						Column: 4,
					},
				}
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := createScanner(t, tc.input).Tokens()
			expected := tc.tokens(tc.input)

			// Compare expected and actual token by token
			var i int
			for ; i < len(expected) && i < len(actual); i++ {
				a, e := actual[i], expected[i]
				if a != e {
					t.Errorf("Expected token #%v to be %v but got %v", i+1, e, a)
				}
			}

			// Report on any tokens not iterated on
			switch {
			case i < len(actual):
				if notSeen := actual[i:]; len(notSeen) > 0 {
					t.Errorf("The following tokens were not expected: %v", notSeen)
				}
			case i < len(expected):
				if notSeen := expected[i:]; len(notSeen) > 0 {
					t.Errorf("The following tokens were expected but not found: %v", notSeen)
				}
			}
		})
	}
}

func TestFloatIdNewlineIdId(t *testing.T) {
	t.Parallel()

	tokens := []token.Token{
		{
			Id:     token.FLOATNUM,
			Lexeme: "1.0",
			Line:   1,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "example_id",
			Line:   1,
			Column: 5,
		},
		{
			Id:     token.ID,
			Lexeme: "Id2",
			Line:   2,
			Column: 2,
		},
		{
			Id:     token.ID,
			Lexeme: "ID3",
			Line:   2,
			Column: 6,
		},
	}

	src := "1.0 example_id\n Id2 ID3"
	s := createScanner(t, src)

	for _, expected := range tokens {
		actual := mustNextToken(t, s)
		if actual != expected {
			t.Errorf("Expected token %v but got %v", expected, actual)
		}
	}
}

func TestFloatIdIdId(t *testing.T) {
	t.Parallel()

	tokens := []token.Token{
		{
			Id:     token.FLOATNUM,
			Lexeme: "1.0",
			Line:   1,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "example_id",
			Line:   1,
			Column: 5,
		},
		{
			Id:     token.ID,
			Lexeme: "Id2",
			Line:   1,
			Column: 16,
		},
		{
			Id:     token.ID,
			Lexeme: "ID3",
			Line:   1,
			Column: 20,
		},
	}

	var src string
	for i, token := range tokens {
		if i > 0 {
			src += " "
		}
		src += string(token.Lexeme)
	}

	s := createScanner(t, src)

	for _, expected := range tokens {
		actual := mustNextToken(t, s)
		if actual != expected {
			t.Errorf("Expected token %v but got %v", expected, actual)
		}
	}
}

func TestDoubleBackup(t *testing.T) {
	t.Parallel()

	expectedToken1 := token.Token{
		Id:     token.FLOATNUM,
		Lexeme: "1.0",
		Line:   1,
		Column: 1,
	}

	expectedToken2 := token.Token{
		Id:     token.ID,
		Lexeme: "example_id",
		Line:   1,
		Column: 4,
	}

	src := string(expectedToken1.Lexeme + expectedToken2.Lexeme)

	// TEST
	t.Run(src, func(t *testing.T) {
		t.Parallel()
		s := createScanner(t, src)

		actualToken1 := mustNextToken(t, s)
		actualToken2 := mustNextToken(t, s)

		if actualToken1 != expectedToken1 {
			t.Errorf("Expected first token to be %v but got %v", expectedToken1, actualToken1)
		}

		if actualToken2 != expectedToken2 {
			t.Errorf("Expected second token to be %v but got %v", expectedToken2, actualToken2)
		}
	})
}

// Tests the prof-provided src data from `lexnegativegrading.src`. Note that the
// data has been inlined into this test file to reduce external dependencies
func TestLexNegativeGrading(t *testing.T) {
	t.Parallel()
	s := createScanner(t, testutils.LEX_NEGATIVE_GRADING_SRC)
	for i, expected := range []token.Token{
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "@",
			Line:   1,
			Column: 1,
		},
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "#",
			Line:   1,
			Column: 3,
		},
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "$",
			Line:   1,
			Column: 5,
		},
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "'",
			Line:   1,
			Column: 7,
		},
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "\\",
			Line:   1,
			Column: 9,
		},
		{
			Id:     token.INVALIDCHAR,
			Lexeme: "~",
			Line:   1,
			Column: 11,
		},

		{
			Id:     token.INVALIDNUM,
			Lexeme: "00",
			Line:   3,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "01",
			Line:   4,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "010",
			Line:   5,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "0120",
			Line:   6,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "01230",
			Line:   7,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "0123450",
			Line:   8,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "01.23",
			Line:   10,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "012.34",
			Line:   11,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "12.340",
			Line:   12,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "012.340",
			Line:   13,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "012.34e10",
			Line:   15,
			Column: 1,
		},
		{
			Id:     token.INVALIDNUM,
			Lexeme: "12.34e010",
			Line:   16,
			Column: 1,
		},

		{
			Id:     token.INVALIDID,
			Lexeme: "_abc",
			Line:   18,
			Column: 1,
		},
		{
			Id:     token.INVALIDID,
			Lexeme: "1abc",
			Line:   19,
			Column: 1,
		},
		{
			Id:     token.INVALIDID,
			Lexeme: "_1abc",
			Line:   20,
			Column: 1,
		},
	} {
		t.Run(fmt.Sprintf("Token-%v[%v]", i+1, expected), func(t *testing.T) {
			expected.Line++ // We added a newline compared to the real file
			actual := mustNextToken(t, s)
			if actual != expected {
				t.Errorf("Expected token %v but got %v", expected, actual)
			}
		})
	}
}

// Tests the prof-provided src data from `lexpositivegrading.src`. Note that the
// data has been inlined into this test file to reduce external dependencies
func TestLexPositiveGrading(t *testing.T) {
	t.Parallel()
	s := createScanner(t, testutils.LEX_POSITIVE_GRADING_SRC)
	for i, expected := range []token.Token{
		{
			Id:     token.EQ,
			Lexeme: "==",
			Line:   1,
			Column: 1,
		},
		{
			Id:     token.PLUS,
			Lexeme: "+",
			Line:   1,
			Column: 4,
		},
		{
			Id:     token.OR,
			Lexeme: "|",
			Line:   1,
			Column: 6,
		},
		{
			Id:     token.OPENPAR,
			Lexeme: "(",
			Line:   1,
			Column: 8,
		},
		{
			Id:     token.SEMI,
			Lexeme: ";",
			Line:   1,
			Column: 10,
		},
		{
			Id:     token.IF,
			Lexeme: "if",
			Line:   1,
			Column: 12,
		},
		{
			Id:     token.PUBLIC,
			Lexeme: "public",
			Line:   1,
			Column: 16,
		},
		{
			Id:     token.READ,
			Lexeme: "read",
			Line:   1,
			Column: 23,
		},
		{
			Id:     token.NOTEQ,
			Lexeme: "<>",
			Line:   2,
			Column: 1,
		},
		{
			Id:     token.MINUS,
			Lexeme: "-",
			Line:   2,
			Column: 4,
		},
		{
			Id:     token.AND,
			Lexeme: "&",
			Line:   2,
			Column: 6,
		},
		{
			Id:     token.CLOSEPAR,
			Lexeme: ")",
			Line:   2,
			Column: 8,
		},
		{
			Id:     token.COMMA,
			Lexeme: ",",
			Line:   2,
			Column: 10,
		},
		{
			Id:     token.THEN,
			Lexeme: "then",
			Line:   2,
			Column: 12,
		},
		{
			Id:     token.PRIVATE,
			Lexeme: "private",
			Line:   2,
			Column: 17,
		},
		{
			Id:     token.WRITE,
			Lexeme: "write",
			Line:   2,
			Column: 25,
		},
		{
			Id:     token.LT,
			Lexeme: "<",
			Line:   3,
			Column: 1,
		},
		{
			Id:     token.MULT,
			Lexeme: "*",
			Line:   3,
			Column: 3,
		},
		{
			Id:     token.NOT,
			Lexeme: "!",
			Line:   3,
			Column: 5,
		},
		{
			Id:     token.OPENCUBR,
			Lexeme: "{",
			Line:   3,
			Column: 7,
		},
		{
			Id:     token.DOT,
			Lexeme: ".",
			Line:   3,
			Column: 9,
		},
		{
			Id:     token.ELSE,
			Lexeme: "else",
			Line:   3,
			Column: 11,
		},
		{
			Id:     token.FUNC,
			Lexeme: "func",
			Line:   3,
			Column: 16,
		},
		{
			Id:     token.RETURN,
			Lexeme: "return",
			Line:   3,
			Column: 21,
		},
		{
			Id:     token.GT,
			Lexeme: ">",
			Line:   4,
			Column: 1,
		},
		{
			Id:     token.DIV,
			Lexeme: "/",
			Line:   4,
			Column: 3,
		},
		{
			Id:     token.CLOSECUBR,
			Lexeme: "}",
			Line:   4,
			Column: 6,
		},
		{
			Id:     token.COLON,
			Lexeme: ":",
			Line:   4,
			Column: 8,
		},
		{
			Id:     token.INTEGER,
			Lexeme: "integer",
			Line:   4,
			Column: 10,
		},
		{
			Id:     token.VAR,
			Lexeme: "var",
			Line:   4,
			Column: 18,
		},
		{
			Id:     token.SELF,
			Lexeme: "self",
			Line:   4,
			Column: 22,
		},
		{
			Id:     token.LEQ,
			Lexeme: "<=",
			Line:   5,
			Column: 1,
		},
		{
			Id:     token.ASSIGN,
			Lexeme: "=",
			Line:   5,
			Column: 4,
		},
		{
			Id:     token.OPENSQBR,
			Lexeme: "[",
			Line:   5,
			Column: 7,
		},
		{
			Id:     token.COLONCOLON,
			Lexeme: "::",
			Line:   5,
			Column: 9,
		},
		{
			Id:     token.FLOAT,
			Lexeme: "float",
			Line:   5,
			Column: 12,
		},
		{
			Id:     token.STRUCT,
			Lexeme: "struct",
			Line:   5,
			Column: 18,
		},
		{
			Id:     token.INHERITS,
			Lexeme: "inherits",
			Line:   5,
			Column: 25,
		},
		{
			Id:     token.GEQ,
			Lexeme: ">=",
			Line:   6,
			Column: 1,
		},
		{
			Id:     token.CLOSESQBR,
			Lexeme: "]",
			Line:   6,
			Column: 6,
		},
		{
			Id:     token.ARROW,
			Lexeme: "->",
			Line:   6,
			Column: 8,
		},
		{
			Id:     token.VOID,
			Lexeme: "void",
			Line:   6,
			Column: 11,
		},
		{
			Id:     token.WHILE,
			Lexeme: "while",
			Line:   6,
			Column: 16,
		},
		{
			Id:     token.LET,
			Lexeme: "let",
			Line:   6,
			Column: 22,
		},
		{
			Id:     token.FUNC,
			Lexeme: "func",
			Line:   7,
			Column: 7,
		},
		{
			Id:     token.IMPL,
			Lexeme: "impl",
			Line:   7,
			Column: 12,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "0",
			Line:   13,
			Column: 1,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "1",
			Line:   14,
			Column: 1,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "10",
			Line:   15,
			Column: 1,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "12",
			Line:   16,
			Column: 1,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "123",
			Line:   17,
			Column: 1,
		},
		{
			Id:     token.INTNUM,
			Lexeme: "12345",
			Line:   18,
			Column: 1,
		},
		{
			Id:     token.FLOATNUM,
			Lexeme: "1.23",
			Line:   20,
			Column: 1,
		},
		{
			Id:     token.FLOATNUM,
			Lexeme: "12.34",
			Line:   21,
			Column: 1,
		},
		{
			Id:     token.FLOATNUM,
			Lexeme: "120.34e10",
			Line:   22,
			Column: 1,
		},
		{
			Id:     token.FLOATNUM,
			Lexeme: "12345.6789e-123",
			Line:   23,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "abc",
			Line:   25,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "abc1",
			Line:   26,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "a1bc",
			Line:   27,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "abc_1abc",
			Line:   28,
			Column: 1,
		},
		{
			Id:     token.ID,
			Lexeme: "abc1_abc",
			Line:   29,
			Column: 1,
		},
		{
			Id:     token.INLINECMT,
			Lexeme: "// this is an inline comment\n",
			Line:   31,
			Column: 1,
		},
		{
			Id:     token.BLOCKCMT,
			Lexeme: "/* this is a single line block comment */",
			Line:   33,
			Column: 1,
		},
		{
			Id:     token.BLOCKCMT,
			Lexeme: "/* this is a\nmultiple line\nblock comment\n*/",
			Line:   35,
			Column: 1,
		},
		{
			Id:     token.BLOCKCMT,
			Lexeme: "/* this is an imbricated\n/* block comment\n*/\n*/",
			Line:   40,
			Column: 1,
		},
	} {
		t.Run(fmt.Sprintf("Token-%v[%v]", i+1, expected.Lexeme), func(t *testing.T) {
			expected.Line++ // We added a newline compared to the real file
			actual := mustNextToken(t, s)
			if actual != expected {
				t.Errorf("Expected token %v but got %v", expected, actual)
			}
		})
	}
}

// Tests a single scan on inputs containing only one token
func TestSingleScans(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   token.Kind
		input  string
		output token.Token
	}{
		{
			name:  token.INLINECMT,
			input: "// asdasdasd \n",
			output: token.Token{
				Id:     token.INLINECMT,
				Lexeme: "// asdasdasd \n",
			},
		},
		{
			name:  token.BLOCKCMT,
			input: "/* \n asdasd \n asdasd \n -123=*/",
			output: token.Token{
				Id:     token.BLOCKCMT,
				Lexeme: "/* \n asdasd \n asdasd \n -123=*/",
			},
		},
		{
			name:  token.DIV,
			input: "/",
			output: token.Token{
				Id:     token.DIV,
				Lexeme: "/",
			},
		},
		{
			name:  token.MULT,
			input: "*",
			output: token.Token{
				Id:     token.MULT,
				Lexeme: "*",
			},
		},
		{
			name:  token.MINUS,
			input: "-",
			output: token.Token{
				Id:     token.MINUS,
				Lexeme: "-",
			},
		},
		{
			name:  token.ARROW,
			input: "->",
			output: token.Token{
				Id:     token.ARROW,
				Lexeme: "->",
			},
		},
		{
			name:  token.ASSIGN,
			input: "=",
			output: token.Token{
				Id:     token.ASSIGN,
				Lexeme: "=",
			},
		},
		{
			name:  token.EQ,
			input: "==",
			output: token.Token{
				Id:     token.EQ,
				Lexeme: "==",
			},
		},
		{
			name:  token.LT,
			input: "<",
			output: token.Token{
				Id:     token.LT,
				Lexeme: "<",
			},
		},
		{
			name:  token.GT,
			input: ">",
			output: token.Token{
				Id:     token.GT,
				Lexeme: ">",
			},
		},
		{
			name:  token.NOTEQ,
			input: "<>",
			output: token.Token{
				Id:     token.NOTEQ,
				Lexeme: "<>",
			},
		},
		{
			name:  token.LEQ,
			input: "<=",
			output: token.Token{
				Id:     token.LEQ,
				Lexeme: "<=",
			},
		},
		{
			name:  token.GEQ,
			input: ">=",
			output: token.Token{
				Id:     token.GEQ,
				Lexeme: ">=",
			},
		},
		{
			name:  token.OR,
			input: "|",
			output: token.Token{
				Id:     token.OR,
				Lexeme: "|",
			},
		},
		{
			name:  token.AND,
			input: "&",
			output: token.Token{
				Id:     token.AND,
				Lexeme: "&",
			},
		},
		{
			name:  token.NOT,
			input: "!",
			output: token.Token{
				Id:     token.NOT,
				Lexeme: "!",
			},
		},

		{
			name:  token.OPENPAR,
			input: "(",
			output: token.Token{
				Id:     token.OPENPAR,
				Lexeme: "(",
			},
		},
		{
			name:  token.CLOSEPAR,
			input: ")",
			output: token.Token{
				Id:     token.CLOSEPAR,
				Lexeme: ")",
			},
		},
		{
			name:  token.OPENSQBR,
			input: "[",
			output: token.Token{
				Id:     token.OPENSQBR,
				Lexeme: "[",
			},
		},
		{
			name:  token.CLOSESQBR,
			input: "]",
			output: token.Token{
				Id:     token.CLOSESQBR,
				Lexeme: "]",
			},
		},
		{
			name:  token.OPENCUBR,
			input: "{",
			output: token.Token{
				Id:     token.OPENCUBR,
				Lexeme: "{",
			},
		},
		{
			name:  token.CLOSECUBR,
			input: "}",
			output: token.Token{
				Id:     token.CLOSECUBR,
				Lexeme: "}",
			},
		},
		{
			name:  token.DOT,
			input: ".",
			output: token.Token{
				Id:     token.DOT,
				Lexeme: ".",
			},
		},
		{
			name:  token.COMMA,
			input: ",",
			output: token.Token{
				Id:     token.COMMA,
				Lexeme: ",",
			},
		},
		{
			name:  token.SEMI,
			input: ";",
			output: token.Token{
				Id:     token.SEMI,
				Lexeme: ";",
			},
		},
		{
			name:  token.COLON,
			input: ":",
			output: token.Token{
				Id:     token.COLON,
				Lexeme: ":",
			},
		},
		{
			name:  token.COLONCOLON,
			input: "::",
			output: token.Token{
				Id:     token.COLONCOLON,
				Lexeme: "::",
			},
		},
		{
			name:  token.IF,
			input: "if",
			output: token.Token{
				Id:     token.IF,
				Lexeme: "if",
			},
		},
		{
			name:  token.THEN,
			input: "then",
			output: token.Token{
				Id:     token.THEN,
				Lexeme: "then",
			},
		},
		{
			name:  token.ELSE,
			input: "else",
			output: token.Token{
				Id:     token.ELSE,
				Lexeme: "else",
			},
		},
		{
			name:  token.INTEGER,
			input: "integer",
			output: token.Token{
				Id:     token.INTEGER,
				Lexeme: "integer",
			},
		},
		{
			name:  token.FLOAT,
			input: "float",
			output: token.Token{
				Id:     token.FLOAT,
				Lexeme: "float",
			},
		},
		{
			name:  token.VOID,
			input: "void",
			output: token.Token{
				Id:     token.VOID,
				Lexeme: "void",
			},
		},
		{
			name:  token.PUBLIC,
			input: "public",
			output: token.Token{
				Id:     token.PUBLIC,
				Lexeme: "public",
			},
		},
		{
			name:  token.PRIVATE,
			input: "private",
			output: token.Token{
				Id:     token.PRIVATE,
				Lexeme: "private",
			},
		},
		{
			name:  token.FUNC,
			input: "func",
			output: token.Token{
				Id:     token.FUNC,
				Lexeme: "func",
			},
		},
		{
			name:  token.VAR,
			input: "var",
			output: token.Token{
				Id:     token.VAR,
				Lexeme: "var",
			},
		},
		{
			name:  token.STRUCT,
			input: "struct",
			output: token.Token{
				Id:     token.STRUCT,
				Lexeme: "struct",
			},
		},
		{
			name:  token.WHILE,
			input: "while",
			output: token.Token{
				Id:     token.WHILE,
				Lexeme: "while",
			},
		},
		{
			name:  token.READ,
			input: "read",
			output: token.Token{
				Id:     token.READ,
				Lexeme: "read",
			},
		},
		{
			name:  token.WRITE,
			input: "write",
			output: token.Token{
				Id:     token.WRITE,
				Lexeme: "write",
			},
		},
		{
			name:  token.RETURN,
			input: "return",
			output: token.Token{
				Id:     token.RETURN,
				Lexeme: "return",
			},
		},
		{
			name:  token.SELF,
			input: "self",
			output: token.Token{
				Id:     token.SELF,
				Lexeme: "self",
			},
		},
		{
			name:  token.INHERITS,
			input: "inherits",
			output: token.Token{
				Id:     token.INHERITS,
				Lexeme: "inherits",
			},
		},
		{
			name:  token.LET,
			input: "let",
			output: token.Token{
				Id:     token.LET,
				Lexeme: "let",
			},
		},
		{
			name:  token.IMPL,
			input: "impl",
			output: token.Token{
				Id:     token.IMPL,
				Lexeme: "impl",
			},
		},
		{
			name:  token.ID,
			input: "asd",
			output: token.Token{
				Id:     token.ID,
				Lexeme: "asd",
			},
		},
		{
			name:  token.INTNUM,
			input: "99",
			output: token.Token{
				Id:     token.INTNUM,
				Lexeme: "99",
			},
		},
		{
			name:  token.FLOATNUM,
			input: "1.0",
			output: token.Token{
				Id:     token.FLOATNUM,
				Lexeme: "1.0",
			},
		},
		{
			name:  token.INVALIDNUM + "[00]",
			input: "00",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "00",
			},
		},
		{
			name:  token.INVALIDNUM + "[01]",
			input: "01",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "01",
			},
		},
		{
			name:  token.INVALIDNUM + "[010]",
			input: "010",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "010",
			},
		},
		{
			name:  token.INVALIDNUM + "[0120]",
			input: "0120",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "0120",
			},
		},
		{
			name:  token.INVALIDNUM + "[01230]",
			input: "01230",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "01230",
			},
		},
		{
			name:  token.INVALIDNUM + "[0123450]",
			input: "0123450",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "0123450",
			},
		},
		{
			name:  token.INVALIDNUM + "[01.23]",
			input: "01.23",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "01.23",
			},
		},
		{
			name:  token.INVALIDNUM + "[012.34]",
			input: "012.34",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "012.34",
			},
		},
		{
			name:  token.INVALIDNUM + "[012.340]",
			input: "012.340",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "012.340",
			},
		},
		{
			name:  token.INVALIDNUM + "[012.34e10]",
			input: "012.34e10",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "012.34e10",
			},
		},
		{
			name:  token.INVALIDNUM + "[12.34e010]",
			input: "12.34e010",
			output: token.Token{
				Id:     token.INVALIDNUM,
				Lexeme: "12.34e010",
			},
		},
		{
			name:  token.INVALIDID + "[_abc]",
			input: "_abc",
			output: token.Token{
				Id:     token.INVALIDID,
				Lexeme: "_abc",
			},
		},
		{
			name:  token.INVALIDID + "[1abc]",
			input: "1abc",
			output: token.Token{
				Id:     token.INVALIDID,
				Lexeme: "1abc",
			},
		},
		{
			name:  token.INVALIDID + "[_1abc]",
			input: "_1abc",
			output: token.Token{
				Id:     token.INVALIDID,
				Lexeme: "_1abc",
			},
		},
		{
			name:  token.INVALIDCHAR + "[@]",
			input: "@",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "@",
			},
		},
		{
			name:  token.INVALIDCHAR + "[#]",
			input: "#",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "#",
			},
		},
		{
			name:  token.INVALIDCHAR + "[$]",
			input: "$",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "$",
			},
		},
		{
			name:  token.INVALIDCHAR + "[']",
			input: "'",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "'",
			},
		},
		{
			name:  token.INVALIDCHAR + "[\\]",
			input: "\\",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "\\",
			},
		},
		{
			name:  token.INVALIDCHAR + "[~]",
			input: "~",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "~",
			},
		},
		{
			name:  token.INVALIDCHAR + "[']",
			input: "'",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "'",
			},
		},
		{
			name:  token.INVALIDCHAR + "[%]",
			input: "%",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "%",
			},
		},
		{
			name:  token.INVALIDCHAR + "[^]",
			input: "^",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "^",
			},
		},
		{
			name:  token.INVALIDCHAR + "[`]",
			input: "`",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "`",
			},
		},
		{
			name:  token.INVALIDCHAR + "[\"]",
			input: "\"",
			output: token.Token{
				Id:     token.INVALIDCHAR,
				Lexeme: "\"",
			},
		},
	} {
		tc := tc
		t.Run(string(tc.name), func(t *testing.T) {
			t.Parallel()
			tc.output.Line = 1
			tc.output.Column = 1
			assertScan(t, tc.input, tc.output)
		})
	}
}

// Grabs the next token.Token from the scanner and asserts that there were no
// errors. Returns the token. If scanner.NextToken() returns an error, this
// function will immediately fail your test with the error.
func mustNextToken(t *testing.T, s scanner.Scanner) token.Token {
	tt, err := s.NextToken()
	if err != nil {
		t.Fatalf("NextToken should succeed: %v", err)
	}
	return tt
}

// Asserts the first token present in the input. Note that this function creates
// a new scanner everytime containing the provided input everytime you call it.
// It is used for testing inputs that contain a single token.
func assertScan(t *testing.T, input string, expected token.Token) {
	s := createScanner(t, input)
	if actual := mustNextToken(t, s); actual != expected {
		t.Fatalf("Expected token %v but got %v", expected, actual)
	}
}

// Creates a scanner with a char source containing the provided contents
func createScanner(t *testing.T, contents string) scanner.LoadableScanner {
	chars := createCharSource(t, contents)
	return scanner.NewLoadableScanner(tabledrivenscanner.NewScanner(chars, table()))
}

// Creates a char source containing the provided contents
func createCharSource(t *testing.T, contents string) *chuggingcharsource.ChuggingCharSource {
	chars := new(chuggingcharsource.ChuggingCharSource)
	err := chars.ChugReader(bytes.NewBufferString(contents))
	if err != nil {
		t.Fatalf("ChugReader should succeed here: %v", err)
	}
	return chars
}

func createParser(contents string) *tabledrivenparser.TableDrivenParser {
	return tabledrivenparser.NewParserNoComments(
		tabledrivenscanner.NewScanner(
			chuggingcharsource.MustChuggingReader(bytes.NewBufferString(contents)),
			scannertable.TABLE()),
		parsertable.TABLE(), nil, nil, token.Comments()...)
}

// Creates a parser that consumes `contents` source code, returns the parser and
// a callback for retrieving the error log
func createErrorLoggingParser(contents string) (
	*tabledrivenparser.TableDrivenParser,
	func() string,
) {
	errs := make([]error, 0, 1024)
	prsr := tabledrivenparser.NewParserNoComments(
		tabledrivenscanner.NewScanner(
			chuggingcharsource.MustChuggingReader(bytes.NewBufferString(contents)),
			scannertable.TABLE()), parsertable.TABLE(),
		func(e *tabledrivenparser.ParserError) { errs = append(errs, e) },
		nil, token.Comments()...)

	return prsr, func() string {
		out := new(bytes.Buffer)
		for _, e := range errs {
			fmt.Fprintln(out, e)
		}
		return out.String()
	}
}

func assertParse(t *testing.T, contents string, result bool) {
	if createParser(contents).Parse() != result {
		t.Fatalf("Parse should succeed, but it returned false")
	}
}

func errSpool() (chan<- tabledrivenparser.ParserError, func(), func() string) {
	out := new(bytes.Buffer)
	errc := make(chan tabledrivenparser.ParserError, 1024)
	donec := make(chan struct{}, 1)
	go func() {
		for err := range errc {
			fmt.Fprintf(
				out,
				"Syntax error on line %v, column %v: %v\n",
				err.Tok.Line, err.Tok.Column, err.Err)
		}
		donec <- struct{}{}
	}()

	return errc, func() { close(errc) }, func() string {
		<-donec
		return out.String()
	}
}

func clean(in string, linePrefix string) string {
	return strings.TrimRight(
		strings.TrimSuffix(
			strings.TrimPrefix(
				testutils.TrimLeading(in, linePrefix),
				"\n"),
			"\n"),
		"\t")
}
