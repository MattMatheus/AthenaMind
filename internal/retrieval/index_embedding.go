package retrieval

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func IndexEntryEmbedding(root, entryID, embeddingEndpoint, sessionID string) (string, error) {
	warnings, err := IndexEntriesEmbeddingBatch(root, []string{entryID}, embeddingEndpoint, sessionID)
	if err != nil {
		return "", err
	}
	if len(warnings) > 0 {
		return warnings[0], nil
	}
	return "", nil
}

func IndexEntriesEmbeddingBatch(root string, entryIDs []string, embeddingEndpoint, sessionID string) ([]string, error) {
	if len(entryIDs) == 0 {
		return nil, nil
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return nil, err
	}

	bodies := make([]string, 0, len(entryIDs))
	validIDs := make([]string, 0, len(entryIDs))
	contentHashes := make([]string, 0, len(entryIDs))
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
		body := string(data)
		bodies = append(bodies, body)
		validIDs = append(validIDs, id)
		contentHashes = append(contentHashes, sha256Hex(body))
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

	profile := ActiveEmbeddingProfile(embeddingEndpoint)
	commitSHA := currentCommitSHA()
	generatedAt := time.Now().UTC().Format(time.RFC3339)
	for i, id := range validIDs {
		if err := index.UpsertEmbeddingRecord(root, types.EmbeddingRecord{
			EntryID:     id,
			Vector:      vecs[i],
			ModelID:     profile.ModelID,
			Provider:    profile.Provider,
			Dim:         len(vecs[i]),
			ContentHash: contentHashes[i],
			CommitSHA:   commitSHA,
			SessionID:   strings.TrimSpace(sessionID),
			GeneratedAt: generatedAt,
		}); err != nil {
			return nil, err
		}
	}

	return warnings, nil
}

func sha256Hex(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}

func currentCommitSHA() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
