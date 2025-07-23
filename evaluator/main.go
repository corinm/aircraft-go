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

	"github.com/caarlos0/env/v11"
	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

type Config struct {
	NatsHost string `env:"AIRCRAFT_NATS_HOST"`
	NatsPort string `env:"AIRCRAFT_NATS_PORT"`
}

func main() {
	log.Println("Evaluator starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file: %w", err)
		panic(err)
	}

	var config Config
	err := env.Parse(&config)

	if err != nil {
		log.Fatal("Error loading configuration:", err)
		panic(err)
	}

	natsUrl := config.NatsHost + ":" + config.NatsPort

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
		isInteresting := evaluateAircraft(aircraft)

		if !isInteresting {
			return
		}

		log.Printf("Aircraft %s meets criteria, republishing...\n", aircraft.IcaoHexCode)

		if err := m.Publish("aircraft.interesting", msg.Data); err != nil {
			log.Println("Error republishing aircraft:", err)
		}
	})

	// Catch interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // blocks until a signal is received
	fmt.Println("Shutting down gracefully...")
	m.Drain()
}

func evaluateAircraft(aircraft data.EnrichedAircraft) bool {
	return aircraft.PlaneAlertDbCategory != ""
}
