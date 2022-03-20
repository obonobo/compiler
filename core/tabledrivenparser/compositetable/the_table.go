package compositetable

import (
	"github.com/obonobo/esac/core/token"
)

func TABLE() *CompositeTable {
	return &CompositeTable{
		StartTerminal: token.START,
		Rules:         token.RULES(),
		Terminals:     token.TERMINALS(),
		Nonterminals:  token.NONTERMINALS(),
		Firsts:        token.FIRSTS(),
		Follows:       token.FOLLOWS(),
		TT:            token.TABLE(),
		NoPush:        token.KindSet{token.EPSILON: {}},
	}
}
