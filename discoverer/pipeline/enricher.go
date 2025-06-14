package pipeline

import "github.com/corinm/aircraft/discovery/data"

type Enricher interface {
	Enrich(a *data.Aircraft) error
}
