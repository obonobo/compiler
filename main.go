package main

import (
	"bytes"
	"fmt"

	"github.com/obonobo/compiler/core/chuggingcharsource"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner/compositetable"
)

func main() {

	x := 1.
	fmt.Println(x)


	// Create a char
	charSource := new(chuggingcharsource.ChuggingCharSource)
	err := charSource.ChugReader(bytes.NewBufferString("/"))
	if err != nil {
		fmt.Println(err)
	}

	scanner := tabledrivenscanner.NewTableDrivenScanner(charSource, compositetable.TABLE)
	token, err := scanner.NextToken()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(token)
}
