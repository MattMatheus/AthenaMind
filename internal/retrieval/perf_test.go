package retrieval

import (
	"fmt"
	"testing"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func seedBenchmarkCorpus(b *testing.B, root string, n int) {
	b.Helper()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "bench",
		Risk:     "low",
	}
	for i := 0; i < n; i++ {
		if err := index.UpsertEntry(root, types.UpsertEntryInput{
			ID:     fmt.Sprintf("entry-%03d", i),
			Title:  fmt.Sprintf("Entry %03d", i),
			Type:   "prompt",
			Domain: "perf",
			Body:   fmt.Sprintf("retrieval benchmark body %03d with repeated tokens network rollback incident", i),
			Stage:  "pm",
		}, policy); err != nil {
			b.Fatalf("seed entry %d: %v", i, err)
		}
	}
}

func BenchmarkRetrieveClassicColdCache(b *testing.B) {
	root := b.TempDir()
	seedBenchmarkCorpus(b, root, 200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clearRetrievalCaches()
		if _, _, err := RetrieveWithOptionsAndEndpointAndSession(
			root,
			"incident rollback network",
			"perf",
			"http://127.0.0.1:1",
			"",
			RetrieveOptions{Mode: "classic", Backend: "sqlite", TopK: 5},
		); err != nil {
			b.Fatalf("retrieve failed: %v", err)
		}
	}
}

func BenchmarkRetrieveClassicWarmCache(b *testing.B) {
	root := b.TempDir()
	seedBenchmarkCorpus(b, root, 200)
	clearRetrievalCaches()
	if _, _, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		"incident rollback network",
		"perf",
		"http://127.0.0.1:1",
		"",
		RetrieveOptions{Mode: "classic", Backend: "sqlite", TopK: 5},
	); err != nil {
		b.Fatalf("warmup retrieve failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := RetrieveWithOptionsAndEndpointAndSession(
			root,
			"incident rollback network",
			"perf",
			"http://127.0.0.1:1",
			"",
			RetrieveOptions{Mode: "classic", Backend: "sqlite", TopK: 5},
		); err != nil {
			b.Fatalf("retrieve failed: %v", err)
		}
	}
}
