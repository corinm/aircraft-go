package enricher

type MetadataFetcher struct {}

type HexDbResponse struct {
	// ModeS            string `json:"ModeS"`
	Registration     string `json:"Registration"`
	Manufacturer     string `json:"Manufacturer"`
	ICAOTypeCode     string `json:"ICAOTypeCode"`
	Type             string `json:"Type"`
	RegisteredOwners string `json:"RegisteredOwners"`
	OperatorFlagCode string `json:"OperatorFlagCode"`
}

var hexToHexDbResponse = map[string]HexDbResponse{}

func (m MetadataFetcher) GetRegistrationByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.Registration, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.Registration, nil
}

func (m MetadataFetcher) GetManufacturerByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.Manufacturer, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.Manufacturer, nil
}

func (m MetadataFetcher) GetICAOTypeCodeByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.ICAOTypeCode, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.ICAOTypeCode, nil
}

func (m MetadataFetcher) GetAircraftTypeByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.Type, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.Type, nil
}

func (m MetadataFetcher) GetRegisteredOwnersByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.RegisteredOwners, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.RegisteredOwners, nil
}

func (m MetadataFetcher) GetIcaoAirlineCodeByHex(hex string) (string, error) {
	if hexDbResponse, exists := hexToHexDbResponse[hex]; exists {
		return hexDbResponse.OperatorFlagCode, nil
	}

	r, err := HexDbGetAircraftInformation(hex)
	if err != nil {
		return "", err
	}

	hexToHexDbResponse[hex] = *r

	return r.OperatorFlagCode, nil
}

// func (m MetadataFetcher) IsMilitaryByHex(hex string) (bool, error) {}

// func (m MetadataFetcher) IsInterestingByHex(hex string) (bool, error) {}
