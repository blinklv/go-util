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
	return toJson(v, false, "")
}

// ToPrettyJson returns the pretty-print JSON encoding of v. It also
// won't escape special characters and doesn't return error.
func ToPrettyJson(v interface{}) []byte {
	return toJson(v, false, "\t")
}

// toJson is the underlying implementation of ToJson and ToPrettyJson.
func toJson(v interface{}, escape bool, indent string) []byte {
	var (
		buf = &bytes.Buffer{}
		enc = json.NewEncoder(buf)
	)

	enc.SetEscapeHTML(escape)
	enc.SetIndent("", indent)

	if err := enc.Encode(v); err != nil {
		return nil
	}

	b := buf.Bytes()
	if last := len(b) - 1; len(b) > 0 && b[last] == '\n' {
		// json.Encoder.Encode will add a newline character at the
		// end, so we need to remove it make this function consistent
		// with json.Marshal.
		return b[:last]
	}
	return b
}
