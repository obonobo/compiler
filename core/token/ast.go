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
	Type     Kind
	Token    Token
	Children []*ASTNode
}

func (n *ASTNode) StringSubtree(depth int) string {
	out := new(bytes.Buffer)
	n.PrintSubtree(out, depth)
	return out.String()
}

func (n *ASTNode) PrintSubtree(fh io.Writer, depth int) {
	PRINT_TOKEN_ONLY_IF_0_CHILDREN_ENABLED := true
	// PRINT_TOKEN_ONLY_IF_0_CHILDREN_ENABLED := false
	PRINT_TOKEN_ENABLED := true
	// PRINT_TOKEN_ENABLED := false

	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "| "
	}

	printToken := (len(n.Children) == 0 ||
		!PRINT_TOKEN_ONLY_IF_0_CHILDREN_ENABLED) &&
		n.Token.Id != "" &&
		PRINT_TOKEN_ENABLED

	fmt.Fprintf(fh, "%v%v", prefix, n.Type)
	if printToken {
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

func (n *ASTNode) String() string {
	out := fmt.Sprintf("ASTNode[Type=%v, ", n.Type)
	if n.Token.Id != "" {
		out += fmt.Sprintf("Token=%v, ", n.Token)
	}
	out += fmt.Sprintf("Children=%v]", n.Children)
	return out
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
