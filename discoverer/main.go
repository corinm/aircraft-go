package main

import (
	"log"
	"os"
	"time"

	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/corinm/aircraft/discovery/messaging"
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

	log.Printf("Connecting to NATS at %s...\n", natsUrl)
	m, err := messaging.NewNatsMessaging(natsUrl)
	if err != nil {
		log.Fatal("Error creating messaging client:", err)
		return
	}
	defer m.Close()
	log.Println("Connected to NATS")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fetchAndPublishAircraft(f, m)
			case <-done:
				log.Println("Stopping fetcher...")
				return
			}
		}
	}()

	select {}
}

func fetchAndPublishAircraft(f fetcher.Tar1090AdsbFetcher, m *messaging.NatsMessaging) {
	log.Println("Fetching aircraft data...")

	aircraft, err := f.FetchAircraft()
	if err != nil {
		log.Println("Error fetching aircraft:", err)
		return
	}

	log.Printf("Found %d aircraft\n", len(aircraft))

	for _, a := range aircraft {
		err := m.Publish("aircraft.raw", []byte(a.AiocHexCode))
		if err != nil {
			log.Println("Error publishing aircraft:", err)
		}
	}
}
