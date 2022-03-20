package token

import (
	"bytes"
	"fmt"
	"io"
	"strings"
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

// This struct can be used to extend the ASTNode with any extra information
// introduced by a Visitor, e.g.: attach a symbol table here, or a symbol table
// entry
//
// All fields are optional and they may be nil on a given node, you have to
// check this yourself in your code
type Meta struct {
	Record      *SymbolTableRecord
	SymbolTable SymbolTable
}

func (m Meta) String() string {
	return fmt.Sprintf(
		`Meta[Record=%v, SymbolTable="%v"]`,
		*m.Record, m.SymbolTable.Id())
}

// A single node of the AST
type ASTNode struct {
	Type     Kind
	Token    Token
	Children []*ASTNode

	// A place to attach extra information about the node, e.g. a Symbol Table
	Meta Meta
}

func (n *ASTNode) Accept(v Visitor) {
	for _, child := range n.Children {
		child.Accept(v)
	}
	v.Visit(n)
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

	var prefix strings.Builder
	for i := 0; i < depth; i++ {
		prefix.WriteString("| ")
	}

	pref := prefix.String()
	printToken := (len(n.Children) == 0 ||
		!PRINT_TOKEN_ONLY_IF_0_CHILDREN_ENABLED) &&
		n.Token.Id != "" &&
		PRINT_TOKEN_ENABLED

	fmt.Fprintf(fh, "%v%v", pref, n.Type)
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
	var out strings.Builder
	fmt.Fprintf(&out, "ASTNode[Type=%v, ", n.Type)
	if n.Token.Id != "" {
		fmt.Fprintf(&out, "Token=%v, ", n.Token)
	}
	if len(n.Children) > 0 {
		fmt.Fprint(&out, "Children=[...]]")
	} else {
		fmt.Fprint(&out, "Children=[]]")
	}
	return out.String()
}
