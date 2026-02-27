package governance

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"athenamind/internal/types"
)

const (
	researchSigningKeyEnv       = "SEDA_RESEARCH_CONTRACT_SIGNING_KEY"
	researchModeEnv             = "SEDA_RESEARCH_MODE"
	researchRequireSignatureEnv = "SEDA_RESEARCH_REQUIRE_SIGNATURE"
)

func EnforceResearchContract(input types.ResearchContractCheckInput) error {
	enabled := input.Enabled || IsTrue(os.Getenv(researchModeEnv))
	if !enabled {
		return nil
	}

	contractPath := strings.TrimSpace(input.ContractFile)
	if contractPath == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_REQUIRED: --research-contract-file is required when research mode is enabled")
	}

	payload, err := os.ReadFile(contractPath)
	if err != nil {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_READ: %w", err)
	}

	var contract types.ResearchContract
	if err := json.Unmarshal(payload, &contract); err != nil {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_PARSE: %w", err)
	}

	if err := validateResearchContract(contract, input); err != nil {
		return err
	}
	return nil
}

func validateResearchContract(contract types.ResearchContract, input types.ResearchContractCheckInput) error {
	if strings.TrimSpace(contract.ModeID) == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_INVALID: mode_id is required")
	}
	if strings.TrimSpace(contract.AbstractionLayerID) == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_INVALID: abstraction_layer_id is required")
	}
	if strings.TrimSpace(contract.ConstraintSpineVersion) == "" || strings.TrimSpace(contract.ConstraintSpineHash) == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_INVALID: constraint_spine_version and constraint_spine_hash are required")
	}
	if strings.TrimSpace(contract.SessionIntentFingerprint) == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_INVALID: session_intent_fingerprint is required")
	}
	if strings.TrimSpace(contract.ContractHash) == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_INVALID: contract_hash is required")
	}

	operation := strings.ToLower(strings.TrimSpace(input.Operation))
	if operation != "" && len(contract.AllowedOperations) > 0 {
		allowed := false
		for _, op := range contract.AllowedOperations {
			if strings.EqualFold(strings.TrimSpace(op), operation) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("ERR_RESEARCH_CONTRACT_OPERATION_FORBIDDEN: operation=%s is not allowed by contract", operation)
		}
	}

	requestedMode := strings.TrimSpace(input.ModeID)
	if requestedMode != "" && !strings.EqualFold(requestedMode, strings.TrimSpace(contract.ModeID)) {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_MODE_MISMATCH: requested_mode=%s contract_mode=%s", requestedMode, contract.ModeID)
	}

	requestedLayer := strings.TrimSpace(input.AbstractionLayerID)
	if requestedLayer != "" && !strings.EqualFold(requestedLayer, strings.TrimSpace(contract.AbstractionLayerID)) {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_LAYER_MISMATCH: requested_layer=%s contract_layer=%s", requestedLayer, contract.AbstractionLayerID)
	}

	for _, ref := range contract.DecisionLedgerRefs {
		status := strings.ToLower(strings.TrimSpace(ref.Status))
		if status == "invalidated" || status == "revoked" {
			return fmt.Errorf("ERR_RESEARCH_CONTRACT_DECISION_INVALIDATED: decision_id=%s status=%s", ref.ID, ref.Status)
		}
	}

	expectedHash, err := calculateContractHash(contract)
	if err != nil {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_HASH: %w", err)
	}
	if !strings.EqualFold(expectedHash, strings.TrimSpace(contract.ContractHash)) {
		return errors.New("ERR_RESEARCH_CONTRACT_HASH_MISMATCH: contract hash does not match payload")
	}

	return verifyContractSignature(strings.TrimSpace(contract.ContractHash), strings.TrimSpace(contract.Signature))
}

func calculateContractHash(contract types.ResearchContract) (string, error) {
	normalized := contract
	normalized.ContractHash = ""
	normalized.Signature = ""
	for i := range normalized.AllowedOperations {
		normalized.AllowedOperations[i] = strings.ToLower(strings.TrimSpace(normalized.AllowedOperations[i]))
	}
	sort.Strings(normalized.AllowedOperations)
	sort.Slice(normalized.DecisionLedgerRefs, func(i, j int) bool {
		return strings.TrimSpace(normalized.DecisionLedgerRefs[i].ID) < strings.TrimSpace(normalized.DecisionLedgerRefs[j].ID)
	})
	for i := range normalized.DecisionLedgerRefs {
		normalized.DecisionLedgerRefs[i].ID = strings.TrimSpace(normalized.DecisionLedgerRefs[i].ID)
		normalized.DecisionLedgerRefs[i].Status = strings.ToLower(strings.TrimSpace(normalized.DecisionLedgerRefs[i].Status))
		normalized.DecisionLedgerRefs[i].InvalidationCondition = strings.TrimSpace(normalized.DecisionLedgerRefs[i].InvalidationCondition)
	}
	for i := range normalized.EntropyPolicy.ArtifactClasses {
		normalized.EntropyPolicy.ArtifactClasses[i] = strings.ToLower(strings.TrimSpace(normalized.EntropyPolicy.ArtifactClasses[i]))
	}
	sort.Strings(normalized.EntropyPolicy.ArtifactClasses)

	raw, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func verifyContractSignature(contractHash, signature string) error {
	signingKey := strings.TrimSpace(os.Getenv(researchSigningKeyEnv))
	requireSignature := IsTrue(os.Getenv(researchRequireSignatureEnv)) || signingKey != ""
	if !requireSignature {
		return nil
	}
	if signingKey == "" {
		return fmt.Errorf("ERR_RESEARCH_CONTRACT_SIGNATURE_CONFIG: %s is required when signature verification is enabled", researchSigningKeyEnv)
	}
	if signature == "" {
		return errors.New("ERR_RESEARCH_CONTRACT_SIGNATURE_REQUIRED: signature is required")
	}
	mac := hmac.New(sha256.New, []byte(signingKey))
	_, _ = mac.Write([]byte(contractHash))
	expected := hex.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(strings.ToLower(signature)), []byte(expected)) != 1 {
		return errors.New("ERR_RESEARCH_CONTRACT_SIGNATURE_INVALID: signature does not match contract hash")
	}
	return nil
}
