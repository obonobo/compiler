package compositetable

import (
	"fmt"

	"github.com/obonobo/esac/core/tabledrivenscanner"
)

type UnrecognizedStateError tabledrivenscanner.State

func (u UnrecognizedStateError) Error() string {
	return fmt.Sprintf("unrecognized state '%v'", tabledrivenscanner.State(u))
}
