package pipeline

import (
	"context"
	"enricher/data"
	"log"
)

type Pipeline struct {
	Enrichers []Enricher
}

func (p *Pipeline) Enrich(ctx context.Context, a *data.EnrichedAircraft) error {
	for _, enricher := range p.Enrichers {
		if err := enricher.Enrich(ctx, a); err != nil {
			log.Printf("enriching aircraft %s with %T: %v", a.IcaoHexCode, enricher, err)
		}
	}

	return nil
}
