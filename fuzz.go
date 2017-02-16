package fuzzy

import (
	"math"
)

// Ratio computes a score of how close two unicode strings are
// based on their Levenshtein edit distance.
// Returns an integer score [0,100], higher score indicates
// that strings are closer.
func Ratio(s1, s2 string) int {
	return int(round(floatRatio(s1, s2)))
}

// PartialRatio computes a score of how close a string is with
// the most similar substring from another string.
// Order of arguments does not matter.
// Returns an integer score [0,100], higher score indicates
// that the string and substring are closer.
func PartialRatio(s1, s2 string) int {
	shorter, longer := s1, s2
	if len(s1) > len(s2) {
		longer, shorter = s1, s2
	}

	matchingBlocks := getMatchingBlocks(shorter, longer)

	bestScore := 0.0
	for _, block := range matchingBlocks {
		longStart := block.dpos - block.spos
		if longStart <= 0 {
			longStart = 0
		}
		longEnd := longStart + len(shorter)
		longSubStr := string([]rune(longer)[longStart:longEnd])

		r := floatRatio(shorter, longSubStr)
		if r > .995 {
			return 100
		} else if r > bestScore {
			bestScore = r
		}
	}

	return int(round(100 * bestScore))
}

func floatRatio(s1, s2 string) float64 {
	lenSum := len(s1) + len(s2)
	editDistance := LevEditDistance(s1, s2, 1)
	return float64(lenSum-editDistance) / float64(lenSum)
}

// QRatio computes a score similar to Ratio, except both strings are trimmed,
// cleansed of non-ASCII characters, and case-standardized.
func QRatio(s1, s2 string) int {
	return quickRatioHelper(s1, s2, true)
}

// UQRatio computes a score similar to Ratio, except both strings are trimmed
// and case-standardized.
func UQRatio(s1, s2 string) int {
	return quickRatioHelper(s1, s2, false)
}

func quickRatioHelper(s1, s2 string, asciiOnly bool) int {
	c1 := Cleanse(s1, asciiOnly)
	c2 := Cleanse(s2, asciiOnly)

	if len(c1) == 0 || len(c2) == 0 {
		return 0
	}
	return Ratio(c1, c2)
}

func round(x float64) float64 {
	if x < 0 {
		return math.Ceil(x - 0.5)
	}
	return math.Floor(x + 0.5)
}
