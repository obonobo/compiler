package tabledrivenscanner

import "github.com/obonobo/compiler/core/scanner"

type TableDrivenScanner struct {
	chars scanner.CharSource
}

func NewTableDrivenScanner(chars scanner.CharSource) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
	}
}

func (t *TableDrivenScanner) NextToken() scanner.Token {
	panic("not implemented") // TODO: Implement
}
