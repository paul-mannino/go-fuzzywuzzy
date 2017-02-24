package fuzzy

import "testing"

func TestExtractOne(t *testing.T) {
	query1 := "new york mets at atlanta braves"
	best, _ := ExtractOne(query1, games)
	if best.Match != "braves vs mets" {
		t.Fatal()
	}
}
