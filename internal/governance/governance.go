package governance

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"athenamind/internal/types"
)

func EnforceConstraintChecks(operation, sessionID, scenarioID, traceID string) error {
	if err := enforceCostConstraint(operation); err != nil {
		return err
	}
	if err := enforceTraceabilityConstraint(sessionID, scenarioID, traceID); err != nil {
		return err
	}
	if err := enforceReliabilityConstraint(); err != nil {
		return err
	}
	return nil
}

func enforceCostConstraint(operation string) error {
	maxPerRun := 0.50
	if v := strings.TrimSpace(os.Getenv("MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD")); v != "" {
		if _, err := fmt.Sscanf(v, "%f", &maxPerRun); err != nil {
			return errors.New("ERR_CONSTRAINT_COST_CONFIG_INVALID: MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD must be numeric")
		}
	}
	estimated := map[string]float64{
		"write":    0.08,
		"retrieve": 0.02,
		"evaluate": 0.30,
	}[operation]
	if estimated > maxPerRun {
		return fmt.Errorf("ERR_CONSTRAINT_COST_BUDGET_EXCEEDED: estimated_%s_cost_usd=%.2f exceeds max_per_run_usd=%.2f", operation, estimated, maxPerRun)
	}
	return nil
}

func enforceTraceabilityConstraint(sessionID, scenarioID, traceID string) error {
	if IsTrue(os.Getenv("MEMORY_CONSTRAINT_FORCE_TRACE_MISSING")) {
		return errors.New("ERR_CONSTRAINT_TRACEABILITY_INCOMPLETE: trace policy forced missing required fields")
	}
	if strings.TrimSpace(sessionID) == "" || strings.TrimSpace(scenarioID) == "" || strings.TrimSpace(traceID) == "" {
		return errors.New("ERR_CONSTRAINT_TRACEABILITY_INCOMPLETE: session_id, scenario_id, and trace_id are required")
	}
	return nil
}

func enforceReliabilityConstraint() error {
	if IsTrue(os.Getenv("MEMORY_CONSTRAINT_RELIABILITY_FREEZE")) {
		return errors.New("ERR_CONSTRAINT_RELIABILITY_FREEZE_ACTIVE: autonomous promotion paths are frozen")
	}
	return nil
}

func IsLatencyDegraded(elapsedMs int64) bool {
	if IsTrue(os.Getenv("MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED")) {
		return true
	}
	threshold := int64(700)
	if v := strings.TrimSpace(os.Getenv("MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS")); v != "" {
		var parsed int64
		if _, err := fmt.Sscanf(v, "%d", &parsed); err == nil && parsed > 0 {
			threshold = parsed
		}
	}
	return elapsedMs > threshold
}

func EnforceWritePolicy(in types.WritePolicyInput) (types.WritePolicyDecision, error) {
	if IsTrue(os.Getenv("AUTONOMOUS_RUN")) {
		return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_NOT_ALLOWED_DURING_AUTONOMOUS_RUN: writes are blocked during autonomous runs")
	}
	allowed := map[string]struct{}{"planning": {}, "architect": {}, "pm": {}}
	if _, ok := allowed[in.Stage]; !ok {
		return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_STAGE_INVALID: --stage must be planning, architect, or pm")
	}
	if strings.TrimSpace(in.Reviewer) == "" {
		return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: --reviewer is required")
	}

	decision := strings.TrimSpace(strings.ToLower(in.Decision))
	if decision == "" {
		if in.ApprovedFlag {
			decision = "approved"
		} else {
			return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: provide --decision=approved|rejected")
		}
	}
	if decision != "approved" && decision != "rejected" {
		return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: --decision must be approved or rejected")
	}
	if strings.TrimSpace(in.Reason) == "" || strings.TrimSpace(in.Risk) == "" || strings.TrimSpace(in.Notes) == "" {
		return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_EVIDENCE_REQUIRED: --reason --risk and --notes are required")
	}
	if decision == "rejected" {
		if strings.TrimSpace(in.ReworkNotes) == "" || strings.TrimSpace(in.ReReviewedBy) == "" {
			return types.WritePolicyDecision{}, errors.New("ERR_MUTATION_REJECTION_EVIDENCE_REQUIRED: rejected decisions require --rework-notes and --re-reviewed-by")
		}
	}

	return types.WritePolicyDecision{
		Decision:     decision,
		Reviewer:     in.Reviewer,
		Notes:        in.Notes,
		Reason:       in.Reason,
		Risk:         in.Risk,
		ReworkNotes:  in.ReworkNotes,
		ReReviewedBy: in.ReReviewedBy,
	}, nil
}

func IsTrue(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes"
}
