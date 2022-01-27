package compositetable

import (
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
)

const (
	ANY    rune = 0  // Represents any character
	LETTER rune = -1 // Represents expression [aA-zZ]
	// DIGIT   rune = -2 // Represents expression [0-9]
	// ONENINE rune = -3 // Represents expression [1-9]
)

// This is the actual implementation table used by the app
var TABLE *CompositeTable = &CompositeTable{
	Start: 1,

	Transitions: map[Key]tabledrivenscanner.State{

		// COMMENTS
		{1, '/'}: 2,

		{2, '/'}: 5,
		{2, '*'}: 4,
		{2, ANY}: 3,

		{4, ANY}: 4,
		{4, '*'}: 7,

		{7, ANY}: 4,
		{7, '/'}: 8,

		{5, ANY}:  5,
		{5, '\n'}: 6,

		// OPERATORS AND PUNCTUATION
		{1, '='}: 9,
		{9, ANY}: 10,
		{9, '='}: 11,

		{1, '-'}:  13,
		{13, ANY}: 14,
		{13, '>'}: 15,

		{1, '<'}:  20,
		{20, ANY}: 21,
		{20, '>'}: 22,
		{20, '='}: 23,

		{1, '>'}:  24,
		{24, ANY}: 25,
		{24, '='}: 26,

		{1, ':'}:  36,
		{36, ANY}: 37,
		{36, ':'}: 38,

		{1, '+'}: 12,
		{1, '|'}: 16,
		{1, '&'}: 17,
		{1, '!'}: 18,
		{1, '*'}: 19,
		{1, '('}: 27,
		{1, ')'}: 28,
		{1, '{'}: 29,
		{1, '}'}: 30,
		{1, '['}: 31,
		{1, ']'}: 32,
		{1, ';'}: 33,
		{1, ','}: 34,
		{1, '.'}: 35,

		// ID AND RESERVED WORDS
		{1, LETTER}:  39,
		{39, LETTER}: 39,
		{39, DIGIT}:  39,
		{39, '_'}:    39,
		{39, ANY}:    40,

		// INTS AND FLOATS
		{1, ONENINE}:  41,
		{1, '0'}:      42,
		{41, ANY}:     43,
		{42, ANY}:     43,
		{41, '.'}:     44,
		{41, '.'}:     44,
		{44, DIGIT}:   45,
		{45, ONENINE}: 45,
		{45, ANY}:     46,
		{45, '0'}:     47,
		{47, ONENINE}: 45,
		{47, ANY}:     48,

		{45, 'e'}:     49,
		{49, '0'}:     51,
		{49, '-'}:     50,
		{49, '+'}:     50,
		{49, ONENINE}: 54,
		{50, '0'}:     51,
		{50, ONENINE}: 54,
		{51, DIGIT}:   52,
		{51, ANY}:     53,
		{52, DIGIT}:   52,
		{52, ANY}:     55,
		{54, DIGIT}:   54,
		{54, ANY}:     53,
	},

	// STATES THAT NEED BACKUP
	NeedBackup: map[tabledrivenscanner.State]struct{}{
		3:  {},
		10: {},
		14: {},
		21: {},
		25: {},
		37: {},
		40: {},
		43: {},
		46: {},
		48: {},
		55: {},
		53: {},
	},

	// STATE TO TOKEN MAPPING
	Tokens: map[tabledrivenscanner.State]scanner.Symbol{
		3:  scanner.DIV,
		6:  scanner.INLINECMT,
		8:  scanner.BLOCKCMT,
		10: scanner.ASSIGN,
		11: scanner.EQ,
		12: scanner.PLUS,
		14: scanner.MINUS,
		15: scanner.ARROW,
		16: scanner.OR,
		17: scanner.AND,
		18: scanner.NOT,
		19: scanner.MULT,
		21: scanner.LT,
		22: scanner.NOTEQ,
		23: scanner.LEQ,
		25: scanner.GT,
		26: scanner.GEQ,
		27: scanner.OPENPAR,
		28: scanner.CLOSEPAR,
		29: scanner.OPENCUBR,
		30: scanner.CLOSECUBR,
		31: scanner.OPENSQBR,
		32: scanner.CLOSESQBR,
		33: scanner.SEMI,
		34: scanner.COMMA,
		35: scanner.DOT,
		40: scanner.ID,
		43: scanner.INTNUM,
		46: scanner.FLOATNUM,
		53: scanner.FLOATNUM,
		48: scanner.INVALIDNUM,
		55: scanner.INVALIDNUM,
	},
}
