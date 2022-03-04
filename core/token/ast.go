package token

import (
	"bytes"
	"fmt"
	"io"
)

// Abstract Syntax Tree created by a parser
type AST struct {
	Root *ASTNode
}

func (a AST) TreeString() string {
	out := new(bytes.Buffer)
	a.Print(out)
	return out.String()
}

func (a AST) Print(fh io.Writer) {
	a.Root.PrintSubtree(fh, 0)
}

// A single node of the AST
type ASTNode struct {

	// TODO: for now we are storing the Type as a string, normally we should
	// TODO: "subclass" this struct
	Type Kind

	Token        Token
	Parent       *ASTNode
	Children     []*ASTNode
	SiblingLeft  *ASTNode
	SiblingRight *ASTNode
}

func (n *ASTNode) StringSubtree(depth int) string {
	out := new(bytes.Buffer)
	n.PrintSubtree(out, depth)
	return out.String()
}

func (n *ASTNode) PrintSubtree(fh io.Writer, depth int) {
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "| "
	}

	// Print me
	fmt.Fprintf(fh, "%v%v", prefix, n.Type)
	if len(n.Children) == 0 && n.Token.Id != "" {
		// If we have a leaf node, then print the token as well
		fmt.Fprintf(fh, ": %v\n", n.Token)
	} else {
		fmt.Fprintln(fh)
	}

	// Print my children
	for _, child := range n.Children {
		child.PrintSubtree(fh, depth+1)
	}
}

// TODO: This interface is as of yet unused
type ASTNodeInterface interface {
	// The semantic action symbol of this node
	Type() Kind

	// The token, if any, for this node
	Token() Token

	Parent() ASTNodeInterface
	Children() []ASTNodeInterface
	LeftSibling() ASTNodeInterface
	RightSibling() ASTNodeInterface
}
