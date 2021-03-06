// error.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-10-28

// Package package contains some utility functions and types.
package util

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
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

	// File represents the file name of the source code
	// that generates the error.
	File string

	// Line represents the line number of the source code
	// at which the error occurs.
	Line int

	// Func represents the name of the function that generates the error.
	Func string
}

// Unwrap returns the underlying raw error.
func (e *Error) Unwrap() error {
	return e.error
}

// Errorf formats according to a format specifier and returns the string as
// a value that satisfies error. Cause the underlying type of the error is
// *util.Error, you need to specify the error code (first parameter).
func Errorf(code int, format string, args ...interface{}) error {
	pc, fn, line, _ := runtime.Caller(1)
	return &Error{
		fmt.Errorf(format, args...),
		code,
		filepath.Base(fn),
		line,
		filepath.Base(runtime.FuncForPC(pc).Name()),
	}
}

// WrapError wraps the raw error with an additional error code. The underlying
// type of the returned error is *util.Error.
func WrapError(code int, err error) error {
	pc, fn, line, _ := runtime.Caller(1)
	return &Error{
		err,
		code,
		filepath.Base(fn),
		line,
		filepath.Base(runtime.FuncForPC(pc).Name()),
	}
}

// ErrorFormatter specifies an interface that can convert the list of errors into an error.
type ErrorFormatter func([]error) error

// ListErrorFormatter is a basic formatter that outputs the number of errors
// in the form of a list. If there is only one error, returns the error itself.
// If there're more than 16 errors, the remaining errors are indicated by ellipsis.
func ListErrorFormatter(es []error) error {
	if len(es) == 0 {
		return nil
	} else if len(es) == 1 {
		return es[0]
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

	return errors.New(b.String())
}

// CommaErrorFormatter outputs errors separated by a comma. If there is only
// one error, returns the error itself.
func CommaErrorFormatter(es []error) error {
	if len(es) == 0 {
		return nil
	} else if len(es) == 1 {
		return es[0]
	}

	strs := make([]string, 0, len(es))
	for _, e := range es {
		strs = append(strs, e.Error())
	}
	return errors.New(strings.Join(strs, ","))
}
