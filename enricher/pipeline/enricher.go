package pipeline

import (
	"context"
	"enricher/data"
)

type Enricher interface {
	Enrich(ctx context.Context, a *data.EnrichedAircraft) error
}
