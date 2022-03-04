package tabledrivenparser

import (
	"errors"
	"fmt"
	"io"

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
	ast   token.AST          // Intermediate Representation created by the parser
	err   error              // An error registered by the parser

	stack    []token.Kind     // Nonterminal stack
	semStack []*token.ASTNode // Semantic stack

	// A callback used to control how the parser closes the channels that are
	// provided to it. This function will be called at the end of parser.Parse()
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
		scnr:  scnr,
		table: table,
		errc:  errc,
		rulec: rulec,

		stack:    make([]token.Kind, 0, 1024),
		semStack: make([]*token.ASTNode, 0, 1024),

		// Default channel closer, if you want to change this, you'll have to
		// modify it after instantiating the parser
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

// Same as NewTableDrivenParserIgnoringComments except uses the comment tokens
// declared in token package
func NewTableDrivenParserIgnoringDefaultComments(
	scnr scanner.Scanner,
	table Table,
	errc chan<- ParserError,
	rulec chan<- token.Rule,
) *TableDrivenParser {
	return NewTableDrivenParserIgnoringComments(scnr, table, errc, rulec, token.Comments()...)
}

func (t *TableDrivenParser) AST() token.AST {
	if t.ast.Root == nil {
		l := len(t.semStack)
		if l != 1 {
			panic(fmt.Errorf("stack should be [PROG] but got %v", t.semStack))
		}
		top := t.semStack[l-1]
		t.ast.Root = top
	}
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

	// for prev := a; !t.empty(); prev = a {
	for prev := a; !t.empty(); {
		x := t.top()

		// Check semantic action first
		if t.isSemAction(x) {
			t.executeSemanticAction(x, prev)
			t.pop()

		} else if t.table.IsTerminal(x) { // Check terminals second
			if x == a.Id {
				t.pop()
				prev = a
				a, err = t.scnr.NextToken()
				if err != nil {
					// This error will probably be EOF, in any case we can't
					// continue with no tokens. EOF does not need to be
					// registered on the parser, eat the error
					if errors.Is(err, io.EOF) {
						// Pop stack symbols that have EPSILON rules, and
						// process remaining semantic actions
						for !t.empty() {
							x = t.top()
							if t.isSemAction(x) {
								t.executeSemanticAction(x, prev)
							} else if !t.table.HasEpsilonRule(x) {
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
				t.emitError(&UnexpectedTokenExpectedInsteadError{Token: a, Instead: x}, a)
				aa, err := t.skipErrors3(a)
				a = aa
				if err != nil {
					break
				}
			}

		} else { // Check nonterminals last
			if l, err := t.table.Lookup(x, a.Id); err == nil {
				t.emitRule(l)
				t.pop()
				t.inverseRHSMultiplePush(l.RHS...)
			} else {
				t.emitError(&UnexpectedTokenError{Token: a, Err: err}, a)
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
	// statementCloserEnabled := false
	statementCloserEnabled := STATEMENT_CLOSER_ENABLED

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
	var hasStatement func(token.Token) bool
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
			if hasStatement(lookahead) {
				// Then we can trim the entire statement, minus the SEMI
				for i := len(statement) - 1; i > 0; i-- {

					// ? Maybe we should do something with the discarded symbols
					// t.emitError(fmt.Errorf(
					// 	"skipErrors() closing statement: expected %v", statement[i]),
					// 	token.Token{})

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

// The statement closer is extra functionality to avoid eating a great many
// tokens when skipErrors() is called. Its utility is for identifying "missing
// token" errors - those errors were there are tokens removed from a construct,
// tokens that are expected to be seen but are missing.
//
// Normally, skipErrors() will drop tokens until it finds a token in FIRST or
// FOLLOW set of the stack-top symbol.
//
// With so-called "missing token" errors, this could lead to dropping way too
// many tokens, so my solution is the statementCloser() function.
//
// How it works is it does a short search through the stack for specific
// construct ending symbols like ';' or '{' or '}'. If we have a token of this
// symbol in hand and there is this same token near the top of the stack,
// perhaps this is a missing-token error? Then the caller of statementCloser()
// can skip the missing symbols near the top of the stack and resume parsing at
// the statement closer symbol e.g.: ';'
//
// This function is essentially the reverse of skipErrors(). Where skipErrors()
// drops tokens, statementCloser() drops stack symbols. Where skipErrors()
// freezes the stack-top in place, statementCloser freezes the lookahead
// instead.
func (t *TableDrivenParser) statementCloser() (
	statement []token.Kind,
	hasStatement func(lookahead token.Token) bool,
) {
	l := len(t.stack)
	if l == 0 {
		return []token.Kind{}, func(lookahead token.Token) bool { return false }
	}

	// These tokens are closable statements, hardcoding this for now. If the
	// feature is useful, then I'll make it configurable
	closers := []token.Kind{
		token.SEMI,
		token.OPENCUBR,
	}

	// We can match against FIRST set (or FOLLOW set if FIRST contains epsilon)
	isCloser := func(symbol token.Kind) (token.Kind, bool) {
		frst := t.first(symbol)
		if contains(frst, token.EPSILON) {
			frst = union(frst, t.follow(symbol))
		}
		for _, c := range closers {
			if contains(frst, c) {
				return c, true
			}
		}
		return "", false
	}

	var found bool
	var foundSymbol token.Kind
	var i int
	for i = l - 1; i >= 0; i-- {
		s := t.stack[i]
		if closer, contains := isCloser(s); contains {
			found = true
			foundSymbol = closer
			break
		}
	}

	hasStatement = func(lookahead token.Token) bool {
		return found && lookahead.Id == foundSymbol
	}

	if i == -1 {
		// Then we've read the entire stack
		return t.stack, hasStatement
	}
	return t.stack[i:], hasStatement
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

func (t *TableDrivenParser) semEmpty() bool {
	return len(t.semStack) == 0
}

func (t *TableDrivenParser) semTop() *token.ASTNode {
	if t.semEmpty() {
		return nil
	}
	return t.semStack[len(t.semStack)-1]
}

func (t *TableDrivenParser) semPop() *token.ASTNode {
	var top *token.ASTNode
	if top = t.semTop(); top == nil {
		return top
	}
	t.semStack = t.semStack[:len(t.semStack)-1]
	return top
}

func (t *TableDrivenParser) semPush(record *token.ASTNode) {
	t.semStack = append(t.semStack, record)
}

func (t *TableDrivenParser) isSemAction(symbol token.Kind) bool {
	return token.IsSemAction(symbol)
}

func (t *TableDrivenParser) executeSemanticAction(x token.Kind, previousToken token.Token) {
	if action, ok := token.SEM_DISPATCH[x]; ok {
		action(x, previousToken, &t.semStack)
	} else {
		// TODO: remove this panic
		panic(fmt.Errorf("no semantic action found, x = %v, previousToken = %v", x, previousToken))
	}
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
	t.ChannelCloser(t.errc, t.rulec)
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
