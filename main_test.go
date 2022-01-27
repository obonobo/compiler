package main

import (
	"bytes"
	"testing"

	"github.com/obonobo/compiler/core/chuggingcharsource"
	"github.com/obonobo/compiler/core/scanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner/compositetable"
)

func TestInlineComment(t *testing.T) {
	t.Parallel()

	comment := "// asdasdasd \n"

	charSource := new(chuggingcharsource.ChuggingCharSource)
	err := charSource.ChugReader(bytes.NewBufferString(comment))
	if err != nil {
		t.Fatalf("ChugReader should succeed here: %v", err)
	}

	scan := tabledrivenscanner.NewTableDrivenScanner(charSource, compositetable.TABLE)
	actual, err := scan.NextToken()
	if err != nil {
		t.Fatalf("NextToken should succeed: %v", err)
	}

	expected := scanner.Token{
		Id:     scanner.INLINECMT,
		Lexeme: scanner.Lexeme(comment),
		Line:   1,
		Column: 1,
	}

	if actual != expected {
		t.Fatalf("Expected token %v but got %v", expected, actual)
	}
}
