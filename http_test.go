// http_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-05-06
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2021-10-19

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
		{
			r:  &http.Request{RemoteAddr: "183.91.1.19"},
			IP: "183.91.1.19",
		},
		{
			r:  &http.Request{RemoteAddr: "Foo, Bar"},
			IP: "Foo, Bar",
		},
	} {
		t.Run(encodeCase(cs), func(t *testing.T) {
			assert.Equal(t, cs.IP, GetClientIP(cs.r))
		})
	}
}

func TestSetCookie(t *testing.T) {
	for _, cs := range []struct {
		OldCookie string `yaml:"old_cookie"`
		Key       string `yaml:"key"`
		Value     string `yaml:"value"`
		NewCookie string `yaml:"new_cookie"`
	}{
		{"", "a", "b", "a=b"},
		{"foo=bar;hello=world", "what", "are", "foo=bar;hello=world;what=are"},
		{"foo=bar;hello=world;what=are", "foo", "foo", "foo=foo;hello=world;what=are"},
		{"foo=bar;hello=world;what=are", "hello", "a", "foo=bar;hello=a;what=are"},
		{"foo=bar;hello=world;what=are", "what", "you", "foo=bar;hello=world;what=you"},
		{"foo=bar;a=b;hello=world;a=b;what=are", "a", "haha", "foo=bar;a=haha;hello=world;a=haha;what=are"},
	} {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Set("Cookie", cs.OldCookie)
		SetCookie(r, cs.Key, cs.Value)
		assert.Equal(t, cs.NewCookie, r.Header.Get("Cookie"))
	}
}

func TestDelCookie(t *testing.T) {
	for _, cs := range []struct {
		OldCookie string `yaml:"old_cookie"`
		DeleteKey string `yaml:"delete_key"`
		NewCookie string `yaml:"new_cookie"`
	}{
		{"", "foo", ""},
		{"a=b", "a", ""},
		{"a=b", "b", "a=b"},
		{"a=b;c=d;hello=world;foo=bar", "hello", "a=b;c=d;foo=bar"},
		{"a=b;foo=bar;c=d;hello=world;foo=bar;x=y", "foo", "a=b;c=d;hello=world;x=y"},
	} {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Set("Cookie", cs.OldCookie)
		DelCookie(r, cs.DeleteKey)
		assert.Equal(t, cs.NewCookie, r.Header.Get("Cookie"))
	}
}

func TestSetURLQuery(t *testing.T) {
	for _, ts := range []struct {
		URL      string
		Query    map[string]string
		FinalURL string
	}{
		{"https::////foo.bar.com", nil, "https::////foo.bar.com"},
		{
			"https://foo.bar.example",
			map[string]string{
				"foo": "bar",
			},
			"https://foo.bar.example?foo=bar",
		},
		{
			"https://foo.bar.example?foo=hello",
			map[string]string{
				"foo": "bar",
			},
			"https://foo.bar.example?foo=bar",
		},
	} {
		query := SetURLQuery(ts.URL, ts.Query)
		t.Logf("query: %s", query)
		assert.Equal(t, ts.FinalURL, query)
	}
}
