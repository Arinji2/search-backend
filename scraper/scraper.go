package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Scraper struct {
	BaseURL string
}

func NewScraper(baseURL string) *Scraper {
	return &Scraper{BaseURL: baseURL}
}

func (s *Scraper) Start() {

	if !checkRobots(s.BaseURL) {
		return
	}

	crawlerClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	doc, err := fetchAndParse(crawlerClient, s.BaseURL)
	if err != nil {
		fmt.Printf("Error fetching page: %v\n", err)
		return
	}

	stopWords, err := loadStopWords()
	if err != nil {
		fmt.Printf("Error loading stop words: %v\n", err)
		return
	}

	lemmatizer, err := loadLemmatizer()
	if err != nil {
		fmt.Printf("Error loading lemmatizer: %v\n", err)
		return
	}

	content := extractContent(doc)
	links := extractLinks(doc)

	for _, link := range links {
		if !strings.HasPrefix(link, "http") {
			link = fmt.Sprintf("%s%s", s.BaseURL, link)
		}
		fmt.Printf("Link: %s\n", link)
	}

	wordFreq := processWords(content, stopWords)

	topWords := getTopWords(wordFreq, 5)

	finalWords := processLemmatization(topWords, lemmatizer)

	for _, wordData := range finalWords {
		fmt.Printf("%s (count: %d)\n", wordData.Word, wordData.Count)
	}

}

func checkRobots(siteURL string) bool {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	parts := strings.Split(siteURL, "/")
	baseURL := strings.TrimSpace((strings.Join(parts[1:3], "")))

	urlBuilder := url.URL{Scheme: "https", Host: baseURL, Path: "/robots.txt"}
	req, err := http.NewRequest("GET", urlBuilder.String(), nil)
	if err != nil {
		return true
	}

	resp, err := client.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true
	}

	rules, err := io.ReadAll(resp.Body)
	if err != nil {
		return true
	}

	if string(rules) == "User-agent: *\nDisallow: /" {
		return false
	}
	return true

}
