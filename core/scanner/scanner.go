package scanner

// The Scanner interface is used to represent a parser that tokenizes a
// character stream
type Scanner interface {
	// Extract the next token in the program. This function is called by the
	// syntactic analyzer
	NextToken() (Token, error)
}
