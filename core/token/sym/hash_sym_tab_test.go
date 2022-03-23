package sym

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/obonobo/esac/core/token"
)

const (
	TARGET    = "TARGET"
	GLOBAL    = "GLOBAL"
	CHILD     = "CHILD"
	INHERITED = "INHERITED"
)

// Some table sizes to test with
var lengths = []int{1, 2, 10, 1 << 5, 1 << 10}

func TestEntries(t *testing.T) {
	t.Parallel()
	for _, tc := range lengths {
		n := tc
		t.Run(fmt.Sprintf("%v", n), func(t *testing.T) {
			t.Parallel()
			expectedRecords, table := createTableAndEntryList(n)
			assertSliceEqual(t, expectedRecords, table.order)
			assertSliceEqual(t, expectedRecords, table.Entries())
		})
	}
}

func TestDeletion(t *testing.T) {
	t.Parallel()

	n := 10
	expectedEntries, table := createTableAndEntryList(n)
	delete := func(i int) {
		rec := expectedEntries[i]
		table.Delete(token.SymbolTableRecord{
			Name: rec.Name,
			Kind: rec.Kind,
			Type: rec.Type,
		})
		expectedEntries = append(expectedEntries[:i], expectedEntries[i+1:]...)
	}

	delete(1)
	delete(3)
	delete(6)

	assertSliceEqual(t, expectedEntries, table.order)
}

func TestDeepSearch_MemberPresentInFirstTable(t *testing.T) {
	assertDeepLookupSearch(t,
		[]token.SymbolTableRecord{{Name: "x", Kind: "Float", Type: token.Type{Type: GLOBAL}}},
		[][]token.SymbolTableRecord{{{Name: "x", Kind: "Float", Type: token.Type{Type: INHERITED}}}},
		[]token.SymbolTableRecord{{Name: "x", Kind: "Float", Type: token.Type{Type: TARGET}}},
		"x")
}

func TestDeepSearch_MemberPresentInParentTable(t *testing.T) {
	assertDeepLookupSearch(t,
		[]token.SymbolTableRecord{{Name: "x", Kind: "Float", Type: token.Type{Type: TARGET}}},
		[][]token.SymbolTableRecord{{{Name: "y", Kind: "Float", Type: token.Type{Type: INHERITED}}}},
		[]token.SymbolTableRecord{{Name: "y", Kind: "Float", Type: token.Type{Type: CHILD}}},
		"x")
}

func TestDeepSearch_MemberPresentInInheritedTable(t *testing.T) {
	assertDeepLookupSearch(t,
		[]token.SymbolTableRecord{{Name: "x", Kind: "Float", Type: token.Type{Type: GLOBAL}}},
		[][]token.SymbolTableRecord{{{Name: "x", Kind: "Float", Type: token.Type{Type: TARGET}}}},
		[]token.SymbolTableRecord{{Name: "y", Kind: "Float", Type: token.Type{Type: CHILD}}},
		"x")
}

func assertDeepLookupSearch(
	t *testing.T,
	global []token.SymbolTableRecord,
	inherited [][]token.SymbolTableRecord,
	me []token.SymbolTableRecord,
	search string,
) {
	var globall token.SymbolTable

	// Create the primary table
	child := NewHashSymTab("Me", nil)
	for _, rec := range me {
		child.Insert(rec)
	}

	// Create global table
	if global != nil {
		globall = NewHashSymTab("Global", nil)
		child.SetParent(globall)
		for _, rec := range global {
			globall.Insert(rec)
		}
	}

	// Create inherited tables
	for i, inherited := range inherited {
		inheritedt := NewHashSymTab(fmt.Sprintf("Inherited-%v", i), globall)
		child.AddInherited(inheritedt)
		for _, rec := range inherited {
			inheritedt.Insert(rec)
		}
	}

	// Perform the search
	got := token.DeepLookup(child, search)
	if len(got) == 0 {
		t.Fatalf("" +
			"DeepLookup() returned no results, " +
			"but there should be a TARGET token returned")
	}

	gott := got[0]
	if expected, actual := TARGET, string(gott.Type.Type); expected != actual {
		t.Fatalf(""+
			"DeepLookup() returned the wrong token. "+
			"Expected token type to be %v but got %v. "+
			"The token that was returned is %v",
			expected, actual, got)
	}

}

func createTableAndEntryList(n int) ([]token.SymbolTableRecord, *HashSymTab) {
	table := NewHashSymTab("Global", nil)
	expectedEntries := make([]token.SymbolTableRecord, 0, 10)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("entry-%v", i)
		record := token.SymbolTableRecord{Name: name}
		expectedEntries = append(expectedEntries, record)
		table.Insert(record)
	}
	return expectedEntries, table
}

func assertSliceEqual[T any](t *testing.T, expected, actual []T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected slice to be %v but got %v", expected, actual)
	}
}
