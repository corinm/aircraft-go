package pipeline

import "github.com/corinm/aircraft/discovery/data"

func EnrichAircraft(a *data.Aircraft, enrichers []Enricher) error {
	for _, enricher := range enrichers {
		if err := enricher.Enrich(a); err != nil {
			return err
		}
	}

	return nil
}
