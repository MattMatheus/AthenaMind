package telemetry

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"athenamind/internal/types"
)

const (
	EventSchema         = "1.0"
	TelemetryRel        = "telemetry/events.jsonl"
	RetrievalMetricsRel = "telemetry/retrieval-metrics.jsonl"
)

func Emit(root, telemetryPath string, ev types.TelemetryEvent) error {
	path := strings.TrimSpace(telemetryPath)
	if path == "" {
		path = filepath.Join(root, filepath.FromSlash(TelemetryRel))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func TelemetryErrorCode(err error) string {
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "ERR_UNKNOWN"
	}
	first := strings.Fields(msg)
	code := strings.TrimSuffix(first[0], ":")
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "ERR_") {
		return code
	}
	return "ERR_COMMAND_FAILED"
}

func NormalizeMemoryType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "procedural", "state", "semantic":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "semantic"
	}
}

func NormalizeOperatorVerdict(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "correct", "partially_correct", "incorrect", "not_scored":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "not_scored"
	}
}

func NormalizeTelemetryValue(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

type RetrievalMetricRate struct {
	SemanticHitRate float64
	FallbackRate    float64
}

func EmitRetrievalMetric(root string, result types.RetrieveResult) (RetrievalMetricRate, error) {
	path := filepath.Join(root, filepath.FromSlash(RetrievalMetricsRel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return RetrievalMetricRate{}, err
	}

	semanticHits := 0
	fallbackHits := 0
	total := 0
	if existing, err := os.Open(path); err == nil {
		scanner := bufio.NewScanner(existing)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			var row struct {
				SemanticHit  bool `json:"semantic_hit"`
				FallbackUsed bool `json:"fallback_used"`
			}
			if err := json.Unmarshal([]byte(line), &row); err != nil {
				continue
			}
			total++
			if row.SemanticHit {
				semanticHits++
			}
			if row.FallbackUsed {
				fallbackHits++
			}
		}
		_ = existing.Close()
	}

	semantic := result.SelectionMode == "semantic" || result.SelectionMode == "embedding_semantic"
	fallback := strings.HasPrefix(result.SelectionMode, "fallback_")
	total++
	if semantic {
		semanticHits++
	}
	if fallback {
		fallbackHits++
	}

	semanticRate := 0.0
	fallbackRate := 0.0
	if total > 0 {
		semanticRate = float64(semanticHits) / float64(total)
		fallbackRate = float64(fallbackHits) / float64(total)
	}
	metric := map[string]any{
		"timestamp_utc":     time.Now().UTC().Format(time.RFC3339),
		"selection_mode":    result.SelectionMode,
		"selected_id":       result.SelectedID,
		"semantic_hit":      semantic,
		"fallback_used":     fallback,
		"precision_proxy":   result.PrecisionHint,
		"semantic_hit_rate": semanticRate,
		"fallback_rate":     fallbackRate,
	}
	data, err := json.Marshal(metric)
	if err != nil {
		return RetrievalMetricRate{}, fmt.Errorf("marshal retrieval metric: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return RetrievalMetricRate{}, err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return RetrievalMetricRate{}, err
	}
	return RetrievalMetricRate{
		SemanticHitRate: semanticRate,
		FallbackRate:    fallbackRate,
	}, nil
}
