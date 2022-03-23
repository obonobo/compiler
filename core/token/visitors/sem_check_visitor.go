package visitors

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/token"
)

type record struct {
	Node *token.ASTNode
}

// Performs semantic checks including type checking
type SemCheckVisitor struct {
	token.DispatchVisitor
	errout func(e *VisitorError)
}

func NewSemCheckVisitor(errout func(e *VisitorError)) *SemCheckVisitor {
	vis := &SemCheckVisitor{errout: errout}
	vis.DispatchVisitor = token.DispatchVisitor{Dispatch: map[token.Kind]token.Visit{
		token.FINAL_FUNC_DEF: vis.typeCheckFunction,
		token.FINAL_VAR_DECL: vis.attachVarDecl,
	}}
	return vis
}

func (vis *SemCheckVisitor) typeCheckFunction(node *token.ASTNode) {
	table := node.Meta.SymbolTable
	vis.attachReturnTypeTable(node)
	body := childrenWithoutVarDecls(node.Children[3].Children)
	vis.typeCheckBlock(table, body)
}

func (vis *SemCheckVisitor) doCheck(table token.SymbolTable, node *token.ASTNode) {
	switch node.Type {
	case token.FINAL_RETURN:
		vis.typeCheckReturn(table, node)
	case token.FINAL_ASSIGN:
		vis.typeCheckAssign(table, node)
	case token.FINAL_IF:
		vis.typeCheckIf(table, node)
	case token.FINAL_WHILE:
		vis.typeCheckWhile(table, node)
	case token.FINAL_READ:
		vis.typeCheckRead(table, node)
	case token.FINAL_WRITE:
		vis.typeCheckWrite(table, node)
	case token.FINAL_FUNC_CALL:
		vis.typeCheckFunctionCall(table, node)
	default:
		vis.typeCheck(table, node)
	}
}

func (vis *SemCheckVisitor) typeCheck(table token.SymbolTable, node *token.ASTNode) token.Type {
	switch child := node.Children[0]; child.Type {
	case token.FINAL_FACTOR:
		return vis.typeCheck(table, child)
	case token.FINAL_ARITH_EXPR:
		return vis.typeCheck(table, child)
	case token.FINAL_VARIABLE:
		return vis.typeCheckVariable(table, child)
	case token.FINAL_FUNC_CALL:
		return vis.typeCheckFunctionCall(table, child)
	case token.FINAL_SUBJECT:
		return vis.typeCheckSubject(table, child)
	case token.FINAL_MULT,
		token.FINAL_DIV,
		token.FINAL_AND,
		token.FINAL_OR,
		token.FINAL_PLUS,
		token.FINAL_MINUS:
		return vis.typeCheckBinaryOperator(table, child)
	case token.FINAL_EQ,
		token.FINAL_NEQ,
		token.FINAL_LEQ,
		token.FINAL_LT,
		token.FINAL_GT,
		token.FINAL_GEQ:
		return vis.typeCheckComparison(table, child)
	case token.FINAL_NOT:
		return vis.typeCheckUnaryOperator(table, child)
	case token.FINAL_NEGATIVE, token.FINAL_POSITIVE:
		return vis.typeCheckSigned(table, node)
	case token.FINAL_INTNUM:
		return token.Type{Type: token.FINAL_INTEGER, Token: child.Token}
	case token.FINAL_FLOATNUM:
		return token.Type{Type: token.FINAL_FLOAT, Token: child.Token}
	default:
		panic(fmt.Errorf("not implemented! %v", child))
	}
}

func (vis *SemCheckVisitor) typeCheckSigned(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	// This will be a Factor with two children
	value := node.Children[1]
	return vis.typeCheck(table, value)
}

// If statement has 3 parts: relExpr, statBlock, statBlock
func (vis *SemCheckVisitor) typeCheckIf(table token.SymbolTable, node *token.ASTNode) {
	relExpr := node.Children[0]
	statBlock1 := node.Children[1].Children
	statBlock2 := node.Children[2].Children

	vis.assertRelExpr(relExpr, fmt.Sprintf("if expression (line %v): ", node.Token.Line))
	vis.typeCheck(table, relExpr)
	vis.typeCheckBlock(table, statBlock1)
	vis.typeCheckBlock(table, statBlock2)
}

// While has 2 parts: relExpr, statBlock
func (vis *SemCheckVisitor) typeCheckWhile(table token.SymbolTable, node *token.ASTNode) {
	relExpr := node.Children[0]
	statBlock := node.Children[1].Children

	vis.assertRelExpr(relExpr, fmt.Sprintf("while expression (line %v): ", node.Token.Line))
	vis.typeCheck(table, relExpr)
	vis.typeCheckBlock(table, statBlock)
}

func (vis *SemCheckVisitor) typeCheckRead(table token.SymbolTable, node *token.ASTNode) {
	variable := node.Children[0]
	vis.assertVariable(variable)
	vis.typeCheck(table, variable)
}

func (vis *SemCheckVisitor) assertVariable(node *token.ASTNode, msgPrefix ...string) {
	if !isTypeNode(node, token.FINAL_VARIABLE) {
		vis.logTypeCheckError(fmt.Sprintf("%vexpected %v but got %v",
			strings.Join(msgPrefix, ""), token.FINAL_VARIABLE, node.Type))
	}
}

func (vis *SemCheckVisitor) typeCheckWrite(table token.SymbolTable, node *token.ASTNode) {
	expr := node.Children[0]
	vis.assertExpr(expr, fmt.Sprintf("write expression (line %v): ", node.Token.Line))
	vis.typeCheck(table, expr)
}

func (vis *SemCheckVisitor) assertExpr(node *token.ASTNode, msgPrefix ...string) {
	or := strings.Join([]string{
		string(token.FINAL_REL_EXPR),
		string(token.FINAL_ARITH_EXPR),
	}, ", ")

	if !isTypeNode(node, token.FINAL_REL_EXPR, token.FINAL_ARITH_EXPR) {
		vis.logTypeCheckError(fmt.Sprintf("%vexpected %v but got %v",
			strings.Join(msgPrefix, ""), or, node.Type))
	}
}

func (vis *SemCheckVisitor) assertRelExpr(node *token.ASTNode, msgPrefix ...string) {
	if !isTypeNode(node, token.FINAL_REL_EXPR) {
		vis.logTypeCheckError(fmt.Sprintf("%vexpected %v but got %v",
			strings.Join(msgPrefix, ""), token.FINAL_REL_EXPR, node.Type))
	}
}

// For operators like <, >, ==, <>, etc.
func (vis *SemCheckVisitor) typeCheckComparison(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	ret := token.Type{
		Type:    token.FINAL_INTEGER, // We will say that booleans are integers
		Token:   node.Token,
		Dimlist: []int{},
	}

	// For this expression to be valid, the operands must be of the same type,
	// but they also must be of a specific set of types that are supported in
	// comparison operations. There is no operator overloading, and custom types
	// cannot be used in comparison operations
	left := vis.typeCheck(table, node.Children[0])
	right := vis.typeCheck(table, node.Children[1])
	if !left.EqualsNoPrivacy(right) {
		vis.emitBinaryOperatorTypeMismatchError(node, left, right)
	}
	return ret
}

// For operators like NOT
func (vis *SemCheckVisitor) typeCheckUnaryOperator(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	t := vis.typeCheck(table, node.Children[0])
	return replaceToken(t, node)
}

// The type of a binary operator expression is equal to the type of both
// operands. Used for operators like AND, OR, PLUS, MINUS, etc.
func (vis *SemCheckVisitor) typeCheckBinaryOperator(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	left := vis.typeCheck(table, node.Children[0])
	right := vis.typeCheck(table, node.Children[1])
	if !left.EqualsNoPrivacy(right) {
		vis.emitBinaryOperatorTypeMismatchError(node, left, right)
	}
	return replaceToken(left, node)
}

func (vis *SemCheckVisitor) emitBinaryOperatorTypeMismatchError(
	node *token.ASTNode,
	left, right token.Type,
) {
	vis.logTypeCheckError(fmt.Sprintf(""+
		"typecheck: type mismatch for operator %v, "+
		"left operands has type %v and right operand has type %v",
		node.Type, left.Type, right.Type))
}

// Type checking the function call is very similar to type checking the
// variable; we need to check the function call subject and possibly do a lookup
// on his table. If there is no subject for the function call, then we need to
// do a lookup for this function within the current scope
func (vis *SemCheckVisitor) typeCheckFunctionCall(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	id := node.Children[1].Token.Lexeme
	switch subjectType := vis.typeCheck(table, node); subjectType.Type {
	case "":
		return vis.functionLookup(table, table, node, string(id))
	default:
		subject := vis.obtainSubjectRecord(table, subjectType, node)
		if subject == nil {
			return token.Type{}
		}
		return vis.functionLookup(table, subject.Link, node, string(id))
	}
}

func (vis *SemCheckVisitor) typeCheckVariable(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	id := node.Children[1].Token.Lexeme
	switch subjectType := vis.typeCheck(table, node); subjectType.Type {
	case "":
		return vis.variableLookup(table, node, string(id))
	default:
		subject := vis.obtainSubjectRecord(table, subjectType, node)
		if subject == nil {
			return token.Type{}
		}
		return vis.variableLookup(subject.Link, node, string(id))
	}
}

func (vis *SemCheckVisitor) obtainSubjectRecord(
	table token.SymbolTable,
	typee token.Type,
	node *token.ASTNode,
) *token.SymbolTableRecord {
	subId := string(typee.Token.Lexeme)
	sub := token.DeepLookup(table, subId)

	if len(sub) == 0 {
		vis.emitLookupError(node.Children[1])
		return nil
	}

	return sub[0]
}

func (vis *SemCheckVisitor) functionLookup(
	table token.SymbolTable,
	paramsTable token.SymbolTable,
	node *token.ASTNode,
	id string,
) token.Type {
	found := token.DeepLookup(paramsTable, string(id))
	if len(found) == 0 {
		vis.emitLookupError(node.Children[1])
		return token.Type{}
	}

	// Gather the parameters of this function call
	callParams := vis.funcCallParams(table, node)
	funcDefCalled := vis.matchFuncCallWithFuncDef(table, found, callParams)
	if funcDefCalled == nil {
		params := make([]string, 0, len(callParams))
		for _, param := range callParams {
			params = append(params, string(param.Type))
		}
		paramsString := strings.Join(params, ", ")
		vis.logTypeCheckError(fmt.Sprintf(""+
			"typecheck: none of the in-scope definitions for function '%v' "+
			"match function call %v(%v) (line %v)",
			id, id, paramsString, node.Children[1].Token.Line))
		return token.Type{}
	}

	return token.Type{
		Type:    funcDefCalled.Type.Type,
		Token:   node.Children[1].Token,
		Dimlist: funcDefCalled.Type.Dimlist,
		Privacy: funcDefCalled.Type.Privacy,
	}
}

// Attempts to unify a function call with a matching function overload. To
// support function/method overloading, we need to search all
func (vis *SemCheckVisitor) matchFuncCallWithFuncDef(
	table token.SymbolTable,
	funcDefs []*token.SymbolTableRecord,
	callParams []token.Type,
) *token.SymbolTableRecord {
loop:
	for _, funcDef := range funcDefs {
		defParams := vis.funcDefparams(funcDef)
		if len(callParams) != len(defParams) {
			continue // This is not the right function
		}
		for i, defParam := range defParams {
			callParam := callParams[i]
			if !defParam.Type.EqualsNoPrivacy(token.Type{
				Type:    callParam.Type,
				Token:   callParam.Token,
				Dimlist: callParam.Dimlist,
			}) {
				continue loop // This is not the right function
			}
		}
		return funcDef // We've found the function we are looking for
	}
	return nil
}

func (vis *SemCheckVisitor) funcDefparams(
	funcDef *token.SymbolTableRecord,
) []token.SymbolTableRecord {
	entries := funcDef.Link.Entries()
	out := make([]token.SymbolTableRecord, 0, len(entries))
	for _, entry := range entries {
		if entry.Kind == token.FINAL_FUNC_DEF_PARAM {
			out = append(out, entry)
		}
	}
	return out
}

func (vis *SemCheckVisitor) funcCallParams(
	table token.SymbolTable,
	node *token.ASTNode,
) []token.Type {
	paramList := node.Children[2].Children
	out := make([]token.Type, 0, len(paramList))
	for _, param := range paramList {
		out = append(out, vis.typeCheck(table, param.Children[0]))
	}
	return out
}

func (vis *SemCheckVisitor) variableLookup(
	table token.SymbolTable,
	node *token.ASTNode,
	id string,
) token.Type {
	found := token.DeepLookup(table, string(id))
	if l := len(found); l == 0 {
		vis.emitLookupError(node.Children[1])
		return token.Type{}
	}

	rec := found[0]
	subscriptedType := vis.typeCheckDimensions(table, node, rec)
	return replaceToken(subscriptedType, node.Children[1])
}

func (vis *SemCheckVisitor) typeCheckDimensions(
	table token.SymbolTable,
	node *token.ASTNode,
	record *token.SymbolTableRecord,
) token.Type {
	indexList := node.Children[2].Children
	dimList := record.Type.Dimlist

	// Make sure that the dimension list is not too long
	if len(indexList) > len(dimList) {
		var variableToken token.Token
		if node.Type == token.FINAL_VARIABLE {
			variableToken = node.Children[1].Token
		} else {
			variableToken = node.Token
		}

		vis.logTypeCheckError(fmt.Sprintf(""+
			"typecheck: invalid index list for variable %v on line %v, "+
			"variable has %v dimensions but index list has %v subscript(s)",
			record.Name, variableToken.Line, len(dimList), len(indexList)))
		return token.Type{}
	}

	// Gather type of each index
	indexListTypes := make([]token.Type, 0, len(indexList))
	for i, index := range indexList {
		indexType := vis.typeCheck(table, index.Children[0])

		// Indexes should always have integer type, let's check that
		if !isType(indexType, token.FINAL_INTEGER) {
			vis.logTypeCheckError(fmt.Sprintf(""+
				"typecheck: index expressions must be integers, "+
				"index #%v for variable '%v' on line %v has type %v",
				i+1,
				node.Children[1].Token.Lexeme,
				node.Children[1].Token.Line,
				indexType))
		}
		indexListTypes = append(indexListTypes, indexType)
	}

	// We have to determine the type of the variable from how many indexes have
	// been provided. E.g. if the dimList has integer[2][3][4] and the index
	// list is [<expr>][<expr>], then the type of the of the variable is
	// integer[4]

	typeDimList := dimList[len(indexListTypes):]
	return replaceToken(token.Type{
		Type:    record.Type.Type,
		Dimlist: typeDimList,
	}, node)
}

func (vis *SemCheckVisitor) typeCheckAssign(table token.SymbolTable, node *token.ASTNode) {
	lhs := vis.typeCheck(table, node)
	rhs := vis.typeCheck(table, node.Children[1])
	if !lhs.EqualsNoPrivacy(rhs) {
		vis.logTypeCheckError(fmt.Sprintf(""+
			"typecheck: mismatched return type for assignment statement "+
			"in function '%v::%v' line %v "+
			"left-hand side has type %v while right-hand side has type %v",
			table.Parent().Id(), table.Id(), node.Children[0].Token.Line,
			lhs.Type, rhs.Type))
	}
}

func (vis *SemCheckVisitor) typeCheckReturn(table token.SymbolTable, node *token.ASTNode) {
	expectedReturnType := functionReturnType(table)
	actualReturnType := vis.typeCheck(table, node.Children[0])
	if !expectedReturnType.EqualsNoPrivacy(actualReturnType) {
		vis.logTypeCheckError(fmt.Sprintf(
			"typecheck: mismatched return type for '%v::%v', expected %v but found %v",
			table.Parent().Id(), table.Id(), expectedReturnType, actualReturnType))
	}
}

func (vis *SemCheckVisitor) typeCheckSubject(
	table token.SymbolTable,
	node *token.ASTNode,
) token.Type {
	if len(node.Children) > 0 {
		return vis.typeCheck(table, node)
	}
	return token.Type{}
}

func (vis *SemCheckVisitor) emitLookupError(node *token.ASTNode) {
	vis.logTypeCheckError(fmt.Sprintf(""+
		"typecheck: id %v was not found within the current scope (line %v)",
		node.Token.Lexeme, node.Token.Line))
}

func (vis *SemCheckVisitor) logErr(e *VisitorError) {
	if vis.errout != nil {
		vis.errout(e)
	}
}

func (vis *SemCheckVisitor) logTypeCheckError(msg string) {
	vis.logErr(&VisitorError{Wrap: &TypeCheckError{Msg: msg}})
}

func (vis *SemCheckVisitor) typeCheckBlock(table token.SymbolTable, statements []*token.ASTNode) {
	for _, statement := range statements {
		vis.doCheck(table, statement)
	}
}

// Attaches symbol tables to VarDecls with custom struct types
func (vis *SemCheckVisitor) attachVarDecl(node *token.ASTNode) {
	typee := node.Children[1].Children[0]
	if typee.Type != token.FINAL_ID {
		return
	}
	customType := typee.Token.Lexeme
	found := token.DeepLookup(node.Meta.Record.Parent, string(customType))
	if len(found) == 0 {
		vis.logErr(&VisitorError{
			Msg: fmt.Sprintf(""+
				"type '%v' is not valid, no such type has been defined (line %v)",
				customType, typee.Token.Line),
		})
	}

	// Update the record
	id := node.Children[0].Token.Lexeme
	varRecord := token.DeepLookup(node.Meta.Record.Parent, string(id))
	varRecord[0].Link = found[0].Link
	node.Meta.Record.Link = found[0].Link
	// typee.Meta.Record.Link = found[0].Link
}

// Attaches a symbol table to function return types
func (vis *SemCheckVisitor) attachReturnTypeTable(node *token.ASTNode) {
	returnType := node.Children[2].Children[0]
	if returnType.Type != token.FINAL_ID {
		return
	}
	customType := returnType.Token.Lexeme
	found := token.DeepLookup(node.Meta.SymbolTable, string(customType))
	if len(found) == 0 {
		vis.logErr(&VisitorError{
			Msg: fmt.Sprintf(""+
				"type '%v' is not valid, no such type has been defined (line %v)",
				customType, returnType.Token.Line),
		})
	}
	// returnType.Meta.Record.Link =  found[0].Link
	// returnType.Meta.SymbolTable = found[0].Link
}

func functionReturnType(functionTable token.SymbolTable) token.Type {
	// Find the entry for this table in its parent
	for _, entry := range functionTable.Parent().Entries() {
		if entry.Link == functionTable {
			// This is our guy
			return entry.Type
		}
	}
	return token.Type{}
}

func childrenWithoutVarDecls(children []*token.ASTNode) []*token.ASTNode {
	ret := make([]*token.ASTNode, 0, len(children))
	for _, child := range children {
		if child.Type != token.FINAL_VAR_DECL {
			ret = append(ret, child)
		}
	}
	return ret
}

func replaceToken(typee token.Type, node *token.ASTNode) token.Type {
	return token.Type{
		Type:    typee.Type,
		Token:   node.Token,
		Dimlist: typee.Dimlist,
		Privacy: typee.Privacy,
	}
}

// Returns true if the type of the specified node matches one of the provided types
func isTypeNode(node *token.ASTNode, types ...token.Kind) bool {
	for _, t := range types {
		if node.Type == t {
			return true
		}
	}
	return false
}

func isType(typee token.Type, types ...token.Kind) bool {
	for _, t := range types {
		if typee.Type == t {
			return true
		}
	}
	return false
}
