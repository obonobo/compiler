package visitors

import "github.com/obonobo/esac/core/token"

// A visitor that annotates symbol tables with memory size information. In this
// project, all memory allocations are known statically. This visitor determines
// the size of all memory allocations that will be emitted by the codegen
// visitor.
type MemSizeVisitor struct {
	token.DispatchVisitor
}

func NewMemSizeVisitor() *MemSizeVisitor {
	vis := &MemSizeVisitor{}
	vis.DispatchVisitor = token.DispatchVisitor{Dispatch: map[token.Kind]token.Visit{}}
	return vis
}
