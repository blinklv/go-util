// error.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-08-12

// Package package contains some utility functions and types.
package util

import (
	"fmt"
	"strings"
)

// Error is an implementation of the error interface, and an additional
// error code is attached to it compared to the raw error. It's useful
// for us to distinguish error types accurately.
type Error struct {
	error // The underlying raw error.

	// Code represents an error code that can help us to
	// locate a particular error quickly.
	Code int
}

// Unwrap returns the underlying raw error.
func (e *Error) Unwrap() error {
	return e.error
}

// Errorf formats according to a format specifier and returns the string as
// a value that satisfies error. Cause the underlying type of the error is
// *util.Error, you need to specify the error code (first parameter).
func Errorf(code int, format string, args ...interface{}) error {
	return &Error{fmt.Errorf(format, args...), code}
}

// WrapError wraps the raw error with an additional error code. The underlying
// type of the returned error is *util.Error.
func WrapError(code int, err error) error {
	return &Error{err, code}
}

// ErrorFormatter specifies an interface that can convert the list of errors into a string.
type ErrorFormatter func([]error) string

// ListErrorFormatter is a basic formatter that outputs the number of errors
// in the form of a list. If there is only one error, returns the error itself.
// If there're more than 16 errors, the remaining errors are indicated by ellipsis.
func ListErrorFormatter(es []error) string {
	if len(es) == 0 {
		return "<nil>"
	} else if len(es) == 1 {
		return es[0].Error()
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("multiple (%d) errors:\n", len(es)))

	for i, e := range es {
		if i > 0 {
			b.WriteString("\n")
		}

		if i >= 16 {
			b.WriteString("      ...")
			break
		}

		b.WriteString(fmt.Sprintf("%4d. %s", i+1, e))
	}

	return b.String()
}
