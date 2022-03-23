package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/obonobo/esac/util"
)

var (
	// Toggle to print all links
	PRINT_ALL_LINKS = false
)

const (
	GLOBAL = "Global"
)

// A SymbolTable lists identifiers that may be referred to within a scope.
//
// WARNING: SymbolTables are recursive and their records may link other
// SymbolTables. Implementations should be careful to avoid cycles. The utility
// functions within this package have been designed with non-cyclic SymbolTables
// in mind. They will break if a cyclic data structure is presented, no checks
// are performed.
type SymbolTable interface {
	// Returns the ID of this table
	Id() string

	// Renames the SymbolTable
	Rename(name string)

	// Adds a record to the SymbolTable
	Insert(record SymbolTableRecord)

	// Prepends a record to the SymbolTable
	Prepend(record ...SymbolTableRecord)

	// Searches for an identifier in the symbol table
	Search(id string) []*SymbolTableRecord

	// Deletes the first entry that match all record fields except
	// SymbolTableRecord.Link
	Delete(like SymbolTableRecord)

	// Deletes all entries matching the provided name
	DeleteAll(id string)

	// Deletes the entry at index i
	DeleteIndex(i int)

	// Iterate over the entries in a SymbolTable, this method may spawn a
	// goroutine
	Entries() []SymbolTableRecord

	// Retrieves the parent symbol table - the symbol table of the enclosing
	// scope, usually this will be a table called "Global". Not to be confused
	// with Inherited().
	Parent() SymbolTable

	// Sets the parent of this table to be the provided table
	SetParent(parent SymbolTable)

	// Returns a list of SymbolTables for structs in the `inherits` list of a
	// struct
	Inherited() []SymbolTable
	ChangeInherited(func(*[]SymbolTable))
}

type Type struct {
	Type    Kind
	Token   Token // Optional
	Dimlist []int // List of
	Privacy Kind  // `public` or `private`. This field is used only on struct members
}

// Fills only the type and token fields, you'll have to fill in the rest
// yourself if you desire those fields
func TypeFromNode(node *ASTNode) Type {
	return Type{
		Type:  node.Type,
		Token: node.Token,
	}
}

func (t Type) String() string {
	return t.StringPrivacy(true)
}

func (t Type) StringPrivacy(includePrivacy bool) string {
	if t.Type == "" {
		return ""
	}
	builder := new(strings.Builder)
	if includePrivacy && t.Privacy != "" {
		fmt.Fprintf(builder, "(%v) ", t.Privacy)
	}
	builder.WriteString(string(t.Token.Lexeme))
	for _, dim := range t.Dimlist {
		fmt.Fprintf(builder, "[%v]", dim)
	}
	return builder.String()

}

func (t Type) dimEquals(t2 Type) bool {
	if len(t.Dimlist) != len(t2.Dimlist) {
		return false
	}
	for i, d := range t.Dimlist {
		d2 := t2.Dimlist[i]
		if d != d2 {
			return false
		}
	}
	return true
}

func (t Type) Equals(t2 Type) bool {
	return t.dimEquals(t2) && t.Type == t2.Type && t.Privacy == t2.Privacy
}

func (t Type) EqualsNoPrivacy(t2 Type) bool {
	return t.dimEquals(t2) && t.Type == t2.Type
}

type SymbolTableRecord struct {
	Name   string // Search key for the record
	Kind   Kind
	Type   Type
	Link   SymbolTable
	Parent SymbolTable
}

func (r SymbolTableRecord) Equal(r2 SymbolTableRecord) bool {
	return r.Name == r2.Name && r.Kind == r2.Kind && r.Type.EqualsNoPrivacy(r2.Type)
}

func (r SymbolTableRecord) String() string {
	link := "nil"
	if r.Link != nil {
		link = r.Link.Id()
	}
	return fmt.Sprintf(
		"SymbolTableRecord[Name=%v, Kind=%v, Type=%v, Link=%v]",
		r.Name, r.Kind, r.Type, link)
}

func (r SymbolTableRecord) ToSimpleJson() map[string]any {
	ret := map[string]any{
		"Name": r.Name,
		"Kind": r.Kind,
		"Type": fmt.Sprintf("%v", r.Type),
	}
	if r.Link != nil {
		ret["Link"] = r.Link.Id()
	}
	return ret
}

func (t SymbolTableRecord) ToSummary() string {
	placeholder := "___"

	var builder strings.Builder
	fmt.Fprintf(&builder, "%v %v %v",
		stror(t.Name, placeholder),
		stror(t.Kind, placeholder),
		stror(fmt.Sprintf("%v", t.Type), placeholder))

	if t.Link != nil {
		fmt.Fprintf(&builder, " ⊙---> SymbolTable[%v]", t.Link.Id())
	}
	return builder.String()
}

func stror[T1, T2 ~string](s T1, or T2) string {
	return util.Or(string(s), string(or))
}

func SymbolTableToJsonMap(table SymbolTable) map[string]any {
	return convertTable(table,
		func(e SymbolTableRecord) map[string]any { return e.ToSimpleJson() },
		func(e SymbolTableRecord) string { return e.Name },
		func(t SymbolTableRecord) map[string]any { return SymbolTableToJsonMap(t.Link) })
}

func SymbolTableToJsonMapBrief(table SymbolTable) map[string]any {
	return convertTable(table,
		func(e SymbolTableRecord) string { return e.ToSummary() },
		func(e SymbolTableRecord) string { return e.ToSummary() },
		func(t SymbolTableRecord) map[string]any { return SymbolTableToJsonMapBrief(t.Link) })
}

func SymbolTableToJson(table SymbolTable) string {
	return printTable(table, func(table SymbolTable) map[string]any {
		return SymbolTableToJsonMap(table)
	})
}

func SymbolTableToJsonBrief(table SymbolTable) string {
	return printTable(table, func(table SymbolTable) map[string]any {
		return SymbolTableToJsonMapBrief(table)
	})
}

func JsonifyAST(node *ASTNode) string {
	return SymbolTableToJson(node.Meta.SymbolTable)
}

func JsonifyASTBrief(node *ASTNode) string {
	return SymbolTableToJsonBrief(node.Meta.SymbolTable)
}

func printTable[T any](table SymbolTable, mapper func(table SymbolTable) T) string {
	data, err := json.MarshalIndent(mapper(table), "", "  ")
	if err != nil {
		return ""
	}

	// Reverse the replacements that json.Marshal does
	return strings.ReplaceAll(string(data), `\u003e`, ">")
}

func convertTable[T, T2 any](
	table SymbolTable,
	bareEntryMapper func(e SymbolTableRecord) T,
	linkKeyEntryMapper func(e SymbolTableRecord) string,
	linkEntryMapper func(e SymbolTableRecord) T2,
) map[string]any {
	m := make(map[string]any, 256)
	for _, e := range table.Entries() {
		if e.Link != nil {
			// Then we have a nested SymbolTable
			m[linkKeyEntryMapper(e)] = linkEntryMapper(e)
		} else {
			// Otherwise, we have an entry with no link; don't recurse
			m[e.Name] = bareEntryMapper(e)
		}
	}
	return m
}

// Pretty prints the table into a string
func PrettySymbolTable(t SymbolTable) string {
	var s bytes.Buffer
	WritePrettySymbolTable(&s, t)
	return s.String()
}

// Pretty prints the symbol table to the provided writer
func WritePrettySymbolTable(w io.Writer, t SymbolTable) {
	// Hardcoded formatting params
	const (
		PREFIX      = "	"
		DEFAULT     = "____"
		MIN_PADDING = 4
	)

	var recurse func(w io.Writer, t SymbolTable, depth int)
	recurse = func(w io.Writer, t SymbolTable, depth int) {
		var prefix strings.Builder
		for i := 0; i < depth; i++ {
			prefix.WriteString(PREFIX)
		}

		entries := t.Entries()

		printouts := make([][4]string, 0, len(entries))
		for _, record := range entries {
			printouts = append(printouts, [4]string{
				util.Or(record.Name, DEFAULT),
				util.Or(string(record.Kind), DEFAULT),
				util.Or(record.Type.String(), DEFAULT),
				util.Or(formatLink(record.Link), DEFAULT),
			})
		}

		// Compute table column widths
		cols := [4]int{4, 4, 4, 4}
		for _, record := range printouts {
			cols[0] = util.Max(cols[0], len(record[0]))
			cols[1] = util.Max(cols[1], len(record[1]))
			cols[2] = util.Max(cols[2], len(record[2]))
			cols[3] = util.Max(cols[3], len(record[3]))
		}

		line := fmt.Sprintf("%v+%v+\n",
			prefix.String(),
			strings.Repeat("-", util.Sum(cols[:]...)+11))

		printLine := func() { fmt.Fprintf(w, "%v", line) }

		centerPad := func(n int, s string) string {
			return strings.Repeat(" ", n/2-len(s)/2-1) + s
		}

		pad := func(n int, s string) string {
			return fmt.Sprintf("%-*v", util.Max(n, MIN_PADDING), s)
		}

		printRow := func(c1, c2, c3, c4 string) {
			fmt.Fprintf(w,
				"%v| %v | %v | %v | %v |\n",
				prefix.String(),
				pad(cols[0], c1),
				pad(cols[1], c2),
				pad(cols[2], c3),
				pad(cols[3], c4))
		}

		// Print the table headers
		fmt.Fprintf(w, "%v%v\n", prefix.String(), centerPad(len(line), t.Id()))
		printLine()
		printRow("Name", "Kind", "Type", "Link")
		printLine()

		// Toggle below if you want to print all
		PRINT_ALL_LINKS := PRINT_ALL_LINKS
		if PRINT_ALL_LINKS {
			nested := make([]SymbolTable, 0, len(entries))
			for i, record := range entries {
				printout := printouts[i]
				if record.Link != nil {
					nested = append(nested, record.Link)
				}
				printRow(printout[0], printout[1], printout[2], printout[3])
			}
			printLine()

			// Recurse on all nested tables
			for _, table := range nested {
				fmt.Fprintln(w)
				recurse(w, table, depth+1)
			}
		} else {
			nested := make([]SymbolTable, 0, len(entries))
			for i, record := range entries {
				printout := printouts[i]
				if record.Link != nil {
					var contains bool
					for _, t := range nested {
						if t.Id() == record.Link.Id() {
							contains = true
							break
						}
					}
					if !contains {
						nested = append(nested, record.Link)
					}
				}
				printRow(printout[0], printout[1], printout[2], printout[3])
			}
			printLine()

			// Recurse on all nested tables
			for _, table := range nested {
				fmt.Fprintln(w)
				recurse(w, table, depth+1)
			}
		}
	}
	recurse(w, t, 0)
}

func formatLink(t SymbolTable) string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("⊙---> %v", t.Id())
}
