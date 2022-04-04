package codegen

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
)



func (v *TagsBasedCodeGenVisitor) log(s string) {
	if v.out != nil {
		v.out(v.logPrefix + s)
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

func (v *TagsBasedCodeGenVisitor) prefix(prefix string) {
	v.logPrefix = prefix
}

func (v *TagsBasedCodeGenVisitor) useDefaultPrefix() {
	v.prefix(PREFIX)
}
