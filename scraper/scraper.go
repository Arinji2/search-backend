package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Arinji2/search-backend/sql"
	"github.com/Arinji2/search-backend/types"
	"github.com/Arinji2/search-backend/utils"
)

type keyword struct {
	ID   string
	Word string
	Freq int
}
type Scraper struct {
	BaseURL string
}

func NewScraper(baseURL string) *Scraper {
	return &Scraper{BaseURL: baseURL}
}

func (s *Scraper) Start() (string, int, error) {
	if !checkRobots(s.BaseURL) {
		return "", 0, fmt.Errorf("robots.txt disallows scraping of %s", s.BaseURL)
	}

	crawlerClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	doc, err := fetchAndParse(crawlerClient, s.BaseURL)
	if err != nil {
		return "", 0, fmt.Errorf("error fetching page: %w", err)
	}

	stopWords, err := loadStopWords()
	if err != nil {
		return "", 0, fmt.Errorf("error loading stop words: %w", err)
	}

	lemmatizer, err := loadLemmatizer()
	if err != nil {
		return "", 0, fmt.Errorf("error loading lemmatizer: %w", err)
	}

	content, firstH1, firstP := extractContent(doc)
	links := extractLinks(doc)

	wordFreq, totalWords := processWords(content, stopWords)
	topWords := getTopWords(wordFreq, 5)
	finalWords := processLemmatization(topWords, lemmatizer)
	title, description, metaImage, favicon := extractMetaInfo(doc)

	if title == "" {
		title = firstH1
	}

	if description == "" {
		description = firstP
	}

	pageData := types.SQLPage{
		URL:         s.BaseURL,
		Title:       title,
		MetaImage:   metaImage,
		Description: description,
		TotalWords:  totalWords,
		Favicon:     favicon,
	}

	var indexingWg sync.WaitGroup
	keywordsIDChan := make(chan keyword, len(finalWords))
	pageIDChan := make(chan string, 1)

	for _, words := range finalWords {
		indexingWg.Add(1)
		go func(word string, count int) {
			defer indexingWg.Done()
			id, err := keywordIndexer(word)
			if err != nil {
				fmt.Printf("Error indexing keyword %s: %v\n", word, err)
				return
			}
			if id != "" {
				keywordsIDChan <- keyword{ID: id, Word: word, Freq: count}
			}
		}(words.Word, words.Count)
	}

	indexingWg.Add(1)
	go func() {
		defer indexingWg.Done()
		id, err := pageIndexer(pageData)
		if err != nil {
			fmt.Printf("Error indexing page %s: %v\n", pageData.URL, err)
			return
		}
		if id != "" {
			pageIDChan <- id
		}
	}()

	go func() {

		if err := addLinksToIndexList(links); err != nil {
			fmt.Printf("Error adding links to index: %v\n", err)
		}
	}()

	var keywords []keyword
	done := make(chan struct{})
	go func() {
		for keyword := range keywordsIDChan {
			keywords = append(keywords, keyword)
		}
		done <- struct{}{}
	}()

	indexingWg.Wait()
	close(keywordsIDChan)
	<-done

	pageID := <-pageIDChan
	close(pageIDChan)

	if pageID != "" && len(keywords) > 0 {
		if err := linkPageKeywords(keywords, pageID); err != nil {
			return "", 0, fmt.Errorf("error linking page keywords: %w", err)
		}
	}

	return title, len(links), nil
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

	return string(rules) != "User-agent: *\nDisallow: /"
}

func keywordIndexer(keyword string) (string, error) {
	keywordExists, id, err := sql.KeywordExists(keyword)
	if err != nil {
		return "", fmt.Errorf("error checking keyword existence %s: %w", keyword, err)
	}

	if keywordExists {
		existingKeywordData, err := sql.GetKeyword(id)
		if err != nil {
			return "", fmt.Errorf("error getting keyword: %w", err)
		}

		existingKeywordData.DocCount++
		existingKeywordData.IDF, err = utils.CalculateIDF(existingKeywordData.DocCount)
		if err != nil {
			return "", fmt.Errorf("error calculating IDF (existing): %w", err)
		}

		if err := sql.UpdateKeyword(id, existingKeywordData); err != nil {
			return "", fmt.Errorf("error updating keyword: %w", err)
		}
		return id, nil

	}

	idf, err := utils.CalculateIDF(1)
	if err != nil {
		return "", fmt.Errorf("error calculating IDF (new): %w", err)
	}

	keywordData := types.SQLKeyword{
		Keyword:  keyword,
		DocCount: 1,
		IDF:      idf,
	}

	id, err = sql.CreateKeyword(keywordData)
	if err != nil {
		return "", fmt.Errorf("error creating keyword %s: %w", keyword, err)
	}
	return id, nil

}

func pageIndexer(pageData types.SQLPage) (string, error) {
	pageExists, existingID, err := sql.PageExists(pageData.URL)
	if err != nil {
		return "", fmt.Errorf("error checking page existence %s: %w", pageData.URL, err)
	}

	if pageExists {

		if err := sql.DeletePageIndex(pageData.URL); err != nil {
			return "", fmt.Errorf("error deleting page index: %w", err)
		}
		return existingID, nil

	}

	id, err := sql.CreatePage(pageData)
	if err != nil {
		return "", fmt.Errorf("error creating page: %w", err)
	}

	if err := sql.DeletePageIndex(pageData.URL); err != nil {
		return "", fmt.Errorf("error deleting page index: %w", err)
	}
	return id, nil

}

func addLinksToIndexList(links []string) error {

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	errorChan := make(chan error, len(links))

	for _, link := range links {
		if !strings.HasPrefix(link, "https") {
			continue
		}

		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := sql.AddIndexList(link); err != nil {
				errorChan <- fmt.Errorf("error adding link %s to index: %w", link, err)
			}
		}(link)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	var errs []error
	for err := range errorChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func linkPageKeywords(keywords []keyword, pageID string) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(keywords))
	semaphore := make(chan struct{}, 10)

	for _, k := range keywords {
		wg.Add(1)
		go func(id string, freq int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := sql.LinkPageKeyword(id, pageID, freq); err != nil {
				errorChan <- fmt.Errorf("error linking page %s to keyword %s: %w", pageID, id, err)
			}
		}(k.ID, k.Freq)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	var errs []error
	for err := range errorChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
