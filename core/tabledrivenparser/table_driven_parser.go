package tabledrivenparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/obonobo/esac/core/parser"
	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/token"
)

// If this option causes problems, disabled it here
var STATEMENT_CLOSER_ENABLED = true

// A parser that is driven by a table
type TableDrivenParser struct {
	scnr  scanner.Scanner    // The lexer, produces a token stream
	table Table              // The parser Table used in the Parse algorithm
	errc  chan<- ParserError // Channel for emitting errors
	rulec chan<- token.Rule  // Channel for emitting which rules are chosen
	stack []token.Kind       // Nonterminal stack
	ast   parser.AST         // Intermediate Representation created by the parser
	err   error              // An error registered by the parser

	// A callback used to control how the parser closes the channels that are
	// provided to it
	ChannelCloser func(errc chan<- ParserError, rulec chan<- token.Rule)
}

// Creates a new TableDriverParser loaded with the scnr, table, and error
// channel errc. If the scnr or table is nil, then this function returns nil.
func NewTableDrivenParser(
	scnr scanner.Scanner,
	table Table,
	errc chan<- ParserError,
	rulec chan<- token.Rule,
) *TableDrivenParser {
	if scnr == nil || table == nil {
		return nil // Nil scanner or table is not workable
	}
	return &TableDrivenParser{
		scnr:          scnr,
		table:         table,
		errc:          errc,
		rulec:         rulec,
		ChannelCloser: CloseChannels,
	}
}

// Wraps your scanner in a scanner.CommentlessScanner
func NewTableDrivenParserIgnoringComments(
	scnr scanner.Scanner,
	table Table,
	errc chan<- ParserError,
	rulec chan<- token.Rule,
	comments ...token.Kind,
) *TableDrivenParser {
	return NewTableDrivenParser(
		scanner.IgnoringComments(scnr, comments...),
		table, errc, rulec)
}

func (t *TableDrivenParser) AST() parser.AST {
	return t.ast
}

// Parses the token stream that is loaded in the Parser. Returns true if the
// parse was successful, false otherwise. The TableDrivenParser's errc will be
// close at the end of this method.
func (t *TableDrivenParser) Parse() bool {
	defer t.closeChannels()

	if t.err != nil {
		return false
	}

	t.push(t.table.Start())
	a, err := t.scnr.NextToken()
	if err != nil {
		t.err = err
		t.emitError(t.err, token.Token{})
		return false
	}

	for prev := a; !t.empty(); prev = a {
		x := t.top()
		if t.table.IsTerminal(x) {
			if x == a.Id {
				t.pop()
				a, err = t.scnr.NextToken()
				if err != nil {
					// This error will probably be EOF, in any case we can't
					// continue with no tokens. EOF does not need to be
					// registered on the parser, eat the error
					if errors.Is(err, io.EOF) {
						// Pop stack symbols that have EPSILON rules
						for !t.empty() {
							if !t.table.HasEpsilonRule(t.top()) {
								t.err = t.unterminatedSentenceError()
								t.emitError(t.err, prev)
								break
							}
							t.pop()
						}
					} else {
						t.err = err
						t.emitError(t.err, prev)
					}
					break
				}
			} else {
				t.emitError(fmt.Errorf("expected to find symbol %v but got %v", x, a.Id), a)
				aa, err := t.skipErrors3(a)
				a = aa
				if err != nil {
					break
				}
			}
		} else {
			if l, err := t.table.Lookup(x, a.Id); err == nil {
				t.emitRule(l)
				t.pop()
				t.inverseRHSMultiplePush(l.RHS...)
			} else {
				// This branch indicates that the lookup failed to return a
				// table cell - i.e. that TT[x, a] is empty and thus the parser
				// was not expecting to encounter this token at this point
				t.emitError(fmt.Errorf("no rule for nonterminal %v, token %v: %w", x, a.Id, err), a)
				aa, err := t.skipErrors3(a)
				a = aa
				if err != nil {
					break
				}
			}
		}
	}

	return t.err == nil && t.empty()
}

// SKIPERRORS VARIATION 3: this version of skip errors scans up to FIRST(top)
// (no pop) or FOLLOW(top) (yes pop)
//
// Small change: STATEMENT CLOSER
//     If the stack contains a semi-colon terminated slice on top and skipErrors
//     encounters a semi colon, skipErrors may trim the entire slice from the
//     top, thereby discarding the remainder of the statement. STATEMENT CLOSER
//     must be checked last, after checking for FIRST(top) and FOLLOW(top). If
//     skipErrors discards symbols via STATEMENT CLOSER, then skipErrors must
//     report all discarded symbols as errors. STATEMENT CLOSER may be checked
//     only once.
func (t *TableDrivenParser) skipErrors3(lookahead token.Token) (token.Token, error) {
	statementCloserEnabled := false
	// statementCloserEnabled := STATEMENT_CLOSER_ENABLED

	if contains(t.follow(t.top()), lookahead.Id) || lookahead.Id == "" {
		t.pop()
		return lookahead, nil
	}

	// If we encounter a symbol in FOLLOW(top), then we will pop and continue
	// from that symbol
	syncTokensPop := t.follow(t.top())

	// If we encounter a symbol in FIRST(top), then we will continue starting
	// from that symbol
	syncTokensNoPop := t.first(t.top())
	if contains(syncTokensNoPop, token.EPSILON) {
		syncTokensNoPop = union(syncTokensNoPop, syncTokensPop)
	}

	// STATEMENT CLOSER
	var statement []token.Kind
	var hasStatement bool
	if statementCloserEnabled {
		statement, hasStatement = t.statementCloser()
	}

	for {
		// Check the no pops first
		if contains(syncTokensNoPop, lookahead.Id) {
			break
		}

		// Check the pops
		if contains(syncTokensPop, lookahead.Id) {
			t.pop()
			break
		}

		// Check statement closer
		if statementCloserEnabled {
			statementCloserEnabled = false // Check only once
			if lookahead.Id == token.SEMI && hasStatement {
				// Then we can trim the entire statement, minus the SEMI
				for i := len(statement) - 1; i > 0; i-- {
					t.emitError(fmt.Errorf(
						"skipErrors() closing statement: expected %v", statement[i]),
						token.Token{})
					t.pop()
				}
				return lookahead, nil
			}
		}

		l, err := t.scnr.NextToken()
		if err != nil {
			t.emitError(fmt.Errorf("skipErrors() failed to scan: %w", err), lookahead)
			return l, t.err
		}
		lookahead = l
	}

	return lookahead, nil
}

func (t *TableDrivenParser) statementCloser() ([]token.Kind, bool) {
	var found bool
	var i int
	for i = len(t.stack) - 1; i >= 0; i-- {
		if s := t.stack[i]; s == token.SEMI {
			found = true
			break
		}
	}
	return t.stack[i:], found
}

func (t *TableDrivenParser) top() token.Kind {
	if t.empty() {
		return ""
	}
	return t.stack[len(t.stack)-1]
}

func (t *TableDrivenParser) pop() token.Kind {
	if t.empty() {
		return ""
	}
	top := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return top
}

func (t *TableDrivenParser) push(symbol token.Kind) {
	t.stack = append(t.stack, symbol)
}

// Pushes all the symbols onto the stack in the reverse order in which they are
// provided
func (t *TableDrivenParser) inverseRHSMultiplePush(symbols ...token.Kind) {
	for i := len(symbols) - 1; i >= 0; i-- {
		t.push(symbols[i])
	}
}

func (t *TableDrivenParser) empty() bool {
	return len(t.stack) == 0
}

func (t *TableDrivenParser) emitError(err error, tok token.Token) {
	t.err = err
	if t.errc != nil {
		t.errc <- ParserError{
			Err: t.err,
			Tok: tok,
			Sym: t.top(),
		}
	}
}

func (t *TableDrivenParser) emitRule(rule token.Rule) {
	if t.rulec != nil {
		t.rulec <- rule
	}
}

func (t *TableDrivenParser) unterminatedSentenceError() error {
	remainingSymbols := make([]token.Kind, 0, len(t.stack))
	for i := len(t.stack) - 1; i >= 0; i-- {
		remainingSymbols = append(remainingSymbols, t.stack[i])
	}
	return &UnterminatedSentence{remainingSymbols}
}

func (t *TableDrivenParser) first(symbol token.Kind) token.KindSet {
	fs, _ := t.table.First(symbol)
	return fs
}

func (t *TableDrivenParser) follow(symbol token.Kind) token.KindSet {
	fs, _ := t.table.Follow(symbol)
	return fs
}

func (t *TableDrivenParser) closeChannels() {
	CloseChannels(t.errc, t.rulec)
}

func contains(haystack token.KindSet, needle token.Kind) bool {
	_, ok := haystack[needle]
	return ok
}

func CloseChannels(errc chan<- ParserError, rulec chan<- token.Rule) {
	if errc != nil {
		close(errc)
	}
	if rulec != nil {
		close(rulec)
	}
}

func union(s1, s2 token.KindSet) token.KindSet {
	s := make(token.KindSet, len(s1)+len(s2))
	for k := range s1 {
		s[k] = struct{}{}
	}
	for k := range s2 {
		s[k] = struct{}{}
	}
	return s
}
