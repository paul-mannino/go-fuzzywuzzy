package fuzzy

type StringSet struct {
	elements map[string]bool
}

func NewStringSet(slice []string) *StringSet {
	sliceStringSet := make(map[string]bool)
	for _, b := range slice {
		sliceStringSet[b] = true
	}
	s := new(StringSet)
	s.elements = sliceStringSet
	return s
}

// Difference returns the set of strings that are present in this set
// but not the other set
func (s *StringSet) Difference(other *StringSet) *StringSet {
	diff := new(StringSet)
	diff.elements = make(map[string]bool)
	for k, v := range s.elements {
		if _, ok := other.elements[k]; !ok {
			diff.elements[k] = v
		}
	}
	return diff
}

// Intersection returns the set of strings that are contained in
// both sets
func (s *StringSet) Intersect(other *StringSet) *StringSet {
	intersection := new(StringSet)
	intersection.elements = make(map[string]bool)
	for k, v := range s.elements {
		if _, ok := other.elements[k]; ok {
			intersection.elements[k] = v
		}
	}
	return intersection
}

// Equals returns true if two sets contain the same elements
func (s *StringSet) Equals(other *StringSet) bool {
	if len(s.elements) != len(other.elements) {
		return false
	}

	for k, _ := range s.elements {
		if _, ok := other.elements[k]; !ok {
			return false
		}
	}
	return true
}

// ToSlice produces a string slice from the set
func (s *StringSet) ToSlice() []string {
	keys := make([]string, len(s.elements))

	i := 0
	for k := range s.elements {
		keys[i] = k
		i++
	}
	return keys
}
