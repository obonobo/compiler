package scanner

// The CharSource interface is kind of like the io.RuneScanner interface except
// that it doesn't report sizes and the `UnreadRune()` method (named
// `BackupChar` in this interface) also returns the unread rune
//
// It also is able to report the current line and column number that it is on.
//
// CharSource should assume UTF-8 encoding
//
type CharSource interface {
	// Reads the next character in the input
	NextChar() (rune, error)

	// Back up one character in the input in case we have just read the next
	// character in order to resolve ambiguity
	BackupChar() (rune, error)

	// Reports the current line number
	Line() int

	// Reports the current column number
	Column() int
}
