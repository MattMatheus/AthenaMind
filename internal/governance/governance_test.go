package governance

import (
	"testing"

	"athenamind/internal/types"
)

func TestEnforceWritePolicyRequiresReviewer(t *testing.T) {
	_, err := EnforceWritePolicy(types.WritePolicyInput{
		Stage:    "planning",
		Decision: "approved",
		Reason:   "r",
		Risk:     "low",
		Notes:    "n",
	})
	if err == nil {
		t.Fatal("expected reviewer requirement error")
	}
}

func TestIsLatencyDegradedDefaultsTo700Ms(t *testing.T) {
	t.Setenv("MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS", "")
	t.Setenv("MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED", "")
	if IsLatencyDegraded(699) {
		t.Fatal("expected 699ms to stay below default latency degradation threshold")
	}
	if !IsLatencyDegraded(701) {
		t.Fatal("expected 701ms to exceed default latency degradation threshold")
	}
}

func TestIsLatencyDegradedUsesConfiguredThreshold(t *testing.T) {
	t.Setenv("MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS", "1500")
	t.Setenv("MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED", "")
	if IsLatencyDegraded(1499) {
		t.Fatal("expected 1499ms to stay below configured threshold")
	}
	if !IsLatencyDegraded(1501) {
		t.Fatal("expected 1501ms to exceed configured threshold")
	}
}

func TestIsLatencyDegradedZeroDisablesLatencyFallback(t *testing.T) {
	t.Setenv("MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS", "0")
	t.Setenv("MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED", "")
	if IsLatencyDegraded(100000) {
		t.Fatal("expected latency degradation to be disabled when threshold is 0")
	}
}
