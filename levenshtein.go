package fuzzy

// this code is a port of python-Levenshtein,
// which is a highly efficient (and obfuscated)
// implementation of Levenshtein distanceS

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

func findEditOps(s1, s2 string) []levEditOp {
	chrs1, chrs2 := []rune(s1), []rune(s2)
	len1, len2 := len(chrs1), len(chrs2)
	return findEditOpsHelper(chrs1, len1, chrs2, len2)
}

func findEditOpsHelper(chrs1 []rune, len1 int, chrs2 []rune, len2 int) []levEditOp {
	p1, p2 := 0, 0
	len1o := 0
	for len1 > 0 && len2 > 0 && chrs1[p1] == chrs2[p2] {
		len1--
		len2--

		p1++
		p2++

		len1o++
	}
	len2o := len1o

	for len1 > 0 && len2 > 0 && chrs1[p1+len1-1] == chrs2[p2+len2-1] {
		len1--
		len2--
	}
	len1++
	len2++

	matrix := make([]int, len2*len1)

	for i := 0; i < len2; i++ {
		matrix[i] = i
	}
	for i := 1; i < len1; i++ {
		matrix[len2*i] = i
	}

	for i := 1; i < len1; i++ {
		prev := (i - 1) * len2
		p := i * len2
		end := p + len2 - 1
		char1 := chrs1[p1+i-1]
		p2copy := p2
		x := i
		p++
		for p <= end {
			prevOpCost := 1
			if char1 == chrs2[p2copy] {
				prevOpCost = 0
			}
			c3 := matrix[prev] + prevOpCost
			p2copy++
			prev++
			x++

			if x > c3 {
				x = c3
			}
			c3 = matrix[prev] + 1
			if x > c3 {
				x = c3
			}

			matrix[p] = x
			p++
		}
	}

	return editOpsFromCostMatrix(len1, chrs1, p1, len1o, len2, chrs2, p2, len2o, matrix)
}

func editOpsFromCostMatrix(len1 int, chrs1 []rune, p1, o1 int, len2 int, chrs2 []rune, p2, o2 int, matrix []int) []levEditOp {
	dir := 0
	pos := matrix[len1*len2-1]
	ops := make([]levEditOp, pos)
	i, j := len1-1, len2-1
	ptr := len1*len2 - 1

	for i > 0 || j > 0 {
		if dir < 0 && j > 0 && matrix[ptr] == matrix[ptr-1]+1 {
			pos--
			j--
			ops[pos] = levEditOp{spos: i + o1, dpos: j + o2, editType: levEditInsert}
			ptr--
			continue
		}

		if dir > 0 && i > 0 && matrix[ptr] == matrix[ptr-len2]+1 {
			pos--
			i--
			ops[pos] = levEditOp{spos: i + o1, dpos: j + o2, editType: levEditDelete}
			ptr -= len2
			continue
		}

		if i > 0 && j > 0 && matrix[ptr] == matrix[ptr-len2-1] &&
			chrs1[p1+i-1] == chrs2[p2+j-1] {
			i--
			j--
			ptr -= len2 + 1
			dir = 0
			continue
		}

		if i > 0 && j > 0 && matrix[ptr] == matrix[ptr-len2-1]+1 {
			pos--
			i--
			j--
			ops[pos] = levEditOp{spos: i + o1, dpos: j + o2, editType: levEditReplace}
			ptr -= len2 + 1
			dir = 0
			continue
		}

		if dir == 0 && j > 0 && matrix[ptr] == matrix[ptr-1]+1 {
			pos--
			j--
			ops[pos] = levEditOp{spos: i + o1, dpos: j + o2, editType: levEditInsert}
			ptr--
			dir = -1
			continue
		}

		if dir == 0 && i > 0 && matrix[ptr] == matrix[ptr-len2]+1 {
			pos--
			i--
			ops[pos] = levEditOp{spos: i + o1, dpos: j + o2, editType: levEditDelete}
			ptr -= len2
			dir = 1
			continue
		}
	}

	return ops
}

func editOpsToOpCodes(ops []levEditOp, len1, len2 int) []levOpCode {
	n := len(ops)
	nBlocks := 0 // number of blocks
	opIdx, spos, dpos := 0, 0, 0
	var editType levEditType

	for i := n; i > 0; {
		for ops[opIdx].editType == levEditKeep && i > 0 {
			i--
			opIdx++
		}

		if i == 0 {
			break
		}

		if spos < ops[opIdx].spos || dpos < ops[opIdx].dpos {
			nBlocks++
			spos = ops[opIdx].spos
			dpos = ops[opIdx].dpos
		}

		nBlocks++
		editType = ops[opIdx].editType

		switch editType {
		case levEditReplace:
			// emulate do...while loop
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				dpos++
				i--
				opIdx++
			}
		case levEditDelete:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				i--
				opIdx++
			}
		case levEditInsert:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				dpos++
				i--
				opIdx++
			}
		}
	}
	if spos < len1 || dpos < len2 {
		nBlocks++
	}

	opCodes := make([]levOpCode, nBlocks)
	opIdx, spos, dpos = 0, 0, 0
	codeIdx := 0

	for i := n; i != 0; {
		for ops[opIdx].editType == levEditKeep {
			i--
			if i <= 0 {
				break
			}
			opIdx++
		}

		if i == 0 {
			break
		}

		oc := levOpCode{sbeg: spos, dbeg: dpos}
		opCodes[codeIdx] = oc
		if spos < ops[opIdx].spos || dpos < ops[opIdx].dpos {
			oc.editType = levEditKeep
			oc.send = ops[opIdx].spos
			oc.dend = ops[opIdx].dpos
			spos = oc.send
			dpos = oc.dend

			codeIdx++
			oc2 := levOpCode{sbeg: spos, dbeg: dpos}
			opCodes[codeIdx] = oc2
		}
		editType = ops[opIdx].editType

		switch editType {
		case levEditReplace:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				dpos++
				i--
				opIdx++
			}
		case levEditDelete:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				i--
				opIdx++
			}
		case levEditInsert:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				dpos++
				i--
				opIdx++
			}
		}

		opCodes[codeIdx].editType = editType
		opCodes[codeIdx].send = spos
		opCodes[codeIdx].dend = dpos
		codeIdx++
	}

	if spos < len1 || dpos < len2 {
		opCodes[codeIdx].editType = levEditKeep
		opCodes[codeIdx].sbeg = spos
		opCodes[codeIdx].dbeg = dpos
		opCodes[codeIdx].send = len1
		opCodes[codeIdx].dend = len2
	}

	return opCodes
}

// emulate do...while loop
func shouldContinue(i int, editOps []levEditOp, opIdx int, editType levEditType, spos, dpos int) bool {
	return i > 0 && editOps[opIdx].editType == editType &&
		editOps[opIdx].dpos == dpos && editOps[opIdx].spos == spos
}

func getMatchingBlocks(s1, s2 string) []levMatchingBlock {
	chrs1, chrs2 := []rune(s1), []rune(s2)
	len1, len2 := len(chrs1), len(chrs2)

	return getMatchingBlocksHelper(len1, len2, findEditOpsHelper(chrs1, len1, chrs2, len2))
}

func getMatchingBlocksHelper(len1, len2 int, ops []levEditOp) []levMatchingBlock {
	n := len(ops)
	nMatchingBlocks := 0
	opIdx, spos, dpos := 0, 0, 0
	var editType levEditType
	for i := n; i > 0; {
		for ops[opIdx].editType == levEditKeep {
			i--
			if i <= 0 {
				break
			}
			opIdx++
		}

		if i == 0 {
			break
		}

		if spos < ops[opIdx].spos || dpos < ops[opIdx].dpos {
			nMatchingBlocks++
			spos = ops[opIdx].spos
			dpos = ops[opIdx].dpos
		}

		editType = ops[opIdx].editType

		switch editType {
		case levEditReplace:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				dpos++
				i--
				opIdx++
			}
		case levEditDelete:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				i--
				opIdx++
			}
		case levEditInsert:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				dpos++
				i--
				opIdx++
			}
		}
	}

	if spos < len1 || dpos < len2 {
		nMatchingBlocks++
	}

	matchingBlocks := make([]levMatchingBlock, nMatchingBlocks+1)

	opIdx = 0
	spos, dpos = 0, 0
	blockIdx := 0
	for i := n; i > 0; {
		for ops[opIdx].editType == levEditKeep {
			i--
			if i <= 0 {
				break
			}
			opIdx++
		}

		if i <= 0 {
			break
		}

		if spos < ops[opIdx].spos || dpos < ops[opIdx].dpos {
			mb := levMatchingBlock{spos: spos, dpos: dpos, length: ops[opIdx].spos - spos}
			spos = ops[opIdx].spos
			dpos = ops[opIdx].dpos
			matchingBlocks[blockIdx] = mb
			blockIdx++
		}

		editType = ops[opIdx].editType
		switch editType {
		case levEditReplace:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				dpos++
				i--
				opIdx++
			}
		case levEditDelete:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				spos++
				i--
				opIdx++
			}
		case levEditInsert:
			for ok := true; ok; ok = shouldContinue(i, ops, opIdx, editType, spos, dpos) {
				dpos++
				i--
				opIdx++
			}
		}
	}
	if spos < len1 || dpos < len2 {
		mb := levMatchingBlock{spos: spos, dpos: dpos, length: len1 - spos}
		matchingBlocks[blockIdx] = mb
		blockIdx++
	}
	lastBlock := levMatchingBlock{spos: len1, dpos: len2, length: 0}
	matchingBlocks[blockIdx] = lastBlock

	return matchingBlocks
}

func getMatchingBlocksFromOpCodes(len1, len2 int, ops []levOpCode) []levMatchingBlock {
	n := len(ops)
	nMB := 0
	codeIdx := 0

	for i := n; i > 0; codeIdx++ {
		i--
		if ops[codeIdx].editType == levEditKeep {
			nMB++
			for i > 0 && ops[codeIdx].editType == levEditKeep {
				i--
				codeIdx++
			}

			if i == 0 {
				break
			}
		}
	}

	matchingBlocks := make([]levMatchingBlock, nMB+1)
	codeIdx = 0
	mbIdx := 0

	for i := n; i > 0; i, codeIdx = i-1, codeIdx+1 {
		if ops[codeIdx].editType == levEditKeep {
			matchingBlocks[mbIdx].spos = ops[codeIdx].sbeg
			matchingBlocks[mbIdx].dpos = ops[codeIdx].dbeg

			for i > 0 && ops[codeIdx].editType == levEditKeep {
				i--
				codeIdx++
			}

			if i == 0 {
				matchingBlocks[mbIdx].length = len1 - matchingBlocks[mbIdx].spos
				mbIdx++
				break
			}

			matchingBlocks[mbIdx].length = ops[codeIdx].sbeg - matchingBlocks[mbIdx].spos
			mbIdx++
		}
	}

	//final matching block
	matchingBlocks[mbIdx].spos = len1
	matchingBlocks[mbIdx].dpos = len2
	matchingBlocks[mbIdx].length = 0

	return matchingBlocks
}

// EditDistance omputes the Levenshtein distance between two strings,
// weighting replacements the same as insertions and deletions.
func EditDistance(s1, s2 string) int {
	return LevEditDistance(s1, s2, 1)
}

// LevEditDistance computes Levenshtein distance between 2 strings
func LevEditDistance(s1, s2 string, xcost int) int {
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
		}
		return len2 - runeContained(chrs1[idx1], chrs2)
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
	}
	return 0
}

func indexOfRune(a rune, list []rune) int {
	for idx, b := range list {
		if b == a {
			return idx
		}
	}
	return -1
}
