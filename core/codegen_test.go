package core

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/obonobo/esac/internal/testutils"
	"github.com/obonobo/esac/util/compile"
)

var (
	moon = "../moon"                   // Path to the moon binary
	libs = []string{"../stdlib/lib.m"} // Libraries to pass to moon
)

func TestAdd(t *testing.T) {
	assertMoon(t, `
	func main() -> void {
		let x: integer;
		let y: integer;
		write(1 + 5);
	}
	`, `
	Loading scratch.moon.
	Loading stdlib/lib.m.
	6
	574 cycles.
	`)
}

// Compiles a program from source and runs it via the moon interpreter. Asserts
// program output. Optionally assert that an error has occured with the
// subprocess
func assertMoon(t *testing.T, src, expectedOutput string, expectedError ...error) {
	var expectedErr error
	if len(expectedError) > 0 {
		expectedErr = expectedError[0]
	}

	compiled, err := compile.TagsBased(bytes.NewBufferString(src))
	if err != nil {
		if err != expectedErr {
			t.Fatalf("Expected compile to return err: %v\nBut got: %v", expectedErr, err)
		}
		expectedErr = nil // If we match this error, then remove it for later
	}

	file, delete := testutils.TempFile("moon_progam", compiled)
	defer delete()

	expectedOutput = cleanMoonOutput(expectedOutput)
	actualOutput, err := runMoon(file)
	actualOutput = cleanMoonOutput(actualOutput)

	// Assert output
	if actual, expected := actualOutput, expectedOutput; actual != expected {
		t.Errorf("Expected moon output:\n%v\nBut got:\n%v", expected, actual)
	}

	// Assert error
	if actual, expected := err, expectedErr; actual != expected {
		t.Errorf("Expected moon error: %v\nBut got: %v", expected, actual)
	}
}

// Runs the moon interpreter on the specified source file
func runMoon(file string) (string, error) {
	out, err := exec.Command(moon, append([]string{file}, libs...)...).Output()
	return string(out), err
}

func cleanMoonOutput(out string) string {
	return strings.TrimRight(testutils.Tail(clean(out, "\t"), -len(libs)-1), "\n")
}
