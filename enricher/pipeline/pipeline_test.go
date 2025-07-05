package pipeline

import (
	"context"
	"enricher/data"
	"testing"
)

func TestPipeline(t *testing.T) {
	tests := []struct {
		name     string
		enrichers []Enricher
		input data.EnrichedAircraft
		expectedOutput data.EnrichedAircraft
	}{
		{
			name: "No enrichers",
			enrichers: []Enricher{},
			input: data.EnrichedAircraft{ IcaoHexCode: "000000" },
			expectedOutput: data.EnrichedAircraft{IcaoHexCode: "000000"},
		},
		{
			name: "With enrichers",
			enrichers: []Enricher{
				&MockRegistrationEnricher{},
				&MockManufacturerEnricher{},
			},
			input: data.EnrichedAircraft{ IcaoHexCode: "000000" },
			expectedOutput: data.EnrichedAircraft{
				IcaoHexCode: "000000",
				Registration: "G-MOCK",
				Manufacturer: "Mock Ltd.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pipeline{
				Enrichers: tt.enrichers,
			}

			aircraft := tt.input

			err := p.Enrich(context.Background(), &aircraft)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if aircraft.IcaoHexCode != tt.expectedOutput.IcaoHexCode {
				t.Errorf("Expected aircraft IcaoHexCode to be %s, got %s", tt.expectedOutput.IcaoHexCode, aircraft.IcaoHexCode)
			}

			if aircraft.Registration != tt.expectedOutput.Registration {
				t.Errorf("Expected aircraft Registration to be %s, got %s", tt.expectedOutput.Registration, aircraft.Registration)
			}

			if aircraft.Manufacturer != tt.expectedOutput.Manufacturer {
				t.Errorf("Expected aircraft Manufacturer to be %s, got %s", tt.expectedOutput.Manufacturer, aircraft.Manufacturer)
			}
		})
	}
}

type MockRegistrationEnricher struct{}
func (m *MockRegistrationEnricher) Enrich(ctx context.Context, aircraft *data.EnrichedAircraft) error {
	aircraft.Registration = "G-MOCK"
	return nil
}

type MockManufacturerEnricher struct{}
func (m *MockManufacturerEnricher) Enrich(ctx context.Context, aircraft *data.EnrichedAircraft) error {
	aircraft.Manufacturer = "Mock Ltd."
	return nil
}
