package fuzzy

import (
	"math"
	"sort"
	"strings"
)

// Ratio computes a score of how close two unicode strings are
// based on their Levenshtein edit distance.
// Returns an integer score [0,100], higher score indicates
// that strings are closer.
func Ratio(s1, s2 string) int {
	return int(round(100 * floatRatio(s1, s2)))
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
		if longStart < 0 {
			longStart = 0
		}
		longEnd := longStart + len(shorter)
		if longEnd > len(longer) {
			longEnd = len(longer)
		}
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
	if lenSum == 0 {
		return 0.0
	}
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

// WRatio computes a score with the following steps:
// 1. Cleanse both strings, remove non-ASCII characters.
// 2. Take Ratio as baseline score.
// 3. Run a few heuristics to determine whether partial ratios
//    should be taken.
// 4. If partial ratios were determined to be necessary,
//    compute PartialRatio, PartialTokenSetRatio, and PartialTokenSortRatio.
//    Otherwise, compute TokenSortRatio and TokenSetRatio.
// 5. Return the max of all computed ratios.
func WRatio(s1, s2 string) int {
	return weightedRatioHelper(s1, s2, true)
}

// UWRatio computes a score similar to WRatio, except non-ASCII
// characters are allowed.
func UWRatio(s1, s2 string) int {
	return weightedRatioHelper(s1, s2, false)
}

func weightedRatioHelper(s1, s2 string, asciiOnly bool) int {
	c1 := Cleanse(s1, asciiOnly)
	c2 := Cleanse(s2, asciiOnly)

	if len(c1) == 0 || len(c2) == 0 {
		return 0
	}

	unbaseScale := .95
	partialScale := .9
	baseScore := float64(Ratio(c1, c2))
	lengthRatio := float64(len(c1)) / float64(len(c2))
	if lengthRatio < 1 {
		lengthRatio = 1 / lengthRatio
	}

	tryPartial := true
	if lengthRatio < 1.5 {
		tryPartial = false
	}

	if lengthRatio > 8 {
		partialScale = .6
	}

	if tryPartial {
		partialScore := float64(PartialRatio(c1, c2)) * partialScale
		tokenSortScore := float64(PartialTokenSortRatio(c1, c2, asciiOnly, false)) *
			unbaseScale * partialScale
		tokenSetScore := float64(PartialTokenSetRatio(c1, c2, asciiOnly, false)) *
			unbaseScale * partialScale
		return int(round(max(baseScore, partialScore, tokenSortScore, tokenSetScore)))
	}
	tokenSortScore := float64(TokenSortRatio(c1, c2, asciiOnly, false)) * unbaseScale
	tokenSetScore := float64(TokenSetRatio(c1, c2, asciiOnly, false)) * unbaseScale
	return int(round(max(baseScore, tokenSortScore, tokenSetScore)))
}

func max(args ...float64) float64 {
	maxVal := args[0]
	for _, arg := range args {
		if arg > maxVal {
			maxVal = arg
		}
	}
	return maxVal
}

// TokenSortRatio computes a score similar to Ratio, except tokens
// are sorted and (optionally) cleansed prior to comparison.
func TokenSortRatio(s1, s2 string, opts ...bool) int {
	return tokenSortRatioHelper(s1, s2, false, opts...)
}

// PartialTokenSortRatio computes a score similar to PartialRatio, except tokens
// are sorted and (optionally) cleansed prior to comparison.
func PartialTokenSortRatio(s1, s2 string, opts ...bool) int {
	return tokenSortRatioHelper(s1, s2, true, opts...)
}

func tokenSortRatioHelper(s1, s2 string, partial bool, opts ...bool) int {
	asciiOnly, cleanse := false, false
	for i, val := range opts {
		switch i {
		case 0:
			asciiOnly = val
		case 1:
			cleanse = val
		}
	}

	sorted1 := tokenSort(s1, asciiOnly, cleanse)
	sorted2 := tokenSort(s2, asciiOnly, cleanse)

	if partial {
		return PartialRatio(sorted1, sorted2)
	}
	return Ratio(sorted1, sorted2)
}

func tokenSort(s string, asciiOnly, cleanse bool) string {
	if cleanse {
		s = Cleanse(s, asciiOnly)
	} else if asciiOnly {
		s = ASCIIOnly(s)
	}

	tokens := strings.Fields(s)
	sort.Strings(tokens)
	return strings.Join(tokens, " ")
}

// TokenSetRatio extracts tokens from each input string, adds
// them to a set, construct strings of the form
// <sorted intersection><sorted remainder>, takes the ratios
// of those two strings, and returns the max.
func TokenSetRatio(s1, s2 string, opts ...bool) int {
	return tokenSetRatioHelper(s1, s2, false, opts...)
}

// PartialTokenSetRatio extracts tokens from each input string, adds
// them to a set, construct two strings of the form
// <sorted intersection><sorted remainder>, takes the partial ratios
// of those two strings, and returns the max.
func PartialTokenSetRatio(s1, s2 string, opts ...bool) int {
	return tokenSetRatioHelper(s1, s2, true, opts...)
}

func tokenSetRatioHelper(s1, s2 string, partial bool, opts ...bool) int {
	asciiOnly, cleanse := false, false
	for i, val := range opts {
		switch i {
		case 0:
			asciiOnly = val
		case 1:
			cleanse = val
		}
	}

	if cleanse {
		s1 = Cleanse(s1, asciiOnly)
		s2 = Cleanse(s2, asciiOnly)
	} else if asciiOnly {
		s1 = ASCIIOnly(s1)
		s2 = ASCIIOnly(s2)
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	set1 := NewStringSet(strings.Fields(s1))
	set2 := NewStringSet(strings.Fields(s2))
	intersection := set1.Intersect(set2).ToSlice()
	diff1to2 := set1.Difference(set2).ToSlice()
	diff2to1 := set2.Difference(set1).ToSlice()

	sort.Strings(intersection)
	sort.Strings(diff1to2)
	sort.Strings(diff2to1)

	sortedIntersect := strings.TrimSpace(strings.Join(intersection, " "))
	combined1to2 := strings.TrimSpace(sortedIntersect + " " + strings.Join(diff1to2, " "))
	combined2to1 := strings.TrimSpace(sortedIntersect + " " + strings.Join(diff2to1, " "))

	var ratioFunction func(string, string) int
	if partial {
		ratioFunction = PartialRatio
	} else {
		ratioFunction = Ratio
	}

	score := ratioFunction(sortedIntersect, combined1to2)
	if alt1 := ratioFunction(sortedIntersect, combined2to1); alt1 > score {
		score = alt1
	}
	if alt2 := ratioFunction(combined1to2, combined2to1); alt2 > score {
		score = alt2
	}

	return score
}

func round(x float64) float64 {
	if x < 0 {
		return math.Ceil(x - 0.5)
	}
	return math.Floor(x + 0.5)
}
