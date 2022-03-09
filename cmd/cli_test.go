package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/obonobo/esac/internal/testutils"
)

const CLI_OUTPUT_LEX_POSITIVE = `
[eq, ==, 1] [plus, +, 1] [or, |, 1] [openpar, (, 1] [semi, ;, 1] [if, if, 1] [public, public, 1] [read, read, 1]
[noteq, <>, 2] [minus, -, 2] [and, &, 2] [closepar, ), 2] [comma, ,, 2] [then, then, 2] [private, private, 2] [write, write, 2]
[lt, <, 3] [mult, *, 3] [not, !, 3] [opencubr, {, 3] [dot, ., 3] [else, else, 3] [func, func, 3] [return, return, 3]
[gt, >, 4] [div, /, 4] [closecubr, }, 4] [colon, :, 4] [integer, integer, 4] [var, var, 4] [self, self, 4]
[leq, <=, 5] [assign, =, 5] [opensqbr, [, 5] [coloncolon, ::, 5] [float, float, 5] [struct, struct, 5] [inherits, inherits, 5]
[geq, >=, 6] [closesqbr, ], 6] [arrow, ->, 6] [void, void, 6] [while, while, 6] [let, let, 6]
[func, func, 7] [impl, impl, 7]
[intnum, 0, 13]
[intnum, 1, 14]
[intnum, 10, 15]
[intnum, 12, 16]
[intnum, 123, 17]
[intnum, 12345, 18]
[floatnum, 1.23, 20]
[floatnum, 12.34, 21]
[floatnum, 120.34e10, 22]
[floatnum, 12345.6789e-123, 23]
[id, abc, 25]
[id, abc1, 26]
[id, a1bc, 27]
[id, abc_1abc, 28]
[id, abc1_abc, 29]
[inlinecmt, // this is an inline comment\n, 31]
[blockcmt, /* this is a single line block comment */, 33]
[blockcmt, /* this is a\nmultiple line\nblock comment\n*/, 35]
[blockcmt, /* this is an imbricated\n/* block comment\n*/\n*/, 40]
`

const CLI_OUTPUT_LEX_NEGATIVE = `
[invalidchar, @, 1] [invalidchar, #, 1] [invalidchar, $, 1] [invalidchar, ', 1] [invalidchar, \, 1] [invalidchar, ~, 1]
[invalidnum, 00, 3]
[invalidnum, 01, 4]
[invalidnum, 010, 5]
[invalidnum, 0120, 6]
[invalidnum, 01230, 7]
[invalidnum, 0123450, 8]
[invalidnum, 01.23, 10]
[invalidnum, 012.34, 11]
[invalidnum, 12.340, 12]
[invalidnum, 012.340, 13]
[invalidnum, 012.34e10, 15]
[invalidnum, 12.34e010, 16]
[invalidid, _abc, 18]
[invalidid, 1abc, 19]
[invalidid, _1abc, 20]
`

const CLI_FILE_OUTPUT_LEX_POSITIVE_ERRORS = ``

const CLI_FILE_OUTPUT_LEX_NEGATIVE_TOKENS = ``
const CLI_FILE_OUTPUT_LEX_NEGATIVE_ERRORS = `
Lexical error: Invalid character: "@": line 1.
Lexical error: Invalid character: "#": line 1.
Lexical error: Invalid character: "$": line 1.
Lexical error: Invalid character: "'": line 1.
Lexical error: Invalid character: "\": line 1.
Lexical error: Invalid character: "~": line 1.
Lexical error: Invalid number: "00": line 3.
Lexical error: Invalid number: "01": line 4.
Lexical error: Invalid number: "010": line 5.
Lexical error: Invalid number: "0120": line 6.
Lexical error: Invalid number: "01230": line 7.
Lexical error: Invalid number: "0123450": line 8.
Lexical error: Invalid number: "01.23": line 10.
Lexical error: Invalid number: "012.34": line 11.
Lexical error: Invalid number: "12.340": line 12.
Lexical error: Invalid number: "012.340": line 13.
Lexical error: Invalid number: "012.34e10": line 15.
Lexical error: Invalid number: "12.34e010": line 16.
Lexical error: Invalid identifier: "_abc": line 18.
Lexical error: Invalid identifier: "1abc": line 19.
Lexical error: Invalid identifier: "_1abc": line 20.
`

const LEX_HELLOWORLD_SRC_ERRORS = ``
const LEX_HELLOWORLD_SRC_TOKENS = `
[blockcmt, /*\nThis is an imaginary program with a made up syntax.\n\nLet us see how the parser handles it...\n*/, 1]
[inlinecmt, // C-style struct\n, 7]
[struct, struct, 8] [id, Student, 8] [opencubr, {, 8]
[float, float, 9] [id, age, 9] [semi, ;, 9]
[integer, integer, 10] [id, id, 10] [semi, ;, 10]
[closecubr, }, 11] [semi, ;, 11]
[public, public, 13] [func, func, 13] [id, main, 13] [openpar, (, 13] [closepar, ), 13] [opencubr, {, 13]
[inlinecmt, // x is my integer variable\n, 15]
[let, let, 16] [id, x, 16] [assign, =, 16] [intnum, 10, 16] [semi, ;, 16]
[blockcmt, /*\n    y is equal to x\n    */, 18]
[var, var, 21] [id, y, 21] [assign, =, 21] [id, x, 21] [semi, ;, 21]
[inlinecmt, // Equality check\n, 23]
[if, if, 24] [openpar, (, 24] [id, y, 24] [eq, ==, 24] [id, x, 24] [closepar, ), 24] [then, then, 24] [opencubr, {, 24]
[var, var, 25] [id, out, 25] [integer, integer, 25] [opensqbr, [, 25] [intnum, 10, 25] [closesqbr, ], 25] [assign, =, 25] [opencubr, {, 25] [id, x, 25] [comma, ,, 25] [id, y, 25] [comma, ,, 25] [intnum, 69, 25] [comma, ,, 25] [intnum, 200, 25] [comma, ,, 25] [intnum, 89, 25] [closecubr, }, 25] [semi, ;, 25]
[write, write, 26] [openpar, (, 26] [id, out, 26] [closepar, ), 26] [semi, ;, 26] [inlinecmt, // Assume we have a function called 'write'\n, 26]
[closecubr, }, 27]
[closecubr, }, 28]
`

const LEX_STRINGS_SRC_TOKENS = `
[var, var, 1] [id, x, 1] [assign, =, 1] [id, this, 1] [id, is, 1] [id, not, 1] [id, valid, 1] [semi, ;, 1]
`
const LEX_STRINGS_SRC_ERRORS = `
Lexical error: Invalid character: """: line 1.
Lexical error: Invalid character: """: line 1.
`

const LEX_STRINGS_SRC_TOKENS_AND_ERRORS = `
[var, var, 1] [id, x, 1] [assign, =, 1] [invalidchar, ", 1] [id, this, 1] [id, is, 1] [id, not, 1] [id, valid, 1] [invalidchar, ", 1] [semi, ;, 1]
`

const LEX_SOMETHINGELSE_TOKENS = `
[id, package, 1] [id, main, 1]
[id, import, 3] [openpar, (, 3]
[id, encoding, 4] [div, /, 4] [id, json, 4]
[id, fmt, 5]
[id, io, 6]
[id, log, 7]
[id, net, 8] [div, /, 8] [id, http, 8]
[closepar, ), 9]
[func, func, 11] [id, main, 11] [openpar, (, 11] [closepar, ), 11] [opencubr, {, 11]
[id, resp, 12] [comma, ,, 12] [id, err, 12] [colon, :, 12] [assign, =, 12] [id, http, 12] [dot, ., 12] [id, Get, 12] [openpar, (, 12] [id, https, 12] [colon, :, 12] [inlinecmt, //www.google.com")\n, 12]
[if, if, 13] [id, err, 13] [not, !, 13] [assign, =, 13] [id, nil, 13] [opencubr, {, 13]
[id, log, 14] [dot, ., 14] [id, Fatalf, 14] [openpar, (, 14] [id, Request, 14] [id, failed, 14] [colon, :, 14] [id, v, 14] [comma, ,, 14] [id, err, 14] [closepar, ), 14]
[closecubr, }, 15]
[id, headers, 17] [comma, ,, 17] [id, err, 17] [colon, :, 17] [assign, =, 17] [id, json, 17] [dot, ., 17] [id, MarshalIndent, 17] [openpar, (, 17] [id, resp, 17] [dot, ., 17] [id, Header, 17] [comma, ,, 17] [comma, ,, 17] [closepar, ), 17]
[if, if, 18] [id, err, 18] [not, !, 18] [assign, =, 18] [id, nil, 18] [opencubr, {, 18]
[id, log, 19] [dot, ., 19] [id, Fatalf, 19] [openpar, (, 19] [id, Failed, 19] [id, to, 19] [id, serialize, 19] [id, response, 19] [id, headers, 19] [colon, :, 19] [id, v, 19] [comma, ,, 19] [id, err, 19] [closepar, ), 19]
[closecubr, }, 20]
[id, fmt, 21] [dot, ., 21] [id, Println, 21] [openpar, (, 21] [id, string, 21] [openpar, (, 21] [id, headers, 21] [closepar, ), 21] [closepar, ), 21]
[id, bod, 23] [comma, ,, 23] [id, err, 23] [colon, :, 23] [assign, =, 23] [id, io, 23] [dot, ., 23] [id, ReadAll, 23] [openpar, (, 23] [id, resp, 23] [dot, ., 23] [id, Body, 23] [closepar, ), 23]
[if, if, 24] [id, err, 24] [not, !, 24] [assign, =, 24] [id, nil, 24] [opencubr, {, 24]
[id, log, 25] [dot, ., 25] [id, Fatalf, 25] [openpar, (, 25] [id, Failed, 25] [id, to, 25] [read, read, 25] [id, body, 25] [closepar, ), 25]
[closecubr, }, 26]
[id, defer, 27] [id, resp, 27] [dot, ., 27] [id, Body, 27] [dot, ., 27] [id, Close, 27] [openpar, (, 27] [closepar, ), 27]
[id, fmt, 29] [dot, ., 29] [id, Println, 29] [openpar, (, 29] [id, string, 29] [openpar, (, 29] [id, bod, 29] [closepar, ), 29] [closepar, ), 29]
[closecubr, }, 30]
`

const LEX_SOMETHINGELSE_ERRORS = `
Lexical error: Invalid character: """: line 4.
Lexical error: Invalid character: """: line 4.
Lexical error: Invalid character: """: line 5.
Lexical error: Invalid character: """: line 5.
Lexical error: Invalid character: """: line 6.
Lexical error: Invalid character: """: line 6.
Lexical error: Invalid character: """: line 7.
Lexical error: Invalid character: """: line 7.
Lexical error: Invalid character: """: line 8.
Lexical error: Invalid character: """: line 8.
Lexical error: Invalid character: """: line 12.
Lexical error: Invalid character: """: line 14.
Lexical error: Invalid character: "%": line 14.
Lexical error: Invalid character: """: line 14.
Lexical error: Invalid character: """: line 17.
Lexical error: Invalid character: """: line 17.
Lexical error: Invalid character: """: line 17.
Lexical error: Invalid character: """: line 17.
Lexical error: Invalid character: """: line 19.
Lexical error: Invalid character: "%": line 19.
Lexical error: Invalid character: """: line 19.
Lexical error: Invalid character: """: line 25.
Lexical error: Invalid character: """: line 25.
`

const LEX_SOMETHINGELSE_TOKENS_AND_ERRORS = `
[id, package, 1] [id, main, 1]
[id, import, 3] [openpar, (, 3]
[invalidchar, ", 4] [id, encoding, 4] [div, /, 4] [id, json, 4] [invalidchar, ", 4]
[invalidchar, ", 5] [id, fmt, 5] [invalidchar, ", 5]
[invalidchar, ", 6] [id, io, 6] [invalidchar, ", 6]
[invalidchar, ", 7] [id, log, 7] [invalidchar, ", 7]
[invalidchar, ", 8] [id, net, 8] [div, /, 8] [id, http, 8] [invalidchar, ", 8]
[closepar, ), 9]
[func, func, 11] [id, main, 11] [openpar, (, 11] [closepar, ), 11] [opencubr, {, 11]
[id, resp, 12] [comma, ,, 12] [id, err, 12] [colon, :, 12] [assign, =, 12] [id, http, 12] [dot, ., 12] [id, Get, 12] [openpar, (, 12] [invalidchar, ", 12] [id, https, 12] [colon, :, 12] [inlinecmt, //www.google.com")\n, 12]
[if, if, 13] [id, err, 13] [not, !, 13] [assign, =, 13] [id, nil, 13] [opencubr, {, 13]
[id, log, 14] [dot, ., 14] [id, Fatalf, 14] [openpar, (, 14] [invalidchar, ", 14] [id, Request, 14] [id, failed, 14] [colon, :, 14] [invalidchar, %, 14] [id, v, 14] [invalidchar, ", 14] [comma, ,, 14] [id, err, 14] [closepar, ), 14]
[closecubr, }, 15]
[id, headers, 17] [comma, ,, 17] [id, err, 17] [colon, :, 17] [assign, =, 17] [id, json, 17] [dot, ., 17] [id, MarshalIndent, 17] [openpar, (, 17] [id, resp, 17] [dot, ., 17] [id, Header, 17] [comma, ,, 17] [invalidchar, ", 17] [invalidchar, ", 17] [comma, ,, 17] [invalidchar, ", 17] [invalidchar, ", 17] [closepar, ), 17]
[if, if, 18] [id, err, 18] [not, !, 18] [assign, =, 18] [id, nil, 18] [opencubr, {, 18]
[id, log, 19] [dot, ., 19] [id, Fatalf, 19] [openpar, (, 19] [invalidchar, ", 19] [id, Failed, 19] [id, to, 19] [id, serialize, 19] [id, response, 19] [id, headers, 19] [colon, :, 19] [invalidchar, %, 19] [id, v, 19] [invalidchar, ", 19] [comma, ,, 19] [id, err, 19] [closepar, ), 19]
[closecubr, }, 20]
[id, fmt, 21] [dot, ., 21] [id, Println, 21] [openpar, (, 21] [id, string, 21] [openpar, (, 21] [id, headers, 21] [closepar, ), 21] [closepar, ), 21]
[id, bod, 23] [comma, ,, 23] [id, err, 23] [colon, :, 23] [assign, =, 23] [id, io, 23] [dot, ., 23] [id, ReadAll, 23] [openpar, (, 23] [id, resp, 23] [dot, ., 23] [id, Body, 23] [closepar, ), 23]
[if, if, 24] [id, err, 24] [not, !, 24] [assign, =, 24] [id, nil, 24] [opencubr, {, 24]
[id, log, 25] [dot, ., 25] [id, Fatalf, 25] [openpar, (, 25] [invalidchar, ", 25] [id, Failed, 25] [id, to, 25] [read, read, 25] [id, body, 25] [invalidchar, ", 25] [closepar, ), 25]
[closecubr, }, 26]
[id, defer, 27] [id, resp, 27] [dot, ., 27] [id, Body, 27] [dot, ., 27] [id, Close, 27] [openpar, (, 27] [closepar, ), 27]
[id, fmt, 29] [dot, ., 29] [id, Println, 29] [openpar, (, 29] [id, string, 29] [openpar, (, 29] [id, bod, 29] [closepar, ), 29] [closepar, ), 29]
[closecubr, }, 30]
`

const LEX_UNTERMINATEDCOMMENTS_TOKENS_AND_ERRORS = `
[inlinecmt, // this is an inline comment\n, 1]
[unterminatedcomment, /* this is a single line block comment, 3]
`

const LEX_UNTERMINATEDCOMMENTS_TOKENS = `
[inlinecmt, // this is an inline comment\n, 1]
`

const LEX_UNTERMINATEDCOMMENTS_ERRORS = `
Lexical error: Unterminated comment: "/* this is a single line block comment": line 3.
`

const LEX_UNTERMINATEDCOMMENTS2_TOKENS_AND_ERRORS = `
[unterminatedcomment, /* this is an imbricated\n/* block comment, 1]
`

const LEX_UNTERMINATEDCOMMENTS2_TOKENS = ``

const LEX_UNTERMINATEDCOMMENTS2_ERRORS = `
Lexical error: Unterminated comment: "/* this is an imbricated\n/* block comment": line 1.
`

func TestLexUnterminatedComments2Stdout(t *testing.T) {
	assertCliStdout(t,
		"TestLexUnterminatedComments2Stdout",
		testutils.UNTERMINATEDCOMMENTS2_SRC,
		LEX_UNTERMINATEDCOMMENTS2_TOKENS_AND_ERRORS)
}

func TestLexUnterminatedComments2NormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexUnterminatedComments2NormalOutput",
		testutils.UNTERMINATEDCOMMENTS2_SRC,
		LEX_UNTERMINATEDCOMMENTS2_TOKENS,
		LEX_UNTERMINATEDCOMMENTS2_ERRORS)
}

func TestLexUnterminatedCommentsStdout(t *testing.T) {
	assertCliStdout(t,
		"TestLexUnterminatedCommentsStdout",
		testutils.UNTERMINATEDCOMMENTS_SRC,
		LEX_UNTERMINATEDCOMMENTS_TOKENS_AND_ERRORS)
}

func TestLexUnterminatedCommentsNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexUnterminatedCommentsNormalOutput",
		testutils.UNTERMINATEDCOMMENTS_SRC,
		LEX_UNTERMINATEDCOMMENTS_TOKENS,
		LEX_UNTERMINATEDCOMMENTS_ERRORS)
}

func TestLexSomethingElseNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexSomethingElseCliStdout",
		testutils.LEX_SOMETHING_ELSE_SRC,
		LEX_SOMETHINGELSE_TOKENS,
		LEX_SOMETHINGELSE_ERRORS)
}

func TestLexSomethingElseCliStdout(t *testing.T) {
	assertCliStdout(t,
		"TestLexSomethingElseCliStdout",
		testutils.LEX_SOMETHING_ELSE_SRC,
		LEX_SOMETHINGELSE_TOKENS_AND_ERRORS)
}

func TestLexStringsCliStdout(t *testing.T) {
	assertCliStdout(t,
		"TestLexStringsNormalOutput",
		testutils.LEX_STRINGS_SRC,
		LEX_STRINGS_SRC_TOKENS_AND_ERRORS)
}

func TestLexStringsNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexStringsNormalOutput",
		testutils.LEX_STRINGS_SRC,
		LEX_STRINGS_SRC_TOKENS,
		LEX_STRINGS_SRC_ERRORS)
}

func TestLexHelloWorldNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexHelloWorldCli",
		testutils.LEX_HELLOWORLD_SRC,
		LEX_HELLOWORLD_SRC_TOKENS,
		LEX_HELLOWORLD_SRC_ERRORS)
}

func TestLexHelloWorldCli(t *testing.T) {
	assertCliStdout(t,
		"TestLexHelloWorldCli",
		testutils.LEX_HELLOWORLD_SRC,
		LEX_HELLOWORLD_SRC_TOKENS)
}

func TestLexNegativeGradingNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexPositiveGradingNormalOutput",
		testutils.LEX_NEGATIVE_GRADING_SRC,
		CLI_FILE_OUTPUT_LEX_NEGATIVE_TOKENS,
		CLI_FILE_OUTPUT_LEX_NEGATIVE_ERRORS)
}

func TestLexPositiveGradingNormalOutput(t *testing.T) {
	assertCliNormal(t,
		"TestLexPositiveGradingNormalOutput",
		testutils.LEX_POSITIVE_GRADING_SRC,
		CLI_OUTPUT_LEX_POSITIVE,
		CLI_FILE_OUTPUT_LEX_POSITIVE_ERRORS)
}

// Tests CLI parsing the lexnegative source
func TestLexNegativeGrading(t *testing.T) {
	assertCliStdout(t,
		"TestLexNegativeGrading",
		testutils.LEX_NEGATIVE_GRADING_SRC,
		CLI_OUTPUT_LEX_NEGATIVE)
}

// Tests CLI parsing the lexpositive source
func TestLexPositiveGrading(t *testing.T) {
	assertCliStdout(t,
		"TestLexPositiveGrading",
		testutils.LEX_POSITIVE_GRADING_SRC,
		CLI_OUTPUT_LEX_POSITIVE)
}

func assertCliNormal(
	t *testing.T,
	testName string,
	inputFiledata string,
	expectedTokensOutput string,
	expectedErrorsOutput string,
) {
	tmp, rm := createTempFile(t, "tmp-"+testName+"*.src", strings.Trim(inputFiledata, "\r\t\n"))
	outTokensName := inputFileNameToOutputFileName(tmp.Name(), OUT_LEX_TOKENS)
	outErrorsName := inputFileNameToOutputFileName(tmp.Name(), OUT_LEX_ERRORS)
	defer func() {
		rm()
		if _, err := os.Stat(outTokensName); !os.IsNotExist(err) {
			os.Remove(outTokensName)
		}
		if _, err := os.Stat(outErrorsName); !os.IsNotExist(err) {
			os.Remove(outErrorsName)
		}
	}()

	args := []string{"esacc", "lex", tmp.Name()}
	if exit := Run(args); exit != 0 {
		t.Fatalf("CLI should have returned '0' exit code but got code '%v'", exit)
	}

	outTokens := readFile(t, outTokensName)
	outErrors := readFile(t, outErrorsName)

	actual := strings.Trim(outTokens, "\r\n\t")
	expected := strings.Trim(expectedTokensOutput, "\r\n\t")
	if actual != expected {
		t.Fatalf(""+
			"Expected file '%v' to contain '%v' but got '%v'",
			outTokensName, expected, actual)
	}

	actual = strings.Trim(outErrors, "\r\n\t")
	expected = strings.Trim(expectedErrorsOutput, "\r\n\t")
	if actual != expected {
		t.Fatalf(""+
			"Expected file '%v' to contain '%v' but got '%v'",
			outErrorsName, expected, actual)

	}
}

func readFile(t *testing.T, name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("Failed to read file '%v': %v", name, err)
	}
	return string(data)
}

func assertCliStdout(t *testing.T, testName, inputFileData, expectedOutput string) {
	output := mockStdoutStderr(t)
	tmp, rm := createTempFile(t, "tmp-"+testName, strings.Trim(inputFileData, "\r\t\n"))
	defer rm()

	args := []string{"esacc", "lex", "-o", "-", tmp.Name()}
	exit := Run(args)
	data := output()
	if exit != 0 {
		t.Fatalf("Expected command to succeed, but got exit code '%v'", exit)
	}

	// Chop the first two lines off
	tailed := strings.Trim(testutils.Tail(data, -2), "\r\t\n\x00")
	expected := strings.Trim(expectedOutput, "\t\r\n\x00")
	if tailed != expected {
		t.Fatalf("Expected output '%v' but got '%v'", expected, tailed)
	}
}

func createTempFile(t *testing.T, name, contents string) (file *os.File, removeFile func()) {
	tmp, err := os.CreateTemp("./", name)
	if err != nil {
		t.Fatalf("Temp file is needed for this test, got an error creating: %v", err)
	}

	destroy := func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	_, err = tmp.WriteString(contents)
	if err != nil {
		destroy()
		t.Fatalf("Writing to temporary file failed: %v", err)
	}

	return tmp, destroy
}

func mockStdoutStderr(t *testing.T) (output func() string) {
	output, err := testutils.MockStdoutStderr()
	if err != nil {
		t.Fatal(err)
	}
	return output
}
