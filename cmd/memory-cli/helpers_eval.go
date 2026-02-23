package main

import "athenamind/internal/retrieval"

func evaluateRetrieval(root string, queries []evaluationQuery, corpusID, querySetID, configID string) (evaluationReport, error) {
	return retrieval.EvaluateRetrieval(root, queries, corpusID, querySetID, configID)
}

func loadEvaluationQueries(path string) ([]evaluationQuery, error) {
	return retrieval.LoadEvaluationQueries(path)
}
