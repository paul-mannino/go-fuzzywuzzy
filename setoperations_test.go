package fuzzy

import (
	"testing"
)

var equalsTestData = []struct {
	set1  []string // input
	set2  []string
	equal bool // expected result
}{
	{[]string{"alpha", "beta"}[:], []string{"beta", "alpha"}[:], true},
	{[]string{"alpha", "beta", "gamma"}, []string{"beta", "alpha"}, false},
	{[]string{"alpha"}, []string{"beta"}, false},
}

func TestEquals(t *testing.T) {
	for _, testCase := range equalsTestData {
		s1 := NewStringSet(testCase.set1)
		s2 := NewStringSet(testCase.set2)

		if s1.Equals(s2) != testCase.equal {
			t.Fatal()
		}
	}
}

func TestDifference(t *testing.T) {
	s1 := NewStringSet([]string{"ab", "bc", "cd"})
	s2 := NewStringSet([]string{"bc", "cd", "de"})
	expectedDiff := NewStringSet([]string{"ab"})

	actualDiff := s1.Difference(s2)
	if !actualDiff.Equals(expectedDiff) {
		t.Fatal()
	}
}

func TestIntersection(t *testing.T) {
	s1 := NewStringSet([]string{"ab", "bc", "cd"})
	s2 := NewStringSet([]string{"bc", "cd", "de"})
	expectedIntersect := NewStringSet([]string{"bc", "cd"})

	actualIntersect := s1.Intersect(s2)
	if !actualIntersect.Equals(expectedIntersect) {
		t.Fatal()
	}
}
