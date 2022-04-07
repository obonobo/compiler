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

	bufEmitted bool
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
	case token.FINAL_REL_EXPR:
		v.relExpr(node)
	case token.FINAL_PLUS:
		v.plus(node)
	case token.FINAL_MINUS:
		v.minus(node)
	case token.FINAL_MULT:
		v.mult(node)
	case token.FINAL_DIV:
		v.div(node)

	case token.FINAL_EQ:
		v.eq(node)
	case token.FINAL_NEQ:
		v.neq(node)
	case token.FINAL_GEQ:
		v.geq(node)
	case token.FINAL_LEQ:
		v.leq(node)
	case token.FINAL_GT:
		v.gt(node)
	case token.FINAL_LT:
		v.lt(node)
	case token.FINAL_AND:
		v.and(node)
	case token.FINAL_OR:
		v.or(node)

	case token.FINAL_IF:
		v.ifStatement(node)
	case token.FINAL_WHILE:
		v.whileStatement(node)

	case token.FINAL_INTNUM:
		v.intnum(node)
	case token.FINAL_ASSIGN:
		v.assign(node)
	case token.FINAL_FACTOR:
		v.factor(node)
	case token.FINAL_STRUCT_DECL:
	default:
		v.propagate(node)
	}
}

func (v *TagsBasedCodeGenVisitor) whileStatement(node *token.ASTNode) {
	v.headerComment("While statement start")

	// Label loop test, traverse expr
	dowhile := "dowhile" + v.tagPool.temp()[1:]
	endwhile := "endwhile" + dowhile[7:]
	v.nop(dowhile)
	v.propagate(node.Children[0])
	reg := v.RegisterPool.ClaimAny()
	defer v.Free(reg)
	top := v.tagPool.pop()
	v.lw(reg, offR0(top))
	v.bz(reg, endwhile)

	// Loop body
	v.propagate(node.Children[1])
	v.j(dowhile)
	v.nop(endwhile)
	v.headerComment("While statement end")
}

func (v *TagsBasedCodeGenVisitor) ifStatement(node *token.ASTNode) {
	v.headerComment("If statement start")

	// Traverse expr
	v.propagate(node.Children[0])
	reg1 := v.RegisterPool.ClaimAny()
	defer v.Free(reg1)
	top := v.tagPool.pop()
	elseTag := v.tagPool.elseTag()
	endIfTag := "endIf" + elseTag[4:]
	v.lw(reg1, offR0(top))
	v.bz(reg1, elseTag)

	// Traverse the 'then' block
	v.propagate(node.Children[1])
	v.j(endIfTag)

	// Emit the 'else' block, starting with a tagged nop
	v.nop(elseTag)
	v.propagate(node.Children[2])
	v.nop(endIfTag)
	v.headerComment("If statement end")
}

func (v *TagsBasedCodeGenVisitor) eq(node *token.ASTNode) {
	v.twoOpPrintHeader("EQ(==)", node, v.equal)
}

func (v *TagsBasedCodeGenVisitor) neq(node *token.ASTNode) {
	v.twoOpPrintHeader("Not equal <>", node, v.notEqual)
}

func (v *TagsBasedCodeGenVisitor) geq(node *token.ASTNode) {
	v.twoOpPrintHeader("Greater than or equal >=", node, v.greaterOrEqual)
}

func (v *TagsBasedCodeGenVisitor) leq(node *token.ASTNode) {
	v.twoOpPrintHeader("Less than or equal <=", node, v.lessOrEqual)
}

func (v *TagsBasedCodeGenVisitor) gt(node *token.ASTNode) {
	v.twoOpPrintHeader("Greater than >", node, v.greater)
}

func (v *TagsBasedCodeGenVisitor) lt(node *token.ASTNode) {
	v.twoOpPrintHeader("Less than <", node, v.less)
}

func (v *TagsBasedCodeGenVisitor) and(node *token.ASTNode) {
	v.propagate(node)
}

func (v *TagsBasedCodeGenVisitor) or(node *token.ASTNode) {
	v.propagate(node)
}

func (v *TagsBasedCodeGenVisitor) relExpr(node *token.ASTNode) {
	v.propagate(node)
}

func (v *TagsBasedCodeGenVisitor) assign(node *token.ASTNode) {
	v.propagate(node)
	lhs := node.Children[0].Children[1].Token.Lexeme
	rhs := v.tagPool.pop()
	v.headerComment(fmt.Sprintf("ASSIGN %v = %v", lhs, rhs))
	reg := v.RegisterPool.ClaimAny()
	defer v.Free(reg)
	v.lw(reg, offR0(rhs))
	v.sw(offR0(lhs), reg)
}

func (v *TagsBasedCodeGenVisitor) factor(node *token.ASTNode) {
	v.propagate(node)
	switch child := node.Children[0]; child.Type {
	case token.FINAL_VARIABLE:
		id := string(child.Children[1].Token.Lexeme)
		v.tagPool.push(id)
	}
}

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
	dimensions := node.Children[2].Children
	size := v.sizeof(typee, dimensions, node)
	v.reserveWord(string(id), size)
}

// Returns the size of some data type pointed to by a node
func (v *TagsBasedCodeGenVisitor) sizeof(
	typee token.Kind,
	dimensions []*token.ASTNode,
	node *token.ASTNode,
) int {
	size := 4
	switch typee {
	case token.FINAL_INTEGER:
	case token.FINAL_FLOAT:
	case token.FINAL_ID:
		symbolTable := node.Meta.Record.Link
		size = computeSymboltableSize(symbolTable)
	default:
	}
	multiplyByDimlist(&size, dimensions)
	return size
}

func computeSymboltableSize(table token.SymbolTable) int {
	size := 0
	for _, entry := range table.Entries() {
		size += sizeofRecord(&entry)
	}
	return size
}

func sizeofRecord(record *token.SymbolTableRecord) int {
	if record.Kind != token.FINAL_VAR_DECL {
		return 0
	}
	size := 0
	switch record.Type.Type {
	case token.FINAL_INTEGER, token.FINAL_FLOAT:
		size = 4
	case token.FINAL_ID:
		size = computeSymboltableSize(record.Link)
	}
	for _, dim := range record.Type.Dimlist {
		size *= dim
	}
	return size
}

func multiplyByDimlist(size *int, dimlist []*token.ASTNode) {
	for _, dim := range dimlist {
		if dimension, err := strconv.Atoi(string(dim.Token.Lexeme)); err == nil {
			*size *= dimension
		}
	}
}

func (v *TagsBasedCodeGenVisitor) write(node *token.ASTNode) {
	v.propagate(node)

	// Reserve some data for a buffer
	buf := "wbuf"
	if !v.bufEmitted {
		bufsize := 32
		v.bufEmitted = true
		v.emitDataf("%v	res	%v		%v Buffer for printing", buf, bufsize, token.MOON_COMMENT)
	}

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
