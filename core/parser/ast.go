package parser

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
)

// Abstract Syntax Tree created by a parser
type AST struct {
	Root *ASTNode
}

func (ast *AST) String() string {
	return fmt.Sprintf("AST[%v]", ast.Root)
}

// A single node of the AST
type ASTNode struct {
	Symbol   token.Kind
	Parent   *ASTNode
	Children []*ASTNode
}

func (n *ASTNode) AddChildren(children ...*ASTNode) {
	for _, child := range children {
		if child == nil {
			continue
		}
		child.Parent = n
		n.Children = append(n.Children, child)
	}
}

func (n *ASTNode) String() string {
	return fmt.Sprintf(
		"ASTNode[Parent=%v, Symbol=%v, Children=%v]",
		n.Parent, n.Symbol, n.Children)
}
