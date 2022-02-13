package scanner

import "github.com/obonobo/esac/core/token"

// The Scanner interface is used to represent a lexer that tokenizes a
// character stream
type Scanner interface {
	// Extract the next token in the program. This function is called by the
	// syntactic analyzer
	NextToken() (token.Token, error)
}
