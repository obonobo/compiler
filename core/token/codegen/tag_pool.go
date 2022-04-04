package codegen

import "fmt"

// Used for keeping track of tags used in expressions. Package-private because
// tags are only used in this particular visitor implementation.
type tagPool struct {
	tempc  int      // Count for temporary tags
	active []string // Stack of tags that are in active usage in the program
}

func newTagPool() *tagPool {
	return &tagPool{active: make([]string, 0, 1024)}
}

// Returns a new temporary tag from the pool
func (p *tagPool) temp() string {
	next := fmt.Sprintf("t%v", p.tempc)
	p.tempc++
	p.active = append(p.active, next)
	return next
}

func (p *tagPool) top() string {
	if len(p.active) == 0 {
		return ""
	}
	return p.active[len(p.active)-1]
}

func (p *tagPool) pop() string {
	got := p.popn(1)
	if len(got) == 0 {
		return ""
	}
	return got[0]
}

func (p *tagPool) pop2() (string, string) {
	popped := p.popn(2)
	if len(popped) == 0 {
		return "", ""
	}
	return popped[0], popped[1]
}

func (p *tagPool) popn(n int) []string {
	l := len(p.active)
	if n <= 0 || l < n {
		return []string{}
	}
	to := len(p.active) - n
	ret := p.active[to:]
	p.active = p.active[:to]
	return ret
}
