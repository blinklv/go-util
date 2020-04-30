// error_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-04-30

package util

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	for _, cs := range []struct {
		Code   int           `json:"code"`
		Format string        `json:"format"`
		Args   []interface{} `json:"args"`
	}{
		{-1, "You're wrong!", nil},
		{-2, "Hello, %s!", []interface{}{"World"}},
		{-3, "Foo is %d, %s is 2", []interface{}{1, "Bar"}},
	} {
		t.Run(encodeCase(cs), func(t *testing.T) {
			e := Errorf(cs.Code, cs.Format, cs.Args...)
			t.Logf("%s", e)

			err, ok := e.(*Error)
			expectedErrMsg := fmt.Sprintf(cs.Format, cs.Args...)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", e))
			assert.True(t, ok)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", err))
			assert.Equal(t, cs.Code, err.Code)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", err.Raw()))
		})
	}
}

func encodeCase(cs interface{}) string {
	var (
		strs []string
		v    = reflect.ValueOf(cs)
		t    = v.Type()
	)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		strs = append(strs, fmt.Sprintf("%s=%v", f.Tag.Get("json"), v.Field(i).Interface()))
	}

	return strings.Join(strs, ";")
}
