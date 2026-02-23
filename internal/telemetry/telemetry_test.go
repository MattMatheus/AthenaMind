package telemetry

import (
	"errors"
	"testing"
)

func TestTelemetryErrorCode(t *testing.T) {
	if got := TelemetryErrorCode(errors.New("ERR_SAMPLE: boom")); got != "ERR_SAMPLE" {
		t.Fatalf("expected ERR_SAMPLE, got %s", got)
	}
}
