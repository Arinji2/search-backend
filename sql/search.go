package sql

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"

	"github.com/Arinji2/search-backend/types"
)

type SearchResult struct {
	Keyword     types.SQLKeyword
	Page        types.SQLPage
	PageKeyword types.SQLPageKeyword
}

type orderedSearchResult struct {
	score float64
	page  types.SQLPage
}

func search(word string, offsets ...int) ([]SearchResult, error) {
	var result []SearchResult
	offset := 0
	if len(offsets) > 0 {
		offset = offsets[0]
	}
	sql := `SELECT * FROM keywords 
	        JOIN page_keywords ON keywords.id = page_keywords.keyword_id 
			JOIN pages ON page_keywords.page_id = pages.id 
			WHERE keywords.keyword = ?
			ORDER BY page_keywords.frequency DESC
			LIMIT 5
			OFFSET ?;
			`
	rows, err := getDB().Query(sql, word, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var keyword types.SQLKeyword
		var page types.SQLPage
		var pageKeywords types.SQLPageKeyword
		err := rows.Scan(&keyword.ID, &keyword.Keyword, &keyword.DocCount, &keyword.IDF, &pageKeywords.PageID, &pageKeywords.KeywordID, &pageKeywords.WeightedFreq, &page.ID, &page.URL, &page.Title, &page.MetaImage, &page.Description, &page.Favicon, &page.TotalWords)
		if err != nil {
			return result, err
		}

		result = append(result, SearchResult{Keyword: keyword, Page: page, PageKeyword: pageKeywords})
	}

	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func SQLSearch(query string) ([]types.SQLPage, error) {
	words := strings.Split(query, " ")
	if len(words) > 15 {
		return nil, errors.New("too many words")
	}

	workerCount := int(math.Min(float64(len(words)), 5))
	wordChan := make(chan string, workerCount)
	resultChan := make(chan []SearchResult, len(words))
	errorsChan := make(chan error, len(words))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range wordChan {
				result, err := search(word)
				if err != nil {
					errorsChan <- err
					continue
				}
				resultChan <- result
			}
		}()
	}

	for _, word := range words {
		wordChan <- word
	}
	close(wordChan)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorsChan)
	}()

	if len(errorsChan) > 0 {
		for err := range errorsChan {
			fmt.Println(err)
		}

		return nil, errors.New("error while searching")
	}

	orderedResults := make([]orderedSearchResult, len(resultChan))

	for result := range resultChan {
		for _, r := range result {
			score := r.PageKeyword.WeightedFreq * r.Keyword.IDF
			orderedResults = append(orderedResults, orderedSearchResult{score: score, page: r.Page})
		}
	}

	if len(orderedResults) > 0 {
		slices.SortFunc(orderedResults, func(a, b orderedSearchResult) int {
			if a.score < b.score {
				return 1
			} else if a.score > b.score {
				return -1
			}
			return 0
		})
	}
	var result []types.SQLPage

	for _, orderedResult := range orderedResults {
		result = append(result, orderedResult.page)
	}

	return result, nil

}
