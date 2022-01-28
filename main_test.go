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
		name   scanner.Symbol
		input  string
		output scanner.Token
	}{
		{
			name:  scanner.INLINECMT,
			input: "// asdasdasd \n",
			output: scanner.Token{
				Id:     scanner.INLINECMT,
				Lexeme: "// asdasdasd \n",
			},
		},
		{
			name:  scanner.BLOCKCMT,
			input: "/* \n asdasd \n asdasd \n -123=*/",
			output: scanner.Token{
				Id:     scanner.BLOCKCMT,
				Lexeme: "/* \n asdasd \n asdasd \n -123=*/",
			},
		},
		{
			name:  scanner.DIV,
			input: "/",
			output: scanner.Token{
				Id:     scanner.DIV,
				Lexeme: "/",
			},
		},
		{
			name:  scanner.MULT,
			input: "*",
			output: scanner.Token{
				Id:     scanner.MULT,
				Lexeme: "*",
			},
		},
		{
			name:  scanner.MINUS,
			input: "-",
			output: scanner.Token{
				Id:     scanner.MINUS,
				Lexeme: "-",
			},
		},
		{
			name:  scanner.ARROW,
			input: "->",
			output: scanner.Token{
				Id:     scanner.ARROW,
				Lexeme: "->",
			},
		},
		{
			name:  scanner.ASSIGN,
			input: "=",
			output: scanner.Token{
				Id:     scanner.ASSIGN,
				Lexeme: "=",
			},
		},
		{
			name:  scanner.EQ,
			input: "==",
			output: scanner.Token{
				Id:     scanner.EQ,
				Lexeme: "==",
			},
		},
		{
			name:  scanner.LT,
			input: "<",
			output: scanner.Token{
				Id:     scanner.LT,
				Lexeme: "<",
			},
		},
		{
			name:  scanner.GT,
			input: ">",
			output: scanner.Token{
				Id:     scanner.GT,
				Lexeme: ">",
			},
		},
		{
			name:  scanner.NOTEQ,
			input: "<>",
			output: scanner.Token{
				Id:     scanner.NOTEQ,
				Lexeme: "<>",
			},
		},
		{
			name:  scanner.LEQ,
			input: "<=",
			output: scanner.Token{
				Id:     scanner.LEQ,
				Lexeme: "<=",
			},
		},
		{
			name:  scanner.GEQ,
			input: ">=",
			output: scanner.Token{
				Id:     scanner.GEQ,
				Lexeme: ">=",
			},
		},
		{
			name:  scanner.OR,
			input: "|",
			output: scanner.Token{
				Id:     scanner.OR,
				Lexeme: "|",
			},
		},
		{
			name:  scanner.AND,
			input: "&",
			output: scanner.Token{
				Id:     scanner.AND,
				Lexeme: "&",
			},
		},
		{
			name:  scanner.NOT,
			input: "!",
			output: scanner.Token{
				Id:     scanner.NOT,
				Lexeme: "!",
			},
		},

		{
			name:  scanner.OPENPAR,
			input: "(",
			output: scanner.Token{
				Id:     scanner.OPENPAR,
				Lexeme: "(",
			},
		},
		{
			name:  scanner.CLOSEPAR,
			input: ")",
			output: scanner.Token{
				Id:     scanner.CLOSEPAR,
				Lexeme: ")",
			},
		},
		{
			name:  scanner.OPENSQBR,
			input: "[",
			output: scanner.Token{
				Id:     scanner.OPENSQBR,
				Lexeme: "[",
			},
		},
		{
			name:  scanner.CLOSESQBR,
			input: "]",
			output: scanner.Token{
				Id:     scanner.CLOSESQBR,
				Lexeme: "]",
			},
		},
		{
			name:  scanner.OPENCUBR,
			input: "{",
			output: scanner.Token{
				Id:     scanner.OPENCUBR,
				Lexeme: "{",
			},
		},
		{
			name:  scanner.CLOSECUBR,
			input: "}",
			output: scanner.Token{
				Id:     scanner.CLOSECUBR,
				Lexeme: "}",
			},
		},
		{
			name:  scanner.DOT,
			input: ".",
			output: scanner.Token{
				Id:     scanner.DOT,
				Lexeme: ".",
			},
		},
		{
			name:  scanner.COMMA,
			input: ",",
			output: scanner.Token{
				Id:     scanner.COMMA,
				Lexeme: ",",
			},
		},
		{
			name:  scanner.SEMI,
			input: ";",
			output: scanner.Token{
				Id:     scanner.SEMI,
				Lexeme: ";",
			},
		},
		{
			name:  scanner.COLON,
			input: ":",
			output: scanner.Token{
				Id:     scanner.COLON,
				Lexeme: ":",
			},
		},
		{
			name:  scanner.COLONCOLON,
			input: "::",
			output: scanner.Token{
				Id:     scanner.COLONCOLON,
				Lexeme: "::",
			},
		},
		{
			name:  scanner.IF,
			input: "if",
			output: scanner.Token{
				Id:     scanner.IF,
				Lexeme: "if",
			},
		},
		{
			name:  scanner.ELSE,
			input: "else",
			output: scanner.Token{
				Id:     scanner.ELSE,
				Lexeme: "else",
			},
		},
		{
			name:  scanner.INTEGER,
			input: "integer",
			output: scanner.Token{
				Id:     scanner.INTEGER,
				Lexeme: "integer",
			},
		},
		{
			name:  scanner.FLOAT,
			input: "float",
			output: scanner.Token{
				Id:     scanner.FLOAT,
				Lexeme: "float",
			},
		},
		{
			name:  scanner.VOID,
			input: "void",
			output: scanner.Token{
				Id:     scanner.VOID,
				Lexeme: "void",
			},
		},
		{
			name:  scanner.PUBLIC,
			input: "public",
			output: scanner.Token{
				Id:     scanner.PUBLIC,
				Lexeme: "public",
			},
		},
		{
			name:  scanner.PRIVATE,
			input: "private",
			output: scanner.Token{
				Id:     scanner.PRIVATE,
				Lexeme: "private",
			},
		},
		{
			name:  scanner.FUNC,
			input: "func",
			output: scanner.Token{
				Id:     scanner.FUNC,
				Lexeme: "func",
			},
		},
		{
			name:  scanner.VAR,
			input: "var",
			output: scanner.Token{
				Id:     scanner.VAR,
				Lexeme: "var",
			},
		},
		{
			name:  scanner.STRUCT,
			input: "struct",
			output: scanner.Token{
				Id:     scanner.STRUCT,
				Lexeme: "struct",
			},
		},
		{
			name:  scanner.WHILE,
			input: "while",
			output: scanner.Token{
				Id:     scanner.WHILE,
				Lexeme: "while",
			},
		},
		{
			name:  scanner.READ,
			input: "read",
			output: scanner.Token{
				Id:     scanner.READ,
				Lexeme: "read",
			},
		},
		{
			name:  scanner.WRITE,
			input: "write",
			output: scanner.Token{
				Id:     scanner.WRITE,
				Lexeme: "write",
			},
		},
		{
			name:  scanner.RETURN,
			input: "return",
			output: scanner.Token{
				Id:     scanner.RETURN,
				Lexeme: "return",
			},
		},
		{
			name:  scanner.SELF,
			input: "self",
			output: scanner.Token{
				Id:     scanner.SELF,
				Lexeme: "self",
			},
		},
		{
			name:  scanner.INHERITS,
			input: "inherits",
			output: scanner.Token{
				Id:     scanner.INHERITS,
				Lexeme: "inherits",
			},
		},
		{
			name:  scanner.LET,
			input: "let",
			output: scanner.Token{
				Id:     scanner.LET,
				Lexeme: "let",
			},
		},
		{
			name:  scanner.IMPL,
			input: "impl",
			output: scanner.Token{
				Id:     scanner.IMPL,
				Lexeme: "impl",
			},
		},
		{
			name:  scanner.ID,
			input: "asd",
			output: scanner.Token{
				Id:     scanner.ID,
				Lexeme: "asd",
			},
		},
		{
			name:  scanner.INTNUM,
			input: "99",
			output: scanner.Token{
				Id:     scanner.INTNUM,
				Lexeme: "99",
			},
		},
		{
			name:  scanner.FLOATNUM,
			input: "1.0",
			output: scanner.Token{
				Id:     scanner.FLOATNUM,
				Lexeme: "1.0",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[00]",
			input: "00",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "00",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[01]",
			input: "01",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "01",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[010]",
			input: "010",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "010",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[0120]",
			input: "0120",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "0120",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[01230]",
			input: "01230",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "01230",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[0123450]",
			input: "0123450",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "0123450",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[01.23]",
			input: "01.23",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "01.23",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[012.34]",
			input: "012.34",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "012.34",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[012.340]",
			input: "012.340",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "012.340",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[012.34e10]",
			input: "012.34e10",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "012.34e10",
			},
		},
		{
			name:  scanner.INVALIDNUM + "[12.34e010]",
			input: "12.34e010",
			output: scanner.Token{
				Id:     scanner.INVALIDNUM,
				Lexeme: "12.34e010",
			},
		},
		{
			name:  scanner.INVALIDIDENTIFIER + "[_abc]",
			input: "_abc",
			output: scanner.Token{
				Id:     scanner.INVALIDIDENTIFIER,
				Lexeme: "_abc",
			},
		},
		{
			name:  scanner.INVALIDIDENTIFIER + "[1abc]",
			input: "1abc",
			output: scanner.Token{
				Id:     scanner.INVALIDIDENTIFIER,
				Lexeme: "1abc",
			},
		},
		{
			name:  scanner.INVALIDIDENTIFIER + "[_1abc]",
			input: "_1abc",
			output: scanner.Token{
				Id:     scanner.INVALIDIDENTIFIER,
				Lexeme: "_1abc",
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
