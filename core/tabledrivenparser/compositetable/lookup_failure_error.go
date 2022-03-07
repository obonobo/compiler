package compositetable

import (
	"fmt"
	"sort"

	"github.com/obonobo/esac/core/token"
)

type LookupFailureError struct {
	Row token.Kind
	Err error // wrap

	table *CompositeTable
}

func (e *LookupFailureError) Error() string {
	return fmt.Sprintf("CompositeTable.Lookup: %v", e.Err)
}

func (e *LookupFailureError) Unwrap() error {
	return e.Err
}

// Returns the possible composite key column values for this row
func (e *LookupFailureError) Possibilities() []token.Kind {
	if e.table == nil {
		return []token.Kind{}
	}
	possibilities := make([]token.Kind, 0, 10)
	for k := range e.table.TT {
		if k.Nonterminal == e.Row {
			possibilities = append(possibilities, k.Terminal)
		}
	}

	// Sort them so that our error messages are deterministic
	sort.Slice(possibilities, func(i, j int) bool {
		return possibilities[i] < possibilities[j]
	})

	return possibilities
}

// Creates a LookupFailureError for this CompositeTable
func (t *CompositeTable) lookupFailure(row token.Kind, wrap error) *LookupFailureError {
	return &LookupFailureError{table: t, Row: row, Err: wrap}
}
