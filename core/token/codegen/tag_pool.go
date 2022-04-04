package codegen

import "fmt"

// Used for keeping track of tags used in expressions. Package-private because
// tags are only used in this particular visitor implementation.
type tagPool struct {
	tempc  int      // Count for temporary tags
	active []string // Stack of tags that are in active usage in the program
}

func newTagPool() *tagPool {
	return &tagPool{}
}

// Returns a new temporary tag from the pool
func (p *tagPool) temp() string {
	next := fmt.Sprintf("%v", p.tempc)
	p.tempc++
	p.active = append(p.active, next)
	return next
}
