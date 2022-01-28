package tabledrivenscanner

import (
	"fmt"

	"github.com/obonobo/compiler/core/scanner"
)

type lexemeSpec struct {
	s    scanner.Lexeme
	line int
	col  int
}

type TableDrivenScanner struct {
	chars scanner.CharSource // A source of characters
	table Table              // A table for performing transitions

	// The lexeme that is being built
	lexeme lexemeSpec
}

func NewTableDrivenScanner(chars scanner.CharSource, table Table) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
		table: table,
		lexeme: lexemeSpec{
			line: chars.Line(),
			col:  chars.Column(),
		},
	}
}

// Scans for the next token present in the character source
func (t *TableDrivenScanner) NextToken() (scanner.Token, error) {
	var token *scanner.Token
	state := t.table.Initial()
	for {
		lookup, err := t.nextChar()
		if err != nil {
			// We are out of input, if there is an ANY transition available,
			// then we can take it, otherwise return the error
			state = t.table.Next(state, ANY)
			if state == NOSTATE {
				return scanner.Token{}, fmt.Errorf("TableDrivenScanner.NextToken(): %w", err)
			}
		} else {
			state = t.table.Next(state, lookup)
		}

		if state == NOSTATE {
			return scanner.Token{},
				fmt.Errorf("TableDrivenScanner: no possible transition")
		}

		if t.table.IsFinal(state) {
			doubleBackTrack := t.table.NeedsDoubleBackup(state)
			backtrack := t.table.NeedsBackup(state)
			if !backtrack && !doubleBackTrack {
				t.pushLexeme(lookup)
			} else if doubleBackTrack {
				t.popLexeme()
			}

			token, err = t.createToken(state)
			if err != nil {
				return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
			}

			if backtrack {
				t.backup()
			} else if doubleBackTrack {
				t.backup()
				t.backup()
			}
		}

		if token != nil {
			break
		}
		t.pushLexeme(lookup)
	}
	return *token, nil
}

func (t *TableDrivenScanner) pushLexeme(char rune) {
	isWhiteSpace := (char == ' ' || char == '\n') && len(t.lexeme.s) == 0
	if !isWhiteSpace {
		t.lexeme.s += scanner.Lexeme(char)
	} else {
		t.resetLexeme()
	}
}

func (t *TableDrivenScanner) popLexeme() {
	if len(t.lexeme.s) > 0 {
		t.lexeme.s = t.lexeme.s[:len(t.lexeme.s)-1]
	}
}

func (t *TableDrivenScanner) resetLexeme() {
	t.lexeme.s = ""
	t.lexeme.col = t.chars.Column()
	t.lexeme.line = t.chars.Line()
}

func (t *TableDrivenScanner) backup() error {
	_, err := t.chars.BackupChar()
	t.lexeme.col = t.chars.Column()
	t.lexeme.line = t.chars.Line()
	return err
}

func (t *TableDrivenScanner) nextChar() (rune, error) {
	r, err := t.chars.NextChar()
	return r, err
}

func (t *TableDrivenScanner) createToken(
	state State,
) (*scanner.Token, error) {
	tt, err := t.table.CreateToken(state, t.lexeme.s, t.lexeme.line, t.lexeme.col)
	t.resetLexeme()
	return &tt, err
}
