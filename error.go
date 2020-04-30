// error.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-04-30

// Package package contains some utility functions and types.
package util

import (
	"fmt"
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

// Raw returns the underlying raw error.
func (e *Error) Raw() error {
	return e.error
}

// Errorf formats according to a format specifier and returns the string as
// a value that satisfies error. Cause the underlying type of the error is
// *util.Error, you need to specify the error code (first parameter).
func Errorf(code int, format string, args ...interface{}) error {
	return &Error{fmt.Errorf(format, args...), code}
}
