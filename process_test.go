package fuzzy

import "testing"

func TestExtractWithoutOrder(t *testing.T) {
	p := func(s string) string {
		return "f"
	}
	ExtractWithoutOrder("test", []string{"test", "test"}, p)
}
