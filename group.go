// group.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-06-03
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-10-28

package util

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

// Group is a collection of goroutines which usually run simultaneously.
type Group struct {
	wg     sync.WaitGroup
	locker sync.Mutex
	result []interface{}

	// Logger specifies an optional logger for unexpected behaviors, which
	// lead to panic, from callback functions.
	Logger *log.Logger
}

// Go method is similar to 'go' statement which starts the execution of
// a function call in a new goroutine except for the function can't be
// arbitrary. The input function doesn't accept any parameter (but you
// can use the closure feature to inject parameters) and returns a value
// of any type. If the return value is (no type) nil, it will be ignored.
func (g *Group) Go(f func() interface{}) {
	g.wg.Add(1)
	go func() {
		defer func() {
			if x := recover(); x != nil {
				g.result = append(g.result, fmt.Errorf("%v", x))
				if g.Logger != nil {
					g.Logger.Printf("panic: %v\n%s", x, debug.Stack())
				}
			}

			// We need to place Done operation at here instead of the end
			// of this anonymous function, cause the custom function 'f'
			// might be panic.
			g.wg.Done()
		}()

		if res := f(); res != nil {
			// Only when the return value of the function is not nil, the
			// value will be added to the collection of results.
			g.locker.Lock()
			g.result = append(g.result, res)
			g.locker.Unlock()
		}
	}()
}

// Result method blocks until all function calls from the Go method have
// returned, then returns all resutls since the last time Result method was
// called, which means results will be cleared after calling this method.
func (g *Group) Result() []interface{} {
	g.wg.Wait()

	// Do we need to lock the result field after g.wg.Wait() has returned?
	// Yes, cause the Go method might be called in another goroutine at
	// this point.
	g.locker.Lock()
	result := g.result
	g.result = nil // Clear results.
	g.locker.Unlock()
	return result
}

// Error calls Result method at first, which means results will be cleared.
// Then it will extract all error results (the underlying type is error)
// and merge them into a single error. The message of the returned error is
// formatted by the ErrorFormatter argument, if you don't specify it (pass nil),
// ListErrorFormatter will be used. If there is no any error result, returns nil.
// In fact, if there are different types of return values, I don't recommend
// you use this method, cause it will ignore some non-error values.
func (g *Group) Error(ef ErrorFormatter) error {
	res := g.Result()
	if res == nil {
		return nil
	}

	es := make([]error, 0, len(res))
	for _, x := range res {
		if e, ok := x.(error); ok {
			es = append(es, e)
		}
	}

	if len(es) == 0 {
		return nil
	}

	if ef == nil {
		ef = ListErrorFormatter
	}
	return ef(es)
}
