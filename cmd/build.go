package main

import (
	"flag"
	"fmt"
)

const BUILD = "build"

type BuildParams struct{}

func buildCmd(config *Config) (usage func(), action func(args []string) int) {
	buildCmd := flag.NewFlagSet(BUILD, flag.ExitOnError)
	buildCmd.Usage = func() {
		fmt.Printf(`The "%v" subcommand has not yet been implemented...`+"\n", BUILD)
	}

	return buildCmd.Usage, func(args []string) int {
		var params BuildParams
		buildCmd.Parse(args)
		fmt.Printf(`The '%v build' subcommand has not yet been implemented...`+"\n", config.Command)
		fmt.Printf(`Run '%v help' for usage.`+"\n", config.Command)
		return Build(params)
	}
}

func Build(params BuildParams) (exit int) {
	return 1
}
