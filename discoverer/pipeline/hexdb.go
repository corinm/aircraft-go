package pipeline

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/corinm/aircraft/discovery/data"
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
	r, err1 := http.Get("https://hexdb.io/api/v1/aircraft/" + hex)
	if err1 != nil {
		return nil, err1
	}
	defer r.Body.Close()
	
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from HexDB, status code: " + r.Status)
	}

	body, err2 := io.ReadAll(r.Body)
	if err2 != nil {
		return nil, err2
	}

	response := hexDbResponse{}

	err3 := json.Unmarshal(body, &response)
	if err3 != nil {
		return nil, err3
	}

	return &response, nil
}

func (e *HexDbEnricher) Enrich(a *data.Aircraft) error {
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
