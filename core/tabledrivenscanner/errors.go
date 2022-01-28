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
