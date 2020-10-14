// error_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-04-30
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-10-14

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
			t.Logf("%s (%s:%d %s)", err, err.File, err.Line, err.Func)
			expectedErrMsg := fmt.Sprintf(cs.Format, cs.Args...)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", e))
			assert.True(t, ok)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", err))
			assert.Equal(t, cs.Code, err.Code)
			assert.Equal(t, expectedErrMsg, fmt.Sprintf("%s", errors.Unwrap(err)))
		})
	}
}

func TestListErrorFormatter(t *testing.T) {
	for _, cs := range []struct {
		Errors []error `json:"errors"`
		Str    string  `json:"str"`
	}{
		{nil, "<nil>"},
		{[]error{}, "<nil>"},
		{[]error{errors.New("foo")}, "foo"},
		{
			[]error{
				errors.New("hello"),
				errors.New("world"),
			},
			`multiple (2) errors:
   1. hello
   2. world`,
		},
		{
			[]error{
				errors.New("aaaaa"),
				errors.New("bbbbb"),
				errors.New("ccccc"),
				errors.New("ddddd"),
				errors.New("eeeee"),
				errors.New("fffff"),
				errors.New("ggggg"),
				errors.New("hhhhh"),
				errors.New("iiiii"),
				errors.New("jjjjj"),
				errors.New("kkkkk"),
				errors.New("lllll"),
				errors.New("mmmmm"),
				errors.New("nnnnn"),
				errors.New("ooooo"),
				errors.New("ppppp"),
			},
			`multiple (16) errors:
   1. aaaaa
   2. bbbbb
   3. ccccc
   4. ddddd
   5. eeeee
   6. fffff
   7. ggggg
   8. hhhhh
   9. iiiii
  10. jjjjj
  11. kkkkk
  12. lllll
  13. mmmmm
  14. nnnnn
  15. ooooo
  16. ppppp`,
		},
		{
			[]error{
				errors.New("aaaaa"),
				errors.New("bbbbb"),
				errors.New("ccccc"),
				errors.New("ddddd"),
				errors.New("eeeee"),
				errors.New("fffff"),
				errors.New("ggggg"),
				errors.New("hhhhh"),
				errors.New("iiiii"),
				errors.New("jjjjj"),
				errors.New("kkkkk"),
				errors.New("lllll"),
				errors.New("mmmmm"),
				errors.New("nnnnn"),
				errors.New("ooooo"),
				errors.New("ppppp"),
				errors.New("qqqqq"),
				errors.New("rrrrr"),
				errors.New("sssss"),
			},
			`multiple (19) errors:
   1. aaaaa
   2. bbbbb
   3. ccccc
   4. ddddd
   5. eeeee
   6. fffff
   7. ggggg
   8. hhhhh
   9. iiiii
  10. jjjjj
  11. kkkkk
  12. lllll
  13. mmmmm
  14. nnnnn
  15. ooooo
  16. ppppp
      ...`,
		},
	} {
		str := ListErrorFormatter(cs.Errors)
		t.Logf("str: %s", str)
		assert.Equal(t, cs.Str, str)
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
