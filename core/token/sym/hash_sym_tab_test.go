package sym

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/obonobo/esac/core/token"
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
