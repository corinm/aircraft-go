package pipeline

import (
	"enricher/data"
)

type Enricher interface {
	Enrich(a *data.EnrichedAircraft) error
}
