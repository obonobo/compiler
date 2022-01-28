package main

import (
	"bytes"
	"fmt"

	"github.com/obonobo/compiler/core/chuggingcharsource"
	"github.com/obonobo/compiler/core/tabledrivenscanner"
	"github.com/obonobo/compiler/core/tabledrivenscanner/compositetable"
)

func main() {

	// Create a char
	charSource := new(chuggingcharsource.ChuggingCharSource)
	err := charSource.ChugReader(bytes.NewBufferString("1.0 example_id\n id2 id3"))
	if err != nil {
		fmt.Println(err)
	}

	scanner := tabledrivenscanner.NewTableDrivenScanner(charSource, compositetable.TABLE)

	for i := 0; i < 4; i++ {
		token, err := scanner.NextToken()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(token)
	}
}
