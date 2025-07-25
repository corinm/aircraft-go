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

	"github.com/caarlos0/env/v11"
	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

type Config struct {
	NatsHost string `env:"AIRCRAFT_NATS_HOST"`
	NatsPort string `env:"AIRCRAFT_NATS_PORT"`
	HexDbUrl          string `env:"HEXDB_URL"`
}

func main() {
	log.Println("Enricher starting...")

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
		panic(err)
	}
	defer m.Close()
	log.Println("Connected to NATS")

	planeAlertDbEnricher, err := pipeline.NewPlaneAlertDbEnricher("./lib/plane-alert-db/plane-alert-db.csv")
	if err != nil {
		log.Fatal("Error creating PlaneAlertDbEnricher:", err)
		panic(err)
	}

	enrichers := []pipeline.Enricher{
		&pipeline.HexDbEnricher{HexDbUrl: config.HexDbUrl},
		planeAlertDbEnricher,
	}

	p := &pipeline.Pipeline{Enrichers: enrichers}

	m.Subscribe("aircraft.raw", func(msg *nats.Msg) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
			defer cancel()

			log.Println("Handling aircraft with hex code:", string(msg.Data))

			aircraft := &data.EnrichedAircraft{IcaoHexCode: string(msg.Data)}

			if errs := p.Enrich(ctx, aircraft); err != nil {
				log.Println("Error enriching aircraft:", errs)
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

			log.Println("Aircraft handled successfully:", aircraft.IcaoHexCode)
		}()
	})

	// Catch interrupt signal to gracefully shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan // blocks until a signal is received
    fmt.Println("Shutting down gracefully...")
    m.Drain()
}
