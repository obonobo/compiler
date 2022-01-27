package compositetable

import (
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
)

const (
	ANY rune = 0 // Represents any character
	// LETTER rune = -1 // Represents expression [aA-zZ]
	// DIGIT  rune = -2 // Represents expression [0-9]
	// ONE    rune = -3 // Represents expression [1-9]
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
		// COMMENTS

		// OPERATORS AND PUNCTUATION
		{1, '='}: 9,
		{9, ANY}: 10,
		{9, '='}: 11,

		// OPERATORS AND PUNCTUATION

		// ID AND RESERVED WORDS

		// ID AND RESERVED WORDS

		// INTS AND FLOATS
		// INTS AND FLOATS
	},

	// STATES THAT NEED BACKUP
	NeedBackup: map[tabledrivenscanner.State]struct{}{
		3:  {},
		10: {},
	},

	// STATE TO TOKEN MAPPING
	Tokens: map[tabledrivenscanner.State]scanner.Symbol{
		3: scanner.DIV,
		6: scanner.INLINECMT,
		8: scanner.BLOCKCMT,
	},
}
