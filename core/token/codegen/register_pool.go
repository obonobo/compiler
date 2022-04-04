package codegen

import (
	"strings"

	"github.com/obonobo/esac/util"
)

var (
	r0 = "r0" // Not part of the regular register set, should always contain 0

	// Starting register set
	moonRegisters = []string{
		"r1", "r2", "r3",
		"r4", "r5", "r6",
		"r7", "r8", "r9",
		"r10", "r11", "r12",
		"r13", "r14", "r15",
	}
)

// A simple pooling machine for keeping track of MOON registers in active use
// during a program
type RegisterPool struct {
	registers []string // Free registers
	active    []string // Used registers
}

func NewRegisterPool() *RegisterPool {
	return &RegisterPool{registers: moonRegisters}
}

// Attempts to claim a specific register, returns false if register is not
// available, otherwise returns true if the register was successfully claimed
func (p *RegisterPool) Claim(id string) bool {
	return util.MoveFromSliceToSlice(&p.registers, &p.active, id)
}

// Attemps to free a register returning it to the available register pool.
// Returns true if successful, false otherwise.
func (p *RegisterPool) Free(id string) bool {
	return util.MoveFromSliceToSlice(&p.active, &p.registers, id)
}

// Checks if the specified register is available to be claimed
func (p *RegisterPool) IsAvailable(id string) bool {
	return util.Contains(p.registers, strings.ToLower(id))
}
