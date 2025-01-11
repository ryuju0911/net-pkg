package http

import "testing"

func TestCleanPath(t *testing.T) {
	for _, test := range []struct {
		in, want string
	}{
		{"//", "/"},
		{"/x", "/x"},
		{"//x", "/x"},
		{"x//", "/x/"},
		{"a//b/////c", "/a/b/c"},
		{"/foo/../bar/./..//baz", "/baz"},
	} {
		got := cleanPath(test.in)
		if got != test.want {
			t.Errorf("%s: got %q, want %q", test.in, got, test.want)
		}
	}
}
