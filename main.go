package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Arinji2/search-backend/scraper"
	"github.com/Arinji2/search-backend/sql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		isProduction := os.Getenv("ENVIRONMENT") == "PRODUCTION"
		if !isProduction {
			log.Fatal("Error loading .env file")
		} else {
			fmt.Println("Using Production Environment")
		}
	} else {
		fmt.Println("Using Development Environment")
	}

	fmt.Println("Running initial scans...")
	scraper.StartScrapers()
	sql.UpdateIDFScores()
	startCronjobs()

	select {}
}

func startCronjobs() {
	fmt.Println("Starting Cronjobs")
	scraperTicker := time.NewTicker(time.Hour * 1)
	go func() {
		for range scraperTicker.C {
			fmt.Println("Starting Scraper")
			scraper.StartScrapers()
		}
	}()

	IDFTicker := time.NewTicker(time.Hour * 5)
	go func() {
		fmt.Println("Starting IDF Ticker")
		for range IDFTicker.C {
			sql.UpdateIDFScores()
		}
	}()

}
