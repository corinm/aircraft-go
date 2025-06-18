package main

import (
	"enricher/messaging"
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

	m, err3 := messaging.NewNatsMessaging(natsUrl)
	if err3 != nil {
		log.Fatal("Error creating messaging client:", err3)
		return
	}

	defer m.Close()
	log.Println("Connected to NATS")

	m.Subscribe("aircraft", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Data)
	})

	// Catch interrupt signal to gracefully shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan // blocks until a signal is received
    fmt.Println("Shutting down gracefully...")
    m.Drain()
}