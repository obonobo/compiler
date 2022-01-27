package tabledrivenscanner

import (
	"fmt"

	"github.com/obonobo/compiler/core/scanner"
)

type TableDrivenScanner struct {
	chars scanner.CharSource
	table Table
}

func NewTableDrivenScanner(chars scanner.CharSource) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
	}
}

func NewTableDrivenScannerWithTable(chars scanner.CharSource, table Table) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
		table: table,
	}
}

func (t *TableDrivenScanner) nextChar() (scanner.Symbol, error) {
	r, err := t.chars.NextChar()
	if err != nil {
		return "", fmt.Errorf("TableDriverScanner: failed to retrieve next char: %w", err)
	}
	return scanner.Symbol(r), nil
}

func (t *TableDrivenScanner) backupChar() {
	t.chars.BackupChar()
}

func (t *TableDrivenScanner) NextToken() (scanner.Token, error) {
	var token *scanner.Token
	state := t.table.Initial()
	for {
		lookup, err := t.nextChar()
		if err != nil {
			return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
		}

		state = t.table.Next(state, lookup)
		if t.table.IsFinal(state) {
			tt, err := t.table.CreateToken(state, "", 0, 0)
			if err != nil {
				return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
			}
			token = &tt

			if t.table.NeedsBackup(state) {
				t.backupChar()
			}
		}

		if token != nil {
			break
		}
	}
	return *token, nil
}
