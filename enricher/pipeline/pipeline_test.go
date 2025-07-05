package pipeline

import (
	"context"
	"enricher/data"
	"errors"
	"testing"
)

func TestPipelineResult(t *testing.T) {
	tests := []struct {
		name     string
		enrichers []Enricher
		input data.EnrichedAircraft
		expectedOutput data.EnrichedAircraft
		expectError bool
	}{
		{
			name: "No enrichers",
			enrichers: []Enricher{},
			input: data.EnrichedAircraft{ IcaoHexCode: "000000" },
			expectedOutput: data.EnrichedAircraft{IcaoHexCode: "000000"},
			expectError: false,
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
			expectError: false,
		},
		{
			name: "Error in enricher doesn't stop pipeline",
			enrichers: []Enricher{
				&MockRegistrationEnricher{},
				&MockErrorEnricher{},
				&MockManufacturerEnricher{},
			},
			input: data.EnrichedAircraft{ IcaoHexCode: "000000" },
			expectedOutput: data.EnrichedAircraft{
				IcaoHexCode: "000000",
				Registration: "G-MOCK",
				Manufacturer: "Mock Ltd.",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pipeline{
				Enrichers: tt.enrichers,
			}

			aircraft := tt.input

			err := p.Enrich(context.Background(), &aircraft)
			if err != nil && !tt.expectError {
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

func TestPipelineErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		enrichers []Enricher
		input data.EnrichedAircraft
		expectedErrors []error
	}{
		{
			name: "Error in enricher",
			enrichers: []Enricher{
				&MockErrorEnricher{},
			},
			input: data.EnrichedAircraft{ IcaoHexCode: "000000" },
			expectedErrors: []error{errors.New("MockErrorEnricher failed to enrich")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pipeline{
				Enrichers: tt.enrichers,
			}

			aircraft := tt.input

			errs := p.Enrich(context.Background(), &aircraft)
			if errs == nil {
				t.Errorf("Expected errors, got nil")
				return
			}

			if len(errs) != len(tt.expectedErrors) {
				t.Errorf("Expected %d errors, got %d", len(tt.expectedErrors), len(errs))
				return
			}

			for i, expectedErr := range tt.expectedErrors {
				if errs[i].Error() != expectedErr.Error() {
					t.Errorf("Expected error %d to be %v, got %v", i, expectedErr, errs[i])
				}
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

type MockErrorEnricher struct{}
func (m *MockErrorEnricher) Enrich(ctx context.Context, aircraft *data.EnrichedAircraft) error {
	return errors.New("MockErrorEnricher failed to enrich")
}
