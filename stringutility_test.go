package fuzzy

import (
	"testing"
)

var asciiOnlyData = [][]interface{}{
	{"one", "one"},
	{"Țwo", "wo"},
	{"ǩƱ©", ""},
}

var cleanseData = [][]interface{}{
	{"  OnE ", "one", "one"},
	{"Țw o", "țw o", "w o"},
	{"ǩƱ©", "ǩʊ ", ""},
}

func TestASCIIOnly(t *testing.T) {
	for _, testCase := range asciiOnlyData {
		actual := ASCIIOnly(testCase[0].(string))
		expected := testCase[1]

		if actual != expected {
			t.Errorf("Ascii-only %v: Expected %v, got %v.",
				testCase[0], expected, actual)
		}
	}
}

func TestCleanse(t *testing.T) {
	for _, testCase := range cleanseData {
		actual := Cleanse(testCase[0].(string), false)
		expected := testCase[1]
		if actual != expected {
			t.Errorf("Cleanse %v: Expected %v, got %v.",
				testCase[0], expected, actual)
		}
	}

	for _, testCase := range cleanseData {
		actual := Cleanse(testCase[0].(string), true)
		expected := testCase[2]
		if actual != expected {
			t.Errorf("Cleanse %v: Expected %v, got %v.",
				testCase[0], expected, actual)
		}
	}
}
