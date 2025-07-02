package pipeline

import (
	"context"
	"encoding/csv"
	"enricher/data"
	"fmt"
	"log"
	"os"
)

type PlaneAlertDbEnricher struct {

}

func NewPlaneAlertDbEnricher(pathToCsv string) (*PlaneAlertDbEnricher, error) {
	// Note: Reading a 2.2MB file into memory
	// It's done once on start-up so, while blocking, it may be acceptable
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

	var aircraftList []PlaneAlertDbAircraft

	// Read the entire CSV content into memory
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		
		// TODO: Put this into a map of ICAO -> PlaneAlertDbAircraft
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
		aircraftList = append(aircraftList, aircraft)
	}

	log.Printf("Loaded %d aircraft from PlaneAlertDB CSV File\n", len(aircraftList))
	fmt.Printf("First aircraft: %+v\n", aircraftList[0])

	return &PlaneAlertDbEnricher{}, nil
}

func (e *PlaneAlertDbEnricher) Enrich(ctx context.Context, a *data.EnrichedAircraft) error {
	return nil
}

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