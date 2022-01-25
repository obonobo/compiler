package scanner

// Unique token identifier
type TokenType string

const (
	ASSIGN TokenType = "assign" // Assignment operator
	ARROW  TokenType = "arrow"  // Right-pointing arrow operator

	EQ    TokenType = "eq"    // Arithmetic operator: equality
	PLUS  TokenType = "plus"  // Arithmetic operator: addition
	MINUS TokenType = "minus" // Arithmetic operator: subtraction
	MULT  TokenType = "mult"  // Arithmetic operator: multiplication
	DIV   TokenType = "div"   // Arithmetic operator: division

	LT    TokenType = "lt"    // Comparison operator: less than
	NOTEQ TokenType = "noteq" // Comparison operator: not equal
	LEQ   TokenType = "leq"   // Comparison operator: less than or equal
	GT    TokenType = "gt"    // Comparison operator: greater than
	GEQ   TokenType = "geq"   // Comparison operator: greater than or equal

	OR  TokenType = "or"  // Logical operator: OR
	AND TokenType = "and" // Logical operator: AND
	NOT TokenType = "not" // Logical operator: NOT

	OPENPAR   TokenType = "openpar"   // Bracket: opening parenthesis
	CLOSEPAR  TokenType = "closepar"  // Bracket: closing parenthesis
	OPENCUBR  TokenType = "opencubr"  // Bracket: opening curly bracket
	CLOSECUBR TokenType = "closecubr" // Bracket: closing curly bracket
	OPENSQBR  TokenType = "opensqbr"  // Bracket: opening square bracket
	CLOSESQBR TokenType = "closesqbr" // Bracket: closing square bracket

	DOT        TokenType = "dot"        // Period
	COMMA      TokenType = "comma"      // Comma
	SEMI       TokenType = "semi"       // Semicolon
	COLON      TokenType = "colon"      // Colon
	COLONCOLON TokenType = "coloncolon" // Double colon

	INLINECMT TokenType = "inlinecmt" // Single-line comment
	BLOCKCMT  TokenType = "blockcmt"  // Multi-line comment

	ID       TokenType = "id"       // Identifier
	INTNUM   TokenType = "intnum"   // Integer
	FLOATNUM TokenType = "floatnum" // Floating-point number

	IF       TokenType = "if"       // Reserved word
	ELSE     TokenType = "else"     // Reserved word
	INTEGER  TokenType = "integer"  // Reserved word
	FLOAT    TokenType = "float"    // Reserved word
	VOID     TokenType = "void"     // Reserved word
	PUBLIC   TokenType = "public"   // Reserved word
	PRIVATE  TokenType = "private"  // Reserved word
	FUNC     TokenType = "func"     // Reserved word
	VAR      TokenType = "var"      // Reserved word
	STRUCT   TokenType = "struct"   // Reserved word
	WHILE    TokenType = "while"    // Reserved word
	READ     TokenType = "read"     // Reserved word
	WRITE    TokenType = "write"    // Reserved word
	RETURN   TokenType = "return"   // Reserved word
	SELF     TokenType = "self"     // Reserved word
	INHERITS TokenType = "inherits" // Reserved word
	LET      TokenType = "let"      // Reserved word
	IMPL     TokenType = "impl"     // Reserved word
)

// Set of reserved words (empty structs as values to allocate 0 memory)
var reservedWords = map[TokenType]struct{}{
	IF:       struct{}{},
	ELSE:     struct{}{},
	INTEGER:  struct{}{},
	FLOAT:    struct{}{},
	VOID:     struct{}{},
	PUBLIC:   struct{}{},
	PRIVATE:  struct{}{},
	FUNC:     struct{}{},
	VAR:      struct{}{},
	STRUCT:   struct{}{},
	WHILE:    struct{}{},
	READ:     struct{}{},
	WRITE:    struct{}{},
	RETURN:   struct{}{},
	SELF:     struct{}{},
	INHERITS: struct{}{},
	LET:      struct{}{},
	IMPL:     struct{}{},
}

func IsReservedWord(s string) (TokenType, bool) {
	t := TokenType(s)
	_, ok := reservedWords[t]
	return t, ok
}

// The text representing a token
type Lexeme string

type Token struct {
	Id     TokenType // The unique identifier of this token
	Lexeme Lexeme    // The exact string that was matched as this token
	Line   int       // The line number on which the token was found
	Column int       // The column number on which the token was found
}
