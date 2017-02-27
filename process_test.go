package fuzzy

import (
	"testing"
)

var baseballStrings = []string{
	"new york mets vs chicago cubs",
	"chicago cubs vs chicago white sox",
	"philladelphia phillies vs atlanta braves",
	"braves vs mets",
}

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
}

func assertMatch(t *testing.T, query, expectedMatch, actualMatch string) {
	if expectedMatch != actualMatch {
		t.Errorf("Expecting [%v] to find match of [%v]. Actual match: [%v]", query, expectedMatch, actualMatch)
	}
}
