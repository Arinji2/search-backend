package scraper

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func loadStopWords() (map[string]struct{}, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	filePath := fmt.Sprintf("%s/stop-words.txt", dir)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening stop words file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading stop words file: %w", err)
	}

	stopWords := make(map[string]struct{})
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			stopWords[word] = struct{}{}
		}
	}

	return stopWords, nil
}

func loadLemmatizer() (map[string][]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	keyDir := fmt.Sprintf("%s/lemmatization.txt", dir)
	jsonFile, err := os.Open(keyDir)
	if err != nil {
		return nil, fmt.Errorf("error opening lemmatization file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {

		return nil, fmt.Errorf("error reading lemmatization file: %w", err)
	}

	lemmatizer := make(map[string][]string)
	lines := strings.Split(string(byteValue), "\n")

	for _, line := range lines {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		lemmatizer[key] = parts[1:]
	}

	return lemmatizer, nil
}
