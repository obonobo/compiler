package tabledrivenparser

import (
	"github.com/obonobo/esac/core/token"
)

// The Table used by the TableDrivenParser
type Table interface {

	// Perform a lookup on the table, may return a NoRuleError if
	Lookup(row, col token.Kind) (token.Rule, error)

	// Returns the starting nonterminal symbol
	Start() token.Kind

	// Determine whether the symbol is part of the set of terminal symbols
	IsTerminal(symbol token.Kind) bool

	// Determine whether the symbol is part of the set of nonterminal symbols
	IsNonterminal(symbol token.Kind) bool

	// Determine whether the symbol has a rule of the form: <symbol> -> EPSILON
	HasEpsilonRule(symbol token.Kind) bool

	// Retrieve the FIRST set of a symbol
	First(symbol token.Kind) (token.KindSet, bool)

	// Retrieve the FOLLOW set of a symbol
	Follow(symbol token.Kind) (token.KindSet, bool)
}
