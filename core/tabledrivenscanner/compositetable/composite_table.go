package compositetable

import (
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
)

const INITIAL tabledrivenscanner.State = 1

type Key struct {
	state tabledrivenscanner.State // Current state that the scanner is on
	next  rune                     // The symbol that is being processed
}

// State transition table. Once initialized, it's contents should never be
// changed. The table should never be written to, only read from. CompositeTable
// has composite Key and Values
type CompositeTable struct {
	Start       tabledrivenscanner.State
	Transitions map[Key]tabledrivenscanner.State
	Tokens      map[tabledrivenscanner.State]scanner.Symbol
	NeedBackup  map[tabledrivenscanner.State]struct{}
	// FinalStates map[tabledrivenscanner.State]struct{}
}

// Perform a transition
func (t *CompositeTable) Next(state tabledrivenscanner.State, char rune) tabledrivenscanner.State {
	s, ok := t.Transitions[Key{state, char}]
	if !ok {
		// Can try to see if there is an ANY state
		s = t.Transitions[Key{state, ANY}]
	}
	return s
}

// Check if a state requires the scanner to backup
func (t *CompositeTable) NeedsBackup(state tabledrivenscanner.State) bool {
	_, ok := t.NeedBackup[state]
	return ok
}

// The initial state
func (t *CompositeTable) Initial() tabledrivenscanner.State {
	return t.Start
}

// Check if a state is a final state
func (t *CompositeTable) IsFinal(state tabledrivenscanner.State) bool {
	// _, ok := t.FinalStates[state]
	_, ok := t.Tokens[state]
	return ok
}

// Generates a token given a State
func (t *CompositeTable) CreateToken(
	state tabledrivenscanner.State,
	lexeme scanner.Lexeme,
	line, col int,
) (scanner.Token, error) {
	symbol, ok := t.Tokens[state]

	if !ok {
		return scanner.Token{}, UnrecognizedStateError(state)
	}

	// IDs could actually be RESERVED WORDS
	res, ok := scanner.IsReservedWordString(string(lexeme))
	if symbol == scanner.ID && ok {
		symbol = res
	}

	return scanner.Token{
		Id:     symbol,
		Lexeme: lexeme,
		Line:   line,
		Column: col,
	}, nil
}
