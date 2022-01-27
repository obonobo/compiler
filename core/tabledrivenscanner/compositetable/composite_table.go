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
	start       tabledrivenscanner.State
	transitions map[Key]tabledrivenscanner.State
	needBackup  map[tabledrivenscanner.State]struct{}
	finalStates map[tabledrivenscanner.State]struct{}
	tokens      map[tabledrivenscanner.State]scanner.Symbol
}

func NewCompositeTable(
	transitions map[Key]tabledrivenscanner.State,
	needBackup []tabledrivenscanner.State,
	finalStates []tabledrivenscanner.State,
	tokens map[tabledrivenscanner.State]scanner.Symbol,
) *CompositeTable {
	t := &CompositeTable{
		start:       INITIAL,
		transitions: make(map[Key]tabledrivenscanner.State, len(transitions)),
		needBackup:  make(map[tabledrivenscanner.State]struct{}, len(needBackup)),
	}

	// Transitions
	for k, v := range transitions {
		t.transitions[k] = v
	}

	// Backups
	for _, k := range needBackup {
		t.needBackup[k] = struct{}{}
	}

	// Final states
	for _, k := range finalStates {
		t.finalStates[k] = struct{}{}
	}

	return t
}

// Perform a transition
func (t *CompositeTable) Next(state tabledrivenscanner.State, char rune) tabledrivenscanner.State {
	s, ok := t.transitions[Key{state, char}]
	if !ok {
		// Can try to see if there is an ANY state
		s = t.transitions[Key{state, ANY}]
	}
	return s
}

// Check if a state requires the scanner to backup
func (t *CompositeTable) NeedsBackup(state tabledrivenscanner.State) bool {
	_, ok := t.needBackup[state]
	return ok
}

// The initial state
func (t *CompositeTable) Initial() tabledrivenscanner.State {
	return t.start
}

// Check if a state is a final state
func (t *CompositeTable) IsFinal(state tabledrivenscanner.State) bool {
	_, ok := t.finalStates[state]
	return ok
}

// Generates a token given a State
func (t *CompositeTable) CreateToken(
	state tabledrivenscanner.State,
	lexeme scanner.Lexeme,
	line, col int,
) (scanner.Token, error) {
	symbol, ok := t.tokens[state]

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
