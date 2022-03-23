package visitors

import (
	"fmt"
	"strings"

	"github.com/obonobo/esac/core/token"
)

const (
	MALFORMED_TYPE = "malformed type"
)

type VisitorError struct {
	Msg  string
	Wrap error
}

func (e *VisitorError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Wrap.Error()
}

func (e *VisitorError) Unwrap() error {
	return e.Wrap
}

type DuplicateIdentifierError struct {
	Name   string
	First  token.Token
	Second token.Token
	Wrap   error
}

func (e *DuplicateIdentifierError) Error() string {
	first, second := e.First, e.Second
	if e.Second.Line < e.First.Line {
		first, second = e.Second, e.First
	}

	name := e.Name
	if name == "" {
		name = string(first.Lexeme)
	}

	return fmt.Sprintf(
		"duplicate definition for '%v' (defined on line %v, and again on line %v)",
		name, first.Line, second.Line)
}

func (e *DuplicateIdentifierError) Unwrap() error {
	return e.Wrap
}

type MethodMismatchError struct {
	Method       string
	Struct       *token.ASTNode
	StructMethod token.SymbolTableRecord
	ImplMethod   token.SymbolTableRecord
	Wrap         error
}

func (e *MethodMismatchError) Error() string {
	return fmt.Sprintf(""+
		"method definition mismatch for struct '%v', "+
		"defined in struct as '%v', "+
		"defined in impl as '%v'",
		e.Struct.Token.Lexeme,
		formatMethodDefinition(e.StructMethod),
		formatMethodDefinition(e.ImplMethod))
}

func formatMethodDefinition(record token.SymbolTableRecord) string {
	out := new(strings.Builder)
	if record.Type.Privacy != "" {
		fmt.Fprintf(out, "%v ", record.Type.Privacy)
	}

	fmt.Fprintf(out, "func %v -> %v",
		formatMethodId(record),
		record.Type.StringPrivacy(false))

	return out.String()
}

func formatMethodId(record token.SymbolTableRecord) string {
	params := formatParams(record)
	out := new(strings.Builder)
	fmt.Fprintf(out, "%v(%v)", record.Name, params)
	return out.String()
}

func formatParams(record token.SymbolTableRecord) string {
	var params []string
	if record.Link != nil {
		entries := record.Link.Entries()
		params = make([]string, 0, len(entries))
		for _, e := range entries {
			if e.Kind == token.FINAL_FUNC_DEF_PARAM {
				params = append(params, e.Type.String())
			}
		}
	}
	return strings.Join(params, ", ")
}

type StructMissingMethodFromImplError struct {
	Node   *token.ASTNode
	Method token.SymbolTableRecord
	Wrap   error
}

func (e *StructMissingMethodFromImplError) Error() string {
	return fmt.Sprintf(""+
		"struct '%v' is missing method '%v' defined in impl (line %v)",
		safeId(e.Node),
		formatMethodDefinition(e.Method),
		e.Method.Type.Token.Line)
}

func (e *StructMissingMethodFromImplError) Unwrap() error {
	return e.Wrap
}

type ImplMissingMethodFromStructError struct {
	Node   *token.ASTNode
	Method token.SymbolTableRecord
	Wrap   error
}

func (e *ImplMissingMethodFromStructError) new(
	node *token.ASTNode,
	method token.SymbolTableRecord,
) error {
	return &ImplMissingMethodFromStructError{
		Node:   node,
		Method: method,
	}
}

func (e *ImplMissingMethodFromStructError) Error() string {
	return fmt.Sprintf(""+
		"impl '%v' is missing method '%v' defined in struct (line %v)",
		safeId(e.Node),
		formatMethodDefinition(e.Method),
		e.Method.Type.Token.Line)
}

func (e *ImplMissingMethodFromStructError) Unwrap() error {
	return e.Wrap
}

type ImplMayOnlyContainFuncDefsError struct {
	Member string
	Impl   *token.ASTNode
	Wrap   error
}

func (e *ImplMayOnlyContainFuncDefsError) Error() string {
	return fmt.Sprintf(""+
		"impl '%v': member '%v' is not valid, impls may only contain function definitions",
		e.Impl.Meta.SymbolTable.Id(), e.Member)
}

func (e *ImplMayOnlyContainFuncDefsError) Unwrap() error {
	return e.Wrap
}

type StructMissingImplError struct {
	Struct *token.ASTNode
	Wrap   error
}

func (e *StructMissingImplError) Error() string {
	return fmt.Sprintf(""+
		"%v: no impl found for struct '%v', "+
		"struct methods declared but not defined",
		MALFORMED_TYPE, e.Struct.Meta.Record.Name)
}

type ImplMissingStructError struct {
	Impl *token.ASTNode
	Wrap error
}

func (e *ImplMissingStructError) Error() string {
	return fmt.Sprintf(""+
		"%v: no struct found for impl '%v', "+
		"impl methods must first be declared in a struct",
		MALFORMED_TYPE, e.Impl.Meta.Record.Name)
}

// Warnings are emitted as errors, but they are distinguishable by this common
// wrapper type. Use `errors.As` to check if an error emitted by the visitor is
// really a warning
type Warning struct {
	Msg  string
	Wrap error
}

func (e *Warning) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Wrap.Error()
}

func (e *Warning) Unwrap() error {
	return e.Wrap
}
