package enricher

import (
	"github.com/corinm/aircraft/discovery/data"
)

func EnrichAircraft(aircraft data.Aircraft) (data.AircraftEnriched, error) {
	metdataFetcher := MetadataFetcher{}

	registration, _ := metdataFetcher.GetRegistrationByHex(aircraft.AiocHexCode)
	manufacturer, _ := metdataFetcher.GetManufacturerByHex(aircraft.AiocHexCode)
	icaoTypeCode, _ := metdataFetcher.GetICAOTypeCodeByHex(aircraft.AiocHexCode)
	aircraftType, _ := metdataFetcher.GetAircraftTypeByHex(aircraft.AiocHexCode)
	registeredOwners, _ := metdataFetcher.GetRegisteredOwnersByHex(aircraft.AiocHexCode)
	icaoAirlineCode, _ := metdataFetcher.GetIcaoAirlineCodeByHex(aircraft.AiocHexCode)
	// isMilitary, _ := metdataFetcher.IsMilitaryByHex(aircraft.AiocHexCode)
	// isInteresting, _ := metdataFetcher.IsInterestingByHex(aircraft.AiocHexCode)

	enrichedAircraft := data.AircraftEnriched{
		AiocHexCode:      aircraft.AiocHexCode,
		Registration:     registration,
		Manufacturer:     manufacturer,
		ICAOTypeCode:     icaoTypeCode,
		AircraftType:     aircraftType,
		RegisteredOwners: registeredOwners,
		IcaoAirlineCode:  icaoAirlineCode,
		// IsMilitary:    isMilitary,
		// IsInteresting: isInteresting,
	}

	return enrichedAircraft, nil
}
