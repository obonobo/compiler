package main

import (
	"bytes"
	"testing"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	"github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
)

// This is the transition table that we are using for our scanner
var table tabledrivenscanner.Table = compositetable.TABLE

// INLINED FILE: `lexpositivegrading.src`
const lexpositivegradingsrc = `
==	+	|	(	;	if 	public	read
<>	-	&	)	,	then	private	write
<	*	!	{	.	else	func	return
>	/		}	:	integer	var	self
<=	=		[	::	float	struct	inherits
>=			]	->	void	while	let
						func	impl





0
1
10
12
123
12345

1.23
12.34
120.34e10
12345.6789e-123

abc
abc1
a1bc
abc_1abc
abc1_abc

// this is an inline comment

/* this is a single line block comment */

/* this is a
multiple line
block comment
*/

/* this is an imbricated
/* block comment
*/
*/




`

// INLINED FILE: `lexnegativegrading.src`
const lexnegativegradingsrc = `
@ # $ ' \ ~

00
01
010
0120
01230
0123450

01.23
012.34
12.340
012.340

012.34e10
12.34e010

_abc
1abc
_1abc

`

// Tests the prof-provided src data from `lexpositivegrading.src`. Note that the
// data has been inlined into this test file to reduce external dependencies
func TestLexPositiveGrading(t *testing.T) {
	t.Parallel()
	// TODO: implement test
}

// Tests the prof-provided src data from `lexnegativegrading.src`. Note that the
// data has been inlined into this test file to reduce external dependencies
func TestLexNegativeGrading(t *testing.T) {
	t.Parallel()
	// TODO: implement test
}

func TestFloatIdNewlineIdId(t *testing.T) {
	t.Parallel()

	tokens := []scanner.Token{
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "1.0",
			Line:   1,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "example_id",
			Line:   1,
			Column: 5,
		},
		{
			Id:     scanner.ID,
			Lexeme: "Id2",
			Line:   2,
			Column: 2,
		},
		{
			Id:     scanner.ID,
			Lexeme: "ID3",
			Line:   2,
			Column: 6,
		},
	}

	src := "1.0 example_id\n Id2 ID3"
	s := createScanner(t, src)

	for _, expected := range tokens {
		actual := assertNextTokenSuccess(t, s)
		if actual != expected {
			t.Errorf("Expected token %v but got %v", expected, actual)
		}
	}
}

func TestFloatIdIdId(t *testing.T) {
	t.Parallel()

	tokens := []scanner.Token{
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "1.0",
			Line:   1,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "example_id",
			Line:   1,
			Column: 5,
		},
		{
			Id:     scanner.ID,
			Lexeme: "Id2",
			Line:   1,
			Column: 16,
		},
		{
			Id:     scanner.ID,
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
		actual := assertNextTokenSuccess(t, s)
		if actual != expected {
			t.Errorf("Expected token %v but got %v", expected, actual)
		}
	}
}

func TestDoubleBackup(t *testing.T) {
	expectedToken1 := scanner.Token{
		Id:     scanner.FLOATNUM,
		Lexeme: "1.0",
		Line:   1,
		Column: 1,
	}

	expectedToken2 := scanner.Token{
		Id:     scanner.ID,
		Lexeme: "example_id",
		Line:   1,
		Column: 4,
	}

	src := string(expectedToken1.Lexeme + expectedToken2.Lexeme)

	// TEST
	t.Run(src, func(t *testing.T) {
		t.Parallel()
		s := createScanner(t, src)

		actualToken1 := assertNextTokenSuccess(t, s)
		actualToken2 := assertNextTokenSuccess(t, s)

		if actualToken1 != expectedToken1 {
			t.Errorf("Expected first token to be %v but got %v", expectedToken1, actualToken1)
		}

		if actualToken2 != expectedToken2 {
			t.Errorf("Expected second token to be %v but got %v", expectedToken2, actualToken2)
		}
	})
}

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

// Grabs the next scanner.Token from the scanner and asserts that there were no
// errors. Returns the token
func assertNextTokenSuccess(t *testing.T, s scanner.Scanner) scanner.Token {
	tt, err := s.NextToken()
	if err != nil {
		t.Fatalf("NextToken should succeed: %v", err)
	}
	return tt
}

// Asserts the first token present in the input
func assertScan(t *testing.T, input string, expected scanner.Token) {
	s := createScanner(t, input)
	if actual := assertNextTokenSuccess(t, s); actual != expected {
		t.Fatalf("Expected token %v but got %v", expected, actual)
	}
}

// Creates a scanner with a char source containing the provided contents
func createScanner(t *testing.T, contents string) *tabledrivenscanner.TableDrivenScanner {
	chars := createCharSource(t, contents)
	return tabledrivenscanner.NewTableDrivenScanner(chars, table)
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
