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

	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
		panic(err)
	}

	natsUrl := config.DiscoveryNatsHost + ":" + config.DiscoveryNatsPort

	log.Printf("Connecting to NATS at %s...\n", natsUrl)
	m, err := messaging.NewNatsMessaging(natsUrl)
	if err != nil {
		log.Fatal("Error creating messaging client:", err)
		return
	}
	defer m.Close()
	log.Println("Connected to NATS")
	
	enrichers := []pipeline.Enricher{
		&pipeline.HexDbEnricher{HexDbUrl: config.HexDbUrl},
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

		// Marshal aircraft to JSON
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

type Config struct {
	DiscoveryNatsHost string `env:"DISCOVERY_NATS_HOST"`
	DiscoveryNatsPort string `env:"DISCOVERY_NATS_PORT"`
	HexDbUrl          string `env:"HEXDB_URL"`
}

func loadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	var config Config

	config.DiscoveryNatsHost = os.Getenv("DISCOVERY_NATS_HOST")
	if config.DiscoveryNatsHost == "" {
		return Config{}, fmt.Errorf("DISCOVERY_NATS_HOST environment variable is not set")
	}
	config.DiscoveryNatsPort = os.Getenv("DISCOVERY_NATS_PORT")
	if config.DiscoveryNatsPort == "" {
		return Config{}, fmt.Errorf("DISCOVERY_NATS_PORT environment variable is not set")
	}
	config.HexDbUrl = os.Getenv("HEXDB_URL")
	if config.HexDbUrl == "" {
		return Config{}, fmt.Errorf("HEXDB_URL environment variable is not set")
	}
	
	return config, nil
}