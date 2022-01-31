package reporting

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/scanner"
)

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
	bufSize int,
	splitErrors bool,
) (
	tokens <-chan string,
	errors <-chan string,
) {
	bufSize = intOr1024(bufSize)
	out := make(chan string, bufSize)

	var errs chan string
	if splitErrors {
		errs = make(chan string, bufSize)
	}

	go func() {
		line := struct {
			n     int
			print string
		}{1, ""}

		resetLine := func(next scanner.Token) {
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
			if splitErrors && scanner.IsError(t.Id) {
				errs <- errorify(t)
			} else {
				if t.Line != line.n {
					resetLine(t)
				} else {
					line.print += " " + t.Report()
				}
			}
		}
		resetLine(scanner.Token{})
		close(out)
		if splitErrors {
			close(errs)
		}
	}()

	return out, errs
}

func errorify(token scanner.Token) string {
	if !scanner.IsError(token.Id) {
		return ""
	}

	errTypes := map[scanner.Kind]string{
		scanner.INVALIDID:   "Invalid identifier",
		scanner.INVALIDCHAR: "Invalid character",
		scanner.INVALIDNUM:  "Invalid number",
	}

	return fmt.Sprintf(""+
		"Lexical error: %v: \"%v\": line %v.",
		errTypes[token.Id], token.Lexeme, token.Line)
}

func intOr1024(i int) int {
	if i < 0 {
		return 1024
	}
	return i
}
