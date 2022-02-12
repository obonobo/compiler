package cmd

import (
	"flag"
	"fmt"
	"os"
)

const ROOT = "root"

type Config struct {
	Command    string
	Subcommand string
}

// Runs the CLI with os.Args and exits with the returned exit code
func RunAndExit() {
	os.Exit(Run(os.Args))
}

// Runs the CLI tool and returns the program exit code. The main function will
// need to manually exit with this code.
func Run(args []string) (exitCode int) {
	config := &Config{Command: args[0]}
	rootCmd := flag.NewFlagSet(ROOT, flag.ExitOnError)
	helpFlag := rootCmd.Bool("help", false, "")
	rootCmd.BoolVar(helpFlag, "h", false, "")
	rootCmd.Usage = func() { printHelp(config) }
	rootCmd.Parse(args[1:])
	if *helpFlag || len(args) < 2 {
		rootCmd.Usage()
		return 1
	}

	lexUsage, lex := lexCmd(config)
	buildUsage, build := buildCmd(config)
	help := helpCmd(config, map[string]func(){
		LEX:   lexUsage,
		BUILD: buildUsage,
	})

	config.Subcommand = args[1]
	rest := args[2:]

	switch config.Subcommand {
	case HELP:
		return help(rest)
	case LEX:
		return lex(rest)
	case BUILD:
		return build(rest)
	default:
		fmt.Println(unknownCommand(config.Command, config.Subcommand))
		return 1
	}
}
