package testutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Consumes the stdout and stderr of the current process, dumping them as a
// single string which is produced when you call the returned `close()`
// function.
func MockStdoutStderr() (close func() string, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to mock stdout/stderr: %w", err)
	}

	out := make(chan string, 1)
	ready, done := make(chan struct{}, 1), make(chan struct{}, 1)
	go func() {
		old := os.Stdout
		olde := os.Stderr
		defer func() {
			r.Close()
			os.Stdout = old
			os.Stderr = olde
			done <- struct{}{}
		}()

		os.Stdout = w
		os.Stderr = w
		ready <- struct{}{}

		var buf bytes.Buffer
		io.Copy(&buf, r)
		io.Copy(old, bytes.NewBuffer(buf.Bytes()))
		out <- buf.String()
	}()

	<-ready
	return func() string {
		w.Close()
		<-done
		return <-out
	}, nil
}

// Tails some string output. Works like the `tail` utility; provide positive n
// to take the bottom n lines, provide negative n to trim n lines off the top
func Tail(data string, n int) string {
	lines := strings.Split(data, "\n")
	var from int
	switch {
	case n < 0:
		from = -n
	case n > 0:
		from = len(lines) - n
	}
	return strings.Join(lines[from:], "\n")
}

// Takes the head of some string output. Works like the `head` utility; provide
// positive n to take the top n lines, provide negative n to trim n lines off
// the bottom
func Head(data string, n int) string {
	lines := strings.Split(data, "\n")
	to := len(lines)
	switch {
	case n < 0:
		to += n
	case n > 0:
		to = n
	}
	return strings.Join(lines[:to], "\n")
}

// Trims a prefix from all lines in a string
func TrimLeading(in string, prefix string) string {
	out := bytes.NewBuffer(make([]byte, 0, len(in)))
	for scnr := bufio.NewScanner(bytes.NewBufferString(in)); scnr.Scan(); {
		fmt.Fprintln(out, strings.TrimPrefix(scnr.Text(), prefix))
	}
	return out.String()
}

// Creates a temporary file with the specified contents. Returns the name of the
// file and a function that deletes the file
func TempFile(name, contents string) (string, func()) {
	if name != "" {
		name += "_"
	}
	fh, err := os.CreateTemp(".", fmt.Sprintf("%v*.tmp", name))
	if err != nil {
		return "", func() {}
	}
	defer fh.Close()
	fh.WriteString(contents)
	name = fh.Name()
	return name, func() {
		if err := os.Remove(name); err != nil {
			panic(err)
		}
	}
}
