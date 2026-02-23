package retrieval

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"athenamind/internal/index"
)

func IndexEntryEmbedding(root, entryID, embeddingEndpoint string) (string, error) {
	entryID = strings.TrimSpace(entryID)
	if entryID == "" {
		return "", errors.New("entry id is required")
	}
	idx, err := index.LoadIndex(root)
	if err != nil {
		return "", err
	}

	var sourcePath string
	for _, e := range idx.Entries {
		if e.ID == entryID {
			sourcePath = e.Path
			break
		}
	}
	if sourcePath == "" {
		return "", fmt.Errorf("entry %s not found for embedding", entryID)
	}

	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(sourcePath)))
	if err != nil {
		return "", err
	}
	vec, err := GenerateEmbedding(embeddingEndpoint, string(data))
	if err != nil {
		return fmt.Sprintf("embedding unavailable; entry %s stored without vector (%v)", entryID, err), nil
	}
	if err := index.UpsertEmbedding(root, entryID, vec); err != nil {
		return "", err
	}
	return "", nil
}
