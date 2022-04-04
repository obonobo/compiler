package codegen

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
)

func (v *TagsBasedCodeGenVisitor) reserveWord(forVariable string, size int) {
	v.emitDataf(""+
		"%v	res %v	%v Space for variable %v",
		forVariable, size, token.MOON_COMMENT, forVariable)
}

// Emits an `addi` instruction
func (v *TagsBasedCodeGenVisitor) addi(ri, rj string, k int32) {
	v.emitf("addi %v, %v, %v", ri, rj, k)
}

func (v *TagsBasedCodeGenVisitor) add(ri, rj, rk string) {
	v.emitf("add %v, %v, %v", ri, rj, rk)
}

// E.g.: sw t1(r0), r1
func (v *TagsBasedCodeGenVisitor) sw(krj, ri string) {
	v.emitf("sw %v, %v", krj, ri)
}

// E.g.: lw r1, t1(r0)
func (v *TagsBasedCodeGenVisitor) lw(ri, krj string) {
	v.emitf("lw %v, %v", ri, krj)
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
	v.emitf(fmt.Sprintf("%v %v", token.MOON_COMMENT, format), a...)
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

func offR0(tag string) string {
	return fmt.Sprintf("%v(%v)", tag, r0)
}
