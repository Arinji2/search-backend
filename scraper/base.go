package scraper

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Arinji2/search-backend/sql"
)

func StartScrapers() {
	indexCount := 10
	indexLinks, err := sql.GetIndexList(indexCount)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	workersCount := 2
	workerChan := make(chan string, workersCount)
	errorChan := make(chan error, indexCount*2)

	var printMu sync.Mutex

	printLog := func(format string, args ...interface{}) {
		printMu.Lock()
		fmt.Printf(format+"\n", args...)
		printMu.Unlock()
	}

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range workerChan {
				startTime := time.Now()
				crawler := NewScraper(url)
				urlTitle, indexedLinks, err := crawler.Start()
				if err != nil {
					errorChan <- fmt.Errorf("error processing %s: %w", url, err)

				}

				if err := sql.DeletePageIndex(url); err != nil {
					errorChan <- fmt.Errorf("error deleting index for %s: %w", url, err)
					continue
				}

				duration := time.Since(startTime)
				printLog("Finished Indexing %s at URL %s in %v. Added %d indexable URLS",
					urlTitle, url, duration, indexedLinks)
			}
		}()
	}

	for _, link := range indexLinks {
		workerChan <- link
	}

	close(workerChan)
	wg.Wait()
	close(errorChan)

	errors := []error{}
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		printLog("\nENCOUNTERED ERRORS WHILST INDEXING:")
		for _, err := range errors {
			printLog("- %v", err)
		}
	}
}
