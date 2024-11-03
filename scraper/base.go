package scraper

func StartScrapers() {
	crawler1 := NewScraper("https://dev.to/maxdavis/the-inside-scoop-helping-my-friend-choose-driving-lessons-in-beaconsfield-e54")
	crawler1.Start()

	crawler2 := NewScraper("https://www.capitalone.com/tech/machine-learning/understanding-tf-idf/")
	crawler2.Start()
}
