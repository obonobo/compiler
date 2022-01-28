package main

import (
	"fmt"
	"log"
	"os"

	"github.com/obonobo/esac/core/chuggingcharsource"
	"github.com/obonobo/esac/core/tabledrivenscanner"
	"github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
)

func main() {
	path := "./resources/handout/lexpositivegrading.src"
	fh, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	charSource := new(chuggingcharsource.ChuggingCharSource)
	err = charSource.ChugReader(fh)
	fh.Close()
	if err != nil {
		log.Fatal(err)
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
