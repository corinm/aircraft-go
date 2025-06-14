package main

import (
	"fmt"
	"log"
	"os"

	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Discoverer starting...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tar1090Url := os.Getenv("TAR1090_URL")
	if tar1090Url == "" {
		panic("TAR1090_URL not set")
	}

	f := fetcher.Tar1090AdsbFetcher{
		URL: tar1090Url,
	}

	// TODO: Do this next bit in a loop

	fmt.Println("Fetching aircraft data...")

	_, err2 := f.FetchAircraft()
	if err2 != nil {
		fmt.Println("Error fetching aircraft:", err2)
		return
	}

	fmt.Println("Discoverer finished")
}