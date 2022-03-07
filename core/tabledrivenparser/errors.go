package tabledrivenparser

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/token"
)

var (
	ErrNilScanner = fmt.Errorf("scanner cannot be nil")
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

func (e *ParserError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e *ParserError) Unwrap() error {
	return e.Err
}

type UnexpectedTokenError struct {
	// The unexpected token
	Token token.Token

	// The token that was expected instead of this token
	Instead token.Kind

	// A set of tokens that was expected instead of this token. Overrides
	// UnexpectedTokenExpectedInsteadError.Instead
	InsteadSlice []token.Kind

	Err error // wrap
}

func (e *UnexpectedTokenError) Error() string {
	switch l := len(e.InsteadSlice); l {
	case 1:
		e.Instead = e.InsteadSlice[0]
		fallthrough
	case 0:
		if e.Instead != "" {
			return fmt.Sprintf(
				"unexpected token '%v', should be '%v'",
				e.Token.Id, e.Instead)
		}
		return fmt.Sprintf("unexpected token '%v'", e.Token.Id)
	default:
		s := make([]string, 0, l)
		for _, k := range e.InsteadSlice {
			s = append(s, fmt.Sprintf("'%v'", k))
		}
		return fmt.Sprintf(
			"unexpected token '%v', should be %v, or %v",
			e.Token.Id, strings.Join(s[:l-1], ", "), s[l-1])
	}
}

func (e *UnexpectedTokenError) Unwrap() error {
	return e.Err
}

// Use for determining possible expected tokens, parser tables may return errors
// that implement this interface, in which case the TableDrivenParser may use
// the LookupPossibilities.Possibilities() method to ascertain more information
// about the lookup error
type LookupPossibilities interface {

	// Returns possible tokens that could result in a successful lookup
	Possibilities() []token.Kind
}
