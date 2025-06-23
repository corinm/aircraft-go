package pipeline

import (
	"enricher/data"
)

func EnrichAircraft(a *data.EnrichedAircraft, enrichers []Enricher) error {
	for _, enricher := range enrichers {
		if err := enricher.Enrich(a); err != nil {
			return err
		}
	}

	return nil
}
