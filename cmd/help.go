package cmd

import (
	"flag"
	"fmt"
	"path"
	"strings"
)

const HELP = "help"

const USAGE = `%v compiles LANG source code into MOON assembly

Usage:
	%v <command> [arguments]

The commands are:

	lex	scan input files, convert them to tokens
	parse	parses token stream, converts it to AST
	build	compile code

Use "%v help <command>" for more information about a command.
`

func helpCmd(config *Config, usages map[string]func()) (action func(args []string) int) {
	helpCmd := flag.NewFlagSet(HELP, flag.ExitOnError)
	return func(args []string) int {
		helpCmd.Parse(args)
		if len(args) < 1 {
			printHelp(config)
			return 1
		}
		needHelpWith := helpCmd.Arg(0)
		if printUsage, ok := usages[needHelpWith]; ok {
			printUsage()
			return 1
		}
		c := path.Base(config.Command)
		fmt.Printf(
			"%v %v %v: unknown command. Run '%v %v'.\n",
			c, HELP, needHelpWith, c, HELP)
		return 1
	}
}

func unknownCommand(command, subcommand string) string {
	return fmt.Sprintf(
		"%v %v: unknown command\nRun '%v help' for usage.",
		command, subcommand, command)
}

func usage(command string) string {
	cmd := path.Base(command)
	return fmt.Sprintf(USAGE, strings.ToUpper(cmd), cmd, cmd)
}

func printHelp(config *Config) {
	fmt.Print(usage(config.Command))
}
