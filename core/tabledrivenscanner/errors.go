package tabledrivenscanner

import (
	"fmt"
)

type UnrecognizedStateError State

func (u UnrecognizedStateError) Error() string {
	return fmt.Sprintf("unrecognized state '%v'", State(u))
}

// The table may generate a PartialTokenError indicating to the algorithm that
// it has reached a final state, but a full token has not been generated. The
// algorithm must record the partial token in its lexeme buffer, but continue as
// if it has received a complete token (i.e. restart at state 1)
type PartialTokenError struct{ Msg string }

func (p PartialTokenError) Error() string {
	msg := "partial token"
	if p.Msg != "" {
		msg += ": " + p.Msg
	}
	return msg
}

// This error is here to reveal any bugs in the transition table that is being
// used by the TableDrivenScanner. If the transition table that is provided to
// the TableDrivenScanner returns NOSTATE when the Table.Next() method is
// called, then there is a case unaccounted for in the DFA implemented by the
// transition table.
//
// Transition table implementations should never return NOSTATE
type NoStateError struct {
	State  State
	Lookup rune
}

func (e NoStateError) Error() string {
	return fmt.Sprintf(""+
		"TableDrivenScanner: NOSTATE returned by Table.Next() "+
		"on state='%v', lookup='%v', this is an error in the "+
		"table implementation being used by this scanner",
		e.State, e.Lookup)
}
