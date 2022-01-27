package main

import (
	"bytes"
	"testing"

	"github.com/obonobo/compiler/core/chuggingcharsource"
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner/compositetable"
)

// Tests a single scan on inputs containing only one token
func TestSingleScans(t *testing.T) {
	for _, tc := range []struct {
		symbol scanner.Symbol
		input  string
		output scanner.Token
	}{
		{
			symbol: scanner.INLINECMT,
			input:  "// asdasdasd \n",
			output: scanner.Token{
				Id:     scanner.INLINECMT,
				Lexeme: "// asdasdasd \n",
			},
		},
		{
			symbol: scanner.BLOCKCMT,
			input:  "/* \n asdasd \n asdasd \n -123=*/",
			output: scanner.Token{
				Id:     scanner.BLOCKCMT,
				Lexeme: "/* \n asdasd \n asdasd \n -123=*/",
			},
		},
		{
			symbol: scanner.DIV,
			input:  "/",
			output: scanner.Token{
				Id:     scanner.DIV,
				Lexeme: "/",
			},
		},
		{
			symbol: scanner.MULT,
			input:  "*",
			output: scanner.Token{
				Id:     scanner.MULT,
				Lexeme: "*",
			},
		},
		{
			symbol: scanner.MINUS,
			input:  "-",
			output: scanner.Token{
				Id:     scanner.MINUS,
				Lexeme: "-",
			},
		},
		{
			symbol: scanner.ARROW,
			input:  "->",
			output: scanner.Token{
				Id:     scanner.ARROW,
				Lexeme: "->",
			},
		},
		{
			symbol: scanner.ASSIGN,
			input:  "=",
			output: scanner.Token{
				Id:     scanner.ASSIGN,
				Lexeme: "=",
			},
		},
		{
			symbol: scanner.EQ,
			input:  "==",
			output: scanner.Token{
				Id:     scanner.EQ,
				Lexeme: "==",
			},
		},
		{
			symbol: scanner.LT,
			input:  "<",
			output: scanner.Token{
				Id:     scanner.LT,
				Lexeme: "<",
			},
		},
		{
			symbol: scanner.GT,
			input:  ">",
			output: scanner.Token{
				Id:     scanner.GT,
				Lexeme: ">",
			},
		},
		{
			symbol: scanner.NOTEQ,
			input:  "<>",
			output: scanner.Token{
				Id:     scanner.NOTEQ,
				Lexeme: "<>",
			},
		},
		{
			symbol: scanner.LEQ,
			input:  "<=",
			output: scanner.Token{
				Id:     scanner.LEQ,
				Lexeme: "<=",
			},
		},
		{
			symbol: scanner.GEQ,
			input:  ">=",
			output: scanner.Token{
				Id:     scanner.GEQ,
				Lexeme: ">=",
			},
		},
		{
			symbol: scanner.OR,
			input:  "|",
			output: scanner.Token{
				Id:     scanner.OR,
				Lexeme: "|",
			},
		},
		{
			symbol: scanner.AND,
			input:  "&",
			output: scanner.Token{
				Id:     scanner.AND,
				Lexeme: "&",
			},
		},
		{
			symbol: scanner.NOT,
			input:  "!",
			output: scanner.Token{
				Id:     scanner.NOT,
				Lexeme: "!",
			},
		},

		{
			symbol: scanner.OPENPAR,
			input:  "(",
			output: scanner.Token{
				Id:     scanner.OPENPAR,
				Lexeme: "(",
			},
		},
		{
			symbol: scanner.CLOSEPAR,
			input:  ")",
			output: scanner.Token{
				Id:     scanner.CLOSEPAR,
				Lexeme: ")",
			},
		},
		{
			symbol: scanner.OPENSQBR,
			input:  "[",
			output: scanner.Token{
				Id:     scanner.OPENSQBR,
				Lexeme: "[",
			},
		},
		{
			symbol: scanner.CLOSESQBR,
			input:  "]",
			output: scanner.Token{
				Id:     scanner.CLOSESQBR,
				Lexeme: "]",
			},
		},
		{
			symbol: scanner.OPENCUBR,
			input:  "{",
			output: scanner.Token{
				Id:     scanner.OPENCUBR,
				Lexeme: "{",
			},
		},
		{
			symbol: scanner.CLOSECUBR,
			input:  "}",
			output: scanner.Token{
				Id:     scanner.CLOSECUBR,
				Lexeme: "}",
			},
		},
		{
			symbol: scanner.DOT,
			input:  ".",
			output: scanner.Token{
				Id:     scanner.DOT,
				Lexeme: ".",
			},
		},
		{
			symbol: scanner.COMMA,
			input:  ",",
			output: scanner.Token{
				Id:     scanner.COMMA,
				Lexeme: ",",
			},
		},
		{
			symbol: scanner.SEMI,
			input:  ";",
			output: scanner.Token{
				Id:     scanner.SEMI,
				Lexeme: ";",
			},
		},
		{
			symbol: scanner.COLON,
			input:  ":",
			output: scanner.Token{
				Id:     scanner.COLON,
				Lexeme: ":",
			},
		},
		{
			symbol: scanner.COLONCOLON,
			input:  "::",
			output: scanner.Token{
				Id:     scanner.COLONCOLON,
				Lexeme: "::",
			},
		},
		{
			symbol: scanner.IF,
			input:  "if",
			output: scanner.Token{
				Id:     scanner.IF,
				Lexeme: "if",
			},
		},
		{
			symbol: scanner.ELSE,
			input:  "else",
			output: scanner.Token{
				Id:     scanner.ELSE,
				Lexeme: "else",
			},
		},
		{
			symbol: scanner.INTEGER,
			input:  "integer",
			output: scanner.Token{
				Id:     scanner.INTEGER,
				Lexeme: "integer",
			},
		},
		{
			symbol: scanner.FLOAT,
			input:  "float",
			output: scanner.Token{
				Id:     scanner.FLOAT,
				Lexeme: "float",
			},
		},
		{
			symbol: scanner.VOID,
			input:  "void",
			output: scanner.Token{
				Id:     scanner.VOID,
				Lexeme: "void",
			},
		},
		{
			symbol: scanner.PUBLIC,
			input:  "public",
			output: scanner.Token{
				Id:     scanner.PUBLIC,
				Lexeme: "public",
			},
		},
		{
			symbol: scanner.PRIVATE,
			input:  "private",
			output: scanner.Token{
				Id:     scanner.PRIVATE,
				Lexeme: "private",
			},
		},
		{
			symbol: scanner.FUNC,
			input:  "func",
			output: scanner.Token{
				Id:     scanner.FUNC,
				Lexeme: "func",
			},
		},
		{
			symbol: scanner.VAR,
			input:  "var",
			output: scanner.Token{
				Id:     scanner.VAR,
				Lexeme: "var",
			},
		},
		{
			symbol: scanner.STRUCT,
			input:  "struct",
			output: scanner.Token{
				Id:     scanner.STRUCT,
				Lexeme: "struct",
			},
		},
		{
			symbol: scanner.WHILE,
			input:  "while",
			output: scanner.Token{
				Id:     scanner.WHILE,
				Lexeme: "while",
			},
		},
		{
			symbol: scanner.READ,
			input:  "read",
			output: scanner.Token{
				Id:     scanner.READ,
				Lexeme: "read",
			},
		},
		{
			symbol: scanner.WRITE,
			input:  "write",
			output: scanner.Token{
				Id:     scanner.WRITE,
				Lexeme: "write",
			},
		},
		{
			symbol: scanner.RETURN,
			input:  "return",
			output: scanner.Token{
				Id:     scanner.RETURN,
				Lexeme: "return",
			},
		},
		{
			symbol: scanner.SELF,
			input:  "self",
			output: scanner.Token{
				Id:     scanner.SELF,
				Lexeme: "self",
			},
		},
		{
			symbol: scanner.INHERITS,
			input:  "inherits",
			output: scanner.Token{
				Id:     scanner.INHERITS,
				Lexeme: "inherits",
			},
		},
		{
			symbol: scanner.LET,
			input:  "let",
			output: scanner.Token{
				Id:     scanner.LET,
				Lexeme: "let",
			},
		},
		{
			symbol: scanner.IMPL,
			input:  "impl",
			output: scanner.Token{
				Id:     scanner.IMPL,
				Lexeme: "impl",
			},
		},
		{
			symbol: scanner.ID,
			input:  "asd",
			output: scanner.Token{
				Id:     scanner.ID,
				Lexeme: "asd",
			},
		},
		{
			symbol: scanner.INTNUM,
			input:  "99",
			output: scanner.Token{
				Id:     scanner.INTNUM,
				Lexeme: "99",
			},
		},
		{
			symbol: scanner.FLOATNUM,
			input:  "1.0",
			output: scanner.Token{
				Id:     scanner.FLOATNUM,
				Lexeme: "1.0",
			},
		},
	} {
		tc := tc
		t.Run(string(tc.symbol), func(t *testing.T) {
			t.Parallel()
			tc.output.Line = 1
			tc.output.Column = 1
			assertScan(t, tc.input, tc.output)
		})
	}
}

// Asserts the first token present in the input
func assertScan(t *testing.T, input string, token scanner.Token) {
	charsource := new(chuggingcharsource.ChuggingCharSource)
	err := charsource.ChugReader(bytes.NewBufferString(input))
	if err != nil {
		t.Fatalf("ChugReader should succeed here: %v", err)
	}

	scan := tabledrivenscanner.NewTableDrivenScanner(charsource, compositetable.TABLE)
	actual, err := scan.NextToken()
	if err != nil {
		t.Fatalf("NextToken should succeed: %v", err)
	}

	if actual != token {
		t.Fatalf("Expected token %v but got %v", token, actual)
	}
}
