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
