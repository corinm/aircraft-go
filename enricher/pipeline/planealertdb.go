package pipeline

import (
	"context"
	"encoding/csv"
	"enricher/data"
	"log"
	"os"
	"time"
)

type PlaneAlertDbAircraft struct {
	Icao string
	Registration string
	Operator string
	Type string
	IcaoType string
	Cmpg string
	Tag1 string
	Tag2 string
	Tag3 string
	Category string
	Link string
}

type PlaneAlertDbEnricher struct {
	data map[string]PlaneAlertDbAircraft
}

func NewPlaneAlertDbEnricher(pathToCsv string) (*PlaneAlertDbEnricher, error) {
	start := time.Now()
	
	// Note: Reading a 2.2MB file into memory
	// It's done once on start-up so, while blocking, it is probably acceptable
	// Memory usage will scale with CSV size
	file, err := os.Open(pathToCsv)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var records = make(map[string]PlaneAlertDbAircraft)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		
		aircraft := PlaneAlertDbAircraft{
			Icao:         record[0],
			Registration: record[1],
			Operator:     record[2],
			Type:         record[3],
			IcaoType:     record[4],
			Cmpg:         record[5],
			Tag1:         record[6],
			Tag2:         record[7],
			Tag3:         record[8],
			Category:     record[9],
			Link:         record[10],
		}
		records[aircraft.Icao] = aircraft
	}

	log.Printf("Loaded %d aircraft from PlaneAlertDB CSV File in %s\n", len(records), time.Since(start))

	return &PlaneAlertDbEnricher{data: records}, nil
}

func (e *PlaneAlertDbEnricher) Enrich(ctx context.Context, a *data.EnrichedAircraft) error {
	data := e.data[a.IcaoHexCode]
	if data.Icao == "" {
		return nil
	}

	// Fill in any missing fields
	if (a.Registration == "") {
		a.Registration = data.Registration
	}

	if (a.IcaoTypeCode == "") {
		a.IcaoTypeCode = data.IcaoType
	}

	if (a.RegisteredOwners == "") {
		a.RegisteredOwners = data.Operator
	}


	// Add new fields
	a.PlaneAlertDbCategory = data.Category

	if data.Cmpg != "" {
		switch data.Cmpg {
		case "Civ":
			a.CMPG = "Civilian"
		case "Mil":
			a.CMPG = "Military"
		case "Pol":
			a.CMPG = "Police"
		case "Gov":
			a.CMPG = "Government"
		default:
			a.CMPG = "Unknown" // "Unknown" or ""?
		}
	}

	a.PlaneAlertDbTags = []string{data.Tag1, data.Tag2, data.Tag3}

	// Note: There is also a "Link" field available that is currently unused

	return nil
}
