package compositetable

import (
	"fmt"

	"github.com/obonobo/compiler/core/tabledrivenscanner"
)

type UnrecognizedStateError tabledrivenscanner.State

func (u UnrecognizedStateError) Error() string {
	return fmt.Sprintf("unrecognized state '%v'", tabledrivenscanner.State(u))
}
