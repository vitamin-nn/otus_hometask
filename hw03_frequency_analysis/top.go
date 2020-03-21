package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
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
	dict := make(map[string]counter)

	//words := strings.Split(s, " ")
	var r = regexp.MustCompile(`[\s\.,;"\!]+`)
	for _, word := range r.Split(s, -1) {
		if !isNonWord(word) {
			continue
		}
		word = strings.ToLower(word)
		v, ok := dict[word]
		if !ok {
			v = counter{word: word, count: 0}
		}
		v.count++
		dict[word] = v
	}

	wordList := make([]counter, 0, len(dict))
	for _, v := range dict {
		wordList = append(wordList, v)
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
	result := true
	if _, ok := nonWords[w]; ok {
		result = false
	}
	return result
}
