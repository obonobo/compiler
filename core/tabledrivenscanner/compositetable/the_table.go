package compositetable

import (
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
)

const ANY rune = 0

// This is the actual implementation table used by the app
var TABLE *CompositeTable = &CompositeTable{
	start: 1,

	transitions: map[Key]tabledrivenscanner.State{
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
	},

	// STATES THAT NEED BACKUP
	needBackup: map[tabledrivenscanner.State]struct{}{
		3: struct{}{},
	},

	// FINAL STATES
	finalStates: map[tabledrivenscanner.State]struct{}{
		3: struct{}{},
		6: struct{}{},
		8: struct{}{},
	},

	// STATE TO TOKEN MAPPING
	tokens: map[tabledrivenscanner.State]scanner.Symbol{
		3: scanner.DIV,
		6: scanner.INLINECMT,
		8: scanner.BLOCKCMT,
	},
}
