package data

type Aircraft struct {
	AiocHexCode string
}

type AircraftEnriched struct {
	AiocHexCode      string
	Registration	 string // e.g. "G-EZBZ"
	Manufacturer     string // e.g. "Boeing"
	ICAOTypeCode     string // e.g. "B738"
	AircraftType     string // e.g. "Boeing 737-800"
	RegisteredOwners string // e.g. "easyJet UK"
	IcaoAirlineCode  string // e.g. "EZY"
}
