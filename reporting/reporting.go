package reporting

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/util"
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
