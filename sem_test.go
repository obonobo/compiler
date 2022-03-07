package main

import (
	"strings"
	"testing"

	"github.com/obonobo/esac/core/parser"
)

func TestStructDecl(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	struct Hey inherits Yo {
		private let a: float;
		public let b: integer;
		public func doIt(A: float, B: float) -> Yo;
		public func hey(x: float, y: integer) -> Hey;
	};

	struct Hey2 {
		private let a: float;
		public func doIt(A: float, B: float) -> Yo;
	};
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | StructDecl
	| | | Id: Token[Id=id, Lexeme=Hey, Line=2, Column=9]
	| | | Inherits
	| | | | Id: Token[Id=id, Lexeme=Yo, Line=2, Column=22]
	| | | Members
	| | | | Member
	| | | | | Private: Token[Id=private, Lexeme=private, Line=3, Column=3]
	| | | | | VarDecl
	| | | | | | Id: Token[Id=id, Lexeme=a, Line=3, Column=15]
	| | | | | | Type
	| | | | | | | Float: Token[Id=float, Lexeme=float, Line=3, Column=18]
	| | | | | | DimList
	| | | | Member
	| | | | | Public: Token[Id=public, Lexeme=public, Line=4, Column=3]
	| | | | | VarDecl
	| | | | | | Id: Token[Id=id, Lexeme=b, Line=4, Column=14]
	| | | | | | Type
	| | | | | | | Integer: Token[Id=integer, Lexeme=integer, Line=4, Column=17]
	| | | | | | DimList
	| | | | Member
	| | | | | Public: Token[Id=public, Lexeme=public, Line=5, Column=3]
	| | | | | FuncDecl
	| | | | | | Id: Token[Id=id, Lexeme=doIt, Line=5, Column=15]
	| | | | | | ParamList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=A, Line=5, Column=20]
	| | | | | | | | Type
	| | | | | | | | | Float: Token[Id=float, Lexeme=float, Line=5, Column=23]
	| | | | | | | | DimList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=B, Line=5, Column=30]
	| | | | | | | | Type
	| | | | | | | | | Float: Token[Id=float, Lexeme=float, Line=5, Column=33]
	| | | | | | | | DimList
	| | | | | | ReturnType
	| | | | | | | Id: Token[Id=id, Lexeme=Yo, Line=5, Column=43]
	| | | | Member
	| | | | | Public: Token[Id=public, Lexeme=public, Line=6, Column=3]
	| | | | | FuncDecl
	| | | | | | Id: Token[Id=id, Lexeme=hey, Line=6, Column=15]
	| | | | | | ParamList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=x, Line=6, Column=19]
	| | | | | | | | Type
	| | | | | | | | | Float: Token[Id=float, Lexeme=float, Line=6, Column=22]
	| | | | | | | | DimList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=y, Line=6, Column=29]
	| | | | | | | | Type
	| | | | | | | | | Integer: Token[Id=integer, Lexeme=integer, Line=6, Column=32]
	| | | | | | | | DimList
	| | | | | | ReturnType
	| | | | | | | Id: Token[Id=id, Lexeme=Hey, Line=6, Column=44]
	| | StructDecl
	| | | Id: Token[Id=id, Lexeme=Hey2, Line=9, Column=9]
	| | | Inherits
	| | | Members
	| | | | Member
	| | | | | Private: Token[Id=private, Lexeme=private, Line=10, Column=3]
	| | | | | VarDecl
	| | | | | | Id: Token[Id=id, Lexeme=a, Line=10, Column=15]
	| | | | | | Type
	| | | | | | | Float: Token[Id=float, Lexeme=float, Line=10, Column=18]
	| | | | | | DimList
	| | | | Member
	| | | | | Public: Token[Id=public, Lexeme=public, Line=11, Column=3]
	| | | | | FuncDecl
	| | | | | | Id: Token[Id=id, Lexeme=doIt, Line=11, Column=15]
	| | | | | | ParamList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=A, Line=11, Column=20]
	| | | | | | | | Type
	| | | | | | | | | Float: Token[Id=float, Lexeme=float, Line=11, Column=23]
	| | | | | | | | DimList
	| | | | | | | Param
	| | | | | | | | Id: Token[Id=id, Lexeme=B, Line=11, Column=30]
	| | | | | | | | Type
	| | | | | | | | | Float: Token[Id=float, Lexeme=float, Line=11, Column=33]
	| | | | | | | | DimList
	| | | | | | ReturnType
	| | | | | | | Id: Token[Id=id, Lexeme=Yo, Line=11, Column=43]
	`, "\n"), "\t", ""))
}

func TestImplDef(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	impl MyImplementation {
		func do_something(x: integer[2]) -> void {
			write(x);
		}

		func and_another_one() -> float {
			return (2.9);
		}
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | ImplDef
	| | | Id: Token[Id=id, Lexeme=MyImplementation, Line=2, Column=7]
	| | | FuncDefList
	| | | | FuncDef
	| | | | | Id: Token[Id=id, Lexeme=do_something, Line=3, Column=8]
	| | | | | ParamList
	| | | | | | Param
	| | | | | | | Id: Token[Id=id, Lexeme=x, Line=3, Column=21]
	| | | | | | | Type
	| | | | | | | | Integer: Token[Id=integer, Lexeme=integer, Line=3, Column=24]
	| | | | | | | DimList
	| | | | | | | | Dim: Token[Id=intnum, Lexeme=2, Line=3, Column=32]
	| | | | | ReturnType
	| | | | | | Void: Token[Id=void, Lexeme=void, Line=3, Column=39]
	| | | | | Body
	| | | | | | Write
	| | | | | | | ArithExpr
	| | | | | | | | Factor
	| | | | | | | | | Variable
	| | | | | | | | | | Subject
	| | | | | | | | | | Id: Token[Id=id, Lexeme=x, Line=4, Column=10]
	| | | | | | | | | | IndexList
	| | | | FuncDef
	| | | | | Id: Token[Id=id, Lexeme=and_another_one, Line=7, Column=8]
	| | | | | ParamList
	| | | | | ReturnType
	| | | | | | Float: Token[Id=float, Lexeme=float, Line=7, Column=29]
	| | | | | Body
	| | | | | | Return
	| | | | | | | ArithExpr
	| | | | | | | | Factor
	| | | | | | | | | FloatNum: Token[Id=floatnum, Lexeme=2.9, Line=8, Column=12]
	`, "\n"), "\t", ""))
}

func TestAssign(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other() -> void {
		id3 = 12;
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=18]
	| | | Body
	| | | | Assign(=)
	| | | | | Variable
	| | | | | | Subject
	| | | | | | Id: Token[Id=id, Lexeme=id3, Line=3, Column=3]
	| | | | | | IndexList
	| | | | | ArithExpr
	| | | | | | Factor
	| | | | | | | IntNum: Token[Id=intnum, Lexeme=12, Line=3, Column=9]
	`, "\n"), "\t", ""))
}

func TestAssignChained(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other() -> void {
		id1[1][2][3].id2(1).id3 = 12;
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=18]
	| | | Body
	| | | | Assign(=)
	| | | | | Variable
	| | | | | | Subject
	| | | | | | | FuncCall
	| | | | | | | | Subject
	| | | | | | | | | Variable
	| | | | | | | | | | Subject
	| | | | | | | | | | Id: Token[Id=id, Lexeme=id1, Line=3, Column=3]
	| | | | | | | | | | IndexList
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=7]
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=2, Line=3, Column=10]
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=3, Line=3, Column=13]
	| | | | | | | | Id: Token[Id=id, Lexeme=id2, Line=3, Column=16]
	| | | | | | | | FuncCallParamList
	| | | | | | | | | FuncCallParam
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=20]
	| | | | | | Id: Token[Id=id, Lexeme=id3, Line=3, Column=23]
	| | | | | | IndexList
	| | | | | ArithExpr
	| | | | | | Factor
	| | | | | | | IntNum: Token[Id=intnum, Lexeme=12, Line=3, Column=29]
	`, "\n"), "\t", ""))
}

func TestIfWhileAndAFewOtherThings(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other() -> void {
		// read(id1[1][2][3].id2(1).id3);
		if (1 < 2) then {
			write(1);
		} else {
			read(id1);
			while (1 == 1) {
				// noop
			};
		};
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=18]
	| | | Body
	| | | | If
	| | | | | RelExpr
	| | | | | | Lt(<)
	| | | | | | | ArithExpr
	| | | | | | | | Factor
	| | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=4, Column=7]
	| | | | | | | ArithExpr
	| | | | | | | | Factor
	| | | | | | | | | IntNum: Token[Id=intnum, Lexeme=2, Line=4, Column=11]
	| | | | | StatBlock
	| | | | | | Write
	| | | | | | | ArithExpr
	| | | | | | | | Factor
	| | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=5, Column=10]
	| | | | | StatBlock
	| | | | | | Read
	| | | | | | | Variable
	| | | | | | | | Subject
	| | | | | | | | Id: Token[Id=id, Lexeme=id1, Line=7, Column=9]
	| | | | | | | | IndexList
	| | | | | | While
	| | | | | | | RelExpr
	| | | | | | | | Eq(==)
	| | | | | | | | | ArithExpr
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=8, Column=11]
	| | | | | | | | | ArithExpr
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=8, Column=16]
	| | | | | | | StatBlock
	`, "\n\t"), "\t", ""))
}

func TestFunctionAndVariables(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other() -> void {
		read(id1[1][2][3].id2(1).id3);
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=18]
	| | | Body
	| | | | Read
	| | | | | Variable
	| | | | | | Subject
	| | | | | | | FuncCall
	| | | | | | | | Subject
	| | | | | | | | | Variable
	| | | | | | | | | | Subject
	| | | | | | | | | | Id: Token[Id=id, Lexeme=id1, Line=3, Column=8]
	| | | | | | | | | | IndexList
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=12]
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=2, Line=3, Column=15]
	| | | | | | | | | | | Index
	| | | | | | | | | | | | Factor
	| | | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=3, Line=3, Column=18]
	| | | | | | | | Id: Token[Id=id, Lexeme=id2, Line=3, Column=21]
	| | | | | | | | FuncCallParamList
	| | | | | | | | | FuncCallParam
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=25]
	| | | | | | Id: Token[Id=id, Lexeme=id3, Line=3, Column=28]
	| | | | | | IndexList
	`, "\n\t"), "\t", ""))
}

func TestArithExpr(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other() -> void {
		write(1 + 1 - 1 * 1 + 5 / 5 + 3 - 2);
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=18]
	| | | Body
	| | | | Write
	| | | | | ArithExpr
	| | | | | | Minus(-)
	| | | | | | | Plus(+)
	| | | | | | | | Plus(+)
	| | | | | | | | | Minus(-)
	| | | | | | | | | | Plus(+)
	| | | | | | | | | | | Factor
	| | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=9]
	| | | | | | | | | | | Factor
	| | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=13]
	| | | | | | | | | | Mult(*)
	| | | | | | | | | | | Factor
	| | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=17]
	| | | | | | | | | | | Factor
	| | | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=1, Line=3, Column=21]
	| | | | | | | | | Div(/)
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=5, Line=3, Column=25]
	| | | | | | | | | | Factor
	| | | | | | | | | | | IntNum: Token[Id=intnum, Lexeme=5, Line=3, Column=29]
	| | | | | | | | Factor
	| | | | | | | | | IntNum: Token[Id=intnum, Lexeme=3, Line=3, Column=33]
	| | | | | | | Factor
	| | | | | | | | IntNum: Token[Id=intnum, Lexeme=2, Line=3, Column=37]
	`, "\n\t"), "\t", ""))
}

func TestFunctionParams(t *testing.T) {
	t.Parallel()
	assertParseAndAst(t, `
	func other(x: integer[2][3][4], z: float, y: integer[2][3][4]) -> void {
		// noop
	}
	`, true, strings.ReplaceAll(strings.TrimLeft(`
	Prog
	| StructOrImplOrFuncList
	| | FuncDef
	| | | Id: Token[Id=id, Lexeme=other, Line=2, Column=7]
	| | | ParamList
	| | | | Param
	| | | | | Id: Token[Id=id, Lexeme=x, Line=2, Column=13]
	| | | | | Type
	| | | | | | Integer: Token[Id=integer, Lexeme=integer, Line=2, Column=16]
	| | | | | DimList
	| | | | | | Dim: Token[Id=intnum, Lexeme=2, Line=2, Column=24]
	| | | | | | Dim: Token[Id=intnum, Lexeme=3, Line=2, Column=27]
	| | | | | | Dim: Token[Id=intnum, Lexeme=4, Line=2, Column=30]
	| | | | Param
	| | | | | Id: Token[Id=id, Lexeme=z, Line=2, Column=34]
	| | | | | Type
	| | | | | | Float: Token[Id=float, Lexeme=float, Line=2, Column=37]
	| | | | | DimList
	| | | | Param
	| | | | | Id: Token[Id=id, Lexeme=y, Line=2, Column=44]
	| | | | | Type
	| | | | | | Integer: Token[Id=integer, Lexeme=integer, Line=2, Column=47]
	| | | | | DimList
	| | | | | | Dim: Token[Id=intnum, Lexeme=2, Line=2, Column=55]
	| | | | | | Dim: Token[Id=intnum, Lexeme=3, Line=2, Column=58]
	| | | | | | Dim: Token[Id=intnum, Lexeme=4, Line=2, Column=61]
	| | | ReturnType
	| | | | Void: Token[Id=void, Lexeme=void, Line=2, Column=68]
	| | | Body
	`, "\n\t"), "\t", ""))
}

func assertParseAndAst(t *testing.T, contents string, valid bool, ast string) {
	// Sanity check: TableDrivenParser should conform to parser.Parser interface
	var prsr parser.Parser = createParser(contents)

	// Assert parse
	if actual, expected := prsr.Parse(), valid; actual != expected {
		if valid {
			t.Errorf(
				"Expected parse to succeed, but it failed (expected: %v, got: %v)",
				expected, actual)
		} else {
			t.Errorf(
				"Expected parse to fail, but it succeeded (expected: %v, got: %v)",
				expected, actual)
		}
	}

	// Assert ast output
	if actual, expected := prsr.AST().TreeString(), ast; actual != expected {
		t.Errorf("\nExpected ast to be:\n%v\nbut got:\n%v", expected, actual)
	}
}
