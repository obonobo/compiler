package tabledrivenparser

// The Table used by the TableDrivenParser
type Table interface {

	// Perform a lookup on the table
	Lookup(row, col string) string

	// Returns the starting nonterminal symbol
	Start() string

	// Determine whether the symbol is part of the set of terminal symbols
	IsTerminal(symbol string) bool

	// Determine whether the symbol is part of the set of nonterminal symbols
	IsNonterminal(symbol string) bool

	// Determine whether the symbol is an error symbol for this table
	IsError(symbol string) bool
}
