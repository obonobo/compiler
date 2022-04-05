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

// Adds the value of the top
func (v *TagsBasedCodeGenVisitor) plus(node *token.ASTNode) {
	v.propagate(node)
	v.headerComment("PLUS")

	top1, top2 := v.tagPool.pop2()
	result := v.tagPool.temp()
	resultReg := v.RegisterPool.ClaimAny()
	leftReg := v.RegisterPool.ClaimAny()
	rightReg := v.RegisterPool.ClaimAny()
	defer v.Free(leftReg)
	defer v.Free(rightReg)
	defer v.Free(resultReg)

	v.reserveWord(result, 4)
	v.lw(leftReg, offR0(top1))
	v.lw(rightReg, offR0(top2))
	v.add(resultReg, leftReg, rightReg)
	v.sw(offR0(result), resultReg)
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
	v.lw(reg, offR0(top))
	v.sw(off(-8, R14), reg, fmt.Sprintf("	%v %v arg1", token.MOON_COMMENT, INTSTR))

	// Put address of buffer on stack
	v.addis(reg, R0, buf)
	v.sw(off(-12, R14), reg, fmt.Sprintf("	%v %v arg2", token.MOON_COMMENT, INTSTR))

	// Call `intstr` procedure from `lib.m`. This will convert the integer to a
	// string
	v.jl(R15, INTSTR, fmt.Sprintf("	%v Procedure call %v", token.MOON_COMMENT, INTSTR))

	// Function return is in r13, we'll use that right away
	v.sw(off(-8, R14), R13, fmt.Sprintf("	%v %v arg1", token.MOON_COMMENT, PUTSTR))

	// Call putstr to print the result
	v.jl(R15, PUTSTR, fmt.Sprintf("	%v Procedure call %v", token.MOON_COMMENT, PUTSTR))
}

// Default action used on a node
func (v *TagsBasedCodeGenVisitor) propagate(node *token.ASTNode) {
	for _, child := range node.Children {
		child.AcceptOnce(v)
	}
}
