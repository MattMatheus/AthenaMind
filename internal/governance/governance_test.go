package governance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
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

func TestEnforceResearchContractRequiresContractFileWhenEnabled(t *testing.T) {
	err := EnforceResearchContract(types.ResearchContractCheckInput{
		Enabled:   true,
		Operation: "retrieve",
	})
	if err == nil {
		t.Fatal("expected missing research contract to fail")
	}
}

func TestEnforceResearchContractValidatesHashSignatureAndMode(t *testing.T) {
	t.Setenv(researchSigningKeyEnv, "unit-test-key")
	root := t.TempDir()
	contractPath := filepath.Join(root, "contract.json")

	contract := types.ResearchContract{
		SchemaVersion:            "1.0",
		ContractID:               "rc-1",
		AllowedOperations:        []string{"retrieve", "evaluate"},
		ModeID:                   "Strategic",
		AbstractionLayerID:       "L2",
		ConstraintSpineVersion:   "2026.02.26",
		ConstraintSpineHash:      "abc123",
		SessionIntentFingerprint: "intent-fp",
		EntropyPolicy: types.EntropyPolicy{
			PreDistillRequired: true,
			ArtifactClasses:    []string{"raw_notes"},
		},
		DecisionLedgerRefs: []types.DecisionLedgerRef{{ID: "D-1", Status: "active"}},
	}
	hash, err := calculateContractHash(contract)
	if err != nil {
		t.Fatalf("calculate hash: %v", err)
	}
	contract.ContractHash = hash
	mac := hmac.New(sha256.New, []byte("unit-test-key"))
	_, _ = mac.Write([]byte(hash))
	contract.Signature = hex.EncodeToString(mac.Sum(nil))

	raw, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(contractPath, raw, 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	err = EnforceResearchContract(types.ResearchContractCheckInput{
		Enabled:            true,
		ContractFile:       contractPath,
		Operation:          "retrieve",
		ModeID:             "Strategic",
		AbstractionLayerID: "L2",
	})
	if err != nil {
		t.Fatalf("expected contract to validate, got error: %v", err)
	}
}

func TestEnforceResearchContractRejectsModeMismatch(t *testing.T) {
	t.Setenv(researchSigningKeyEnv, "unit-test-key")
	root := t.TempDir()
	contractPath := filepath.Join(root, "contract.json")

	contract := types.ResearchContract{
		SchemaVersion:            "1.0",
		ContractID:               "rc-2",
		AllowedOperations:        []string{"retrieve"},
		ModeID:                   "Validation",
		AbstractionLayerID:       "L1",
		ConstraintSpineVersion:   "2026.02.26",
		ConstraintSpineHash:      "abc123",
		SessionIntentFingerprint: "intent-fp",
	}
	hash, err := calculateContractHash(contract)
	if err != nil {
		t.Fatalf("calculate hash: %v", err)
	}
	contract.ContractHash = hash
	mac := hmac.New(sha256.New, []byte("unit-test-key"))
	_, _ = mac.Write([]byte(hash))
	contract.Signature = hex.EncodeToString(mac.Sum(nil))

	raw, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(contractPath, raw, 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	err = EnforceResearchContract(types.ResearchContractCheckInput{
		Enabled:            true,
		ContractFile:       contractPath,
		Operation:          "retrieve",
		ModeID:             "Strategic",
		AbstractionLayerID: "L1",
	})
	if err == nil {
		t.Fatal("expected mode mismatch to fail")
	}
}
