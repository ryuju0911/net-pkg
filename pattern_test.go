// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import "testing"

func TestParsePattern(t *testing.T) {
	for _, test := range []struct {
		in   string
		want pattern
	}{
		{"/", pattern{}},
		{
			"example.com/",
			pattern{host: "example.com"},
		},
		{
			"GET /",
			pattern{method: "GET"},
		},
		{
			"POST example.com/foo/{w}",
			pattern{
				method: "POST",
				host:   "example.com",
			},
		},
	} {
		got := mustParsePattern(t, test.in)
		if !got.equal(&test.want) {
			t.Errorf("%q:\ngot  %#v\nwant %#v", test.in, got, &test.want)
		}
	}
}

func (p1 *pattern) equal(p2 *pattern) bool {
	return p1.method == p2.method && p1.host == p2.host
}

func mustParsePattern(tb testing.TB, s string) *pattern {
	tb.Helper()
	p, err := parsePattern(s)
	if err != nil {
		tb.Fatal(err)
	}
	return p
}
