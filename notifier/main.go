package main

import (
	"encoding/json"
	"fmt"
	"log"
	"notifier/data"
	"notifier/messaging"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/gregdel/pushover"
	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

type Config struct {
	NatsHost string `env:"AIRCRAFT_NATS_HOST"`
	NatsPort string `env:"AIRCRAFT_NATS_PORT"`
	PushoverAppToken string `env:"PUSHOVER_APP_TOKEN"`
	PushoverUserKey string `env:"PUSHOVER_USER_KEY"`
}

func main() {
	log.Print("Notifier service is starting...")

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

	p := pushover.New(config.PushoverAppToken)
	recipient := pushover.NewRecipient(config.PushoverUserKey)

	m.Subscribe("aircraft.interesting", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Subject)

		aircraft := data.EnrichedAircraft{}
		err := json.Unmarshal(msg.Data, &aircraft)
		if err != nil {
			log.Println("Error unmarshalling aircraft data:", err)
			return
		}

		message := pushover.NewMessage(formatAircraft(aircraft))

		response, err := p.SendMessage(message, recipient)
		if err != nil {
			log.Println("Error sending message:", err)
		}

		log.Println("Message sent successfully:", response)
	})

	// Catch interrupt signal to gracefully shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan // blocks until a signal is received
    fmt.Println("Shutting down gracefully...")
    m.Drain()
}

func formatAircraft(aircraft data.EnrichedAircraft) string {
	return fmt.Sprintf("Interesting aircraft: %s (%s) %s %s %s", aircraft.Registration, aircraft.IcaoHexCode, aircraft.RegisteredOwners, aircraft.Manufacturer, aircraft.AircraftType)
}