package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const countWords = 10

var regexpEndSymbols = regexp.MustCompile(`^[[:punct:]]*|[[:punct:]]*?$`)

type Rating struct {
	Word  string
	Count uint32
}

type RatingList []Rating

func Top10(text string) []string {
	if text == "" {
		return []string{}
	}

	stringToLower := strings.ToLower(text)

	var result []string

	rate := make(map[string]uint32)

	for _, match := range strings.Fields(stringToLower) {
		if len(match) > 0 {
			if len(match) > 0 {
				match = regexpEndSymbols.ReplaceAllString(match, "")
			}

			if len(match) > 0 {
				rate[match]++
			}
		}
	}

	rez := topRatingWord(rate)

	for ind, Rating := range rez {
		if ind < countWords {
			result = append(result, Rating.Word)
		}
	}

	return result
}

func topRatingWord(wordFrequencies map[string]uint32) RatingList {
	rl := make(RatingList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		rl[i] = Rating{k, v}
		i++
	}
	sort.Sort(rl)

	return rl
}

// Implementing the sorting interface

func (rl RatingList) Len() int {
	return len(rl)
}

func (rl RatingList) Less(i, j int) bool {
	if rl[i].Count == rl[j].Count {
		return rl[i].Word < rl[j].Word
	}
	return rl[i].Count > rl[j].Count
}

func (rl RatingList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
