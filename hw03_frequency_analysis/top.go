package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
	"unicode"
)

const topLength = 10

var nonWords = map[string]struct{}{
	"—": {},
	"-": {},
	" ": {},
}

func Top10(s string) []string {
	result := []string{}

	if s == "" {
		return result
	}
	type counter struct {
		word  string
		count int
	}

	dict := make(map[string]int)

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && (string(c) != "-")
	}
	words := strings.FieldsFunc(s, f)
	for _, word := range words {
		if !isNonWord(word) {
			continue
		}
		word = strings.ToLower(word)

		if _, ok := dict[word]; !ok {
			dict[word] = 0
		}
		dict[word]++
	}

	wordList := make([]counter, 0, len(dict))
	for word, cnt := range dict {
		c := counter{
			word:  word,
			count: cnt,
		}
		wordList = append(wordList, c)
	}

	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i].count > wordList[j].count
	})

	n := topLength
	if len(wordList) < topLength {
		n = len(wordList)
	}
	wordList = wordList[:n]

	for _, v := range wordList {
		result = append(result, v.word)
	}
	return result
}

func isNonWord(w string) bool {
	_, ok := nonWords[w]
	return !ok
}
