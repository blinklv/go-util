// http.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-06
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-06-28

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

// AddCookie adds a cookie to the HTTP request.  AddCookie does not attach
// more than one Cookie header field. That means all cookies, if any, are
// written into the same line, separated by semicolon.
func AddCookie(r *http.Request, name, value string) {
	r.AddCookie(&http.Cookie{Name: name, Value: value})
}

// GetCookie gets the value of a cookie. If the named cookie is not
// found, returns an empty string. If multiple cookies match the given
// name, only one cookie value will be returned.
func GetCookie(r *http.Request, name string) string {
	if cookie, err := r.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

// DelCookie deletes the cookies associated with the given name from the
// HTTP request's Cookie header.
func DelCookie(r *http.Request, name string) {
	var cookies []string
	for _, cookie := range r.Cookies() {
		if cookie.Name != name {
			cookies = append(cookies, cookie.String())
		}
	}

	if len(cookies) != 0 {
		r.Header.Set("Cookie", strings.Join(cookies, ";"))
	} else {
		// No cookies reserved now, so delete Cookie header directly.
		r.Header.Del("Cookie")
	}
}
