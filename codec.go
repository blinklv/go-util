// codec.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-08
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-05-08

package util

import (
	"bytes"
	"encoding/json"
)

// ToJson returns the JSON encoding of v. Compared with json.Marshal,
// it won't escape special characters (&, <, and >) in quoted strings
// to avoid certain safety problems, and doesn't return error.
func ToJson(v interface{}) []byte {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	return b.Bytes()
}

// ToPrettyJson returns the pretty-print JSON encoding of v. It also
// won't escape special characters and doesn't return error.
func ToPrettyJson(v interface{}) []byte {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	enc.Encode(v)
	return b.Bytes()
}
