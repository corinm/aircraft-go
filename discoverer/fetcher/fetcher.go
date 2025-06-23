package fetcher

import "github.com/corinm/aircraft/discovery/data"

type Fetcher interface {
	fetchAircraft() ([]data.RawAircraft, error)
}
