package http

import "golang.org/x/net/http/httpguts"

func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}
