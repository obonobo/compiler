package chuggingcharsource

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

const (
	abcdefg = "abcdefg"
)

func TestChugFile(t *testing.T) {
	t.Parallel()
	tmp, destroy := tempFileWithContents(t, "TestChugFile", abcdefg)
	defer destroy()

	chugger := new(ChuggingCharSource)
	if err := chugger.Chug(tmp.Name()); err != nil {
		t.Fatalf("Chug should have succeeded: %v", err)
	}

	data, _ := io.ReadAll(chugger)
	if contents := string(data); contents != abcdefg {
		t.Fatalf("Expected chugger to contain '%v' but got '%v'", abcdefg, contents)
	}
}

func TestNextChar(t *testing.T) {
	t.Parallel()
	chugger := chuggerWithContents(t, abcdefg)
	for _, expected := range abcdefg {
		actual, err := chugger.NextChar()
		if err != nil || actual != expected {
			t.Fatalf(
				"Expected chugger.NextChar() to return '%s' but got '%s'",
				string(expected), string(actual))
		}
	}
}

func TestBackupChar(t *testing.T) {
	t.Parallel()
	var rev string
	chugger := chuggerWithContents(t, abcdefg)

	// Read all the characters
	for range abcdefg {
		r, err := chugger.NextChar()
		if err != nil {
			t.Fatalf("Failed to grab next char: %v", err)
		}
		rev = string(r) + rev
	}

	// Read backwards and assert
	for _, expected := range rev {
		if actual, err := chugger.BackupChar(); err != nil {
			t.Fatalf("Should have been able to back up here: %v", err)
		} else if actual != expected {
			t.Fatalf(
				"Expected to read '%v' but got '%v'",
				string(expected), string(actual))
		}
	}
}

func TestNextCharEndOfFile(t *testing.T) {
	t.Parallel()
	assertThrowsEndOfCharSourceError(
		t, func(chugger *ChuggingCharSource) (interface{}, error) {
			return chugger.NextChar()
		})
}

func TestBackupCharEndOfFile(t *testing.T) {
	t.Parallel()
	assertThrowsEndOfCharSourceError(
		t, func(chugger *ChuggingCharSource) (interface{}, error) {
			return chugger.BackupChar()
		})
}

func assertThrowsEndOfCharSourceError(
	t *testing.T,
	do func(chugger *ChuggingCharSource) (interface{}, error),
) {
	_, err := do(chuggerWithContents(t, ""))
	if err == nil {
		t.Fatalf("Should have got an error trying to read an empty chugger")
	}
	if e := new(EndOfCharSourceError); !errors.As(err, &e) {
		t.Fatalf("Expected error to be of type EndOfCharSourceError but got '%v'", err)
	}
}

func chuggerWithContents(t *testing.T, contents string) *ChuggingCharSource {
	chugger := new(ChuggingCharSource)
	if err := chugger.ChugReader(bytes.NewBufferString(contents)); err != nil {
		t.Fatalf("Chugger should chug data okay: %v", err)
	}
	return chugger
}

// Creates a temporary file with the provided contents, sets the offset for the
// next read/write to the beginning of the file
func tempFileWithContents(
	t *testing.T,
	filenamePrefix string,
	contents string,
) (fh *os.File, destroy func()) {
	tmp, destroy := tempFile(t, "TestChugFile")
	if _, err := tmp.WriteString(contents); err != nil {
		t.Fatalf("Failed to write data to file: %v", err)
	}
	if _, err := tmp.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek to beginning of file: %v", err)
	}
	return tmp, destroy
}

// Creates a temporary file, provides a destory function to eradicate the file
func tempFile(t *testing.T, filenamePrefix string) (fh *os.File, destroy func()) {
	tmp, err := os.CreateTemp(".", fmt.Sprintf("%v_*", filenamePrefix))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	return tmp, func() {
		tmp.Close()
		for try := 0; try < 3; try++ {
			if err := os.Remove(tmp.Name()); err == nil {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}
