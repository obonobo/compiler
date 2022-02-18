package parser

type Parser interface {
	// Parses the token stream that is loaded in the Parser. Returns true if the
	// parse was successful, false otherwise
	Parse() bool

	// Retrieve the AST that was created by called Parser.Parse()
	AST() AST
}
