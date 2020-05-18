// error_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-05-18

package util

import (
	"errors"
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
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", errors.Unwrap(err)))
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
		ft, fv := t.Field(i), v.Field(i)
		// PkgPath is empty for upper case (exported) field names.
		if ft.PkgPath == "" {
			strs = append(strs, fmt.Sprintf("%s=%v", ft.Tag.Get("json"), fv.Interface()))
		}
	}

	return strings.Join(strs, ";")
}
