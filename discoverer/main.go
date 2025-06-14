package main

import (
	"log"
	"os"

	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/corinm/aircraft/discovery/pipeline"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Discoverer starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
		panic("Error loading .env file")
	}

	tar1090Url := os.Getenv("TAR1090_URL")
	if tar1090Url == "" {
		log.Fatal("TAR1090_URL environment variable is not set")
		panic("TAR1090_URL not set")
	}

	hexDbUrl := os.Getenv("HEXDB_URL")
	if hexDbUrl == "" {
		log.Fatal("HEXDB_URL environment variable is not set")
		panic("HEXDB_URL not set")
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

	log.Printf("Found %d aircraft\n", len(aircraft))

	enrichers := []pipeline.Enricher{
		&pipeline.HexDbEnricher{HexDbUrl: hexDbUrl},
	}

	for _, a := range aircraft {
		pipeline.EnrichAircraft(&a, enrichers)
		log.Printf("Enriched Aircraft: %+v\n", a)
	}

	log.Println("Discoverer finished")
}
