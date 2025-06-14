package data

type Aircraft struct {
	AiocHexCode string
}

type AircraftEnriched struct {
	AiocHexCode string
	Manufacturer string
	ICAOTypeCode string // e.g. "B738", "A320" etc.
	AircraftType  string // e.g. "Boeing 737-800", "Airbus A320" etc.
	Operator	 string
	PlaneAlertDbTags []string // e.g. "Police Squad", "Military Transport" etc.
}
