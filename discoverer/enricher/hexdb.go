package enricher

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func HexDbGetAircraftInformation(hex string) (*HexDbResponse, error) {
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

	response := HexDbResponse{}

	err3 := json.Unmarshal(body, &response)
	if err3 != nil {
		return nil, err3
	}

	return &response, nil
}