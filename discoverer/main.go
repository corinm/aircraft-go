package main

import (
	"log"
	"os"

	"github.com/corinm/aircraft/discovery/enricher"
	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Discoverer starting...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}

	tar1090Url := os.Getenv("TAR1090_URL")
	if tar1090Url == "" {
		log.Fatal("TAR1090_URL environment variable is not set")
		panic("TAR1090_URL not set")
	}

	f := fetcher.Tar1090AdsbFetcher{
		URL: tar1090Url,
	}

	// TODO: Do this next bit in a loop

	log.Println("Fetching aircraft data...")

	aircraft, err2 := f.FetchAircraft()
	if err2 != nil {
		log.Println("Error fetching aircraft:", err2)
		return
	}

	for _, a := range aircraft {
		a2, _ := enricher.EnrichAircraft(a)
		log.Printf("Enriched Aircraft: %+v\n", a2)
	}

	log.Println("Discoverer finished")
}