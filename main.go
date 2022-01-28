package main

import (
	"os"

	"github.com/obonobo/esac/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args))
}
