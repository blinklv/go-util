// group.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-06-03
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-08-06

package util

import "sync"

// Group is a collection of goroutines which usually run simultaneously.
type Group struct {
	wg     sync.WaitGroup
	locker sync.Mutex
	result []interface{}
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
				g.result = append(g.result, x)
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
