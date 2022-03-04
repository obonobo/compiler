package tabledrivenparser

import (
	"reflect"
	"testing"

	"github.com/obonobo/esac/core/token"
)

// Test the statementCloser function with various stacks
func TestStatementCloser(t *testing.T) {
	semiLookahead := token.Token{
		Id:     token.SEMI,
		Lexeme: ";",
		Line:   -1,
		Column: -1,
	}

	openCurlyBracketLookahead := token.Token{
		Id:     token.OPENCUBR,
		Lexeme: "{",
		Line:   -1,
		Column: -1,
	}

	for _, tc := range []struct {
		name      string
		stack     []token.Kind
		lookahead token.Token

		expectedSlice []token.Kind
		expectedOk    bool
	}{
		{
			name:          "empty stack",
			stack:         []token.Kind{},
			lookahead:     semiLookahead,
			expectedSlice: []token.Kind{},
			expectedOk:    false,
		},
		{
			name:          "one element stack with no closer",
			stack:         []token.Kind{token.ID},
			lookahead:     semiLookahead,
			expectedSlice: []token.Kind{token.ID},
			expectedOk:    false,
		},
		{
			name:          "one element stack with semi-colon",
			stack:         []token.Kind{token.SEMI},
			lookahead:     semiLookahead,
			expectedSlice: []token.Kind{token.SEMI},
			expectedOk:    true,
		},
		{
			name:      "normal scenario with multi-symbol statement on the stack",
			lookahead: semiLookahead,
			stack: []token.Kind{
				"<rept-prog0>",
				token.SEMI,
				"closecubr",
				"<rept-structDecl4>",
				token.SEMI,
				"<rept-varDecl4>",
				"<type>",
				"colon",
			},
			expectedSlice: []token.Kind{
				token.SEMI,
				"<rept-varDecl4>",
				"<type>",
				"colon",
			},
			expectedOk: true,
		},
		{
			name:      "lookahead=opencubr normal scenario with multi-symbol statement on the stack",
			lookahead: openCurlyBracketLookahead,
			stack: []token.Kind{
				"<rept-prog0>",
				token.OPENCUBR,
				"closecubr",
				"<rept-structDecl4>",
				token.OPENCUBR,
				"<rept-varDecl4>",
				"<type>",
				"colon",
			},
			expectedSlice: []token.Kind{
				token.OPENCUBR,
				"<rept-varDecl4>",
				"<type>",
				"colon",
			},
			expectedOk: true,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := &TableDrivenParser{stack: tc.stack, table: &MockTable{}}
			actualSlice, ok := parser.statementCloser()
			actualOk := ok(tc.lookahead)
			if actualOk != tc.expectedOk || !reflect.DeepEqual(actualSlice, tc.expectedSlice) {
				t.Errorf(
					"Expected statementCloser() to return (%v, %v) but got (%v, %v)",
					tc.expectedSlice, tc.expectedOk, actualSlice, actualOk)
			}
		})
	}
}

type MockTable struct{}

func (t *MockTable) Lookup(row token.Kind, col token.Kind) (token.Rule, error) {
	panic("not implemented") // TODO: Implement
}

func (t *MockTable) Start() token.Kind {
	panic("not implemented") // TODO: Implement
}

func (t *MockTable) IsTerminal(symbol token.Kind) bool {
	panic("not implemented") // TODO: Implement
}

func (t *MockTable) IsNonterminal(symbol token.Kind) bool {
	panic("not implemented") // TODO: Implement
}

func (t *MockTable) HasEpsilonRule(symbol token.Kind) bool {
	panic("not implemented") // TODO: Implement
}

func (t *MockTable) First(symbol token.Kind) (token.KindSet, bool) {
	return token.KindSet{symbol: {}}, true
}

func (t *MockTable) Follow(symbol token.Kind) (token.KindSet, bool) {
	return token.KindSet{symbol: {}}, true
}
