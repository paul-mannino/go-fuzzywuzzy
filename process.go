package fuzzy

import (
	"errors"
	"fmt"
	"sort"
)

type MatchPair struct {
	Match string
	Score int
}

type MatchPairs []*MatchPair

func (slice MatchPairs) Len() int {
	return len(slice)
}

func (slice MatchPairs) Less(i, j int) bool {
	return slice[i].Score > slice[j].Score
}

func (slice MatchPairs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type alphaLengthSortPairs []*MatchPair

func (slice alphaLengthSortPairs) Len() int {
	return len(slice)
}

func (slice alphaLengthSortPairs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice alphaLengthSortPairs) Less(i, j int) bool {
	if len(slice[i].Match) != len(slice[j].Match) {
		return len(slice[i].Match) > len(slice[j].Match)
	}
	return slice[i].Match < slice[j].Match
}

type optionalArgs struct {
	processor    func(string) string
	scorer       func(string, string) int
	scoreCutoff  int
	processorSet bool
	scorerSet    bool
	cutoffSet    bool
}

func ExtractWithoutOrder(query string, choices []string, args ...interface{}) (MatchPairs, error) {
	noProcess := func(s string) string {
		return s
	}
	if len(args) > 3 {
		return nil, errors.New("Expecting 3 or fewer optional parameters")
	}
	processor := func(s string) string {
		return Cleanse(s, false)
	}
	scorer := func(s1, s2 string) int {
		return WRatio(s1, s2)
	}
	scoreCutoff := 0

	opts, err := parseArgs(args...)
	if err != nil {
		return nil, err
	}

	if opts.scorerSet {
		scorer = opts.scorer
	}
	if opts.cutoffSet {
		scoreCutoff = opts.scoreCutoff
	}
	if opts.processorSet {
		processor = opts.processor
	}

	processedQuery := processor(query)

	if !opts.scorerSet && !opts.processorSet {
		// if using default scorer, set active processor to noProcess
		// to avoid processing twice
		processor = noProcess
	}

	results := MatchPairs{}
	for _, choice := range choices {
		processedChoice := processor(choice)
		score := scorer(processedQuery, processedChoice)
		if score >= scoreCutoff {
			pair := &MatchPair{Match: choice, Score: score}
			results = append(results, pair)
		}
	}
	return results, nil
}

func parseArgs(args ...interface{}) (*optionalArgs, error) {
	var processor func(string) string
	var scorer func(string, string) int
	var scoreCutoff int
	processorSet, scorerSet, cutoffSet := false, false, false

	for _, arg := range args {
		switch v := arg.(type) {
		default:
			return nil, fmt.Errorf("not expecting optional argument of type %T", v)
		case func(string) string:
			if processorSet {
				return nil, errors.New("expecting only one processing function of the form f(string)->string")
			}
			processor = arg.(func(string) string)
			processorSet = true
		case func(string, string) int:
			if scorerSet {
				return nil, errors.New("expecting only one scoring function of the form f(string,string)->int")
			}
			scorer = arg.(func(string, string) int)
			scorerSet = true
		case int:
			if cutoffSet {
				return nil, errors.New("expecting only one integer scoring cutoff")
			}
			scoreCutoff = arg.(int)
			cutoffSet = true
		}

	}
	opts := new(optionalArgs)
	opts.processor = processor
	opts.scorer = scorer
	opts.scoreCutoff = scoreCutoff
	opts.cutoffSet = cutoffSet
	opts.scorerSet = scorerSet
	opts.processorSet = processorSet
	return opts, nil
}

func Extract(query string, choices []string, limit int, args ...interface{}) (MatchPairs, error) {
	pairs, err := ExtractWithoutOrder(query, choices, args...)
	if err != nil {
		return nil, err
	}
	largestKPairs := largestKMatchPairs(pairs, limit)
	return largestKPairs, nil
}

func ExtractOne(query string, choices []string, args ...interface{}) (*MatchPair, error) {
	matches, err := ExtractWithoutOrder(query, choices, args...)
	if err != nil {
		return nil, err
	}
	bestPair, err := bestScoreMatchPair(matches)
	if err != nil {
		return nil, err
	}
	return bestPair, nil
}

func bestScoreMatchPair(pairs MatchPairs) (*MatchPair, error) {
	bestPair := &MatchPair{}
	bestScore := -1
	for _, pair := range pairs {
		if pair.Score > bestScore {
			bestPair = pair
			bestScore = pair.Score
		}
	}
	if bestScore < 0 {
		return nil, errors.New("no matches found between query and provided choices")
	}
	return bestPair, nil
}

func largestKMatchPairs(pairs MatchPairs, k int) MatchPairs {
	//todo: implement priority queue algorithm
	sort.Sort(pairs)
	n := len(pairs)
	if k > n {
		k = n
	}

	if k > 0 {
		return pairs[:k]
	} else if k < 0 {
		return pairs
	}
	return make(MatchPairs, 0)
}

var defaultThreshold = 70

func Dedupe(sliceWithDupes []string, threshold *int, scorer func(string, string) int) ([]string, error) {
	if scorer == nil {
		scorer = func(s1, s2 string) int {
			return TokenSetRatio(s1, s2, true, true)
		}
	}
	if threshold == nil {
		threshold = &defaultThreshold
	}

	extracted := []string{}
	for _, elem := range sliceWithDupes {
		matches, err := Extract(elem, sliceWithDupes, -1, scorer)
		if err != nil {
			return nil, err
		}
		filtered := MatchPairs{}
		for _, m := range matches {
			if m.Score > *threshold {
				filtered = append(filtered, m)
			}
		}
		if len(filtered) == 1 {
			extracted = append(extracted, filtered[0].Match)
		} else if len(filtered) > 0 {
			altPoints := alphaLengthSortPairs(filtered)
			sort.Sort(altPoints)
			extracted = append(extracted, altPoints[0].Match)
		}
	}
	set := NewStringSet(extracted)
	// dedupe extracted slice
	extracted = set.ToSlice()

	if len(extracted) == len(sliceWithDupes) {
		return sliceWithDupes, nil
	}
	return extracted, nil
}
