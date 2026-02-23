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
