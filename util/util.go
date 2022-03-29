package util

import (
	"fmt"
	"io"
	"strings"
)

// Replaces `\n` with `\\n` and `\r` with `\\r`
func SingleLinify(l string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			string(l), "\n", "\\n"),
		"\r", "\\r")
}

// Constrains: any type that you can do '<' or '>' with
type Ordered interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64 |
		~string
}

// Finds the Max amongst some comparable values
func Max[T Ordered](values ...T) T {
	var max T
	for _, val := range values {
		if val > max {
			max = val
		}
	}
	return max
}

func Min[T Ordered](values ...T) T {
	var min T
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}

// Sums all provided values
func Sum[T Ordered](values ...T) T {
	var sum T
	for _, val := range values {
		sum += val
	}
	return sum
}

// Returns the first non-zero value in the list
func Or[T comparable](values ...T) T {
	for _, v := range values {
		if v != *new(T) {
			return v
		}
	}
	return *new(T)
}

func SlicesEqual[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, e1 := range s1 {
		e2 := s2[i]
		if e1 != e2 {
			return false
		}
	}
	return true
}

// A generic map function
func Map[T, R any](slice []T, mapper func(T) R) []R {
	ret := make([]R, 0, len(slice))
	for _, t := range slice {
		ret = append(ret, mapper(t))
	}
	return ret
}

// Deep copies a slice (copies up to capacity)
func Copy[T any](slice []T) []T {
	ret := make([]T, cap(slice))
	copy(ret, slice[:cap(slice)])
	return ret[:len(slice)]
}

func CopySet[T comparable](set map[T]struct{}) map[T]struct{} {
	ret := make(map[T]struct{}, len(set))
	for k := range set {
		ret[k] = struct{}{}
	}
	return ret
}

// Curried logging function - useful for visitor callbacks (the visitors use
// callbacks to log errors/warnings). This function generates callbacks that log
// items to the curried writer.
func Logback[E any](out io.Writer) func(E) {
	return func(e E) { fmt.Fprintln(out, e) }
}

// Curried append function - useful for creating callbacks that collect elements
// emitted from the different parsing phases
func Appendback[E any](slice *[]E) func(E) {
	return func(e E) { *slice = append(*slice, e) }
}
