package compositetable

import "github.com/obonobo/esac/core/token"

type Rule int

// type Key struct {
// 	Nonterminal token.Kind
// 	Terminal    token.Kind
// }

type Key = token.Key

// t *CompositeTable tabledrivenparser.Table
type CompositeTable struct {
	Rules        token.Rules
	Terminals    token.KindSet
	Nonterminals token.KindSet
	Firsts       map[token.Kind]token.KindSet
	Follows      map[token.Kind]token.KindSet
	TT           map[Key]token.Rule
}

// Perform a lookup on the table
func (t *CompositeTable) Lookup(row string, col string) string {
	panic("not implemented") // TODO: Implement
}

// Returns the starting nonterminal symbol
func (t *CompositeTable) Start() string {
	panic("not implemented") // TODO: Implement
}

// Determine whether the symbol is part of the set of terminal symbols
func (t *CompositeTable) IsTerminal(symbol string) bool {
	panic("not implemented") // TODO: Implement
}

// Determine whether the symbol is part of the set of nonterminal symbols
func (t *CompositeTable) IsNonterminal(symbol string) bool {
	panic("not implemented") // TODO: Implement
}

// Determine whether the symbol is an error symbol for this table
func (t *CompositeTable) IsError(symbol string) bool {
	panic("not implemented") // TODO: Implement
}
