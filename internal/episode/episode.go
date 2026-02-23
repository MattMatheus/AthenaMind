package episode

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func Write(root string, in types.WriteEpisodeInput, policy types.WritePolicyDecision) (types.EpisodeRecord, error) {
	repo := normalizeKey(in.Repo)
	sessionID := strings.TrimSpace(in.SessionID)
	cycleID := strings.TrimSpace(in.CycleID)
	storyID := strings.TrimSpace(in.StoryID)
	outcome := strings.ToLower(strings.TrimSpace(in.Outcome))
	if repo == "" || sessionID == "" || cycleID == "" || storyID == "" || outcome == "" {
		return types.EpisodeRecord{}, errors.New("--repo --session-id --cycle-id --story-id --outcome are required")
	}
	if outcome != "success" && outcome != "partial" && outcome != "blocked" {
		return types.EpisodeRecord{}, errors.New("--outcome must be success|partial|blocked")
	}

	summary, err := loadInlineOrFile(in.Summary, in.SummaryFile)
	if err != nil {
		return types.EpisodeRecord{}, err
	}
	if strings.TrimSpace(summary) == "" {
		return types.EpisodeRecord{}, errors.New("episode summary is required via --summary or --summary-file")
	}
	decisions, err := loadInlineOrFile(in.Decisions, in.DecisionsFile)
	if err != nil {
		return types.EpisodeRecord{}, err
	}
	if strings.TrimSpace(decisions) == "" {
		return types.EpisodeRecord{}, errors.New("episode decisions are required via --decisions or --decisions-file")
	}

	files := []string{}
	for _, part := range strings.Split(in.FilesChanged, ",") {
		v := strings.TrimSpace(part)
		if v != "" {
			files = append(files, v)
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	record := types.EpisodeRecord{
		ID:           fmt.Sprintf("episode-%d", time.Now().UTC().UnixNano()),
		Repo:         repo,
		SessionID:    sessionID,
		CycleID:      cycleID,
		StoryID:      storyID,
		Outcome:      outcome,
		Summary:      strings.TrimSpace(summary),
		FilesChanged: files,
		Decisions:    strings.TrimSpace(decisions),
		CreatedAt:    now,
	}

	dir := filepath.Join(root, "episodes", repo)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return types.EpisodeRecord{}, err
	}
	recordPath := filepath.Join(dir, record.ID+".json")
	if err := index.WriteJSONAsYAML(recordPath, record); err != nil {
		return types.EpisodeRecord{}, err
	}

	latest := types.EpisodeContext{
		Repo:      repo,
		Scenario:  "episode",
		CycleID:   record.CycleID,
		StoryID:   record.StoryID,
		Outcome:   record.Outcome,
		Summary:   record.Summary,
		Timestamp: record.CreatedAt,
	}
	if err := index.WriteJSONAsYAML(filepath.Join(dir, "latest.json"), latest); err != nil {
		return types.EpisodeRecord{}, err
	}

	entryBody := strings.TrimSpace(fmt.Sprintf("Outcome: %s\nCycle: %s\nStory: %s\nSummary: %s\nDecisions: %s\nFiles Changed: %s", record.Outcome, record.CycleID, record.StoryID, record.Summary, record.Decisions, strings.Join(record.FilesChanged, ", ")))
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:       record.ID,
		Title:    fmt.Sprintf("Episode %s", record.CycleID),
		Type:     "episode",
		Domain:   repo,
		Body:     entryBody,
		BodyFile: "",
		Stage:    in.Stage,
	}, policy); err != nil {
		return types.EpisodeRecord{}, err
	}

	return record, nil
}

func List(root, repo string) ([]types.EpisodeRecord, error) {
	repo = normalizeKey(repo)
	if strings.TrimSpace(repo) == "" {
		return nil, errors.New("--repo is required")
	}
	dir := filepath.Join(root, "episodes", repo)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return []types.EpisodeRecord{}, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	rows := []types.EpisodeRecord{}
	for _, ent := range entries {
		if ent.IsDir() || ent.Name() == "latest.json" || !strings.HasSuffix(ent.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, ent.Name()))
		if err != nil {
			return nil, err
		}
		var rec types.EpisodeRecord
		if err := json.Unmarshal(data, &rec); err != nil {
			return nil, err
		}
		rows = append(rows, rec)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt > rows[j].CreatedAt })
	return rows, nil
}

func loadInlineOrFile(inline, path string) (string, error) {
	v := strings.TrimSpace(inline)
	if strings.TrimSpace(path) != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		v = strings.TrimSpace(string(data))
	}
	return v, nil
}

func normalizeKey(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	v = strings.ReplaceAll(v, " ", "-")
	return strings.Trim(v, "-")
}
