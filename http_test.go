// http_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-06
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-05-06

package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientIP(t *testing.T) {
	for _, cs := range []struct {
		r  *http.Request
		IP string `json:"client_ip"`
	}{
		{
			r: &http.Request{
				Header: http.Header(map[string][]string{
					"X-Forwarded-For": []string{"192.168.1.1"},
				}),
			},
			IP: "192.168.1.1",
		},
		{
			r: &http.Request{
				Header: http.Header(map[string][]string{
					"X-Forwarded-For": []string{"192.168.1.2, 176.10.1.2, 123.1.1.1"},
				}),
				RemoteAddr: "183.91.1.19:80",
			},
			IP: "192.168.1.2",
		},
		{
			r:  &http.Request{RemoteAddr: "183.91.1.19:80"},
			IP: "183.91.1.19",
		},
	} {
		t.Run(encodeCase(cs), func(t *testing.T) {
			assert.Equal(t, cs.IP, GetClientIP(cs.r))
		})
	}
}
