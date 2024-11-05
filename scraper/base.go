package scraper

import "sync"

func StartScrapers() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		crawler := NewScraper("https://dev.to/maxdavis/the-inside-scoop-helping-my-friend-choose-driving-lessons-in-beaconsfield-e54")
		crawler.Start()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		crawler := NewScraper("https://medium.com/@petertou12/bite-sized-tech-automated-reasoning-and-generative-ai-04f6c0878c41")
		crawler.Start()
	}()
	wg.Wait()
}
