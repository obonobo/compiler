package constructs

import "github.com/obonobo/esac/core/token"

type Leaf struct {
	Symbol token.Kind
	Up     token.ASTNodeInterface
}

// The semantic action symbol of this node
func (n *Leaf) Type() token.Kind {
	panic("not implemented") // TODO: Implement
}

// The token, if any, for this node
func (n *Leaf) Token() token.Token {
	panic("not implemented") // TODO: Implement
}

func (n *Leaf) Parent() token.ASTNodeInterface {
	return n.Up
}

func (n *Leaf) Children() []token.ASTNodeInterface {
	return nil
}

func (n *Leaf) LeftSibling() token.ASTNodeInterface {
	return nil
}

func (n *Leaf) RightSibling() token.ASTNodeInterface {
	return nil
}
