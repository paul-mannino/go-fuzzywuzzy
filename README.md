# go-fuzzywuzzy
This is a port of SeatGeek's [fuzzywuzzy](https://github.com/seatgeek/fuzzywuzzy), a fuzzy string matching library. Also included is a port of [python-levenshtein](https://github.com/miohtama/python-Levenshtein), a wicked-fast implementation of Levenshtein edit distance.
##Usage 
####Levenshtein Edit Distance
```go
fuzzy.EditDistance("bart", "bort")
1
```
####Simple Ratio
```go
fuzzy.Ratio("coolstring", "coooolstring")
91
fuzzy.Ratio("coolstring", "radstring"))
63
```
####Partial Ratio
```go
fuzzy.Ratio("needle", "haystackneedelhaystack")
36
fuzzy.PartialRatio("needle", "haystackneedelhaystack")
83
```
####Token Sort Ratio
```go
fuzzy.Ratio("several tokens arbitrary order", "order arbitrary several tokens")
50
fuzzy.TokenSortRatio("several tokens arbitrary order", "order arbitrary several tokens")
100
```
####Token Set Ratio
```go
fuzzy.TokenSortRatio("several tokens arbitrary order", "order order arbitrary several tokens")
91
fuzzy.TokenSetRatio("several tokens arbitrary order", "order order arbitrary several tokens")
100
```
####Process
```go
choices := []string{"Wayne Shorter", "Jonathan Richman", "Wayne Hancock", "Kate Bush"}
fuzzy.ExtractOne("wayne hancock", choices)
{Match:"Wayne Hancock", Score:100}
fuzzy.Extract("wayne hancock", choices, 2)
[{Match:"Wayne Hancock", Score:100}, {Match:"Wayne Shorter", Score:62}]
```
