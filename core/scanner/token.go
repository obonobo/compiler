package scanner

import "fmt"

// Unique token identifier
type Kind string

// The text representing a token
type Lexeme string

type Token struct {
	Id     Kind   // The unique identifier of this token
	Lexeme Lexeme // The exact string that was matched as this token
	Line   int    // The line number on which the token was found
	Column int    // The column number on which the token was found
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token[Id=%v, Lexeme=%v, Line=%v, Column=%v]",
		t.Id, t.Lexeme, t.Line, t.Column)
}

const (
	ASSIGN Kind = "assign" // Assignment operator `=`
	ARROW  Kind = "arrow"  // Right-pointing arrow operator `->`

	EQ    Kind = "eq"    // Arithmetic operator: equality `==`
	PLUS  Kind = "plus"  // Arithmetic operator: addition `+`
	MINUS Kind = "minus" // Arithmetic operator: subtraction `-`
	MULT  Kind = "mult"  // Arithmetic operator: multiplication `*`
	DIV   Kind = "div"   // Arithmetic operator: division `/`

	LT    Kind = "lt"    // Comparison operator: less than `<`
	NOTEQ Kind = "noteq" // Comparison operator: not equal `<>`
	LEQ   Kind = "leq"   // Comparison operator: less than or equal `<=`
	GT    Kind = "gt"    // Comparison operator: greater than `>`
	GEQ   Kind = "geq"   // Comparison operator: greater than or equal `>=`

	OR  Kind = "or"  // Logical operator: OR `|`
	AND Kind = "and" // Logical operator: AND `&`
	NOT Kind = "not" // Logical operator: NOT `!`

	OPENPAR   Kind = "openpar"   // Bracket: opening parenthesis `(`
	CLOSEPAR  Kind = "closepar"  // Bracket: closing parenthesis `)`
	OPENCUBR  Kind = "opencubr"  // Bracket: opening curly bracket `{`
	CLOSECUBR Kind = "closecubr" // Bracket: closing curly bracket `}`
	OPENSQBR  Kind = "opensqbr"  // Bracket: opening square bracket `[`
	CLOSESQBR Kind = "closesqbr" // Bracket: closing square bracket `]`

	DOT        Kind = "dot"        // Period `.`
	COMMA      Kind = "comma"      // Comma `,`
	SEMI       Kind = "semi"       // Semicolon `;`
	COLON      Kind = "colon"      // Colon `:`
	COLONCOLON Kind = "coloncolon" // Double colon `::`

	INLINECMT   Kind = "inlinecmt"   // Single-line comment `// ... \n`
	BLOCKCMT    Kind = "blockcmt"    // Multi-line comment `/* ... */`
	CLOSEINLINE Kind = "closeinline" // End of an inline comment `\n`
	CLOSEBLOCK  Kind = "closeblock"  // End of a block comment `*/`
	OPENINLINE  Kind = "openinline"  // Start of an inline comment `//`
	OPENBLOCK   Kind = "openblock"   // Start of a block comment `/*`

	ID       Kind = "id"       // Identifier `exampleId_123`
	INTNUM   Kind = "intnum"   // Integer `123`
	FLOATNUM Kind = "floatnum" // Floating-point number `1.23`

	IF       Kind = "if"       // Reserved word `if`
	THEN     Kind = "then"     // Reserved word `then`
	ELSE     Kind = "else"     // Reserved word `else`
	INTEGER  Kind = "integer"  // Reserved word `integer`
	FLOAT    Kind = "float"    // Reserved word `float`
	VOID     Kind = "void"     // Reserved word `void`
	PUBLIC   Kind = "public"   // Reserved word `public`
	PRIVATE  Kind = "private"  // Reserved word `private`
	FUNC     Kind = "func"     // Reserved word `func`
	VAR      Kind = "var"      // Reserved word `var`
	STRUCT   Kind = "struct"   // Reserved word `struct`
	WHILE    Kind = "while"    // Reserved word `while`
	READ     Kind = "read"     // Reserved word `read`
	WRITE    Kind = "write"    // Reserved word `write`
	RETURN   Kind = "return"   // Reserved word `return`
	SELF     Kind = "self"     // Reserved word `self`
	INHERITS Kind = "inherits" // Reserved word `inherits`
	LET      Kind = "let"      // Reserved word `let`
	IMPL     Kind = "impl"     // Reserved word `impl`

	INVALIDID   Kind = "invalidid"   // Error token
	INVALIDNUM  Kind = "invalidnum"  // Error token
	INVALIDCHAR Kind = "invalidchar" // Error token
)

// Set of reserved words (empty structs as values to allocate 0 memory)
var reservedWords = map[Kind]struct{}{
	IF:       {},
	THEN:     {},
	ELSE:     {},
	INTEGER:  {},
	FLOAT:    {},
	VOID:     {},
	PUBLIC:   {},
	PRIVATE:  {},
	FUNC:     {},
	VAR:      {},
	STRUCT:   {},
	WHILE:    {},
	READ:     {},
	WRITE:    {},
	RETURN:   {},
	SELF:     {},
	INHERITS: {},
	LET:      {},
	IMPL:     {},
}

var errorSymbols = map[Kind]struct{}{
	INVALIDNUM:  {},
	INVALIDCHAR: {},
	INVALIDID:   {},
}

func IsReservedWord(s Kind) bool {
	_, ok := reservedWords[s]
	return ok
}

func IsReservedWordString(s string) (Kind, bool) {
	t := Kind(s)
	return t, IsReservedWord(t)
}

func IsError(s Kind) bool {
	_, ok := errorSymbols[s]
	return ok
}
