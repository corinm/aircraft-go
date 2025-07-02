package pipeline

import (
	"enricher/data"
)

type Pipeline struct {
	Enrichers []Enricher
}

func (p *Pipeline) Enrich(a *data.EnrichedAircraft) error {
	for _, enricher := range p.Enrichers {
		if err := enricher.Enrich(a); err != nil {
			return err
		}
	}

	return nil
}
