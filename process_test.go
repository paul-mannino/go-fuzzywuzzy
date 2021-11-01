package fuzzy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var baseballStrings = []string{
	"new york mets vs chicago cubs",
	"chicago cubs vs chicago white sox",
	"philladelphia phillies vs atlanta braves",
	"braves vs mets",
}

var moreBaseballStrings = []string{
	"new york mets vs chicago cubs",
	"chicago cubs at new york mets",
	"atlanta braves vs pittsbugh pirates",
	"new york yankees vs boston red sox",
}

var someEmptyStrings = []string{
	"",
	"new york mets vs chicago cubs",
	"new york yankees vs boston red sox",
	"",
	"",
}

var someNullStrings = []string{}

func TestExtractOne(t *testing.T) {
	query1 := "new york mets at atlanta braves"
	best1, _ := ExtractOne(query1, baseballStrings)
	assertMatch(t, query1, baseballStrings[3], best1.Match)

	query2 := "philadelphia phillies at atlanta braves"
	best2, _ := ExtractOne(query2, baseballStrings)
	assertMatch(t, query2, baseballStrings[2], best2.Match)

	query3 := "atlanta braves at philadelphia phillies"
	best3, _ := ExtractOne(query3, baseballStrings)
	assertMatch(t, query3, baseballStrings[2], best3.Match)

	query4 := "chicago cubs vs new york mets"
	best4, _ := ExtractOne(query4, baseballStrings)
	assertMatch(t, query4, baseballStrings[0], best4.Match)

	query5 := "new york mets at chicago cubs"
	best5, _ := ExtractOne(query5, baseballStrings)
	assertMatch(t, query5, baseballStrings[0], best5.Match)

	customScorer := func(s1, s2 string) int {
		return QRatio(s1, s2)
	}
	best6, _ := ExtractOne(query5, baseballStrings, customScorer)
	assertMatch(t, query5, baseballStrings[0], best6.Match)

	query7 := "los angeles dodgers vs san francisco giants"
	scoreCutoff := 50
	best7, _ := ExtractOne(query7, moreBaseballStrings, scoreCutoff)
	if best7 != nil {
		t.Error("expecting to find no matches")
	}

	query8 := "new york mets vs chicago cubs"
	scoreCutoff = 100
	best8, _ := ExtractOne(query8, moreBaseballStrings, scoreCutoff)
	if best8 == nil {
		t.Error("expecting to find a match")
	}

	query9 := "new york mets at chicago cubs"
	best9, _ := ExtractOne(query9, someEmptyStrings)
	assertMatch(t, query9, someEmptyStrings[1], best9.Match)

	query10 := "a, b"
	choices := []string{query10}
	expectedResult := new(MatchPair)
	expectedResult.Match = query10
	expectedResult.Score = 100
	customScorer = func(s1, s2 string) int {
		return Ratio(s1, s2)
	}
	partialScorer := func(s1, s2 string) int {
		return PartialRatio(s1, s2)
	}

	res, _ := ExtractOne(query10, choices, customScorer)
	partialRes, _ := ExtractOne(query10, choices, partialScorer)
	if *res != *expectedResult {
		t.Error("simple match failed")
	}
	if *partialRes != *expectedResult {
		t.Error("simple partial match failed")
	}
}

func TestDedupe(t *testing.T) {
	sliceWithDupes := []string{"Frodo Baggins", "Tom Sawyer", "Bilbo Baggin", "Samuel L. Jackson", "F. Baggins", "Frody Baggins", "Bilbo Baggins"}
	res, err := Dedupe(sliceWithDupes, nil, nil)
	assert.Nil(t, err)
	if len(res) >= len(sliceWithDupes) {
		t.Error("expecting Dedupe to remove at least one string from slice")
	}

	sliceWithoutDupes := []string{"Tom", "Dick", "Harry"}
	res2, err := Dedupe(sliceWithoutDupes, nil, nil)
	assert.Nil(t, err)
	if len(res2) != len(sliceWithoutDupes) {
		t.Error("not expecting Dedupe to remove any strings from slice")
	}

	lowThreshold := 1
	res3, err := Dedupe(sliceWithDupes, &lowThreshold, nil)
	assert.Nil(t, err)
	if len(res3) != 1 {
		t.Error("expecting low threshold to dedupe all items")
	}

	highThreshold := 99
	res4, err := Dedupe(sliceWithDupes, &highThreshold, nil)
	assert.Nil(t, err)
	if len(res4) != len(sliceWithDupes) {
		t.Error("expecting high threshold to maintain all items")
	}

	threshold := 1
	res5, err := Dedupe(sliceWithDupes, &threshold, func(s1 string, s2 string) int {
		diff := len(s1) - len(s2)
		if diff < 0 {
			diff *= -1
		}
		return diff
	})
	assert.Nil(t, err)
	if len(res5) != 2 {
		t.Error("expecting custom scorer to yield two results")
	}
}

func assertMatch(t *testing.T, query, expectedMatch, actualMatch string) {
	if expectedMatch != actualMatch {
		t.Errorf("expecting [%v] to find match of [%v], actual match was [%v]", query, expectedMatch, actualMatch)
	}
}
