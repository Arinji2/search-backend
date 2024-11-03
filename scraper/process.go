package scraper

import (
	types "github.com/Arinji2/search-backend/types"
)

func processLemmatization(words []types.WordCount, lemmatizer map[string][]string) []types.WordCount {

	var processedWord []types.WordCount

	for _, wc := range words {

		if _, ok := lemmatizer[wc.Word]; ok {
			processedWord = append(processedWord, types.WordCount{Word: wc.Word, Count: wc.Count})
			continue
		}

		for base, variants := range lemmatizer {
			if contains(variants, wc.Word) {
				processedWord = append(processedWord, types.WordCount{Word: base, Count: wc.Count})
				break
			}
		}

		processedWord = append(processedWord, types.WordCount{Word: wc.Word, Count: wc.Count})

	}

	return processedWord
}

// contains checks if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
