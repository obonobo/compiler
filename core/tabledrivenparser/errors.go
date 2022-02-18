package tabledrivenparser

import (
	"errors"
	"fmt"

	"github.com/obonobo/esac/core/token"
)

var (
	ErrNilScanner = errors.New("scanner cannot be nil")
)

type NoRuleError struct {
	Row, Col token.Kind
}

func (e *NoRuleError) Error() string {
	return fmt.Sprintf("no entry for Key{%v, %v}", e.Row, e.Col)
}

type UnterminatedSentence struct {
	expectedSymbols []token.Kind
}

func (e *UnterminatedSentence) Error() string {
	return fmt.Sprintf(
		"got EOF in the middle of a sentence, expected to find symbols for %v",
		e.expectedSymbols)
}

type ParserError struct {
	Err error       // An error message
	Tok token.Token // The most recent token
	Sym token.Kind  // The symbol on top of the stack
}

func (e ParserError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
