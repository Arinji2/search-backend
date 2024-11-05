package scraper

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/Arinji2/search-backend/types"
	"github.com/gertd/go-pluralize"
	"golang.org/x/net/html"
)

func fetchAndParse(client *http.Client, url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return html.Parse(resp.Body)
}

func extractContent(n *html.Node) ([]string, string, string) {
	var content []string

	var firstH1 string
	var firstP string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				text := extractText(n)
				if text != "" {
					if firstH1 == "" {
						firstH1 = text
					}
					content = append(content, text)
				}
			case "p":
				text := extractText(n)
				if text != "" {
					if firstP == "" {
						firstP = text
					}
					content = append(content, text)
				}
			case "h2", "h3", "h4", "h5", "h6":
				text := extractText(n)
				if text != "" {
					content = append(content, text)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return content, firstH1, firstP
}

func extractText(n *html.Node) string {
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			sb.WriteString(c.Data)
		}
	}
	return strings.TrimSpace(sb.String())
}

func processWords(sentences []string, stopWords map[string]struct{}) (map[string]int, int) {
	wordFreq := make(map[string]int)
	totalWords := 0
	pluralizer := pluralize.NewClient()

	for _, sentence := range sentences {
		fields := strings.Fields(strings.ToLower(sentence))
		for _, word := range fields {

			if _, isStopWord := stopWords[word]; isStopWord {
				continue
			}

			word = strings.Trim(word, ".,!?\"';:()")
			if word == "" {
				continue
			}

			if pluralizer.IsPlural(word) {
				word = pluralizer.Singular(word)
			}

			wordFreq[word]++
			totalWords++
		}
	}

	return wordFreq, totalWords
}

func getTopWords(wordFreq map[string]int, n int) []types.ScraperWordCount {

	pairs := make([]types.ScraperWordCount, 0, len(wordFreq))
	for word, count := range wordFreq {
		pairs = append(pairs, types.ScraperWordCount{Word: word, Count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].Word < pairs[j].Word
		}
		return pairs[i].Count > pairs[j].Count
	})

	if len(pairs) > n {
		pairs = pairs[:n]
	}
	return pairs
}
