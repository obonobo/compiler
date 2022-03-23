package visitors

import "github.com/obonobo/esac/core/token"

// Performs semantic checks including type checking
type SemCheckVisitor struct {
	token.DispatchVisitor
}

func NewSemCheckVisitor() *SemCheckVisitor {
	vis := &SemCheckVisitor{}

	vis.DispatchVisitor = token.DispatchVisitor{Dispatch: map[token.Kind]token.Visit{

	}}

	return vis
}
