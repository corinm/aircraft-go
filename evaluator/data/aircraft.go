package data

type EnrichedAircraft struct {
	AiocHexCode      string `json:"aiocHexCode"`      // e.g. "4CA293"
	Registration     string `json:"registration"`     // e.g. "G-EZBZ"
	Manufacturer     string `json:"manufacturer"`     // e.g. "Boeing"
	ICAOTypeCode     string `json:"icaoTypeCode"`     // e.g. "B738"
	AircraftType     string `json:"aircraftType"`     // e.g. "Boeing 737-800"
	RegisteredOwners string `json:"registeredOwners"` // e.g. "easyJet UK"
	IcaoAirlineCode  string `json:"icaoAirlineCode"`  // e.g. "EZY"
}
