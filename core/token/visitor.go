package token

type Visit func(node *ASTNode)

type Visitor interface {
	Visit(node *ASTNode)
}

// A a general-purpose Visitor implementation that can be customized via the
// Dispatch table provided to it. This type is made to be embedded into Visitor
// implementations
type DispatchVisitor struct {
	Dispatch map[Kind]Visit
}

// Default visit function performs a lookup in the dispatch map, if no action is
// present, then this is a noop
func (t *DispatchVisitor) Visit(node *ASTNode) {
	if act, ok := t.Dispatch[node.Type]; ok {
		act(node)
	}
}

// A map containing noop methods for each kind of ASTNode
var DEFAULT_VISITOR_METHODS = map[Kind]Visit{
	FINAL_PROG:                        func(node *ASTNode) {},
	FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST: func(node *ASTNode) {},
	FINAL_FUNC_DEF:                    func(node *ASTNode) {},
	FINAL_STRUCT_DECL:                 func(node *ASTNode) {},
	FINAL_IMPL_DEF:                    func(node *ASTNode) {},
	FINAL_TYPE:                        func(node *ASTNode) {},
	FINAL_ID:                          func(node *ASTNode) {},
	FINAL_VOID:                        func(node *ASTNode) {},
	FINAL_INTEGER:                     func(node *ASTNode) {},
	FINAL_FLOAT:                       func(node *ASTNode) {},
	FINAL_INTNUM:                      func(node *ASTNode) {},
	FINAL_FLOATNUM:                    func(node *ASTNode) {},
	FINAL_FUNC_DEF_PARAM:              func(node *ASTNode) {},
	FINAL_FUNC_DEF_PARAMLIST:          func(node *ASTNode) {},
	FINAL_FUNC_BODY:                   func(node *ASTNode) {},
	FINAL_VAR_DECL:                    func(node *ASTNode) {},
	FINAL_FUNC_DECL:                   func(node *ASTNode) {},
	FINAL_STATEMENT:                   func(node *ASTNode) {},
	FINAL_WRITE:                       func(node *ASTNode) {},
	FINAL_IF:                          func(node *ASTNode) {},
	FINAL_WHILE:                       func(node *ASTNode) {},
	FINAL_READ:                        func(node *ASTNode) {},
	FINAL_RETURN:                      func(node *ASTNode) {},
	FINAL_ASSIGN:                      func(node *ASTNode) {},
	FINAL_EXPR:                        func(node *ASTNode) {},
	FINAL_ARITH_EXPR:                  func(node *ASTNode) {},
	FINAL_FACTOR:                      func(node *ASTNode) {},
	FINAL_TERM:                        func(node *ASTNode) {},
	FINAL_PLUS:                        func(node *ASTNode) {},
	FINAL_MINUS:                       func(node *ASTNode) {},
	FINAL_OR:                          func(node *ASTNode) {},
	FINAL_MULT:                        func(node *ASTNode) {},
	FINAL_DIV:                         func(node *ASTNode) {},
	FINAL_AND:                         func(node *ASTNode) {},
	FINAL_EQ:                          func(node *ASTNode) {},
	FINAL_NEQ:                         func(node *ASTNode) {},
	FINAL_LT:                          func(node *ASTNode) {},
	FINAL_GT:                          func(node *ASTNode) {},
	FINAL_LEQ:                         func(node *ASTNode) {},
	FINAL_GEQ:                         func(node *ASTNode) {},
	FINAL_RETURNTYPE:                  func(node *ASTNode) {},
	FINAL_NOT:                         func(node *ASTNode) {},
	FINAL_NEGATIVE:                    func(node *ASTNode) {},
	FINAL_POSITIVE:                    func(node *ASTNode) {},
	FINAL_REL_EXPR:                    func(node *ASTNode) {},
	FINAL_DIM:                         func(node *ASTNode) {},
	FINAL_DIMLIST:                     func(node *ASTNode) {},
	FINAL_INDEX:                       func(node *ASTNode) {},
	FINAL_INDEXLIST:                   func(node *ASTNode) {},
	FINAL_SUBJECT:                     func(node *ASTNode) {},
	FINAL_VARIABLE:                    func(node *ASTNode) {},
	FINAL_FUNC_CALL:                   func(node *ASTNode) {},
	FINAL_FUNC_CALL_PARAM:             func(node *ASTNode) {},
	FINAL_FUNC_CALL_PARAMLIST:         func(node *ASTNode) {},
	FINAL_STATBLOCK:                   func(node *ASTNode) {},
	FINAL_FUNC_DEF_LIST:               func(node *ASTNode) {},
	FINAL_INHERITS:                    func(node *ASTNode) {},
	FINAL_MEMBER:                      func(node *ASTNode) {},
	FINAL_MEMBERS:                     func(node *ASTNode) {},
	FINAL_PRIVATE:                     func(node *ASTNode) {},
	FINAL_PUBLIC:                      func(node *ASTNode) {},
}
