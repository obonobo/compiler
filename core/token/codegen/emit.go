package codegen

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/util"
)

func (v *TagsBasedCodeGenVisitor) jl(ri, proc string, comments ...string) {
	v.emit(withComments(fmt.Sprintf("jl	%v, %v", ri, proc), comments...))
}

// Emits an `addi` instruction
func (v *TagsBasedCodeGenVisitor) addi(ri, rj string, k int32) {
	addix(v, ri, rj, k)
}

func (v *TagsBasedCodeGenVisitor) addis(ri, rj, k string) {
	addix(v, ri, rj, k)
}

func addix[X ~string | ~int32](v *TagsBasedCodeGenVisitor, ri, rj string, x X) {
	v.emitf("addi	%v, %v, %v", ri, rj, x)
}

func (v *TagsBasedCodeGenVisitor) add(ri, rj, rk string) {
	v.emit3op("add", ri, rj, rk)
}

func (v *TagsBasedCodeGenVisitor) multiply(ri, rj, rk string) {
	v.emit3op("mul", ri, rj, rk)
}

func (v *TagsBasedCodeGenVisitor) divide(ri, rj, rk string) {
	v.emit3op("div", ri, rj, rk)
}

func (v *TagsBasedCodeGenVisitor) sub(ri, rj, rk string) {
	v.emit3op("sub", ri, rj, rk)
}

// E.g.: sw t1(r0), r1
func (v *TagsBasedCodeGenVisitor) sw(krj, ri string, comments ...string) {
	v.emit(withComments(fmt.Sprintf("sw	%v, %v", krj, ri), comments...))
}

// E.g.: lw r1, t1(r0)
func (v *TagsBasedCodeGenVisitor) lw(ri, krj string) {
	v.emitf("lw	%v, %v", ri, krj)
}

func (v *TagsBasedCodeGenVisitor) emit(s string) {
	if v.out != nil {
		v.out(v.logPrefix + s)
	}
}

func (v *TagsBasedCodeGenVisitor) emitf(format string, a ...any) {
	v.emit(fmt.Sprintf(format, a...))
}

func (v *TagsBasedCodeGenVisitor) emitData(s string) {
	if v.dataOut != nil {
		v.dataOut(s)
	}
}

func (v *TagsBasedCodeGenVisitor) emitDataf(format string, a ...any) {
	v.emitData(fmt.Sprintf(format, a...))
}

func (v *TagsBasedCodeGenVisitor) comment(format string, a ...any) {
	v.emit(fmt.Sprintf("%v %v", append([]any{token.MOON_COMMENT, format}, a...)...))
}

func (v *TagsBasedCodeGenVisitor) commentData(format string, a ...any) {
	v.emitDataf(fmt.Sprintf("%v %v", token.MOON_COMMENT, format), a...)
}

func (v *TagsBasedCodeGenVisitor) prefix(prefix string) {
	v.logPrefix = prefix
}

func (v *TagsBasedCodeGenVisitor) useDefaultPrefix() {
	v.prefix(PREFIX)
}

func (v *TagsBasedCodeGenVisitor) reserveWord(forVariable string, size int) {
	v.emitDataf(""+
		"%v	res	%v		%v Space for variable %v",
		forVariable, size, token.MOON_COMMENT, forVariable)
}

func (v *TagsBasedCodeGenVisitor) headerComment(header string) {
	v.emit("")
	v.comment(header)
}

func (v *TagsBasedCodeGenVisitor) emit3op(op, ri, rj, rk string) {
	v.emitf("%v	%v, %v, %v", op, ri, rj, rk)
}

func withComments(s string, comments ...string) string {
	if len(comments) > 0 {
		s += strings.Join(comments, " ")
	}
	return s
}

func off[V1, V2 util.Ordered](outer V2, inner V1) string {
	return fmt.Sprintf("%v(%v)", outer, inner)
}

func offR0[T util.Ordered](tag T) string {
	return off(tag, R0)
}
