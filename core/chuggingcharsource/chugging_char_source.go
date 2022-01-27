package chuggingcharsource

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

var ErrEOF = fmt.Errorf("EOF")

// A character source that is initialized by first chugging an entire file or
// io.Reader into its internal buffer. Assumes the file/io.Reader is UTF-8
// encoded.
type ChuggingCharSource struct {
	buf     []byte
	i       int
	line    int
	columns []int // Need to keep a stack of columns here in case of backtrack
}

// Initializes the ChuggingCharSource by chugging a file
func (c *ChuggingCharSource) Chug(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return &ChuggingError{err}
	}
	defer f.Close()
	return c.ChugReader(f)
}

// Initializes the ChuggingCharSource by chugging from the provided io.Reader.
// If this chugger already contains a buffer, then the old buffer gets replaced
// completely by a new buffer with the contents of the reader
func (c *ChuggingCharSource) ChugReader(reader io.Reader) error {
	b, err := io.ReadAll(reader) // Let io.ReadAll create the buffer
	c.buf = b
	c.i = 0
	if err != nil {
		return &ChuggingError{err}
	}
	return nil
}

// Reads the next character in the input
func (c *ChuggingCharSource) NextChar() (rune, error) {
	r, _, err := c.ReadRune()
	return r, err
}

// Back up one character in the input in case we have just read the next
// character in order to resolve ambiguity
func (c *ChuggingCharSource) BackupChar() (rune, error) {
	r, s, err := c.PeekBack()
	if err != nil {
		return 0, err
	}

	c.i -= s
	if r == '\n' {
		c.uncountNewline()
	} else {
		c.uncountColumn()
	}

	return r, err
}

// Reads the remainder of the buffer, starting from the current character
// position. This method is just so you can use the chugger as a raw buffer to
// read from, similar to bytes.Buffer
//
// WARNING: column and line numbers are not counted when reading this way. You
// must call ChuggingCharSource.Reset() to reset the counts before trying to
// obtain column and line numbers from the chugger
func (c *ChuggingCharSource) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p) && c.i < len(c.buf); i++ {
		p[i] = c.buf[c.i]
		c.i++
		n++
	}
	return n, io.EOF
}

func (c *ChuggingCharSource) ReadRune() (r rune, size int, err error) {
	r, s, err := c.Peek()
	if err != nil {
		return 0, 0, err
	}

	c.i += s
	if r == '\n' {
		c.countNewline()
	} else {
		c.countColumn()
	}

	return r, s, err
}

func (c *ChuggingCharSource) UnreadRune() error {
	_, err := c.BackupChar()
	return err
}

// Returns the closest rune infront of the cursor without advancing the chugger
// forward
func (c *ChuggingCharSource) Peek() (r rune, size int, err error) {
	if c.i > len(c.buf) {
		return 0, 0, &EndOfCharSourceError{io.EOF}
	}
	r, s := utf8.DecodeRune(c.buf[c.i:])
	if r == utf8.RuneError {
		if s == 0 {
			return 0, 0, &EndOfCharSourceError{io.EOF}
		}
		return 0, 0, fmt.Errorf("ChuggingCharSource: RuneError from utf8 lib")
	}
	return r, s, nil
}

// Returns the closest rune behind the cursor without backing up the chugger
func (c *ChuggingCharSource) PeekBack() (r rune, s int, err error) {
	if c.i == 0 {
		return 0, 0, &EndOfCharSourceError{io.EOF}
	}
	r, s = utf8.DecodeLastRune(c.buf[:c.i])
	if r == utf8.RuneError {
		if s == 0 {
			return 0, 0, &EndOfCharSourceError{io.EOF}
		}
		return 0, 0, fmt.Errorf("ChuggingCharSource: RuneError from utf8 lib")
	}
	return r, s, nil
}

// Reports the current line number
func (c *ChuggingCharSource) Line() int {
	return c.line + 1
}

// Reports the current column number
func (c *ChuggingCharSource) Column() int {
	if len(c.columns) < 1 {
		return 1
	}
	return c.columns[len(c.columns)-1] + 1
}

func (c *ChuggingCharSource) Reset() {
	c.columns = make([]int, 0, cap(c.columns))
	c.line = 0
	c.i = 0
}

func (c *ChuggingCharSource) countColumn() {
	if len(c.columns) < 1 {
		c.columns = append(c.columns, 0)
	}
	c.columns[len(c.columns)-1]++
}

func (c *ChuggingCharSource) uncountColumn() {
	c.columns[len(c.columns)-1]--
}

func (c *ChuggingCharSource) countNewline() {
	c.line++
	c.columns = append(c.columns, 0)
}

func (c *ChuggingCharSource) uncountNewline() {
	c.line--
	c.columns = c.columns[:len(c.columns)-1]
}
