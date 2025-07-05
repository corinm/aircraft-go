package pipeline

import (
	"context"
	"enricher/data"
	"testing"
)

func TestPipelineNoEnrichersDoesNothing(t *testing.T) {
	p := Pipeline{}

	aircraft := &data.EnrichedAircraft{
		IcaoHexCode: "000000",
	}

	err := p.Enrich(context.Background(), aircraft)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if aircraft.IcaoHexCode != "000000" {
		t.Errorf("Expected aircraft IcaoHexCode to be unchanged i.e. '000000', got %s", aircraft.IcaoHexCode)
	}
}

type MockRegistrationEnricher struct{}
func (m *MockRegistrationEnricher) Enrich(ctx context.Context, aircraft *data.EnrichedAircraft) error {
	aircraft.Registration = "G-MOCK"
	return nil
}

type MockManufacturerEnricher struct{}
func (m *MockManufacturerEnricher) Enrich(ctx context.Context, aircraft *data.EnrichedAircraft) error {
	aircraft.Manufacturer = "Mock Manufacturer"
	return nil
}

func TestPipelineRunsEnrichers(t *testing.T) {
	p := Pipeline{
		Enrichers: []Enricher{
			&MockRegistrationEnricher{},
			&MockManufacturerEnricher{},
		},
	}

	aircraft := &data.EnrichedAircraft{
		IcaoHexCode: "000000",
	}

	errs := p.Enrich(context.Background(), aircraft)
	if errs != nil {
		t.Errorf("Expected no error, got %v", errs)
	}

	if aircraft.Registration != "G-MOCK" {
		t.Errorf("Expected aircraft Registration to be 'G-MOCK', got %s", aircraft.Registration)
	}

	if aircraft.Manufacturer != "Mock Manufacturer" {
		t.Errorf("Expected aircraft Manufacturer to be 'Mock Manufacturer', got %s", aircraft.Manufacturer)
	}
}
