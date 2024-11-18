package fuzzy

import "testing"

var levEditDistanceTestData = [][]interface{}{
	{"one", "", 3},
	{"", "qwertuiop", 9},
	{"four", "tour", 1},
	{"cupid", "pulpit", 3},
	{"no", "y", 2},
	{"nowish", "n", 5},
	{"JOHNSMITH6211986", "JOHNSMITH6201986", 1},
	{"你好", "你好，世界", 3},
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
