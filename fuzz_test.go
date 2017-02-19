package fuzzy

import "testing"

var teams = []string{
	"new york mets",                       //0
	"new york mets",                       //1
	"new YORK mets",                       //2
	"the wonderful new york mets",         //3
	"new york mets vs atlanta braves",     //4
	"atlanta braves vs new york mets",     //5
	"new york mets - atlanta braves",      //6
	"new york city mets - atlanta braves", //7
}

func TestRatio(t *testing.T) {
	r1 := Ratio(teams[0], teams[1])
	assertRatioIs100(t, "Ratio", teams[0], teams[1], r1)

	r2 := Ratio(teams[1], teams[2])
	if r2 == 100 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be less than 100. Got %v", teams[1], teams[2], r2)
	}

	r3 := Ratio(Cleanse(teams[1], true), Cleanse(teams[2], true))
	if r3 != 100 {
		t.Errorf("Expected Ratio (cleansed) of '%v' and '%v' to be 100. Got %v", teams[1], teams[2], r3)
	}
}

func TestPartialRatio(t *testing.T) {
	r1 := PartialRatio(teams[1], teams[3])
	assertRatioIs100(t, "PartialRatio", teams[1], teams[3], r1)
}

func TestTokenSortRatio(t *testing.T) {
	r1 := PartialRatio(teams[1], teams[0])
	assertRatioIs100(t, "TokenSortRatio", teams[1], teams[0], r1)
}

func TestPartialTokenSortRatio(t *testing.T) {
	r1 := PartialTokenSortRatio(teams[0], teams[1], false, false)
	assertRatioIs100(t, "PartialTokenSortRatio", teams[0], teams[1], r1)
	r2 := PartialTokenSortRatio(teams[4], teams[5], false, false)
	assertRatioIs100(t, "PartialTokenSortRatio", teams[4], teams[5], r2)
}

func TestTokenSetRatio(t *testing.T) {
	r1 := TokenSetRatio(teams[4], teams[5], false, false)
	assertRatioIs100(t, "TokenSetRatio", teams[4], teams[5], r1)
}

func TestPartialTokenSetRatio(t *testing.T) {
	r1 := PartialTokenSetRatio(teams[4], teams[7], false, false)
	assertRatioIs100(t, "PartialTokenSetRatio", teams[4], teams[7], r1)
}

func TestQuickRatio(t *testing.T) {
	r1 := QRatio(teams[0], teams[1])
	assertRatioIs100(t, "QRatio", teams[0], teams[1], r1)
	r2 := QRatio(teams[0], teams[2])
	assertRatioIs100(t, "QRatio", teams[0], teams[2], r2)
	r3 := QRatio(teams[0], teams[3])
	assertRatioIsNot100(t, "QRatio", teams[0], teams[3], r3)
}

func TestWRatio(t *testing.T) {
	r1 := WRatio(teams[0], teams[1])
	assertRatioIs100(t, "WRatio", teams[0], teams[1], r1)
	r2 := WRatio(teams[0], teams[2])
	assertRatioIs100(t, "WRatio", teams[0], teams[2], r2)
	r3 := WRatio(teams[0], teams[3])
	assertRatio(t, "WRatio", teams[0], teams[3], 90, r3)
	r4 := WRatio(teams[4], teams[5])
	assertRatio(t, "WRatio", teams[4], teams[5], 95, r4)
}

func assertRatio(t *testing.T, methodName, s1, s2 string, expectedRatio, actualRatio int) {
	if actualRatio != expectedRatio {
		t.Errorf("Expected %v of %v and %v to be %v. Got %v", methodName, s1, s2, expectedRatio, actualRatio)
	}
}

func assertRatioIs100(t *testing.T, methodName, s1, s2 string, actualRatio int) {
	if actualRatio != 100 {
		t.Errorf("Expected %v of %v and %v to be 100. Got %v", methodName, s1, s2, actualRatio)
	}
}

func assertRatioIsNot100(t *testing.T, methodName, s1, s2 string, actualRatio int) {
	if actualRatio == 100 {
		t.Errorf("Expected %v of %v and %v to be less than 100. Got %v", methodName, s1, s2, actualRatio)
	}
}
