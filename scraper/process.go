package scraper

import (
	types "github.com/Arinji2/search-backend/types"
)

func processLemmatization(words []types.ScraperWordCount, lemmatizer map[string][]string) []types.ScraperWordCount {

	var processedWord []types.ScraperWordCount

	for _, wc := range words {

		if _, ok := lemmatizer[wc.Word]; ok {
			processedWord = append(processedWord, types.ScraperWordCount{Word: wc.Word, Count: wc.Count})
			continue
		}

		for base, variants := range lemmatizer {
			if contains(variants, wc.Word) {
				processedWord = append(processedWord, types.ScraperWordCount{Word: base, Count: wc.Count})
				break
			}
		}

		processedWord = append(processedWord, types.ScraperWordCount{Word: wc.Word, Count: wc.Count})

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
