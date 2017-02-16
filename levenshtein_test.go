package fuzzy

import (
	"testing"
)

var levEditDistanceTestData = [][]interface{}{
	{"one", "", 3},
	{"four", "tour", 1},
	{"cupid", "pulpit", 3},
	{"no", "y", 2},
	{"nowish", "n", 5},
}

func TestLevEditDistance(t *testing.T) {
	for _, test := range levEditDistanceTestData {
		actual := LevEditDistance(test[0].(string), test[1].(string), 0)
		expected := test[2]
		if actual != test[2] {
			t.Errorf("Edit distance from %v to %v is %d; got %d.",
				test[0], test[1], expected, actual)
		}
	}
}
