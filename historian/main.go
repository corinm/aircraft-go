package main

import (
	"context"
	"encoding/json"
	"fmt"
	"historian/data"
	pg "historian/db/sql/generated"
	"historian/messaging"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lpernett/godotenv"
	"github.com/nats-io/nats.go"
)

type Config struct {
	DiscoveryNatsHost string `env:"DISCOVERY_NATS_HOST"`
	DiscoveryNatsPort string `env:"DISCOVERY_NATS_PORT"`
	PostgresHost      string `env:"POSTGRES_HOST"`
	PostgresPort      uint16 `env:"POSTGRES_PORT"`
	PostgresUser      string `env:"POSTGRES_USER"`
	PostgresPassword  string `env:"POSTGRES_PASSWORD"`
	PostgresDB        string `env:"POSTGRES_DB"`
}

func main() {
	log.Println("Historian starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatalln("failed to load .env file:", err)
		panic(err)
	}

	var config Config
	err := env.Parse(&config)
	
	if err != nil {
		log.Fatalln("Error loading configuration:", err)
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

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDB,
	))
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
		return
	}
	defer conn.Close(ctx)

	queries := pg.New(conn)

	aircraft, err := queries.ListAircraft(ctx)
	if err != nil {
		log.Printf("Error listing aircraft: %v", err)
	}

	log.Printf("Found %d aircraft in the database", len(aircraft))
	
	m.Subscribe("aircraft.raw", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Subject)
		log.Printf("Message data: %s\n", string(msg.Data))
		
		aircraft := data.RawAircraft{IcaoHexCode: string(msg.Data)}
		
		log.Printf("Processing aircraft: %s\n", aircraft.IcaoHexCode)
	})
	
	m.Subscribe("aircraft.enriched", func(msg *nats.Msg) {
		log.Println("Received message on subject:", msg.Subject)
		log.Printf("Message data: %s\n", string(msg.Data))

		aircraft := data.EnrichedAircraft{}
		err := json.Unmarshal(msg.Data, &aircraft)
		if err != nil {
			log.Println("Error unmarshalling aircraft data:", err)
			return
		}

		log.Printf("Processing aircraft: %s\n", aircraft.IcaoHexCode)

		err = queries.CreateAircraft(ctx, marshallAircraftForPostgres(aircraft))
		if err != nil {
			log.Printf("Error creating aircraft in database: %v", err)
			return
		}
		log.Printf("Aircraft %s successfully created in the database", aircraft.IcaoHexCode)
	})

	// Catch interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // blocks until a signal is received
	fmt.Println("Shutting down gracefully...")
	m.Drain()
}

func marshallAircraftForPostgres(aircraft data.EnrichedAircraft) pg.CreateAircraftParams {
	cmpg := aircraft.CMPG
	if cmpg == "" {
		cmpg = string(pg.CmpgUnknown)
	}

	return pg.CreateAircraftParams{
		IcaoHexCode: aircraft.IcaoHexCode,
		Registration: pgtype.Text{String: aircraft.Registration, Valid: true},
		Manufacturer: pgtype.Text{String: aircraft.Manufacturer, Valid: true},
		IcaoTypeCode: pgtype.Text{String: aircraft.IcaoTypeCode, Valid: true},
		AircraftType: pgtype.Text{String: aircraft.AircraftType, Valid: true},
		RegisteredOwners: pgtype.Text{String: aircraft.RegisteredOwners, Valid: true},
		IcaoAirlineCode: pgtype.Text{String: aircraft.IcaoAirlineCode, Valid: true},
		Cmpg: pg.Cmpg(cmpg),
		PlaneAlertDbCategory: pgtype.Text{String: aircraft.PlaneAlertDbCategory, Valid: true},
		PlaneAlertDbTags: aircraft.PlaneAlertDbTags,
	}
}