package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/obonobo/esac/internal/testutils"
	"github.com/obonobo/esac/util"
	"github.com/obonobo/esac/util/compile"
)

var (
	moon = "../moon"                   // Path to the moon binary
	libs = []string{"../stdlib/lib.m"} // Libraries to pass to moon
)

const (
	writeSomething = `
	func main() -> void {
		write(%v);
	}
	`
)

func TestArithmeticWithVariables(t *testing.T) {
	t.Parallel()
	testMoon(t, []MoonTestCase{
		{
			name: "basic",
			src: `
			func main() -> void {
				let x: integer;
				let y: integer;
				x = 10;
				y = x;
				y = y * y;
				write(x + y);
			}
			`,
			output: "110",
		},
	})
}

// Tests a bunch of arithmetic expressions
func TestArithmetic(t *testing.T) {
	t.Parallel()
	testArithmeticExpression(t, [][2]string{
		{"10 + 30 / 10", "13"},
		{"10 + 5 * 30 / 10", "25"},
		{"10 + 5 * 30 / 10 - 4", "21"},
	})
}

func TestDiv(t *testing.T) {
	t.Parallel()
	testTwoOp(t, "/", [][3]int{
		{30, 10, 3},
	})
}

func TestMult(t *testing.T) {
	t.Parallel()
	testTwoOp(t, "*", [][3]int{
		{10, 3, 30},
		{1, 2, 2},
		{10, 10, 100},
	})
}

func TestSub(t *testing.T) {
	t.Parallel()
	testTwoOp(t, "-", [][3]int{
		{10, 3, 7},
	})
}

func TestAdd(t *testing.T) {
	t.Parallel()
	testTwoOp(t, "+", [][3]int{
		{1, 5, 6},
		{10, 7, 17},
		{100, 9, 109},
	})
}

// Runs a series of tests cases that check output from a two operands operator
// like `+`, `-`, `*`, etc.
func testTwoOp(t *testing.T, op string, testCases [][3]int) {
	testArithmeticExpression(t, util.Map(testCases, func(tc [3]int) [2]string {
		return [2]string{
			fmt.Sprintf("%v%v%v", tc[0], op, tc[1]),
			fmt.Sprintf("%v", tc[2]),
		}
	}))
}

// Tests some arbitrary arithmetic expressions
func testArithmeticExpression(t *testing.T, testCases [][2]string) {
	testMoon(t, util.Map(testCases, func(tc [2]string) MoonTestCase {
		return MoonTestCase{
			name:   strings.ReplaceAll(fmt.Sprintf("%v=%v", tc[0], tc[1]), " ", ""),
			src:    fmt.Sprintf(writeSomething, tc[0]),
			output: tc[1],
		}
	}))
}

type MoonTestCase struct {
	name   string
	src    string
	output string
}

func testMoon(t *testing.T, testCases []MoonTestCase, linePrefix ...string) {
	prefix := "\t\t\t"
	if len(linePrefix) > 0 {
		prefix = linePrefix[0]
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertMoon(t, tc.src, tc.output, prefix)
		})
	}
}

// Compiles a program from source and runs it via the moon interpreter. Asserts
// program output. Optionally assert that an error has occured with the
// subprocess
func assertMoon(t *testing.T, src, expectedOutput, linePrefix string, expectedError ...error) {
	assertCycles := strings.Contains(expectedOutput, "cycles")
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

	expectedOutput = cleanMoonOutput(expectedOutput, linePrefix)
	actualOutput, err := runMoon(file)
	actualOutput = cleanMoonOutput(actualOutput, linePrefix)
	if !assertCycles {
		actualOutput = testutils.Head(actualOutput, -1)
	}

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

func cleanMoonOutput(out, prefix string) string {
	tail := 0
	if strings.Contains(out, "Loading") {
		tail = -len(libs) - 1
	}
	return strings.TrimRight(testutils.Tail(clean(out, prefix), tail), "\n")
}
