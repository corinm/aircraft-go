package main

import (
	"log"
	"os"

	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/corinm/aircraft/discovery/messaging"
	"github.com/corinm/aircraft/discovery/pipeline"
	"github.com/lpernett/godotenv"
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

	natsHost := os.Getenv("DISCOVERY_NATS_HOST")
	if natsHost == "" {
		log.Fatal("DISCOVERY_NATS_HOST environment variable is not set")
		panic("DISCOVERY_NATS_HOST not set")
	}
	
	natsPort := os.Getenv("DISCOVERY_NATS_PORT")
	if natsPort == "" {
		log.Fatal("DISCOVERY_NATS_PORT environment variable is not set")
		panic("DISCOVERY_NATS_PORT not set")
	}

	natsUrl := natsHost + ":" + natsPort

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

	log.Printf("Connecting to NATS at %s\n", natsUrl)

	m, err3 := messaging.NewNatsMessaging(natsUrl)
	if err3 != nil {
		log.Fatal("Error creating messaging client:", err3)
		return
	}

	defer m.Close()
	log.Println("Connected to NATS")

	for _, a := range aircraft {
		pipeline.EnrichAircraft(&a, enrichers)
		log.Printf("Enriched Aircraft: %+v\n", a)
		err4 := m.Publish("aircraft", []byte(a.AiocHexCode))
		if err4 != nil {
			log.Println("Error publishing aircraft:", err3)
		}
	}

	log.Println("Discoverer finished")
}
