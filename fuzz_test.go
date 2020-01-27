package fuzzy

import "testing"

var games = []string{
	"new york mets",                       //0
	"new york mets",                       //1
	"new YORK mets",                       //2
	"the wonderful new york mets",         //3
	"new york mets vs atlanta braves",     //4
	"atlanta braves vs new york mets",     //5
	"new york mets - atlanta braves",      //6
	"new york city mets - atlanta braves", //7
}

var alphanumeric = []string{
	"JOHNSMITH6211986", //0
	"JOHNSMITH6201986", //1
}

var nonascii = []string{
	"你貴姓大名？",  //0
	"你叫什麼名字？", //1
}

func TestRatio(t *testing.T) {
	r1 := Ratio(games[0], games[1])
	assertRatioIs100(t, "Ratio", games[0], games[1], r1)

	r2 := Ratio(games[1], games[2])
	if r2 == 100 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be less than 100. Got %v", games[1], games[2], r2)
	}

	r3 := Ratio(Cleanse(games[1], true), Cleanse(games[2], true))
	assertRatioIs100(t, "Ratio (cleansed)", games[1], games[2], r3)

	r4 := Ratio("", "")
	assertRatio(t, "Ratio", "[empty string]", "[empty string]", 0, r4)

	r5 := Ratio(alphanumeric[0], alphanumeric[1])
	assertRatio(t, "Ratio", alphanumeric[0], alphanumeric[1], 94, r5)
}

func TestPartialRatio(t *testing.T) {
	r1 := PartialRatio(games[1], games[3])
	assertRatioIs100(t, "PartialRatio", games[1], games[3], r1)

	r2 := PartialRatio("", "")
	assertRatio(t, "PartialRatio", "[empty string]", "[empty string]", 0, r2)

	s1 := "HSINCHUANG"
	s2 := "SINJHUAN"
	s3 := "LSINJHUANG DISTRIC"
	s4 := "SINJHUANG DISTRICT"
	r3 := PartialRatio(s1, s2)
	if r3 <= 75 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be greater than 75. Got %v", s1, s2, r3)
	}
	r4 := PartialRatio(s1, s3)
	if r4 <= 75 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be greater than 75. Got %v", s1, s3, r4)
	}
	r5 := PartialRatio(s1, s4)
	if r5 <= 75 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be greater than 75. Got %v", s1, s4, r5)
	}

	s5, s6 := "栶eeƵ画-ʏĜ橭畏p父«P^艎鹥ʭ攆", "eeǸɁ碳簫S晑=2#父«厄].稍咾靐Ë"
	r6 := PartialRatio(s5, s6)
	assertRatio(t, "Ratio", s5, s6, 21, r6)
}

func TestTokenSortRatio(t *testing.T) {
	r1 := PartialRatio(games[1], games[0])
	assertRatioIs100(t, "TokenSortRatio", games[1], games[0], r1)
}

func TestPartialTokenSortRatio(t *testing.T) {
	r1 := PartialTokenSortRatio(games[0], games[1], false, false)
	assertRatioIs100(t, "PartialTokenSortRatio", games[0], games[1], r1)
	r2 := PartialTokenSortRatio(games[4], games[5], false, false)
	assertRatioIs100(t, "PartialTokenSortRatio", games[4], games[5], r2)
}

func TestTokenSetRatio(t *testing.T) {
	r1 := TokenSetRatio(games[4], games[5], false, false)
	assertRatioIs100(t, "TokenSetRatio", games[4], games[5], r1)
}

func TestPartialTokenSetRatio(t *testing.T) {
	r1 := PartialTokenSetRatio(games[4], games[7], false, false)
	assertRatioIs100(t, "PartialTokenSetRatio", games[4], games[7], r1)
}

func TestQuickRatio(t *testing.T) {
	r1 := QRatio(games[0], games[1])
	assertRatioIs100(t, "QRatio", games[0], games[1], r1)
	r2 := QRatio(games[0], games[2])
	assertRatioIs100(t, "QRatio", games[0], games[2], r2)
	r3 := QRatio(games[0], games[3])
	assertRatioIsNot100(t, "QRatio", games[0], games[3], r3)

	s1, s2 := "XYZ", "XYZÜ"
	r4 := QRatio(s1, s2)
	assertRatio(t, "QRatio", s1, s2, 100, r4)
}

func TestWRatio(t *testing.T) {
	r1 := WRatio(games[0], games[1])
	assertRatioIs100(t, "WRatio", games[0], games[1], r1)
	r2 := WRatio(games[0], games[2])
	assertRatioIs100(t, "WRatio", games[0], games[2], r2)
	r3 := WRatio(games[0], games[3])
	assertRatio(t, "WRatio", games[0], games[3], 90, r3)
	r4 := WRatio(games[4], games[5])
	assertRatio(t, "WRatio", games[4], games[5], 95, r4)

	r5 := WRatio(nonascii[0], nonascii[1])
	if r5 != 0 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be 0. Got %v", nonascii[0], nonascii[1], r5)
	}
}

func TestUWRatio(t *testing.T) {
	r1 := UWRatio(nonascii[0], nonascii[1])
	if r1 == 0 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be greater than 0. Got 0", nonascii[0], nonascii[1])
	}
}

func TestQRatio(t *testing.T) {
	r1 := QRatio(nonascii[0], nonascii[1])
	if r1 != 0 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be 0. Got %v", nonascii[0], nonascii[1], r1)
	}
}

func TestUQRatio(t *testing.T) {
	r1 := UQRatio(nonascii[0], nonascii[1])
	if r1 == 0 {
		t.Errorf("Expected Ratio of '%v' and '%v' to be greater than 0. Got 0", nonascii[0], nonascii[1])
	}
}

func TestReadmeExamples(t *testing.T) {
	s1 := "coolstring"
	s2 := "coooolstring"
	assertRatio(t, "Ratio", s1, s2, 91, Ratio(s1, s2))

	s1 = "coolstring"
	s2 = "radstring"
	assertRatio(t, "Ratio", s1, s2, 63, Ratio(s1, s2))

	s1 = "needle"
	s2 = "haystackneedelhaystack"
	assertRatio(t, "Ratio", s1, s2, 36, Ratio(s1, s2))
	assertRatio(t, "PartialRatio", s1, s2, 83, PartialRatio(s1, s2))

	s1 = "several tokens arbitrary order"
	s2 = "order arbitrary several tokens"
	assertRatio(t, "Ratio", s1, s2, 50, Ratio(s1, s2))
	assertRatio(t, "TokenSortRatio", s1, s2, 100, TokenSortRatio(s1, s2))
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
