//
// This package contains code for working with the symbol table
//
package sym

import (
	"github.com/obonobo/esac/core/token"
)

// Implements:
// t *HashSymTab SymbolTable
type HashSymTab struct {
	id        string
	order     []token.SymbolTableRecord
	parent    token.SymbolTable
	inherited []token.SymbolTable
}

// `id` should be a string that uniquely identifies this symbol table, e.g.:
// "Global"
func NewHashSymTab(id string, parent token.SymbolTable) *HashSymTab {
	return &HashSymTab{
		id:        id,
		order:     make([]token.SymbolTableRecord, 0, 256),
		parent:    parent,
		inherited: make([]token.SymbolTable, 0, 4),
	}
}

func (t *HashSymTab) Id() string {
	return t.id
}

func (t *HashSymTab) Rename(name string) {
	t.id = name
}

// Adds a record to the SymbolTable
func (t *HashSymTab) Insert(record token.SymbolTableRecord) {
	t.order = append(t.order, record)
}

func (t *HashSymTab) Prepend(records ...token.SymbolTableRecord) {
	t.order = append(records, t.order...)
}

// Searches for an identifier in the symbol table
func (t *HashSymTab) Search(id string) []*token.SymbolTableRecord {
	entries := make([]*token.SymbolTableRecord, 0, 64)
	for i, r := range t.order {
		if r.Name == id {
			entries = append(entries, &t.order[i])
		}
	}
	return entries
}

func (t *HashSymTab) Delete(like token.SymbolTableRecord) {
	for i := 0; i < len(t.order); i++ {
		rec := t.order[i]
		if rec.Equal(like) {
			t.DeleteIndex(i)
			i--
		}
	}
}

func (t *HashSymTab) DeleteAll(id string) {
	for i := 0; i < len(t.order); i++ {
		rec := t.order[i]
		if rec.Name == id {
			t.DeleteIndex(i)
			i--
		}
	}
}

func (t *HashSymTab) DeleteIndex(i int) {
	if i < 0 || i > len(t.order)-1 {
		return
	}
	switch l := len(t.order); {
	case i == 0:
		t.order = t.order[1:]
	case i == l-1:
		t.order = t.order[:l-1]
	default:
		t.order = append(t.order[:i], t.order[i+1:]...)
	}
}

func (t *HashSymTab) Entries() []token.SymbolTableRecord {
	return t.order
}

func (t *HashSymTab) Parent() token.SymbolTable {
	return t.parent
}

func (t *HashSymTab) SetParent(parent token.SymbolTable) {
	t.parent = parent
}

func (t *HashSymTab) Inherited() []token.SymbolTable {
	return t.inherited
}

func (t *HashSymTab) AddInherited(inherited token.SymbolTable) {
	t.inherited = append(t.inherited, inherited)
}

func (t *HashSymTab) RemoveInherited(name string) {
	for i, inherited := range t.inherited {
		if inherited.Id() == name {
			t.inherited = append(t.inherited[:i], t.inherited[i+1:]...)
		}
	}
}
