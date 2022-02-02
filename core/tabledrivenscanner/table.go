package tabledrivenscanner

import "github.com/obonobo/esac/core/scanner"

const (
	ANY    rune = -2 // Represents any character
	LETTER rune = -1 // Represents expression [aA-zZ]
)

const (
	START               State = 1    // A suggested starting State for table implementations
	NOSTATE             State = -666 // A State that is not attached to DFA
	UNTERMINATEDCOMMENT State = -999 // A State that represents unterminated comments
)

type State int

// State transition table
type Table interface {
	// The initial state
	Initial() State

	// Perform a transition
	Next(state State, char rune) State

	// Generates a token given a State
	CreateToken(state State, lexeme scanner.Lexeme, line, col int) (scanner.Token, error)

	// Check if a state requires the scanner to backup
	NeedsBackup(state State) bool

	// Check if a state requires the scanner to backup TWICE
	NeedsDoubleBackup(state State) bool

	// Check if a state is a final state
	IsFinal(state State) bool

	// Checks if a symbol is whitespace
	IsWhiteSpace(char rune) bool

	// Report whether the table thinks it is reading a comment
	InCommentMode() bool

	// Returns the state that the table is using to represent unterminated
	// comments
	UnterminatedCommentState() State
}
