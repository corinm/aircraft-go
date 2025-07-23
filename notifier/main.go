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

	"github.com/gregdel/pushover"
	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Print("Notifier service is starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
		panic("Error loading .env file")
	}

	natsHost := os.Getenv("AIRCRAFT_NATS_HOST")
	if natsHost == "" {
		log.Fatal("AIRCRAFT_NATS_HOST environment variable is not set")
		panic("AIRCRAFT_NATS_HOST not set")
	}
	
	natsPort := os.Getenv("AIRCRAFT_NATS_PORT")
	if natsPort == "" {
		log.Fatal("AIRCRAFT_NATS_PORT environment variable is not set")
		panic("AIRCRAFT_NATS_PORT not set")
	}

	natsUrl := natsHost + ":" + natsPort

	pushoverAppToken := os.Getenv("PUSHOVER_APP_TOKEN")
	if pushoverAppToken == "" {
		log.Fatal("PUSHOVER_APP_TOKEN environment variable is not set")
		panic("PUSHOVER_APP_TOKEN not set")
	}

	pushoverUserKey := os.Getenv("PUSHOVER_USER_KEY")
	if pushoverUserKey == "" {
		log.Fatal("PUSHOVER_USER_KEY environment variable is not set")
		panic("PUSHOVER_USER_KEY not set")
	}

	log.Printf("Connecting to NATS at %s...\n", natsUrl)
	m, err := messaging.NewNatsMessaging(natsUrl)
	if err != nil {
		log.Fatal("Error creating messaging client:", err)
		return
	}
	defer m.Close()
	log.Println("Connected to NATS")

	p := pushover.New(pushoverAppToken)
	recipient := pushover.NewRecipient(pushoverUserKey)

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