// http.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-06
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2021-10-19

package util

import (
	"net"
	"net/http"
	"net/url"
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

// SetCookie sets the named cookie to value. If the named cookie is
// not found, adds the new one.
func SetCookie(r *http.Request, name, value string) {
	var (
		cookies []string
		found   bool
	)

	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			cookie.Value = value
			found = true
		}
		cookies = append(cookies, cookie.String())
	}

	if !found {
		cookies = append(cookies, (&http.Cookie{Name: name, Value: value}).String())
	}

	if len(cookies) > 0 {
		r.Header.Set("Cookie", strings.Join(cookies, ";"))
	}
	return
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

// SetURLQuery injects query parameters into the URL and returns the modified one.
// If some query paramenters already exist, they will be updated.
func SetURLQuery(URL string, query map[string]string) string {
	// If the url arguments is invalid, returns it.
	u, err := url.Parse(URL)
	if err != nil {
		return URL
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
