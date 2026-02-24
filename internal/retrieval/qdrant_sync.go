package retrieval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

type QdrantSyncReport struct {
	Root       string   `json:"root"`
	URL        string   `json:"url"`
	Collection string   `json:"collection"`
	Indexed    int      `json:"indexed"`
	Synced     int      `json:"synced"`
	Skipped    int      `json:"skipped"`
	Warnings   []string `json:"warnings,omitempty"`
}

func SyncQdrantCollection(root, qdrantURL, collection string, batchSize int) (QdrantSyncReport, error) {
	if strings.TrimSpace(root) == "" {
		return QdrantSyncReport{}, fmt.Errorf("root is required")
	}
	if strings.TrimSpace(qdrantURL) == "" {
		qdrantURL = strings.TrimSpace(os.Getenv("ATHENA_QDRANT_URL"))
	}
	if strings.TrimSpace(qdrantURL) == "" {
		qdrantURL = "http://localhost:6333"
	}
	if strings.TrimSpace(collection) == "" {
		collection = strings.TrimSpace(os.Getenv("ATHENA_QDRANT_COLLECTION"))
	}
	if strings.TrimSpace(collection) == "" {
		collection = "athena_memories"
	}
	if batchSize <= 0 {
		batchSize = 128
	}

	report := QdrantSyncReport{
		Root:       root,
		URL:        qdrantURL,
		Collection: collection,
		Warnings:   []string{},
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return report, err
	}
	report.Indexed = len(idx.Entries)

	recs, err := index.GetEmbeddingRecords(root, nil)
	if err != nil {
		return report, err
	}
	entryByID := map[string]types.IndexEntry{}
	for _, entry := range idx.Entries {
		entryByID[entry.ID] = entry
	}

	points := make([]map[string]any, 0, len(recs))
	for id, rec := range recs {
		if len(rec.Vector) == 0 {
			report.Skipped++
			continue
		}
		entry := entryByID[id]
		payload := map[string]any{
			"entry_id": id,
			"id":       id,
			"domain":   entry.Domain,
			"title":    entry.Title,
			"path":     entry.Path,
			"type":     entry.Type,
		}
		points = append(points, map[string]any{
			"id":      id,
			"vector":  rec.Vector,
			"payload": payload,
		})
	}

	for start := 0; start < len(points); start += batchSize {
		end := start + batchSize
		if end > len(points) {
			end = len(points)
		}
		if err := qdrantUpsert(qdrantURL, collection, points[start:end]); err != nil {
			return report, err
		}
		report.Synced += end - start
	}
	if report.Synced == 0 {
		report.Warnings = append(report.Warnings, "no embeddings found to sync")
	}
	return report, nil
}

func qdrantUpsert(qdrantURL, collection string, points []map[string]any) error {
	if len(points) == 0 {
		return nil
	}
	payload := map[string]any{"points": points}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := strings.TrimRight(qdrantURL, "/") + "/collections/" + collection + "/points?wait=true"
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if key := strings.TrimSpace(os.Getenv("ATHENA_QDRANT_API_KEY")); key != "" {
		req.Header.Set("api-key", key)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		msg := strings.TrimSpace(string(data))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("qdrant upsert failed: status %d: %s", resp.StatusCode, msg)
	}
	return nil
}
