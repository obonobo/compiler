package scanner

import "github.com/obonobo/esac/core/token"

// An alternative mode of operation for the scanner. Better support for Go
// looping constructs. Works similar to bufio.Scanner
type LoadableScanner interface {
	Scanner

	// Advances the scanner by one token
	Scan() bool

	// Reports the error registered by the scanner
	Err() error

	// Reports the token loaded in the scanner. Use after calling
	// LoadableScanner.Scan() first.
	Token() token.Token

	// Processes all tokens from the scanner
	Tokens() []token.Token
}

type loadableScanner struct {
	Scanner
	err   error
	token token.Token
}

// Wraps your scanner in a LoadableScanner
func NewLoadableScanner(scnr Scanner) LoadableScanner {
	return &loadableScanner{Scanner: scnr}
}

func (s *loadableScanner) Scan() bool {
	s.token, s.err = s.NextToken()
	return s.err == nil
}

func (s *loadableScanner) Err() error {
	return s.err
}

func (s *loadableScanner) Token() token.Token {
	return s.token
}

func (s *loadableScanner) Tokens() []token.Token {
	tokens := make([]token.Token, 0, 512)
	for s.Scan() {
		tokens = append(tokens, s.Token())
	}
	return tokens
}
