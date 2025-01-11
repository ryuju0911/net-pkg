// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"slices"
	"strings"
	"testing"
)

func TestParsePattern(t *testing.T) {
	lit := func(name string) segment {
		return segment{s: name}
	}

	wild := func(name string) segment {
		return segment{s: name, wild: true}
	}

	multi := func(name string) segment {
		s := wild(name)
		s.multi = true
		return s
	}

	for _, test := range []struct {
		in   string
		want pattern
	}{
		{"/", pattern{segments: []segment{multi("")}}},
		{"/a", pattern{segments: []segment{lit("a")}}},
		{
			"/a/",
			pattern{segments: []segment{lit("a"), multi("")}},
		},
		{"/path/to/something", pattern{segments: []segment{
			lit("path"), lit("to"), lit("something"),
		}}},
		{
			"example.com/",
			pattern{host: "example.com", segments: []segment{multi("")}},
		},
		{
			"GET /",
			pattern{method: "GET", segments: []segment{multi("")}},
		},
		{
			"POST example.com/foo/{w}",
			pattern{
				method:   "POST",
				host:     "example.com",
				segments: []segment{lit("foo")},
			},
		},
		{
			"//",
			pattern{segments: []segment{lit(""), multi("")}},
		},
		{
			"/foo///./../bar",
			pattern{segments: []segment{lit("foo"), lit(""), lit(""), lit("."), lit(".."), lit("bar")}},
		},
		{
			"a.com/foo//",
			pattern{host: "a.com", segments: []segment{lit("foo"), lit(""), multi("")}},
		},
		{
			"/%61%62/%7b/%",
			pattern{segments: []segment{lit("ab"), lit("{"), lit("%")}},
		},
		// Allow multiple spaces matching regexp '[ \t]+' between method and path.
		{
			"GET\t  /",
			pattern{method: "GET", segments: []segment{multi("")}},
		},
	} {
		got := mustParsePattern(t, test.in)
		if !got.equal(&test.want) {
			t.Errorf("%q:\ngot  %#v\nwant %#v", test.in, got, &test.want)
		}
	}
}

func TestParsePatternError(t *testing.T) {
	for _, test := range []struct {
		in       string
		contains string
	}{
		{"", "empty pattern"},
		{"A=B /", "at offset 0: invalid method"},
		{" ", "at offset 1: host/path missing /"},
		{"{a}/b", "at offset 0: host contains '{' (missing initial '/'?)"},
		{"GET //", "at offset 4: non-CONNECT pattern with unclean path"},
	} {
		_, err := parsePattern(test.in)
		if err == nil || !strings.Contains(err.Error(), test.contains) {
			t.Errorf("%q:\ngot %v, want error containing %q", test.in, err, test.contains)
		}
	}
}

func (p1 *pattern) equal(p2 *pattern) bool {
	return p1.method == p2.method && p1.host == p2.host &&
		slices.Equal(p1.segments, p2.segments)
}

func mustParsePattern(tb testing.TB, s string) *pattern {
	tb.Helper()
	p, err := parsePattern(s)
	if err != nil {
		tb.Fatal(err)
	}
	return p
}
