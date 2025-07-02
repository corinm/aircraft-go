package pipeline

import (
	"encoding/json"
	"enricher/data"
	"fmt"
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
	ICAOTypeCode     string `json:"ICAOTypeCode"`
	Type             string `json:"Type"`
	RegisteredOwners string `json:"RegisteredOwners"`
	OperatorFlagCode string `json:"OperatorFlagCode"`
}

func hexDbGetAircraftInformation(hex string) (*hexDbResponse, error) {
	url := fmt.Sprintf("https://hexdb.io/api/v1/aircraft/%s", hex)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for HexDB (%s): %w", url, err)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from HexDB (%s): %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from from HexDB (%s), status code: %s", url, resp.Status)
	}

	var result hexDbResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response from HexDB (%s): %w", url, err)
	}

	return &result, nil
}

func (e *HexDbEnricher) Enrich(a *data.EnrichedAircraft) error {
	resp, err := hexDbGetAircraftInformation(a.AiocHexCode)
	if err != nil {
		return err
	}

	a.Registration = resp.Registration
	a.Manufacturer = resp.Manufacturer
	a.ICAOTypeCode = resp.ICAOTypeCode
	a.AircraftType = resp.Type
	a.RegisteredOwners = resp.RegisteredOwners
	a.IcaoAirlineCode = resp.OperatorFlagCode

	return nil
}
