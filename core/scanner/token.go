package scanner

import "fmt"

// Unique token identifier
type Symbol string

// The text representing a token
type Lexeme string

type Token struct {
	Id     Symbol // The unique identifier of this token
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
	ASSIGN Symbol = "assign" // Assignment operator
	ARROW  Symbol = "arrow"  // Right-pointing arrow operator

	EQ    Symbol = "eq"    // Arithmetic operator: equality
	PLUS  Symbol = "plus"  // Arithmetic operator: addition
	MINUS Symbol = "minus" // Arithmetic operator: subtraction
	MULT  Symbol = "mult"  // Arithmetic operator: multiplication
	DIV   Symbol = "div"   // Arithmetic operator: division

	LT    Symbol = "lt"    // Comparison operator: less than
	NOTEQ Symbol = "noteq" // Comparison operator: not equal
	LEQ   Symbol = "leq"   // Comparison operator: less than or equal
	GT    Symbol = "gt"    // Comparison operator: greater than
	GEQ   Symbol = "geq"   // Comparison operator: greater than or equal

	OR  Symbol = "or"  // Logical operator: OR
	AND Symbol = "and" // Logical operator: AND
	NOT Symbol = "not" // Logical operator: NOT

	OPENPAR   Symbol = "openpar"   // Bracket: opening parenthesis
	CLOSEPAR  Symbol = "closepar"  // Bracket: closing parenthesis
	OPENCUBR  Symbol = "opencubr"  // Bracket: opening curly bracket
	CLOSECUBR Symbol = "closecubr" // Bracket: closing curly bracket
	OPENSQBR  Symbol = "opensqbr"  // Bracket: opening square bracket
	CLOSESQBR Symbol = "closesqbr" // Bracket: closing square bracket

	DOT        Symbol = "dot"        // Period
	COMMA      Symbol = "comma"      // Comma
	SEMI       Symbol = "semi"       // Semicolon
	COLON      Symbol = "colon"      // Colon
	COLONCOLON Symbol = "coloncolon" // Double colon

	INLINECMT Symbol = "inlinecmt" // Single-line comment
	BLOCKCMT  Symbol = "blockcmt"  // Multi-line comment

	ID       Symbol = "id"       // Identifier
	INTNUM   Symbol = "intnum"   // Integer
	FLOATNUM Symbol = "floatnum" // Floating-point number

	IF       Symbol = "if"       // Reserved word
	ELSE     Symbol = "else"     // Reserved word
	INTEGER  Symbol = "integer"  // Reserved word
	FLOAT    Symbol = "float"    // Reserved word
	VOID     Symbol = "void"     // Reserved word
	PUBLIC   Symbol = "public"   // Reserved word
	PRIVATE  Symbol = "private"  // Reserved word
	FUNC     Symbol = "func"     // Reserved word
	VAR      Symbol = "var"      // Reserved word
	STRUCT   Symbol = "struct"   // Reserved word
	WHILE    Symbol = "while"    // Reserved word
	READ     Symbol = "read"     // Reserved word
	WRITE    Symbol = "write"    // Reserved word
	RETURN   Symbol = "return"   // Reserved word
	SELF     Symbol = "self"     // Reserved word
	INHERITS Symbol = "inherits" // Reserved word
	LET      Symbol = "let"      // Reserved word
	IMPL     Symbol = "impl"     // Reserved word

	INVALIDNUM        Symbol = "invalidnum"        // Error token
	INVALIDCHAR       Symbol = "invalidchar"       // Error token
	INVALIDIDENTIFIER Symbol = "invalididentifier" // Error token
)

// Set of reserved words (empty structs as values to allocate 0 memory)
var reservedWords = map[Symbol]struct{}{
	IF:       {},
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

var errorSymbols = map[Symbol]struct{}{
	INVALIDNUM:        {},
	INVALIDCHAR:       {},
	INVALIDIDENTIFIER: {},
}

func IsReservedWord(s Symbol) bool {
	_, ok := reservedWords[s]
	return ok
}

func IsReservedWordString(s string) (Symbol, bool) {
	t := Symbol(s)
	return t, IsReservedWord(t)
}

func IsError(s Symbol) bool {
	_, ok := errorSymbols[s]
	return ok
}
