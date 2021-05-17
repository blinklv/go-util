// group_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-06-04
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2021-05-17

package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	g := &Group{}

	// 1. Check whether nil result will be ignored.
	for i := 0; i < 128; i++ {
		g.Go(func() interface{} {
			return nil
		})
	}
	assert.True(t, g.Result() == nil)

	// 2. Check normal case.
	g.Go(func() interface{} { // nil result.
		time.Sleep(2 * time.Second)
		return nil
	})

	g.Go(func() interface{} { // error result.
		return fmt.Errorf("foo")
	})

	var input = struct {
		a int
		b string
	}{
		a: 0,
		b: "Hello, World!",
	}
	g.Go(func() interface{} { // valid result.
		return input
	})

	result := g.Result()
	assert.Equal(t, 2, len(result))
	for _, res := range result {
		t.Logf("result: %#v", res)
	}

	// 3. Check whether the result has been cleared.
	assert.True(t, g.Result() == nil)

	// 4. Check panic case 1.
	var panicStr string
	g.Go(func() interface{} {
		defer func() {
			if x := recover(); x != nil {
				panicStr = fmt.Sprintf("%s", x)
			}
		}()

		panic("bar")
		return "foo"
	})

	assert.True(t, g.Result() == nil)
	assert.Equal(t, "bar", panicStr)

	// 5. Check panic case 2.
	g.Go(func() interface{} {
		panic(errors.New("Hello, Boy!"))
		return "bar"
	})

	res := g.Result()
	e, ok := res[0].(error)
	assert.True(t, ok)
	assert.EqualError(t, e, "Hello, Boy!")

	// 6. Check panic case 3.
	g = &Group{Logger: log.New(os.Stderr, "", log.LstdFlags)}
	g.Go(func() interface{} {
		panic("Hi!")
		return nil
	})
	res = g.Result()
	e, ok = res[0].(error)
	assert.True(t, ok)
	assert.EqualError(t, e, "Hi!")

	// 7. Check Error method case 1.
	g.Go(func() interface{} {
		panic("error A")
		return nil
	})
	g.Go(func() interface{} {
		return errors.New("error B")
	})
	g.Go(func() interface{} {
		return 10
	})
	g.Go(func() interface{} {
		return "error C"
	})
	g.Go(func() interface{} {
		return nil
	})

	e = g.Error()
	t.Logf("%s", e)
	assert.Error(t, e)
	assert.Nil(t, g.Result())

	// 7. Check Error method case 2.
	g.Go(func() interface{} {
		return nil
	})
	g.Go(func() interface{} {
		return 10
	})
	g.Go(func() interface{} {
		return "error A"
	})

	e = g.Error(nil)
	assert.NoError(t, e)
	assert.Nil(t, g.Result())
}
