// http.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-06
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-05-06

package util

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP gets the client ip from a http request. The implementation
// of this function based on some special HTTP headers, so it has security
// implications. Do NOT use this function unless there exists a trusted
// reverse proxy. The returned string might be empty or an illegal IP.
func GetClientIP(r *http.Request) string {
	// Extracts client ip from X-Forwarded-For header at first.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {

		// The value of X-Forwarded-For header is a comma-space separated list
		// of IP addresses, the left-most being the original client, and each
		// successive proxy that passed the request adding the IP address where
		// it received the request from.
		if parts := strings.Split(xff, ","); len(parts) != 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// RemoteAddr field will be set to an 'ip:port' address in most cases.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}
