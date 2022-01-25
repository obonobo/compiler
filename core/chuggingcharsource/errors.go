package chuggingcharsource

import "fmt"

type ChuggingError struct{ Err error }

func (e *ChuggingError) Error() string { return fmt.Sprintf("failed to chug: %v", e.Err) }
func (e *ChuggingError) Unwrap() error { return e.Err }

type EndOfCharSourceError struct{ Err error }

func (e *EndOfCharSourceError) Error() string { return fmt.Sprintf("no more chars: %v", e.Err) }
func (e *EndOfCharSourceError) Unwrap() error { return e.Err }
