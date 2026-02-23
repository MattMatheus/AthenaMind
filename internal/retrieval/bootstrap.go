package retrieval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func Bootstrap(root, repo, sessionID, scenario string) (types.BootstrapPayload, error) {
	payload := types.BootstrapPayload{
		Repo:          strings.TrimSpace(repo),
		SessionID:     strings.TrimSpace(sessionID),
		Scenario:      strings.TrimSpace(scenario),
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		MemoryEntries: []types.BootstrapMemoryEntry{},
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return payload, err
	}

	query := strings.TrimSpace(repo + " " + scenario)
	tokens := tokenSet(query)
	if len(tokens) > 0 {
		type scoredEntry struct {
			entry types.IndexEntry
			score float64
		}
		scored := []scoredEntry{}
		for _, e := range idx.Entries {
			if e.Status != "approved" || e.Type != "instruction" {
				continue
			}
			hay := strings.ToLower(strings.Join([]string{e.ID, e.Title, e.Domain, e.Path}, " "))
			hits := 0
			for tok := range tokens {
				if strings.Contains(hay, tok) {
					hits++
				}
			}
			score := float64(hits) / float64(len(tokens))
			if score > 0 {
				scored = append(scored, scoredEntry{entry: e, score: score})
			}
		}
		sort.SliceStable(scored, func(i, j int) bool {
			if scored[i].score == scored[j].score {
				return scored[i].entry.ID < scored[j].entry.ID
			}
			return scored[i].score > scored[j].score
		})
		for i, s := range scored {
			if i >= 5 {
				break
			}
			payload.MemoryEntries = append(payload.MemoryEntries, types.BootstrapMemoryEntry{
				ID:            s.entry.ID,
				SelectionMode: "semantic",
				SourcePath:    s.entry.Path,
				Confidence:    s.score,
				Reason:        "bootstrap contextual match",
				Type:          s.entry.Type,
				Domain:        s.entry.Domain,
				Title:         s.entry.Title,
			})
		}
	}

	if ep := loadLatestEpisode(root, repo, scenario); ep != nil {
		payload.Episode = ep
	}

	return payload, nil
}

func loadLatestEpisode(root, repo, scenario string) *types.EpisodeContext {
	repoKey := normalizeKey(repo)
	scenarioKey := normalizeKey(scenario)
	path := filepath.Join(root, "episodes", repoKey, scenarioKey, "latest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var ep types.EpisodeContext
	if err := json.Unmarshal(data, &ep); err != nil {
		return nil
	}
	if strings.TrimSpace(ep.Repo) == "" {
		ep.Repo = strings.TrimSpace(repo)
	}
	if strings.TrimSpace(ep.Scenario) == "" {
		ep.Scenario = strings.TrimSpace(scenario)
	}
	return &ep
}

func normalizeKey(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return "unknown"
	}
	re := regexp.MustCompile(`[^a-z0-9._-]+`)
	out := re.ReplaceAllString(v, "-")
	out = strings.Trim(out, "-")
	if out == "" {
		return "unknown"
	}
	return out
}
