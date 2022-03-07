package compositetable

import (
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/token"
)

func NewTableDrivenParser(
	scnr scanner.Scanner,
	errc chan<- tabledrivenparser.ParserError,
	rulec chan<- token.Rule,
) *tabledrivenparser.TableDrivenParser {
	return tabledrivenparser.NewParser(scnr, TABLE(), errc, rulec)
}

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
