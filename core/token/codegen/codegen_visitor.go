package codegen

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
)

const (
	// Format prefix for printed assembly code
	PREFIX = "	"

	// Expression output tag. The output of evaluating an expression must be
	// tagged with this tag in order for it to be picked up by e.g. a `write`
	// statement
	TN = "tn"
)

// A visitor that emits MOON assembly that uses a tags-based approach to
// variables and memory allocation.
//
// For data declarations, use the `dataOut` callback.
//
// For regular assembly instructions, use the `out` callback.
type TagsBasedCodeGenVisitor struct {
	// Assembly instructions will be printed to this callback, one line at a
	// time
	out func(string)

	// Memory reserved and tagged will be printed to this callback, one line at
	// a time
	dataOut func(string)

	// A prefix that is used internally for logging
	logPrefix string

	tagPool
	RegisterPool
}

func NewTagsBasedCodeGenVisitor(out, dataOut func(string)) *TagsBasedCodeGenVisitor {
	vis := &TagsBasedCodeGenVisitor{out: out, dataOut: dataOut, logPrefix: PREFIX}
	return vis
}

func (v *TagsBasedCodeGenVisitor) Visit(node *token.ASTNode) {
	switch node.Type {
	case token.FINAL_PROG:
		v.prog(node)
	case token.FINAL_VAR_DECL:
		v.varDecl(node)
	case token.FINAL_WRITE:
		v.write(node)
	case token.FINAL_ARITH_EXPR:
		v.arithExpr(node)
	case token.FINAL_INTNUM:
		v.intnum(node)
	case token.FINAL_FLOATNUM:
		v.floatnum(node)
	default:
		v.propagate(node)
	}
}

func (v *TagsBasedCodeGenVisitor) floatnum(node *token.ASTNode) {

}

func (v *TagsBasedCodeGenVisitor) intnum(node *token.ASTNode) {

}

func (v *TagsBasedCodeGenVisitor) arithExpr(node *token.ASTNode) {
	// Process all children
	v.propagate(node)
	fmt.Println()
}

func (v *TagsBasedCodeGenVisitor) prog(node *token.ASTNode) {
	v.log("entry")
	v.propagate(node)
	v.log("hlt")
}

func (v *TagsBasedCodeGenVisitor) varDecl(node *token.ASTNode) {
	id := node.Children[0].Token.Lexeme
	typee := node.Children[1].Children[0].Type
	// dimensions := node.Children[2]

	switch typee {

	// TODO: remove hardcoding
	case token.FINAL_INTEGER:
		size := 4
		v.logDataf("%v	res %v	%v Space for variable %v", id, size, token.MOON_COMMENT, id)
	}
}

func (v *TagsBasedCodeGenVisitor) write(node *token.ASTNode) {
	v.propagate(node)
}

// Default action used on a node
func (v *TagsBasedCodeGenVisitor) propagate(node *token.ASTNode) {
	for _, child := range node.Children {
		child.AcceptOnce(v)
	}
}
