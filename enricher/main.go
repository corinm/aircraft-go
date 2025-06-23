package main

import (
	"encoding/json"
	"enricher/data"
	"enricher/messaging"
	"enricher/pipeline"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Enricher starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
		panic("Error loading .env file")
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

	hexDbUrl := os.Getenv("HEXDB_URL")
	if hexDbUrl == "" {
		log.Fatal("HEXDB_URL environment variable is not set")
		panic("HEXDB_URL not set")
	}

	m, err := messaging.NewNatsMessaging(natsUrl)
	if err != nil {
		log.Fatal("Error creating messaging client:", err)
		return
	}
	defer m.Close()
	log.Println("Connected to NATS")
	
	enrichers := []pipeline.Enricher{
		&pipeline.HexDbEnricher{HexDbUrl: hexDbUrl},
	}

	m.Subscribe("aircraft.raw", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Subject)
		// Process message, enrich aircraft, republish to enriched subject
		aircraft := &data.EnrichedAircraft{AiocHexCode: string(msg.Data)}

		err := pipeline.EnrichAircraft(aircraft, enrichers)
		if err != nil {
			log.Println("Error enriching aircraft:", err)
			return
		}
		log.Println("Enriched aircraft successfully, republishing...")

		// Marshall aircraft to JSON
		aircraftData, err := json.Marshal(aircraft)
		if err != nil {
			log.Println("Error marshalling enriched aircraft to JSON:", err)
			return
		}

		if err := m.Publish("aircraft.enriched", aircraftData); err != nil {
			log.Println("Error republishing enriched aircraft:", err)
			return
		}
		log.Println("Enriched aircraft republished successfully")
	})

	// Catch interrupt signal to gracefully shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan // blocks until a signal is received
    fmt.Println("Shutting down gracefully...")
    m.Drain()
}