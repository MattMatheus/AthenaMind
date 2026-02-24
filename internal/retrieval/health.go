package retrieval

import (
	"athenamind/internal/index"
)

type SemanticHealthReport struct {
	IndexedEntries    int
	StoredEmbeddings  int
	MissingEmbeddings int
	CoverageOK        bool
	SelectionMode     string
	SemanticAvailable bool
	Warning           string
	SelectedID        string
	Confidence        float64
	Pass              bool
}

func EvaluateSemanticHealth(root, query, domain, endpoint, sessionID string) (SemanticHealthReport, error) {
	idx, err := index.LoadIndex(root)
	if err != nil {
		return SemanticHealthReport{}, err
	}
	embeddings, err := index.GetEmbeddingRecords(root, nil)
	if err != nil {
		return SemanticHealthReport{}, err
	}

	missing := 0
	for _, e := range idx.Entries {
		if rec, ok := embeddings[e.ID]; !ok || len(rec.Vector) == 0 {
			missing++
		}
	}

	result, warning, err := RetrieveWithEmbeddingEndpointAndSession(root, query, domain, endpoint, sessionID)
	if err != nil {
		return SemanticHealthReport{}, err
	}
	semantic := result.SelectionMode == "embedding_semantic" ||
		result.SelectionMode == "semantic" ||
		result.SelectionMode == "hybrid_rrf"
	coverageOK := missing == 0
	return SemanticHealthReport{
		IndexedEntries:    len(idx.Entries),
		StoredEmbeddings:  len(embeddings),
		MissingEmbeddings: missing,
		CoverageOK:        coverageOK,
		SelectionMode:     result.SelectionMode,
		SemanticAvailable: semantic,
		Warning:           warning,
		SelectedID:        result.SelectedID,
		Confidence:        result.Confidence,
		Pass:              coverageOK && semantic,
	}, nil
}
