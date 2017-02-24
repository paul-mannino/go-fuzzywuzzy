package fuzzy

import (
	"regexp"
	"strings"
)

func Cleanse(s string, forceAscii bool) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if forceAscii {
		s = ASCIIOnly(s)
	}

	r := regexp.MustCompile("[^\\p{L}\\p{N}]")
	s = r.ReplaceAllString(s, " ")
	return s
}

func ASCIIOnly(s string) string {
	runes := []rune(s)
	stripped := make([]rune, len(runes))

	w := 0
	for _, r := range runes {
		if r < 128 {
			stripped[w] = r
			w++
		}
	}
	return string(stripped[0:w])
}
