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
			"/{w1}/lit/{w2}",
			pattern{
				segments: []segment{wild("w1"), lit("lit"), wild("w2")},
			},
		},
		{
			"/{w1}/lit/{w2}/",
			pattern{
				segments: []segment{wild("w1"), lit("lit"), wild("w2"), multi("")},
			},
		},
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
				segments: []segment{lit("foo"), wild("w")},
			},
		},
		{
			"/{$}",
			pattern{segments: []segment{lit("/")}},
		},
		{
			"DELETE example.com/a/{foo12}/{$}",
			pattern{method: "DELETE", host: "example.com", segments: []segment{lit("a"), wild("foo12"), lit("/")}},
		},
		{
			"/foo/{$}",
			pattern{segments: []segment{lit("foo"), lit("/")}},
		},
		{
			"/{a}/foo/{rest...}",
			pattern{segments: []segment{wild("a"), lit("foo"), multi("rest")}},
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
		{
			"POST \t  example.com/foo/{w}",
			pattern{
				method:   "POST",
				host:     "example.com",
				segments: []segment{lit("foo"), wild("w")},
			},
		},
		{
			"DELETE    \texample.com/a/{foo12}/{$}",
			pattern{method: "DELETE", host: "example.com", segments: []segment{lit("a"), wild("foo12"), lit("/")}},
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
		{"/{w}x", "at offset 1: bad wildcard segment"},
		{"/x{w}", "at offset 1: bad wildcard segment"},
		{"/{wx", "at offset 1: bad wildcard segment"},
		{"/a/{/}/c", "at offset 3: bad wildcard segment"},
		{"/a/{%61}/c", "at offset 3: bad wildcard name"}, // wildcard names aren't unescaped
		{"/{a$}", "at offset 1: bad wildcard name"},
		{"/{}", "at offset 1: empty wildcard"},
		{"POST a.com/x/{}/y", "at offset 13: empty wildcard"},
		{"/{...}", "at offset 1: empty wildcard"},
		{"/{$...}", "at offset 1: bad wildcard"},
		{"/{$}/", "at offset 1: {$} not at end"},
		{"/{$}/x", "at offset 1: {$} not at end"},
		{"/abc/{$}/x", "at offset 5: {$} not at end"},
		{"/{a...}/", "at offset 1: {...} wildcard not at end"},
		{"/{a...}/x", "at offset 1: {...} wildcard not at end"},
		{"{a}/b", "at offset 0: host contains '{' (missing initial '/'?)"},
		{"/a/{x}/b/{x...}", "at offset 9: duplicate wildcard name"},
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
