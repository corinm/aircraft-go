package data

type EnrichedAircraft struct {
	IcaoHexCode          string   `json:"icaoHexCode"`          // e.g. "4CA293"

	Registration         string   `json:"registration"`         // e.g. "G-EZBZ"
	Manufacturer         string   `json:"manufacturer"`         // e.g. "Boeing"
	IcaoTypeCode         string   `json:"icaoTypeCode"`         // e.g. "B738"
	AircraftType         string   `json:"aircraftType"`         // e.g. "Boeing 737-800"
	RegisteredOwners     string   `json:"registeredOwners"`     // e.g. "easyJet UK"
	IcaoAirlineCode      string   `json:"icaoAirlineCode"`      // e.g. "EZY"

	CMPG                 string   `json:"cmpg"`                 // e.g. "Civilian"
	
	// The values of these two fields are inherently coupled to PlaneAlertDb
	// so I have included "PlaneAlertDb" in the name
	PlaneAlertDbCategory string   `json:"planeAlertDbCategory"` // e.g. "Special Forces"
	PlaneAlertDbTags     []string `json:"planeAlertDbTags"`     // e.g. ["Special Ops", "Special Air Service", "22 SAS"]
}
