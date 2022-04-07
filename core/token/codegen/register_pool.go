package codegen

import (
	"strings"

	"github.com/obonobo/esac/util"
)

const (
	R0  = "r0"  // Register 0 - reserved as zero value
	R1  = "r1"  // Register 1
	R2  = "r2"  // Register 2
	R3  = "r3"  // Register 3
	R4  = "r4"  // Register 4
	R5  = "r5"  // Register 5
	R6  = "r6"  // Register 6
	R7  = "r7"  // Register 7
	R8  = "r8"  // Register 8
	R9  = "r9"  // Register 9
	R10 = "r10" // Register 10
	R11 = "r11" // Register 11
	R12 = "r12" // Register 12
	R13 = "r13" // Register 13 - reserved for procedure return values
	R14 = "r14" // Register 14 - reserved for stack bottom
	R15 = "r15" // Register 15 - reserved for jump link

	TOPADDR = "topaddr" // `topaddr` MOON intrinsic, points to the start of the stack
)

// Starting register set
var registers = []string{
	// R0, // reserved

	R1, R2, R3,
	R4, R5, R6,
	R7, R8, R9,
	R10, R11, R12,

	// R13, R14, R15, // reserved
}

// A simple pooling machine for keeping track of MOON registers in active use
// during a program
type RegisterPool struct {
	registers []string // Free registers
	active    []string // Used registers
}

func NewRegisterPool() *RegisterPool {
	return &RegisterPool{
		registers: util.Copy(registers),
		active:    make([]string, 0, len(registers)),
	}
}

// Claims any available register
func (p *RegisterPool) ClaimAny() string {
	l := len(p.registers)
	if l == 0 {
		return ""
	}
	return p.Claim(p.registers[l-1])
}

// Attempts to claim a specific register, returns false if register is not
// available, otherwise returns true if the register was successfully claimed
func (p *RegisterPool) Claim(id string) string {
	if util.MoveFromSliceToSlice(&p.registers, &p.active, id) {
		return id
	}
	return ""
}

// Attemps to free a register returning it to the available register pool.
// Returns true if successful, false otherwise.
func (p *RegisterPool) Free(ids ...string) bool {
	var result bool
	if len(ids) == 0 {
		return true
	}
	result = util.MoveFromSliceToSlice(&p.active, &p.registers, ids[0])
	for _, id := range ids[1:] {
		result = result && util.MoveFromSliceToSlice(&p.active, &p.registers, id)
	}
	return result
}

func (p *RegisterPool) FreeAll() bool {
	return p.Free(p.active...)
}

// Checks if the specified register is available to be claimed
func (p *RegisterPool) IsAvailable(id string) bool {
	return util.Contains(p.registers, strings.ToLower(id))
}

func (p *RegisterPool) ClaimN(n int) []string {
	claimed := make([]string, 0, n)
	for range claimed[:n] {
		claimed = append(claimed, p.ClaimAny())
	}
	return claimed
}

func (p *RegisterPool) Claim6() (string, string, string, string, string, string) {
	claimed := p.ClaimN(6)
	return claimed[0], claimed[1], claimed[2], claimed[3], claimed[4], claimed[5]
}

func (p *RegisterPool) Claim3() (string, string, string) {
	claimed := p.ClaimN(3)
	return claimed[0], claimed[1], claimed[2]
}
