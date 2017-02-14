package fuzzywuzzy

type levMatchingBlock struct {
	spos   int // source block pos
	dpos   int // destination block pos
	length int
}

type levEditType int

const (
	levEditKeep levEditType = iota
	levEditReplace
	levEditInsert
	levEditDelete
	levEditLast
)

type levEditOp struct {
	editType levEditType
	spos     int
	dpos     int
}

type levOpCode struct {
	editType   levEditType
	sbeg, send int
	dbeg, dend int
}

func Ratio(s1 string, s2 string) float64 {
	lenSum := len(s1) + len(s2)
	editDistance := LevEditDistance(s1, s2, 1)
	return float64(lenSum-editDistance) / float64(lenSum)
}

func PartialRatio(s1 string, s2 string) float64 {
	shorter, longer := s1, s2
	if len(s1) > len(s2) {
		longer, shorter = s1, s2
	}
	_, _ = shorter, longer
	return 0
}

// function for edit distance ported from python-Levenshtein
func LevEditDistance(s1 string, s2 string, xcost int) int {
	chrs1, chrs2 := []rune(s1), []rune(s2)
	len1, len2 := len(chrs1), len(chrs2)
	return levEditDistance(chrs1, len1, chrs2, len2, xcost)
}

func levEditDistance(chrs1 []rune, len1 int, chrs2 []rune, len2 int, xcost int) int {
	idx1, idx2 := 0, 0
	// strip common prefix
	for len1 > 0 && len2 > 0 && chrs1[idx1] == chrs2[idx2] {
		len1--
		len2--
		idx1++
		idx2++
	}
	// strip suffix
	for len1 > 0 && len2 > 0 && chrs1[idx1+len1-1] == chrs2[idx2+len2-1] {
		len1--
		len2--
	}

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}
	// if s1 is longer than s2, switch them around
	if len1 > len2 {
		idx1, idx2 = idx2, idx1
		len1, len2 = len2, len1
		chrs1, chrs2 = chrs2, chrs1
	}

	if len1 == 1 {
		if xcost != 0 {
			return len2 + 1 - 2*runeContained(chrs1[idx1], chrs2)
		} else {
			return len2 - runeContained(chrs1[idx1], chrs2)
		}
	}

	len1++
	len2++
	half := len1 >> 1

	var row = make([]int, len2)
	end := len2 - 1
	tmp := 0
	if xcost == 0 {
		tmp = half
	}

	for i := 0; i < len2-tmp; i++ {
		row[i] = i
	}

	if xcost != 0 {
		for i := 1; i < len1; i++ {
			p := 1

			char1 := chrs1[idx1+i-1]
			c2p := idx2
			D, x := i, i

			for p <= end {
				if char1 == chrs2[c2p] {
					D--
					x = D
				} else {
					x++
				}
				c2p++
				D = row[p]
				D++

				if x > D {
					x = D
				}
				row[p] = x
				p++
			}
		}
	} else {
		row[0] = len1 - half - 1
		for i := 1; i < len1; i++ {
			var c2p, D, x, p int
			char1 := chrs1[idx1+i-1]
			if i >= len1-half {
				offset := i - (len1 - half)
				c2p = idx2 + offset
				p = offset
				loc := 1
				if char1 == chrs2[c2p] {
					loc = 0
				}
				c3 := row[p] + loc
				p++
				c2p++
				x = row[p]
				x++
				D = x
				if x > c3 {
					x = c3
				}
				row[p] = x
				p++
			} else {
				p = 1
				c2p = idx2
				x, D = i, i
			}

			if i <= half+1 {
				end = len2 + i - half - 2
			}

			for p <= end {
				D--
				tmp := 1
				if char1 == chrs2[c2p] {
					tmp = 0
				}
				c2p++
				c3 := D + tmp
				x++
				if x > c3 {
					x = c3
				}
				D = row[p]
				D++
				if x > D {
					x = D
				}
				row[p] = x
				p++
			}

			if i <= half {
				D--
				loc := 1
				if char1 == chrs2[c2p] {
					loc = 0
				}
				x++
				c3 := D + loc
				if x > c3 {
					x = c3
				}
				row[p] = x
			}
		}
	}

	return row[end]
}

func runeContained(a rune, list []rune) int {
	if indexOfRune(a, list) >= 0 {
		return 1
	} else {
		return 0
	}
}

func indexOfRune(a rune, list []rune) int {
	for idx, b := range list {
		if b == a {
			return idx
		}
	}
	return -1
}
