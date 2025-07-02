package main

import (
	"context"
	"encoding/json"
	"enricher/data"
	"enricher/messaging"
	"enricher/pipeline"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		panic(err)
	}
	defer m.Close()
	log.Println("Connected to NATS")

	enrichers := []pipeline.Enricher{
		&pipeline.HexDbEnricher{HexDbUrl: config.HexDbUrl},
	}

	p := &pipeline.Pipeline{Enrichers: enrichers}

	m.Subscribe("aircraft.raw", func(msg *nats.Msg) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
			defer cancel()

			log.Println("Handling aircraft with hex code:", string(msg.Data))

			aircraft := &data.EnrichedAircraft{AiocHexCode: string(msg.Data)}

			if err := p.Enrich(ctx, aircraft); err != nil {
				log.Println("Error enriching aircraft:", err)
				return
			}

			aircraftData, err := json.Marshal(aircraft)
			if err != nil {
				log.Println("Error marshalling aircraft to JSON:", err)
				return
			}

			if err := m.Publish("aircraft.enriched", aircraftData); err != nil {
				log.Println("Error publishing aircraft:", err)
				return
			}

			log.Println("Aircraft handled successfully:", aircraft.AiocHexCode)
		}()
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