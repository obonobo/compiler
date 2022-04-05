package codegen

import (
	"fmt"
	"strconv"

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

	*tagPool
	*RegisterPool
}

func NewTagsBasedCodeGenVisitor(out, dataOut func(string)) *TagsBasedCodeGenVisitor {
	vis := &TagsBasedCodeGenVisitor{
		out:          out,
		dataOut:      dataOut,
		logPrefix:    PREFIX,
		tagPool:      newTagPool(),
		RegisterPool: NewRegisterPool(),
	}
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
	case token.FINAL_PLUS:
		v.plus(node)
	case token.FINAL_MINUS:
		v.minus(node)
	case token.FINAL_MULT:
		v.mult(node)
	case token.FINAL_DIV:
		v.div(node)
	case token.FINAL_INTNUM:
		v.intnum(node)
	default:
		v.propagate(node)
	}
}

// We need to emit:
//
// t1 res 4
// addi r1, r0, <integer literal>
// sw t1(r0), r1
func (v *TagsBasedCodeGenVisitor) intnum(node *token.ASTNode) {
	value, _ := strconv.Atoi(string(node.Token.Lexeme))
	v.headerComment(fmt.Sprintf("INTNUM %v", value))
	tag := v.tagPool.temp()
	reg := v.RegisterPool.ClaimAny()
	defer v.Free(reg)

	v.reserveWord(tag, 4)
	v.addi(reg, R0, int32(value))
	v.sw(offR0(tag), reg)
}

func (v *TagsBasedCodeGenVisitor) plus(node *token.ASTNode) {
	v.twoOpPrintHeader("PLUS", node, v.add)
}

func (v *TagsBasedCodeGenVisitor) minus(node *token.ASTNode) {
	v.twoOpPrintHeader("SUB", node, v.sub)
}

func (v *TagsBasedCodeGenVisitor) mult(node *token.ASTNode) {
	v.twoOpPrintHeader("MULT", node, v.multiply)
}

func (v *TagsBasedCodeGenVisitor) div(node *token.ASTNode) {
	v.twoOpPrintHeader("DIV", node, v.divide)
}

func (v *TagsBasedCodeGenVisitor) arithExpr(node *token.ASTNode) {
	v.propagate(node)
}

func (v *TagsBasedCodeGenVisitor) prog(node *token.ASTNode) {
	v.comment("PROG")
	v.emit("entry")

	// The library functions that we will call need a stack pointer in r14. In
	// this tags-based approach, we will set this register and then never free
	// that register back to the pool. It will be a constant pointer to the
	// topaddr.
	v.addis(R14, R0, TOPADDR)

	v.propagate(node)
	v.emit("hlt")
}

func (v *TagsBasedCodeGenVisitor) varDecl(node *token.ASTNode) {
	id := node.Children[0].Token.Lexeme
	typee := node.Children[1].Children[0].Type
	// dimensions := node.Children[2]

	switch typee {

	// TODO: remove hardcoding
	case token.FINAL_INTEGER:
		v.reserveWord(string(id), 4)
	}
}

func (v *TagsBasedCodeGenVisitor) write(node *token.ASTNode) {
	v.propagate(node)

	// Reserve some data for a buffer
	buf, bufsize := "buf", 32
	v.emitDataf("%v	res	%v		%v Buffer for printing", buf, bufsize, token.MOON_COMMENT)

	// This is the value to be printed
	top := v.tagPool.pop()
	v.headerComment(fmt.Sprintf("WRITE(%v)", top))

	// Put the value to be printed on top of the stack
	reg := v.RegisterPool.ClaimAny()
	defer v.Free(reg)

	// Assembly
	v.lw(reg, offR0(top))
	v.sw(off(-8, R14), reg, fmt.Sprintf("	%v %v arg1", token.MOON_COMMENT, INTSTR))
	v.addis(reg, R0, buf)
	v.sw(off(-12, R14), reg, fmt.Sprintf("	%v %v arg2", token.MOON_COMMENT, INTSTR))
	v.jl(R15, INTSTR, fmt.Sprintf("	%v Procedure call %v", token.MOON_COMMENT, INTSTR))
	v.sw(off(-8, R14), R13, fmt.Sprintf("	%v %v arg1", token.MOON_COMMENT, PUTSTR))
	v.jl(R15, PUTSTR, fmt.Sprintf("	%v Procedure call %v", token.MOON_COMMENT, PUTSTR))
}

// Default action used on a node
func (v *TagsBasedCodeGenVisitor) propagate(node *token.ASTNode) {
	for _, child := range node.Children {
		child.AcceptOnce(v)
	}
}

func (v *TagsBasedCodeGenVisitor) twoOp(op func(resultReg, leftReg, rightReg string)) {
	left, right := v.tagPool.pop2()
	result := v.tagPool.temp()
	resultReg := v.RegisterPool.ClaimAny()
	leftReg := v.RegisterPool.ClaimAny()
	rightReg := v.RegisterPool.ClaimAny()
	defer v.Free(resultReg)
	defer v.Free(leftReg)
	defer v.Free(rightReg)
	v.reserveWord(result, 4)
	v.lw(leftReg, offR0(left))
	v.lw(rightReg, offR0(right))
	op(resultReg, leftReg, rightReg)
	v.sw(offR0(result), resultReg)
}

func (v *TagsBasedCodeGenVisitor) twoOpPrintHeader(
	header string,
	node *token.ASTNode,
	op func(resultReg, leftReg, rightRef string),
) {
	v.propagate(node)
	v.headerComment(header)
	v.twoOp(op)
}
