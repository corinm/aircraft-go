package pipeline

import (
	"context"
	"enricher/data"
	"log"
)

type Pipeline struct {
	Enrichers []Enricher
}

func (p *Pipeline) Enrich(ctx context.Context, a *data.EnrichedAircraft) []error {
	errs := make([]error, 0)

	for _, enricher := range p.Enrichers {
		if err := enricher.Enrich(ctx, a); err != nil {
			log.Printf("enriching aircraft %s with %T: %v", a.IcaoHexCode, enricher, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
