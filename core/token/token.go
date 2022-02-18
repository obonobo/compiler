package token

import (
	"fmt"

	"github.com/obonobo/esac/util"
)

// The text representing a token
type Lexeme string

type Token struct {
	Id     Kind   // The unique identifier of this token
	Lexeme Lexeme // The exact string that was matched as this token
	Line   int    // The line number on which the token was found
	Column int    // The column number on which the token was found
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token[Id=%v, Lexeme=%v, Line=%v, Column=%v]",
		t.Id, t.Lexeme, t.Line, t.Column)
}

// A premade mapper function to be used with 'reporting.TranformTokenStream'
func (t Token) Report() string {
	return fmt.Sprintf("[%v, %v, %v]", t.Id, util.SingleLinify(string(t.Lexeme)), t.Line)
}

// Set of reserved words (empty structs as values to allocate 0 memory)
var reservedWords = map[Kind]struct{}{
	IF:       {},
	THEN:     {},
	ELSE:     {},
	INTEGER:  {},
	FLOAT:    {},
	VOID:     {},
	PUBLIC:   {},
	PRIVATE:  {},
	FUNC:     {},
	VAR:      {},
	STRUCT:   {},
	WHILE:    {},
	READ:     {},
	WRITE:    {},
	RETURN:   {},
	SELF:     {},
	INHERITS: {},
	LET:      {},
	IMPL:     {},
}

var errorSymbols = map[Kind]struct{}{
	INVALIDNUM:          {},
	INVALIDCHAR:         {},
	INVALIDID:           {},
	UNTERMINATEDCOMMENT: {},
}

func IsReservedWord(s Kind) bool {
	_, ok := reservedWords[s]
	return ok
}

func IsReservedWordString(s string) (Kind, bool) {
	t := Kind(s)
	return t, IsReservedWord(t)
}

// Returns all reserved word token kinds
func ReservedWords() []Kind {
	return setToSlice(reservedWords)
}

func IsError(s Kind) bool {
	_, ok := errorSymbols[s]
	return ok
}

// Returns all error token Kinds
func ErrorTokens() []Kind {
	return setToSlice(errorSymbols)
}

func setToSlice(set map[Kind]struct{}) []Kind {
	ret := make([]Kind, 0, len(set))
	for k := range set {
		ret = append(ret, k)
	}
	return ret
}
