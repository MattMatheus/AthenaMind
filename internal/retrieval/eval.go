package retrieval

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"athenamind/internal/types"
)

func EvaluateRetrieval(root string, queries []types.EvaluationQuery, corpusID, querySetID, configID string) (types.EvaluationReport, error) {
	return EvaluateRetrievalWithEmbeddingEndpoint(root, queries, corpusID, querySetID, configID, DefaultEmbeddingEndpoint)
}

func EvaluateRetrievalWithEmbeddingEndpoint(root string, queries []types.EvaluationQuery, corpusID, querySetID, configID, embeddingEndpoint string) (types.EvaluationReport, error) {
	return EvaluateRetrievalWithOptionsAndEmbeddingEndpoint(
		root,
		queries,
		corpusID,
		querySetID,
		configID,
		embeddingEndpoint,
		RetrieveOptions{},
	)
}

func EvaluateRetrievalWithOptionsAndEmbeddingEndpoint(
	root string,
	queries []types.EvaluationQuery,
	corpusID, querySetID, configID, embeddingEndpoint string,
	options RetrieveOptions,
) (types.EvaluationReport, error) {
	options = normalizeRetrieveOptions(options)
	report := types.EvaluationReport{
		CorpusID:            corpusID,
		QuerySetID:          querySetID,
		ConfigID:            configID,
		Strategy:            fmt.Sprintf("mode=%s,backend=%s,top_k=%d", options.Mode, options.Backend, options.TopK),
		Status:              "FAIL",
		FailingQueries:      []types.QueryMiss{},
		DeterministicReplay: []types.DeterministicReplay{},
	}
	if len(queries) == 0 {
		return report, errors.New("ERR_EVAL_QUERY_SET_INVALID: query set must contain at least one query")
	}

	total := len(queries)
	top1Useful := 0
	selectionModePresent := 0
	sourceTracePresent := 0
	fallbackChecks := 0
	fallbackStable := 0
	latencies := make([]float64, 0, total)
	var latencySum float64

	for _, q := range queries {
		started := time.Now()
		result, _, err := RetrieveWithOptionsAndEndpointAndSession(
			root,
			q.Query,
			q.Domain,
			embeddingEndpoint,
			"",
			options,
		)
		if err != nil {
			return report, err
		}
		elapsed := float64(time.Since(started).Milliseconds())
		latencies = append(latencies, elapsed)
		latencySum += elapsed
		if result.SelectionMode != "" {
			selectionModePresent++
		}
		if result.SelectedID != "" && result.SourcePath != "" {
			sourceTracePresent++
		}
		if (result.SelectionMode == "semantic" || result.SelectionMode == "embedding_semantic" || result.SelectionMode == "hybrid_rrf") &&
			q.ExpectedID != "" &&
			result.SelectedID == q.ExpectedID {
			top1Useful++
		} else if q.ExpectedID != "" && result.SelectedID != q.ExpectedID {
			report.FailingQueries = append(report.FailingQueries, types.QueryMiss{
				Query:      q.Query,
				ExpectedID: q.ExpectedID,
				ActualID:   result.SelectedID,
				Mode:       result.SelectionMode,
			})
		}

		if strings.HasPrefix(result.SelectionMode, "fallback_") {
			fallbackChecks++
			stable := true
			for i := 0; i < 4; i++ {
				again, _, err := RetrieveWithOptionsAndEndpointAndSession(
					root,
					q.Query,
					q.Domain,
					embeddingEndpoint,
					"",
					options,
				)
				if err != nil {
					return report, err
				}
				if again.SelectionMode != result.SelectionMode || again.SelectedID != result.SelectedID || again.SourcePath != result.SourcePath {
					stable = false
					break
				}
			}
			if stable {
				fallbackStable++
			}
			report.DeterministicReplay = append(report.DeterministicReplay, types.DeterministicReplay{
				Query:      q.Query,
				Mode:       result.SelectionMode,
				SelectedID: result.SelectedID,
				SourcePath: result.SourcePath,
				StableRuns: 5,
			})
		}
	}

	report.Top1UsefulRate = metric(top1Useful, total)
	report.SelectionModeReporting = metric(selectionModePresent, total)
	report.SourceTraceCompleteness = metric(sourceTracePresent, total)
	report.AvgLatencyMS = latencySum / float64(total)
	report.LatencyP50MS = percentile(latencies, 0.50)
	report.LatencyP95MS = percentile(latencies, 0.95)
	if fallbackChecks == 0 {
		report.FallbackDeterminism = types.EvaluationMetric{Numerator: 1, Denominator: 1, Rate: 1}
	} else {
		report.FallbackDeterminism = metric(fallbackStable, fallbackChecks)
	}

	maxP95 := float64(1500)
	if env := strings.TrimSpace(os.Getenv("MEMORY_EVAL_MAX_P95_MS")); env != "" {
		if parsed, err := strconv.ParseFloat(env, 64); err == nil && parsed > 0 {
			maxP95 = parsed
		}
	}

	pass := total >= 50 &&
		report.Top1UsefulRate.Rate >= 0.80 &&
		report.FallbackDeterminism.Rate == 1 &&
		report.SelectionModeReporting.Rate == 1 &&
		report.SourceTraceCompleteness.Rate == 1 &&
		report.LatencyP95MS <= maxP95

	if pass {
		report.Status = "PASS"
		if report.Top1UsefulRate.Rate >= 0.80 && report.Top1UsefulRate.Rate <= 0.82 {
			report.Status = "WATCH"
		}
	}
	if report.Status == "PASS" && report.LatencyP95MS <= (maxP95*0.8) && report.Top1UsefulRate.Rate >= 0.85 {
		report.Recommendation = "promote"
	} else if report.Status == "FAIL" {
		report.Recommendation = "eliminate"
	} else {
		report.Recommendation = "iterate"
	}

	return report, nil
}

func metric(numerator, denominator int) types.EvaluationMetric {
	if denominator == 0 {
		return types.EvaluationMetric{Numerator: 0, Denominator: 0, Rate: 0}
	}
	return types.EvaluationMetric{
		Numerator:   numerator,
		Denominator: denominator,
		Rate:        float64(numerator) / float64(denominator),
	}
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	cp := append([]float64(nil), values...)
	sort.Float64s(cp)
	if p <= 0 {
		return cp[0]
	}
	if p >= 1 {
		return cp[len(cp)-1]
	}
	pos := int(float64(len(cp)-1) * p)
	if pos < 0 {
		pos = 0
	}
	if pos >= len(cp) {
		pos = len(cp) - 1
	}
	return cp[pos]
}

func LoadEvaluationQueries(path string) ([]types.EvaluationQuery, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var queries []types.EvaluationQuery
	if err := json.Unmarshal(data, &queries); err != nil {
		return nil, fmt.Errorf("ERR_EVAL_QUERY_SET_INVALID: cannot parse %s: %w", path, err)
	}
	for _, q := range queries {
		if strings.TrimSpace(q.Query) == "" {
			return nil, errors.New("ERR_EVAL_QUERY_SET_INVALID: each query must include query text")
		}
	}
	return queries, nil
}
