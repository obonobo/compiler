package compositetable

import (
	"fmt"

	tdp "github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/token"
)

// IMPLEMENTS:
// t *CompositeTable tabledrivenparser.Table
type CompositeTable struct {
	StartTerminal token.Kind
	Rules         token.Rules
	Terminals     token.KindSet
	Nonterminals  token.KindSet
	Firsts        map[token.Kind]token.KindSet
	Follows       map[token.Kind]token.KindSet
	TT            map[token.Key]token.Rule
	NoPush        token.KindSet // Set of tokens that don't need to be pushed on the stack
}

// Perform a lookup on the table, may return a NoRuleError if
func (t *CompositeTable) Lookup(row token.Kind, col token.Kind) (token.Rule, error) {
	r, ok := t.TT[token.Key{Nonterminal: row, Terminal: col}]
	if !ok {
		return token.Rule{},
			fmt.Errorf("CompositeTable.Lookup: %w", &tdp.NoRuleError{Row: row, Col: col})
	}
	return t.filterNoPush(r), nil
}

// Returns the starting nonterminal symbol
func (t *CompositeTable) Start() token.Kind {
	return t.StartTerminal
}

// Determine whether the symbol is part of the set of terminal symbols
func (t *CompositeTable) IsTerminal(symbol token.Kind) bool {
	_, ok := t.Terminals[symbol]
	return ok
}

// Determine whether the symbol is part of the set of nonterminal symbols
func (t *CompositeTable) IsNonterminal(symbol token.Kind) bool {
	_, ok := t.Nonterminals[symbol]
	return ok
}

// Determine whether the symbol has a rule of the form: <symbol> -> EPSILON
func (t *CompositeTable) HasEpsilonRule(symbol token.Kind) bool {
	if r, ok := t.Rules[symbol]; ok {
		for _, rr := range r {
			if len(rr.RHS) == 1 {
				for tok := range t.NoPush {
					if tok == rr.RHS[0] {
						return true
					}
				}
			}
		}
	}
	return false
}

// Retrieve the FIRST set of a symbol
func (t *CompositeTable) Follow(symbol token.Kind) (token.KindSet, bool) {
	ks, ok := t.Follows[symbol]
	return ks, ok
}

// Retrieve the FOLLOW set of a symbol
func (t *CompositeTable) First(symbol token.Kind) (token.KindSet, bool) {
	ks, ok := t.Firsts[symbol]
	return ks, ok
}

func (t *CompositeTable) filterNoPush(r token.Rule) token.Rule {
	rr := token.Rule{LHS: r.LHS, RHS: make([]token.Kind, 0, len(r.RHS))}
	for _, k := range r.RHS {
		if !t.isNoPush(k) {
			rr.RHS = append(rr.RHS, k)
		}
	}
	return rr
}

func (t *CompositeTable) isNoPush(k token.Kind) bool {
	_, ok := t.NoPush[k]
	return ok
}
