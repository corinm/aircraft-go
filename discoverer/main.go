package main

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/corinm/aircraft/discovery/fetcher"
	"github.com/corinm/aircraft/discovery/messaging"
	"github.com/lpernett/godotenv"
)

type Config struct {
	Tar1090Url 	      string `env:"TAR1090_URL"`
	DiscoveryNatsHost string `env:"DISCOVERY_NATS_HOST"`
	DiscoveryNatsPort string `env:"DISCOVERY_NATS_PORT"`
}

func main() {
	log.Println("Discoverer starting...")

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

	natsUrl := config.DiscoveryNatsHost + ":" + config.DiscoveryNatsPort

	f := fetcher.Tar1090AdsbFetcher{
		URL: config.Tar1090Url,
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
		err := m.Publish("aircraft.raw", []byte(a.IcaoHexCode))
		if err != nil {
			log.Println("Error publishing aircraft:", err)
		}
	}
}
