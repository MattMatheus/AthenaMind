package retrieval

import (
	"fmt"
	"os"
	"path/filepath"

	"athenamind/internal/index"
)

func IndexEntryEmbedding(root, entryID, embeddingEndpoint string) (string, error) {
	warnings, err := IndexEntriesEmbeddingBatch(root, []string{entryID}, embeddingEndpoint)
	if err != nil {
		return "", err
	}
	if len(warnings) > 0 {
		return warnings[0], nil
	}
	return "", nil
}

func IndexEntriesEmbeddingBatch(root string, entryIDs []string, embeddingEndpoint string) ([]string, error) {
	if len(entryIDs) == 0 {
		return nil, nil
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return nil, err
	}

	bodies := make([]string, 0, len(entryIDs))
	validIDs := make([]string, 0, len(entryIDs))
	var warnings []string

	for _, id := range entryIDs {
		var sourcePath string
		for _, e := range idx.Entries {
			if e.ID == id {
				sourcePath = e.Path
				break
			}
		}
		if sourcePath == "" {
			warnings = append(warnings, fmt.Sprintf("entry %s not found for embedding", id))
			continue
		}

		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(sourcePath)))
		if err != nil {
			return nil, err
		}
		bodies = append(bodies, string(data))
		validIDs = append(validIDs, id)
	}

	if len(bodies) == 0 {
		return warnings, nil
	}

	vecs, err := GenerateEmbeddings(embeddingEndpoint, bodies)
	if err != nil {
		for _, id := range validIDs {
			warnings = append(warnings, fmt.Sprintf("embedding unavailable; entry %s stored without vector (%v)", id, err))
		}
		return warnings, nil
	}

	for i, id := range validIDs {
		if err := index.UpsertEmbedding(root, id, vecs[i]); err != nil {
			return nil, err
		}
	}

	return warnings, nil
}
