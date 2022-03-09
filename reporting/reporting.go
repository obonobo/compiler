package reporting

import (
	"fmt"
	"log"
	"strings"

	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/util"
)

// Spools strings from data chan to logger. Prints one string per line.
func LogSpool(data <-chan string, logger *log.Logger) (done <-chan struct{}) {
	donec := make(chan struct{}, 1)
	go func() {
		for s := range data {
			logger.Println(s)
		}
		donec <- struct{}{}
	}()
	return donec
}

// Spools errors reported by the parser, logs them on the provided logger
func ErrSpool(logger *log.Logger) chan<- tabledrivenparser.ParserError {
	errc := make(chan tabledrivenparser.ParserError, 1024)
	go func() {
		for err := range errc {
			if logger != nil {
				logger.Print(ParserErrorPrintout(err))
			}
		}
	}()
	return errc
}

func ParserErrorPrintout(err tabledrivenparser.ParserError) string {
	return fmt.Sprintf(
		"Syntax error on line %v, column %v: %v",
		err.Tok.Line, err.Tok.Column, err.Err)
}

// Spools rules reported by the parser, logs them on the provided logger
func RuleSpool(logger *log.Logger) chan<- token.Rule {
	rulec := make(chan token.Rule, 1024)
	go func() {
		for err := range rulec {
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	return rulec
}

func StreamLinesSplitErrors(
	scnr scanner.Scanner,
	bufSize int,
) (tokens, errors <-chan string) {
	return StreamLinesOptionallySplitErrors(scnr, bufSize, true)
}

// Groups tokens by line and prints them on the output chan
func StreamLines(scnr scanner.Scanner, bufSize int) (lines <-chan string) {
	lines, _ = StreamLinesOptionallySplitErrors(scnr, bufSize, false)
	return lines
}

// Groups tokens by line and prints them on the output chan, optionally prints
// errors to the error chan
func StreamLinesOptionallySplitErrors(
	scnr scanner.Scanner,
	bufsize int,
	splitErrors bool,
) (
	tokens <-chan string,
	errors <-chan string,
) {
	bufsize = intOr1024(bufsize)
	out := make(chan string, bufsize)

	var errs chan string
	if splitErrors {
		errs = make(chan string, bufsize)
	}

	go func() {
		line := struct {
			n     int
			print string
		}{1, ""}

		resetLine := func(next token.Token) {
			if len(line.print) > 0 {
				out <- strings.TrimLeft(line.print, " ")
			}
			line.n = next.Line
			line.print = next.Report()
		}

		for {
			t, err := scnr.NextToken()
			if err != nil {
				break
			}

			// Errors tokens may optionally be split into a separate stream
			if splitErrors && token.IsError(t.Id) {
				errs <- errorify(t)
			} else {
				if t.Line != line.n {
					resetLine(t)
				} else {
					line.print += " " + t.Report()
				}
			}
		}
		resetLine(token.Token{})
		close(out)
		if splitErrors {
			close(errs)
		}
	}()

	return out, errs
}

// Consumes tokens from the token chan, groups them by line and prints them on
// the output chan, prints errors to the error chan
func StreamTokensSplitErrors(tokens <-chan token.Token) (tokc, errc <-chan string) {
	bufsize := 1024
	tokcc, errcc := make(chan string, bufsize), make(chan string, bufsize)

	go func() {
		line := struct {
			n     int
			print string
		}{1, ""}

		resetLine := func(next token.Token) {
			if len(line.print) > 0 {
				tokcc <- strings.TrimLeft(line.print, " ")
			}
			line.n = next.Line
			line.print = next.Report()
		}

		for t := range tokens {
			if token.IsError(t.Id) {
				errcc <- errorify(t)
			} else {
				if t.Line != line.n {
					resetLine(t)
				} else {
					line.print += " " + t.Report()
				}
			}
		}

		resetLine(token.Token{})
		close(tokcc)
		close(errcc)
	}()

	return tokcc, errcc
}

func errorify(tok token.Token) string {
	if !token.IsError(tok.Id) {
		return ""
	}

	errTypes := map[token.Kind]string{
		token.INVALIDID:           "Invalid identifier",
		token.INVALIDCHAR:         "Invalid character",
		token.INVALIDNUM:          "Invalid number",
		token.UNTERMINATEDCOMMENT: "Unterminated comment",
	}

	return fmt.Sprintf(""+
		"Lexical error: %v: \"%v\": line %v.",
		errTypes[tok.Id], util.SingleLinify(string(tok.Lexeme)), tok.Line)
}

func intOr1024(i int) int {
	if i < 0 {
		return 1024
	}
	return i
}
