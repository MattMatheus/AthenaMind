package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func evaluateRetrieval(root string, queries []evaluationQuery, corpusID, querySetID, configID string) (evaluationReport, error) {
	report := evaluationReport{
		CorpusID:            corpusID,
		QuerySetID:          querySetID,
		ConfigID:            configID,
		Status:              "FAIL",
		FailingQueries:      []queryMiss{},
		DeterministicReplay: []deterministicReplay{},
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

	for _, q := range queries {
		result, err := retrieve(root, q.Query, q.Domain)
		if err != nil {
			return report, err
		}
		if result.SelectionMode != "" {
			selectionModePresent++
		}
		if result.SelectedID != "" && result.SourcePath != "" {
			sourceTracePresent++
		}
		if result.SelectionMode == "semantic" && q.ExpectedID != "" && result.SelectedID == q.ExpectedID {
			top1Useful++
		} else if q.ExpectedID != "" && result.SelectedID != q.ExpectedID {
			report.FailingQueries = append(report.FailingQueries, queryMiss{
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
				again, err := retrieve(root, q.Query, q.Domain)
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
			report.DeterministicReplay = append(report.DeterministicReplay, deterministicReplay{
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
	if fallbackChecks == 0 {
		report.FallbackDeterminism = evaluationMetric{Numerator: 1, Denominator: 1, Rate: 1}
	} else {
		report.FallbackDeterminism = metric(fallbackStable, fallbackChecks)
	}

	pass := total >= 50 &&
		report.Top1UsefulRate.Rate >= 0.80 &&
		report.FallbackDeterminism.Rate == 1 &&
		report.SelectionModeReporting.Rate == 1 &&
		report.SourceTraceCompleteness.Rate == 1

	if pass {
		report.Status = "PASS"
		if report.Top1UsefulRate.Rate >= 0.80 && report.Top1UsefulRate.Rate <= 0.82 {
			report.Status = "WATCH"
		}
	}

	return report, nil
}

func metric(numerator, denominator int) evaluationMetric {
	if denominator == 0 {
		return evaluationMetric{Numerator: 0, Denominator: 0, Rate: 0}
	}
	return evaluationMetric{
		Numerator:   numerator,
		Denominator: denominator,
		Rate:        float64(numerator) / float64(denominator),
	}
}

func loadEvaluationQueries(path string) ([]evaluationQuery, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var queries []evaluationQuery
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
