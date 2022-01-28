package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	"github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
)

// This is the transition table that we are using for our scanner
// var table tabledrivenscanner.Table = compositetable.TABLE
var table = compositetable.TABLEE

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
	s := createScanner(t, lexpositivegradingsrc)
	for i, expected := range []scanner.Token{
		{
			Id:     scanner.EQ,
			Lexeme: "==",
			Line:   1,
			Column: 1,
		},
		{
			Id:     scanner.PLUS,
			Lexeme: "+",
			Line:   1,
			Column: 4,
		},
		{
			Id:     scanner.OR,
			Lexeme: "|",
			Line:   1,
			Column: 6,
		},
		{
			Id:     scanner.OPENPAR,
			Lexeme: "(",
			Line:   1,
			Column: 8,
		},
		{
			Id:     scanner.SEMI,
			Lexeme: ";",
			Line:   1,
			Column: 10,
		},
		{
			Id:     scanner.IF,
			Lexeme: "if",
			Line:   1,
			Column: 12,
		},
		{
			Id:     scanner.PUBLIC,
			Lexeme: "public",
			Line:   1,
			Column: 16,
		},
		{
			Id:     scanner.READ,
			Lexeme: "read",
			Line:   1,
			Column: 23,
		},
		{
			Id:     scanner.NOTEQ,
			Lexeme: "<>",
			Line:   2,
			Column: 1,
		},
		{
			Id:     scanner.MINUS,
			Lexeme: "-",
			Line:   2,
			Column: 4,
		},
		{
			Id:     scanner.AND,
			Lexeme: "&",
			Line:   2,
			Column: 6,
		},
		{
			Id:     scanner.CLOSEPAR,
			Lexeme: ")",
			Line:   2,
			Column: 8,
		},
		{
			Id:     scanner.COMMA,
			Lexeme: ",",
			Line:   2,
			Column: 10,
		},
		{
			Id:     scanner.THEN,
			Lexeme: "then",
			Line:   2,
			Column: 12,
		},
		{
			Id:     scanner.PRIVATE,
			Lexeme: "private",
			Line:   2,
			Column: 17,
		},
		{
			Id:     scanner.WRITE,
			Lexeme: "write",
			Line:   2,
			Column: 25,
		},
		{
			Id:     scanner.LT,
			Lexeme: "<",
			Line:   3,
			Column: 1,
		},
		{
			Id:     scanner.MULT,
			Lexeme: "*",
			Line:   3,
			Column: 3,
		},
		{
			Id:     scanner.NOT,
			Lexeme: "!",
			Line:   3,
			Column: 5,
		},
		{
			Id:     scanner.OPENCUBR,
			Lexeme: "{",
			Line:   3,
			Column: 7,
		},
		{
			Id:     scanner.DOT,
			Lexeme: ".",
			Line:   3,
			Column: 9,
		},
		{
			Id:     scanner.ELSE,
			Lexeme: "else",
			Line:   3,
			Column: 11,
		},
		{
			Id:     scanner.FUNC,
			Lexeme: "func",
			Line:   3,
			Column: 16,
		},
		{
			Id:     scanner.RETURN,
			Lexeme: "return",
			Line:   3,
			Column: 21,
		},
		{
			Id:     scanner.GT,
			Lexeme: ">",
			Line:   4,
			Column: 1,
		},
		{
			Id:     scanner.DIV,
			Lexeme: "/",
			Line:   4,
			Column: 3,
		},
		{
			Id:     scanner.CLOSECUBR,
			Lexeme: "}",
			Line:   4,
			Column: 6,
		},
		{
			Id:     scanner.COLON,
			Lexeme: ":",
			Line:   4,
			Column: 8,
		},
		{
			Id:     scanner.INTEGER,
			Lexeme: "integer",
			Line:   4,
			Column: 10,
		},
		{
			Id:     scanner.VAR,
			Lexeme: "var",
			Line:   4,
			Column: 18,
		},
		{
			Id:     scanner.SELF,
			Lexeme: "self",
			Line:   4,
			Column: 22,
		},
		{
			Id:     scanner.LEQ,
			Lexeme: "<=",
			Line:   5,
			Column: 1,
		},
		{
			Id:     scanner.ASSIGN,
			Lexeme: "=",
			Line:   5,
			Column: 4,
		},
		{
			Id:     scanner.OPENSQBR,
			Lexeme: "[",
			Line:   5,
			Column: 7,
		},
		{
			Id:     scanner.COLONCOLON,
			Lexeme: "::",
			Line:   5,
			Column: 9,
		},
		{
			Id:     scanner.FLOAT,
			Lexeme: "float",
			Line:   5,
			Column: 12,
		},
		{
			Id:     scanner.STRUCT,
			Lexeme: "struct",
			Line:   5,
			Column: 18,
		},
		{
			Id:     scanner.INHERITS,
			Lexeme: "inherits",
			Line:   5,
			Column: 25,
		},
		{
			Id:     scanner.GEQ,
			Lexeme: ">=",
			Line:   6,
			Column: 1,
		},
		{
			Id:     scanner.CLOSESQBR,
			Lexeme: "]",
			Line:   6,
			Column: 6,
		},
		{
			Id:     scanner.ARROW,
			Lexeme: "->",
			Line:   6,
			Column: 8,
		},
		{
			Id:     scanner.VOID,
			Lexeme: "void",
			Line:   6,
			Column: 11,
		},
		{
			Id:     scanner.WHILE,
			Lexeme: "while",
			Line:   6,
			Column: 16,
		},
		{
			Id:     scanner.LET,
			Lexeme: "let",
			Line:   6,
			Column: 22,
		},
		{
			Id:     scanner.FUNC,
			Lexeme: "func",
			Line:   7,
			Column: 7,
		},
		{
			Id:     scanner.IMPL,
			Lexeme: "impl",
			Line:   7,
			Column: 12,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "0",
			Line:   13,
			Column: 1,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "1",
			Line:   14,
			Column: 1,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "10",
			Line:   15,
			Column: 1,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "12",
			Line:   16,
			Column: 1,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "123",
			Line:   17,
			Column: 1,
		},
		{
			Id:     scanner.INTNUM,
			Lexeme: "12345",
			Line:   18,
			Column: 1,
		},
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "1.23",
			Line:   20,
			Column: 1,
		},
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "12.34",
			Line:   21,
			Column: 1,
		},
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "120.34e10",
			Line:   22,
			Column: 1,
		},
		{
			Id:     scanner.FLOATNUM,
			Lexeme: "12345.6789e-123",
			Line:   23,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "abc",
			Line:   25,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "abc1",
			Line:   26,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "a1bc",
			Line:   27,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "abc_1abc",
			Line:   28,
			Column: 1,
		},
		{
			Id:     scanner.ID,
			Lexeme: "abc1_abc",
			Line:   29,
			Column: 1,
		},
		{
			Id:     scanner.INLINECMT,
			Lexeme: "// this is an inline comment\n",
			Line:   31,
			Column: 1,
		},
		{
			Id:     scanner.BLOCKCMT,
			Lexeme: "/* this is a single line block comment */",
			Line:   33,
			Column: 1,
		},
		{
			Id:     scanner.BLOCKCMT,
			Lexeme: "/* this is a\nmultiple line\nblock comment\n*/",
			Line:   35,
			Column: 1,
		},
		{
			Id:     scanner.BLOCKCMT,
			Lexeme: "/* this is an imbricated\n/* block comment\n*/\n*/",
			Line:   40,
			Column: 1,
		},
	} {
		t.Run(fmt.Sprintf("Token-%v[%v]", i+1, expected.Lexeme), func(t *testing.T) {
			expected.Line++ // We added a newline compared to the read file
			actual := mustNextToken(t, s)
			if actual != expected {
				t.Errorf("Expected token %v but got %v", expected, actual)
			}
		})
	}
}

// Tests the prof-provided src data from `lexnegativegrading.src`. Note that the
// data has been inlined into this test file to reduce external dependencies
func TestLexNegativeGrading(t *testing.T) {
	t.Parallel()
	// TODO: implement test
}

// Tests the ability to lex nested comments
func TestImbricatedComments(t *testing.T) {
	t.Parallel()

	// This should be a single BLOCKCMT token
	src := `/* this is an imbricated
	/* block comment
	*/
	*/`

	s := createScanner(t, src)
	actual := mustNextToken(t, s)
	expected := scanner.Token{
		Id:     scanner.BLOCKCMT,
		Lexeme: scanner.Lexeme(src),
		Line:   1,
		Column: 1,
	}

	if actual != expected {
		t.Errorf("Expected token %v but got %v", expected, actual)
	}
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
		actual := mustNextToken(t, s)
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
		actual := mustNextToken(t, s)
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
			name:  scanner.THEN,
			input: "then",
			output: scanner.Token{
				Id:     scanner.THEN,
				Lexeme: "then",
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
// errors. Returns the token. If scanner.NextToken() returns an error, this
// function will immediately fail your test with the error.
func mustNextToken(t *testing.T, s scanner.Scanner) scanner.Token {
	tt, err := s.NextToken()
	if err != nil {
		t.Fatalf("NextToken should succeed: %v", err)
	}
	return tt
}

// Asserts the first token present in the input. Note that this function creates
// a new scanner everytime containing the provided input everytime you call it.
// It is used for testing inputs that contain a single token.
func assertScan(t *testing.T, input string, expected scanner.Token) {
	s := createScanner(t, input)
	if actual := mustNextToken(t, s); actual != expected {
		t.Fatalf("Expected token %v but got %v", expected, actual)
	}
}

// Creates a scanner with a char source containing the provided contents
func createScanner(t *testing.T, contents string) *tabledrivenscanner.TableDrivenScanner {
	chars := createCharSource(t, contents)
	return tabledrivenscanner.NewTableDrivenScanner(chars, table())
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
