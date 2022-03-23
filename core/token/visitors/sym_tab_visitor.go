package visitors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/core/token/sym"
)

// A version of the symbol table that knows which node it is attached to
type NodeAwareSymbolTable struct {
	token.SymbolTable
	node *token.ASTNode
}

// Structs must be linked with Impls, if the struct defines at least one method.
//
// The symbol table of a struct is split across two constructs. The `struct`
// definition defines data members as well as methods, but we cannot link the
// method symbol table until we see an `impl`.
//
// How this visitor works is that it processes these two portions of a struct in
// either order and creates the StructTable type below. Once both portions have
// been processed, the StructTable is considered complete. If we reach the
// 'prog' semantic action and discover an incomplete StructTable, then an error
// is emitted from the visitor indicating a partial struct definition, which is
// not allowed.
//
// Further, when an incomplete StructTable is being completed, by either the
// `impl` or `struct` semantic action, that semantic action must match
// one-to-one the newly defined methods with those already defined in the
// incomplete symbol table. If there is a missing or extra method, the visitor
// must emit an error; this is a malformed struct definition.
//
// But where to place a completed StructTable?
//
// The visitor annotates both the `impl` and the `struct` ASTNodes with a
// pointer to the same completed StructTable for the struct jointly defined by
// these constructs.
//
// If a semantic action discovers an already complete StructTable, this
// indicates a multiply defined `struct` or `impl` for this type, which is
// treated as an error emitted from the visitor in this implementation.
type StructTable struct {
	*NodeAwareSymbolTable
	complete bool
}

type key struct {
	parent string
	me     string
}

type SymTabVisitor struct {
	token.DispatchVisitor

	// Keep a running tally of tables constructed. Program symbol tables are
	// uniquely identified by a composite key of {parent-table, table}
	tables map[key]token.SymbolTable

	errout func(e *VisitorError)
}

// TODO: add inheritance
// TODO: shadowed methods, shadowed variables
func NewSymTabVisitor(errout func(e *VisitorError)) *SymTabVisitor {
	vis := &SymTabVisitor{
		tables: make(map[key]token.SymbolTable, 64),
		errout: errout,
	}
	vis.tables[key{"", token.GLOBAL}] = newSymbolTable(token.GLOBAL, nil, nil)

	vis.DispatchVisitor = token.DispatchVisitor{Dispatch: map[token.Kind]token.Visit{
		token.FINAL_PROG: func(node *token.ASTNode) {
			table := vis.tables[key{"", token.GLOBAL}]
			table.(*NodeAwareSymbolTable).node = node
			node.Meta.SymbolTable = table
			addChildren(vis, node, node.Children[0].Children)
			vis.verifyStructTables()

			// Emit warnings for all overloaded methods in the table
			warnOverloads(vis, node.Meta.SymbolTable)
		},

		token.FINAL_STRUCT_DECL: func(node *token.ASTNode) {
			id := id(node)
			partial, ok := vis.tables[key{token.GLOBAL, id}]
			if ok {
				// Then we need to try to complete the table
				completeSymbolTableStruct(vis, id, node, partial)
			} else {
				// We must create a fresh table
				createFreshTableStruct(vis, id, node)
			}
		},

		// This method works almost the same as the 'struct' version except that
		// it must complete the 'impl' portion of the table
		token.FINAL_IMPL_DEF: func(node *token.ASTNode) {
			id := id(node)
			partial, ok := vis.tables[key{token.GLOBAL, id}]
			if ok {
				// Then we need to try to complete the table
				completeSymbolTableImpl(vis, id, node, partial)
			} else {
				// We must create a fresh table
				createFreshTableImpl(vis, id, node)
			}
		},

		// Members are either FuncDecl or VarDecl but with an extra privacy
		// modifier attached to them. This action just adds the privacy
		// annotation to the record
		token.FINAL_MEMBER: func(node *token.ASTNode) {
			node.Children[1].Meta.Record.Type.Privacy = node.Children[0].Token.Id
			node.Meta.Record = node.Children[1].Meta.Record
		},

		token.FINAL_FUNC_DECL: func(node *token.ASTNode) {
			vis.parseFuncHead(node, token.FINAL_FUNC_DECL)
		},

		token.FINAL_FUNC_DEF: func(node *token.ASTNode) {
			vis.parseFuncHead(node, token.FINAL_FUNC_DEF)
			addChildren(vis, node, node.Children[3].Children)
		},

		token.FINAL_FUNC_DEF_PARAM: func(node *token.ASTNode) {
			vis.parseVarRecord(node, token.FINAL_FUNC_DEF_PARAM)
		},

		token.FINAL_VAR_DECL: func(node *token.ASTNode) {
			vis.parseVarRecord(node, token.FINAL_VAR_DECL)
		},
	}}
	return vis
}

// Checks the StructTables of a program
func (v *SymTabVisitor) verifyStructTables() {
	for _, t := range v.tables {
		if tt, ok := t.(*StructTable); ok {
			if !tt.complete {
			loop:
				for _, rec := range tt.Entries() {
					switch kind := rec.Kind; kind {
					case token.FINAL_FUNC_DECL:
						// Then we have a defined struct, but no impl
						v.logErr(&VisitorError{Wrap: &StructMissingImplError{Struct: tt.node}})
						break loop
					case token.FINAL_FUNC_DEF:
						// Then we have a defined impl, but no struct
						v.logErr(&VisitorError{Wrap: &ImplMissingStructError{Impl: tt.node}})
						break loop
					}
				}
			}
		}
	}
}

func (v *SymTabVisitor) logErr(e *VisitorError) {
	if v.errout != nil {
		v.errout(e)
	}
}

func setParent(member, node *token.ASTNode) {
	member.Meta.Record.Parent = node.Meta.SymbolTable
	if childTable := member.Meta.Record.Link; childTable != nil {
		childTable.SetParent(node.Meta.SymbolTable)
	}
}

func methodsMatch(rec1, rec2 token.SymbolTableRecord) bool {
	if len(rec1.Type.Dimlist) != len(rec2.Type.Dimlist) {
		return false
	}
	for i, d1 := range rec1.Type.Dimlist {
		d2 := rec2.Type.Dimlist[i]
		if d1 != d2 {
			return false
		}
	}
	return (rec1.Kind == token.FINAL_FUNC_DECL || rec1.Kind == token.FINAL_FUNC_DEF) &&
		(rec2.Kind == token.FINAL_FUNC_DECL || rec2.Kind == token.FINAL_FUNC_DEF) &&
		rec1.Name == rec1.Name &&
		rec1.Type.Type == rec2.Type.Type
}

func isFuncMember(record token.SymbolTableRecord) bool {
	switch record.Kind {
	case token.FINAL_FUNC_DEF:
		return true
	}
	return false
}

func isDataMember(record token.SymbolTableRecord) bool {
	switch record.Kind {
	case token.FINAL_VAR_DECL:
		return true
	}
	return false
}

func newSymbolTable(
	id string,
	node *token.ASTNode,
	parent token.SymbolTable,
) *NodeAwareSymbolTable {
	return &NodeAwareSymbolTable{
		SymbolTable: sym.NewHashSymTab(id, parent),
		node:        node,
	}
}

func newSymbolTableNoParent(id string, node *token.ASTNode) *NodeAwareSymbolTable {
	return &NodeAwareSymbolTable{
		SymbolTable: sym.NewHashSymTab(id, nil),
		node:        node,
	}
}

// Returns true if the symbol table defines a partial StructTable containing
// only the 'struct' portion of the definition
func isStructTable(table token.SymbolTable) bool {
	// We determine if this is a partial 'struct' StructTable by checking what
	// the method entries look like
	for _, e := range table.Entries() {
		if e.Type.Type == token.FINAL_FUNC_DECL {
			return true
		} else if e.Type.Type == token.FINAL_FUNC_DEF {
			return false // This the 'impl' version of the partial StructTable
		}
	}
	return false
}

// Returns true if the symbol table defines a partial StructTable containing
// only the 'impl' portion of the definition
func isImplTable(table token.SymbolTable) bool {
	for _, e := range table.Entries() {
		if e.Type.Type == token.FINAL_FUNC_DECL {
			return false // This the 'struct' version of the partial StructTable
		} else if e.Type.Type == token.FINAL_FUNC_DEF {
			return true
		}
	}
	return false
}

func (vis *SymTabVisitor) parseFuncHead(node *token.ASTNode, kind token.Kind) {
	id := id(node)
	node.Meta.SymbolTable = newSymbolTableNoParent(id, node)
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: kind,
		// Type: token.TypeFromNode(idNode(node)),
		Type: token.TypeFromNode(node.Children[2].Children[0]),
		Link: node.Meta.SymbolTable,
	}

	// Fill out params
	for _, child := range node.Children[1].Children {
		node.Meta.SymbolTable.Insert(*child.Meta.Record)
	}

	id = formatMethodId(*node.Meta.Record)
	node.Meta.SymbolTable.Rename(id)
	// node.Meta.Record.Name = id
}

func (t *SymTabVisitor) parseVarRecord(node *token.ASTNode, kind token.Kind) {
	id := id(node)
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: kind,
		Type: token.TypeFromNode(node.Children[1].Children[0]),
	}

	// Fill out param dimensions
	for _, child := range node.Children[2].Children {
		if dim := parseDim(child); dim != -1 {
			node.Meta.Record.Type.Dimlist = append(node.Meta.Record.Type.Dimlist, dim)
		}
	}
}

// Parses an int dimension from an ASTNode
func parseDim(node *token.ASTNode) int {
	if node == nil {
		return -1
	}
	ret, err := strconv.Atoi(string(node.Token.Lexeme))
	if err != nil {
		panic(fmt.Errorf("ParseDim(): %w", err))
	}
	return ret
}

func id(node *token.ASTNode) string {
	return string(idNode(node).Token.Lexeme)
}

func safeId(node *token.ASTNode) string {
	if node != nil && len(node.Children) >= 1 {
		return id(node)
	}
	return ""
}

func idNode(node *token.ASTNode) *token.ASTNode {
	return node.Children[0]
}

func lengthOrPanic(children []*token.ASTNode, desiredLength int) {
	l := len(children)
	if l < desiredLength {
		panic(fmt.Errorf(
			"expected stack to be %v elements, but was %v. Here is the stack: %v",
			desiredLength, l, children))
	}
}

func searchRecordByIdAndKindAndType(
	table token.SymbolTable,
	id string,
	kind token.Kind,
	typee token.Type,
) *token.SymbolTableRecord {
	records := table.Search(id)
	for _, r := range records {
		if r.Kind == kind && r.Type.Equals(typee) {
			return r
		}
	}
	return nil
}

// Check that the method definitions match
func checkMismatchedMethods(
	vis *SymTabVisitor,
	structMember, implMember token.SymbolTableRecord,
	node *token.ASTNode,
) bool {
	if !methodsMatch(structMember, implMember) {
		vis.logErr(&VisitorError{Wrap: &MethodMismatchError{
			Method:       structMember.Name,
			Struct:       node,
			ImplMethod:   implMember,
			StructMethod: structMember,
		}})
		return true
	}
	return false
}

// Searches for a method with specified params
func searchForMethod(
	table token.SymbolTable,
	name string,
	typee token.Type,
	params []token.SymbolTableRecord,
) *token.SymbolTableRecord {
	records := table.Search(name)

loop:
	for _, r := range records {
		if typee.EqualsNoPrivacy(r.Type) {
			// Check the params
			if r.Link != nil {
				pars := getParams(*r)
				if len(pars) != len(params) {
					continue
				}

				for i, p1 := range pars {
					p2 := params[i]
					if p2.Name != p1.Name || !p2.Type.EqualsNoPrivacy(p1.Type) {
						continue loop
					}
				}
				return r
			}
		}
	}
	return nil
}

func getParams(record token.SymbolTableRecord) []token.SymbolTableRecord {
	pars := make([]token.SymbolTableRecord, 0, 16)
	for _, p := range record.Link.Entries() {
		if p.Kind == token.FINAL_FUNC_DEF_PARAM {
			pars = append(pars, p)
		}
	}
	return pars
}

func checkPartialStructTable(
	vis *SymTabVisitor,
	partial token.SymbolTable,
	node *token.ASTNode,
	isStruct bool,
) *StructTable {
	partialStructTable, ok := partial.(*StructTable)

	var check bool
	if isStruct {
		check = isStructTable(partial)
	} else {
		check = isImplTable(partial)
	}

	if !ok || check || partialStructTable.complete {
		// Then this is a duplicate identifier
		vis.logErr(&VisitorError{Wrap: &DuplicateIdentifierError{
			First:  idNode(partialStructTable.NodeAwareSymbolTable.node).Token,
			Second: idNode(node).Token,
		}})
		return nil
	}
	return partialStructTable
}

func checkFunctionMember(
	vis *SymTabVisitor,
	implMember token.SymbolTableRecord,
	node *token.ASTNode,
) bool {
	// Error if member is not a function member
	if !isFuncMember(implMember) {
		vis.logErr(&VisitorError{Wrap: &ImplMayOnlyContainFuncDefsError{
			Member: implMember.Name,
			Impl:   node,
		}})
		return true
	}
	return false
}

func checkStructMember(
	vis *SymTabVisitor,
	structMember *token.SymbolTableRecord,
	implMember token.SymbolTableRecord,
	node *token.ASTNode,
) bool {
	if structMember == nil {
		vis.logErr(&VisitorError{Wrap: &StructMissingMethodFromImplError{
			Method: implMember,
			Node:   node,
		}})
		return true
	}
	return false
}

func checkImplMember(
	vis *SymTabVisitor,
	implMember *token.SymbolTableRecord,
	structMember token.SymbolTableRecord,
	node *token.ASTNode,
) bool {
	if implMember == nil {
		vis.logErr(&VisitorError{Wrap: &ImplMissingMethodFromStructError{
			Node:   node,
			Method: structMember,
		}})
		return true
	}
	return false
}

// Completes the `struct` portion of a symbol table
func completeSymbolTableStruct(
	vis *SymTabVisitor,
	id string,
	node *token.ASTNode,
	partial token.SymbolTable,
) {
	partialStructTable := checkPartialStructTable(vis, partial, node, true)
	if partialStructTable == nil {
		return
	}

	dataMembers := make(
		[]token.SymbolTableRecord, 0,
		len(node.Children[1].Children))

	// Collect definitions from this node
	structMethods := make(map[string]token.SymbolTableRecord, 64)
	implDefitions := methods(partialStructTable)
	for _, member := range node.Children[2].Children {
		structMember := *member.Meta.Record

		structKey := structMember.String()
		if found, ok := structMethods[structKey]; ok {
			vis.logErr(&VisitorError{Wrap: &DuplicateIdentifierError{
				Name:   found.Name,
				First:  found.Type.Token,
				Second: structMember.Type.Token,
			}})
			continue
		}

		// If this is a data member, then add it to the table
		if isDataMember(structMember) {
			dataMembers = append(dataMembers, structMember)
			continue
		}

		implMember := searchForMethod(
			partialStructTable,
			structMember.Name,
			structMember.Type,
			getParams(structMember))

		if checkImplMember(vis, implMember, structMember, partialStructTable.node) {
			continue
		}

		if checkMismatchedMethods(vis, structMember, *implMember, partialStructTable.node) {
			continue
		}

		// Mark this method off the list
		implKey := implMember.String()
		delete(implDefitions, implKey)
		structMethods[structKey] = structMember

		// Complete the record by adding the privacy modifier and changing the
		// node that is pointed to (for better error messages)
		implMember.Type.Privacy = member.Meta.Record.Type.Privacy
		partialStructTable.node = node
	}

	// We also need to check that all members from the impl have been accounted for
	logMissingMethods[StructMissingMethodFromImplError](vis, node, implDefitions)

	// Prepend the data members
	partialStructTable.Prepend(dataMembers...)

	// Mark the table as being complete
	partialStructTable.complete = true

	// Make a record for this node
	node.Meta.SymbolTable = partialStructTable
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: token.FINAL_STRUCT_DECL,
		Type: token.Type{Token: idNode(node).Token},
		Link: node.Meta.SymbolTable,
	}

	// Emit warnings for all overloaded methods in the table
	warnOverloads(vis, partialStructTable)
}

func createFreshTableStruct(vis *SymTabVisitor, id string, node *token.ASTNode) {
	node.Meta.SymbolTable = &StructTable{
		NodeAwareSymbolTable: newSymbolTable(id, node, vis.tables[key{"", token.GLOBAL}]),
	}
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: token.FINAL_STRUCT_DECL,
		Type: token.Type{Token: idNode(node).Token},
		Link: node.Meta.SymbolTable,
	}
	addChildren(vis, node, node.Children[2].Children)
	vis.tables[key{token.GLOBAL, id}] = node.Meta.SymbolTable
}

func completeSymbolTableImpl(
	vis *SymTabVisitor,
	id string,
	node *token.ASTNode,
	partial token.SymbolTable,
) {
	partialStructTable := checkPartialStructTable(vis, partial, node, false)
	if partialStructTable == nil {
		return
	}

	// Collect method definitions from this node
	implMethods := make(map[string]token.SymbolTableRecord, 64)
	structMethods := methods(partialStructTable)
	for _, member := range node.Children[1].Children {
		implMember := *member.Meta.Record // Panic if nil

		implKey := implMember.String()
		if found, ok := implMethods[implKey]; ok {
			vis.logErr(&VisitorError{Wrap: &DuplicateIdentifierError{
				Name:   found.Name,
				First:  found.Type.Token,
				Second: implMember.Type.Token,
			}})
			continue
		}

		if checkFunctionMember(vis, implMember, node) {
			continue
		}

		// Check that this method has been defined in the struct
		structMember := searchForMethod(
			partialStructTable,
			implMember.Name,
			implMember.Type,
			getParams(implMember))

		if checkStructMember(vis, structMember, implMember, partialStructTable.node) {
			continue
		}

		if checkMismatchedMethods(vis, *structMember, implMember, partialStructTable.node) {
			continue
		}

		structKey := structMember.String()
		delete(structMethods, structKey)
		implMethods[implKey] = implMember

		// Complete the record with new information from the impl def
		structMember.Kind = implMember.Kind
		structMember.Link = implMember.Link
	}

	logMissingMethods[ImplMissingMethodFromStructError](vis, node, structMethods)

	// Mark the table as being complete
	partialStructTable.complete = true

	// Make a record for this node
	node.Meta.SymbolTable = partialStructTable
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: token.FINAL_IMPL_DEF,
		Type: token.Type{Token: idNode(node).Token},
		Link: node.Meta.SymbolTable,
	}

	// Emit warnings for all overloaded methods in the table
	warnOverloads(vis, partialStructTable)
}

func createFreshTableImpl(vis *SymTabVisitor, id string, node *token.ASTNode) {
	node.Meta.SymbolTable = &StructTable{
		NodeAwareSymbolTable: newSymbolTable(id, node, vis.tables[key{"", token.GLOBAL}]),
	}
	node.Meta.Record = &token.SymbolTableRecord{
		Name: id,
		Kind: token.FINAL_IMPL_DEF,
		Type: token.Type{Token: idNode(node).Token},
		Link: node.Meta.SymbolTable,
	}

	addChildren(vis, node, node.Children[1].Children)
	vis.tables[key{token.GLOBAL, id}] = node.Meta.SymbolTable
}

// WARNING: generics experiment below - Frankenstein's generics
//
// This is generic function that can instantiate one of two error types.
func logMissingMethods[
	ErrorType ImplMissingMethodFromStructError | StructMissingMethodFromImplError,
	ErrorPointer interface {
		*ErrorType
		error
	},
](
	vis *SymTabVisitor,
	node *token.ASTNode,
	methods map[string]token.SymbolTableRecord,
) {
	for _, missing := range methods {
		vis.logErr(&VisitorError{Wrap: ErrorPointer(&ErrorType{
			Node:   node,
			Method: missing,
		})})
	}
}

func methods(table token.SymbolTable) map[string]token.SymbolTableRecord {
	methods := make(map[string]token.SymbolTableRecord, 64)
	for _, r := range table.Entries() {
		if r.Kind == token.FINAL_FUNC_DEF || r.Kind == token.FINAL_FUNC_DECL {
			methods[r.String()] = r
		}
	}
	return methods
}

func alreadyExists(
	table token.SymbolTable,
	record token.SymbolTableRecord,
) *token.SymbolTableRecord {
	for _, r := range table.Search(record.Name) {
		switch {
		case r.Kind != record.Kind:
			continue
		case r.Link == nil && record.Link == nil:
			fallthrough
		case r.Link != nil && record.Link != nil && r.Link.Id() == record.Link.Id():
			return r
		default:
			continue
		}
	}
	return nil
}

// Check if record already exists
func checkAlreadyExists(
	vis *SymTabVisitor,
	node *token.ASTNode,
	add token.SymbolTableRecord,
) bool {
	if r := alreadyExists(node.Meta.SymbolTable, add); r != nil {
		vis.logErr(&VisitorError{Wrap: &DuplicateIdentifierError{
			Name:   add.Name,
			First:  r.Type.Token,
			Second: add.Type.Token,
		}})
		return true
	}
	return false
}

func addChild(vis *SymTabVisitor, node, member *token.ASTNode) {
	add := *member.Meta.Record
	if checkAlreadyExists(vis, node, add) {
		return
	}
	setParent(member, node)
	node.Meta.SymbolTable.Insert(*member.Meta.Record)
}

func addChildren(vis *SymTabVisitor, node *token.ASTNode, children []*token.ASTNode) {
	for _, member := range children {
		if member.Meta.Record != nil {
			addChild(vis, node, member)
		}
	}
}

// Emit warnings for all overloaded methods in the table
func warnOverloads(vis *SymTabVisitor, table token.SymbolTable) {
	overloads := make(map[string][]token.SymbolTableRecord, 32)
	for _, entry := range table.Entries() {
		if entry.Kind == token.FINAL_FUNC_DEF || entry.Kind == token.FINAL_FUNC_DECL {
			overloads[entry.Name] = append(overloads[entry.Name], entry)
		}
	}
	for k, v := range overloads {
		l := len(v)
		if l < 2 {
			continue
		}
		out := make([]string, 0, l)
		for _, overload := range v {
			out = append(out, formatMethodId(overload))
		}
		outt := strings.Join(out, ", ")
		vis.logErr(&VisitorError{Wrap: &Warning{
			Msg: fmt.Sprintf(
				"'%v::%v' has been overloaded %v times: %v",
				table.Id(), k, l, outt),
		}})
	}
}
