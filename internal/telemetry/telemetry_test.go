package telemetry

import (
	"errors"
	"testing"

	"athenamind/internal/types"
)

func TestTelemetryErrorCode(t *testing.T) {
	if got := TelemetryErrorCode(errors.New("ERR_SAMPLE: boom")); got != "ERR_SAMPLE" {
		t.Fatalf("expected ERR_SAMPLE, got %s", got)
	}
}

func TestEmitRetrievalMetricTracksRunningRates(t *testing.T) {
	root := t.TempDir()
	first, err := EmitRetrievalMetric(root, types.RetrieveResult{SelectionMode: "embedding_semantic", SelectedID: "a", PrecisionHint: 1})
	if err != nil {
		t.Fatalf("emit first metric: %v", err)
	}
	if first.SemanticHitRate != 1 || first.FallbackRate != 0 {
		t.Fatalf("unexpected first rates: %+v", first)
	}
	second, err := EmitRetrievalMetric(root, types.RetrieveResult{SelectionMode: "fallback_path_priority", SelectedID: "b", PrecisionHint: 0})
	if err != nil {
		t.Fatalf("emit second metric: %v", err)
	}
	if second.SemanticHitRate != 0.5 || second.FallbackRate != 0.5 {
		t.Fatalf("unexpected second rates: %+v", second)
	}
}
