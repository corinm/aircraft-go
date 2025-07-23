package pipeline

import (
	"context"
	"encoding/json"
	"enricher/data"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HexDbEnricher struct {
	HexDbUrl string
}

type hexDbResponse struct {
	// ModeS            string `json:"ModeS"`
	Registration     string `json:"Registration"`
	Manufacturer     string `json:"Manufacturer"`
	IcaoTypeCode     string `json:"ICAOTypeCode"`
	Type             string `json:"Type"`
	RegisteredOwners string `json:"RegisteredOwners"`
	OperatorFlagCode string `json:"OperatorFlagCode"`
}

func hexDbGetAircraftInformation(ctx context.Context, hex string) (*hexDbResponse, error) {
	url := fmt.Sprintf("https://hexdb.io/api/v1/aircraft/%s", hex)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for HexDB (%s): %w", url, err)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from HexDB (%s): %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("aircraft with hex code %s not found in HexDB (%s)", hex, url)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from from HexDB (%s), status code: %s", url, resp.Status)
	}

	var result hexDbResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response from HexDB (%s): %w", url, err)
	}

	return &result, nil
}

func (e *HexDbEnricher) Enrich(ctx context.Context, a *data.EnrichedAircraft) error {
	resp, err := hexDbGetAircraftInformation(ctx, a.IcaoHexCode)
	if err != nil {
		return err
	}

	a.Registration = resp.Registration
	a.Manufacturer = resp.Manufacturer
	a.IcaoTypeCode = resp.IcaoTypeCode
	a.AircraftType = resp.Type
	a.RegisteredOwners = resp.RegisteredOwners
	a.IcaoAirlineCode = resp.OperatorFlagCode
	// HexDB sometimes provides longer, invalid flag codes that appear to actually be type codes
	if len(a.IcaoAirlineCode) > 3 {
		log.Printf("Warning: IcaoAirlineCode for %s is longer than 3 characters: %s", a.IcaoHexCode, a.IcaoAirlineCode)
		a.IcaoAirlineCode = ""
	}

	return nil
}
