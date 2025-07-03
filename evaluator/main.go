package main

import (
	"encoding/json"
	"evaluator/data"
	"evaluator/messaging"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Evaluator starting...")

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

	m.Subscribe("aircraft.enriched", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Subject)

		aircraft := data.EnrichedAircraft{}
		err := json.Unmarshal(msg.Data, &aircraft)
		if err != nil {
			log.Println("Error unmarshalling aircraft data:", err)
			return
		}

		log.Printf("Evaluating aircraft: %s\n", aircraft.IcaoHexCode)
		shouldNotify, err := evaluateAircraft(aircraft)
		if err != nil {
			log.Println("Error evaluating aircraft:", err)
			return
		}

		if !shouldNotify {
			log.Printf("Aircraft %s does not meet criteria, no notification sent.\n", aircraft.IcaoHexCode)
			return
		}

		log.Printf("Aircraft %s meets criteria, sending notification...\n", aircraft.IcaoHexCode)
		// Print the response if you want
		if err := m.Publish("aircraft.interesting", msg.Data); err != nil {
			log.Println("Error republishing enriched aircraft:", err)
		}
	})

	// Catch interrupt signal to gracefully shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan // blocks until a signal is received
    fmt.Println("Shutting down gracefully...")
    m.Drain()
}

func evaluateAircraft(aircraft data.EnrichedAircraft) (shouldNotify bool, err error) {
	return false, nil
}

type Config struct {
	DiscoveryNatsHost string `env:"DISCOVERY_NATS_HOST"`
	DiscoveryNatsPort string `env:"DISCOVERY_NATS_PORT"`
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
	
	return config, nil
}
