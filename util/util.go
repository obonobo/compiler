package util

import (
	"strings"
)

// Replaces `\n` with `\\n` and `\r` with `\\r`
func SingleLinify(l string) string {
	return strings.ReplaceAll(strings.ReplaceAll(string(l), "\n", "\\n"), "\r", "\\r")
}
