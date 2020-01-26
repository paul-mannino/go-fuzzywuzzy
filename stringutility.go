package fuzzy

import (
	"strings"
	"unicode"
)

func Cleanse(s string, forceASCII bool) string {
	if forceASCII {
		s = ASCIIOnly(s)
	}
	s = strings.TrimSpace(s)
	rs := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			r = ' '
		}
		rs = append(rs, r)
	}
	return strings.ToLower(string(rs))
}

func ASCIIOnly(s string) string {
	b := make([]byte, 0, len(s))
	for _, r := range s {
		if r <= unicode.MaxASCII {
			b = append(b, byte(r))
		}
	}
	return string(b)
}
