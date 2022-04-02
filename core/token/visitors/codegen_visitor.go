package visitors

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
)

// A visitor that emits MOON assembly that uses a tags-based approach to
// variables and memory allocation.
//
// For data declarations, use the `dataOut` callback.
//
// For regular assembly instructions, use the `out` callback.
type TagsBasedCodeGenVisitor struct {
	token.DispatchVisitor

	// Assembly instructions will be printed to this callback, one line at a
	// time
	out func(string)

	// Memory reserved and tagged will be printed to this callback, one line at
	// a time
	dataOut func(string)
}

func NewTagsBasedCodeGenVisitor(out, dataOut func(string)) *TagsBasedCodeGenVisitor {
	vis := &TagsBasedCodeGenVisitor{out: out, dataOut: dataOut}
	vis.DispatchVisitor = token.DispatchVisitor{Dispatch: map[token.Kind]token.Visit{
		token.FINAL_VAR_DECL: vis.varDecl,
	}}
	return vis
}

func (v *TagsBasedCodeGenVisitor) varDecl(node *token.ASTNode) {
	id := node.Children[0].Token.Lexeme
	typee := node.Children[1].Children[0].Type
	// dimensions := node.Children[2]

	switch typee {

	// TODO: remove hardcoding
	case token.FINAL_INTEGER:
		size := 4
		v.logDataf("	%v Space for variable %v", token.MOON_COMMENT, id)
		v.logDataf("%v	res %v", id, size)
	}
}

func (v *TagsBasedCodeGenVisitor) log(s string) {
	if v.out != nil {
		v.out(s)
	}
}

func (v *TagsBasedCodeGenVisitor) logf(format string, a ...any) {
	v.log(fmt.Sprintf(format, a...))
}

func (v *TagsBasedCodeGenVisitor) logData(s string) {
	if v.dataOut != nil {
		v.dataOut(s)
	}
}

func (v *TagsBasedCodeGenVisitor) logDataf(format string, a ...any) {
	v.logData(fmt.Sprintf(format, a...))
}

func (v *TagsBasedCodeGenVisitor) comment(format string, a ...any) {
	v.logf(fmt.Sprintf("%v %v", token.MOON_COMMENT, format), a...)
}

func (v *TagsBasedCodeGenVisitor) commentData(format string, a ...any) {
	v.logDataf(fmt.Sprintf("%v %v", token.MOON_COMMENT, format), a...)
}
