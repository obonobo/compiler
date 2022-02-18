package compositetable

import (
	"log"

	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/token"
)

func NewTableDrivenParser(scnr scanner.Scanner, logger *log.Logger) *tabledrivenparser.TableDrivenParser {
	return tabledrivenparser.NewTableDriverParser(scnr, TABLE(), logger)
}

func TABLE() *CompositeTable {
	return &CompositeTable{
		Rules:        token.RULES(),
		Terminals:    token.TERMINALS(),
		Nonterminals: token.NONTERMINALS(),
		Firsts:       token.FIRSTS(),
		Follows:      token.FOLLOWS(),
		TT:           token.TABLE(),
	}
}
