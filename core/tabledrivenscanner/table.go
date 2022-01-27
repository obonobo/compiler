package tabledrivenscanner

import "github.com/obonobo/compiler/core/scanner"

const (
	ANY    rune = 0  // Represents any character
	LETTER rune = -1 // Represents expression [aA-zZ]
)

const (
	START   State = 1    // A suggested starting State for table implementations
	NOSTATE State = -666 // A State that is not attached to DFA
)

type State int

// State transition table
type Table interface {
	// Perform a transition
	Next(state State, char rune) State

	// Check if a state requires the scanner to backup
	NeedsBackup(state State) bool

	// The initial state
	Initial() State

	// Check if a state is a final state
	IsFinal(state State) bool

	// Generates a token given a State
	CreateToken(state State, lexeme scanner.Lexeme, line, col int) (scanner.Token, error)
}
