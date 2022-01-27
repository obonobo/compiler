package tabledrivenscanner

import (
	"fmt"

	"github.com/obonobo/compiler/core/scanner"
)

type TableDrivenScanner struct {
	chars  scanner.CharSource // A source of characters
	table  Table              // A table for performing transitions
	lexeme scanner.Lexeme     // The lexeme that is being built
}

func NewTableDrivenScanner(chars scanner.CharSource, table Table) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
		table: table,
	}
}

// Scans for the next token present in the character source
func (t *TableDrivenScanner) NextToken() (scanner.Token, error) {
	var token *scanner.Token
	state := t.table.Initial()
	for {
		lookup, err := t.nextChar()
		if err != nil {
			return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
		}

		state = t.table.Next(state, lookup)
		if state == 0 {
			return scanner.Token{},
				fmt.Errorf("TableDrivenScanner: no possible transition (state = 0)")
		}

		if t.table.IsFinal(state) {
			tt, err := t.table.CreateToken(state, t.lexeme, 0, 0)
			if err != nil {
				return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
			}
			token = &tt

			if t.table.NeedsBackup(state) {
				t.backup()
			}
		}

		if token != nil {
			break
		}
	}
	return *token, nil
}

func (t *TableDrivenScanner) backup() error {
	_, err := t.chars.BackupChar()
	t.lexeme = t.lexeme[:len(t.lexeme)-1]
	return err
}

func (t *TableDrivenScanner) nextChar() (rune, error) {
	r, err := t.chars.NextChar()
	t.lexeme += scanner.Lexeme(r)
	return r, err
}
