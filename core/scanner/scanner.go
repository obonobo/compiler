package scanner

// The Scanner interface is used to represent a parser that tokenizes a
// character stream
type Scanner interface {
	// Extract the next token in the program. This function is called by the
	// syntactic analyzer
	NextToken() Token
}

// The CharSource interface is kind of like the io.RuneScanner interface except
// that it doesn't report sizes and the `UnreadRune()` method (named
// `BackupChar` in this interface) also returns the unread rune
type CharSource interface {
	// Reads the next character in the input
	NextChar() (rune, error)

	// Back up one character in the input in case we have just read the next
	// character in order to resolve ambiguity
	BackupChar() (rune, error)
}
