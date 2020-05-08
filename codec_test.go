// codec_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-08
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-05-08

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
	for i, cs := range []struct {
		v      interface{}
		output string
	}{
		{nil, `null`},
		{-123, `-123`},
		{"", `""`},
		{"<foo>hello</foo>", `"<foo>hello</foo>"`},
		{
			map[string]interface{}{
				"url": "https://www.google.com/search?sxsrf=ALeKk03T38aZACvdHVXYVKAQDgcXuqxYcw%3A1588925736603&source=hp&ei=KBW1XsWkIomkmAW8nZv4BQ&q=golang+return+value&oq=golang+return+value&gs_lcp=CgZwc3ktYWIQAzIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywE6BggjECcQEzoECCMQJzoECAAQQzoFCAAQgwE6BwgAEIMBEEM6AggAUM4BWMAVYJgaaABwAHgAgAHIAYgBohCSAQYxMS43LjGYAQCgAQGqAQdnd3Mtd2l6&sclient=psy-ab&ved=0ahUKEwjF2OjQ6aPpAhUJEqYKHbzOBl8Q4dUDCAc&uact=5",
			},
			`{"url":"https://www.google.com/search?sxsrf=ALeKk03T38aZACvdHVXYVKAQDgcXuqxYcw%3A1588925736603&source=hp&ei=KBW1XsWkIomkmAW8nZv4BQ&q=golang+return+value&oq=golang+return+value&gs_lcp=CgZwc3ktYWIQAzIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywEyBQgAEMsBMgUIABDLATIFCAAQywE6BggjECcQEzoECCMQJzoECAAQQzoFCAAQgwE6BwgAEIMBEEM6AggAUM4BWMAVYJgaaABwAHgAgAHIAYgBohCSAQYxMS43LjGYAQCgAQGqAQdnd3Mtd2l6&sclient=psy-ab&ved=0ahUKEwjF2OjQ6aPpAhUJEqYKHbzOBl8Q4dUDCAc&uact=5"}`,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, cs.output, string(ToJson(cs.v)))
		})
	}
}

func TestToPrettyJson(t *testing.T) {
	o := map[string]interface{}{
		"hello": 1,
		"world": "bar",
		"foo": map[string]interface{}{
			"Andy": map[string]interface{}{
				"age": 12,
				"sex": "male",
			},
			"Gina": map[string]interface{}{
				"age": 24,
				"sex": "female",
			},
		},
		"url": "https://www.bar.com/foo/path?a=b&c=1&hello=200",
	}

	var out = &bytes.Buffer{}
	json.Indent(out, ToJson(o), "", "\t")
	assert.Equal(t, out.String(), string(ToPrettyJson(o)))
}
