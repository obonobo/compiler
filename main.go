package main

import (
	"bytes"
	"fmt"
	"log"

	ccs "github.com/obonobo/esac/core/chuggingcharsource"
	tdp "github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/tabledrivenparser/compositetable"
	tds "github.com/obonobo/esac/core/tabledrivenscanner"
	scannertable "github.com/obonobo/esac/core/tabledrivenscanner/compositetable"
)

func main() {
	chrs := ccs.MustChuggingReader(bytes.NewBufferString("asdasd"))
	scnr := tds.NewTableDrivenScanner(chrs, scannertable.TABLE())
	prsr := tdp.NewTableDriverParser(scnr, compositetable.TABLE(), log.Default())

	fmt.Println(prsr)
}

func main2() {
	// cmd.RunAndExit()
}
