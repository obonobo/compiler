package tabledrivenparser

import (
	"errors"
	"log"
	"os"

	"github.com/obonobo/esac/core/scanner"
)

// A parser that is driven by a table
type TableDrivenParser struct {
	scnr   scanner.Scanner // The lexer, produces a token stream
	table  Table           // The parser Table used in the Parse algorithm
	logger *log.Logger     // A logger to print out
	stack  []string        // Nonterminal stack
	err    error           // An error registered by the parser
}

// Creates a new TableDriverParser loaded with the scnr and logger. Will panic
// if scnr is nil, though nil logger is okay
func NewTableDriverParser(
	scnr scanner.Scanner,
	table Table,
	logger *log.Logger,
) *TableDrivenParser {
	if scnr == nil {
		panic(ErrNilScanner) // Nil scanner is not workable
	}
	if logger == nil {
		logger = log.New(os.Stderr, "", 0) // Nil logger prints to stderr
	}
	return &TableDrivenParser{scnr: scnr, table: table, logger: logger}
}

// Parses the token stream that is loaded in the Parser. Returns true if the
// parse was successful, false otherwise
func (t *TableDrivenParser) Parse() bool {
	if t.err != nil {
		return false
	}

	t.push(t.table.Start())
	a, err := t.scnr.NextToken()
	if err != nil {
		// TODO: add appropriate error handling here
		t.err = err
		return false
	}

	for !t.empty() {
		x := t.top()
		if t.table.IsTerminal(x) {

		} else {
			l := t.table.Lookup(x, string(a.Id))
			if !t.table.IsError(l) {
				t.pop()

				// TODO: change the lookup method so that it actually returns a
				// TODO: list of the items that need to be RHS pushed here
				t.inverseRHSMultiplePush(l)
			} else {
				t.skipErrors()

				// TODO: implement errors
				t.err = errors.New("some error")
			}
		}
	}

	return false
}

func (t *TableDrivenParser) skipErrors() {
	// TODO: implement method
	panic("not implemented")
}

func (t *TableDrivenParser) top() string {
	if t.empty() {
		return ""
	}
	return t.stack[len(t.stack)-1]
}

func (t *TableDrivenParser) pop() string {
	if t.empty() {
		return ""
	}
	top := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return top
}

func (t *TableDrivenParser) push(symbol string) {
	t.stack = append(t.stack, symbol)
}

// Pushes all the symbols onto the stack in the reverse order in which they are
// provided
func (t *TableDrivenParser) inverseRHSMultiplePush(symbols ...string) {
	for i := len(symbols) - 1; i >= 0; i-- {
		t.push(symbols[i])
	}
}

func (t *TableDrivenParser) empty() bool {
	return len(t.stack) == 0
}
