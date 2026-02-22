package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	supportedMajor = 1
	defaultSchema  = "1.0"
)

type indexFile struct {
	SchemaVersion string       `json:"schema_version"`
	UpdatedAt     string       `json:"updated_at"`
	Entries       []indexEntry `json:"entries"`
}

type indexEntry struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	MetadataPath string `json:"metadata_path"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updated_at"`
	Title        string `json:"title"`
}

type metadataFile struct {
	SchemaVersion string     `json:"schema_version"`
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	UpdatedAt     string     `json:"updated_at"`
	Review        reviewMeta `json:"review"`
}

type reviewMeta struct {
	ReviewedBy   string `json:"reviewed_by"`
	ReviewedAt   string `json:"reviewed_at"`
	Decision     string `json:"decision"`
	DecisionNote string `json:"decision_notes"`
}

type retrieveResult struct {
	SelectedID    string  `json:"selected_id"`
	SelectionMode string  `json:"selection_mode"`
	SourcePath    string  `json:"source_path"`
	Confidence    float64 `json:"confidence"`
	Reason        string  `json:"reason"`
}

type candidate struct {
	Entry indexEntry
	Body  string
	Score float64
}

func main() {
	if len(os.Args) < 2 {
		exitErr(errors.New("usage: memory-cli <write|retrieve> [flags]"))
	}

	var err error
	switch os.Args[1] {
	case "write":
		err = runWrite(os.Args[2:])
	case "retrieve":
		err = runRetrieve(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		exitErr(err)
	}
}

func runWrite(args []string) error {
	fs := flag.NewFlagSet("write", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	id := fs.String("id", "", "entry id")
	title := fs.String("title", "", "entry title")
	typeValue := fs.String("type", "", "entry type: prompt|instruction")
	domain := fs.String("domain", "", "entry domain")
	body := fs.String("body", "", "entry body")
	bodyFile := fs.String("body-file", "", "path to markdown body")
	stage := fs.String("stage", "", "workflow stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")	
	approved := fs.Bool("approved", false, "set true when change is reviewed/approved")
	notes := fs.String("notes", "", "review notes")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := enforceWritePolicy(*stage, *reviewer, *approved); err != nil {
		return err
	}

	if *id == "" || *title == "" || *typeValue == "" || *domain == "" {
		return errors.New("--id --title --type --domain are required")
	}
	if *typeValue != "prompt" && *typeValue != "instruction" {
		return errors.New("--type must be prompt or instruction")
	}

	entryBody := strings.TrimSpace(*body)
	if *bodyFile != "" {
		data, err := os.ReadFile(*bodyFile)
		if err != nil {
			return err
		}
		entryBody = strings.TrimSpace(string(data))
	}
	if entryBody == "" {
		return errors.New("entry body is required via --body or --body-file")
	}

	idx, err := loadIndex(*root)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	dirName := "prompts"
	if *typeValue == "instruction" {
		dirName = "instructions"
	}

	relContentPath := filepath.ToSlash(filepath.Join(dirName, *domain, *id+".md"))
	relMetaPath := filepath.ToSlash(filepath.Join("metadata", *id+".yaml"))
	contentPath := filepath.Join(*root, filepath.FromSlash(relContentPath))
	metaPath := filepath.Join(*root, filepath.FromSlash(relMetaPath))

	if err := os.MkdirAll(filepath.Dir(contentPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(metaPath), 0o755); err != nil {
		return err
	}

	markdown := fmt.Sprintf("# %s\n\n%s\n", *title, entryBody)
	if err := os.WriteFile(contentPath, []byte(markdown), 0o644); err != nil {
		return err
	}

	meta := metadataFile{
		SchemaVersion: defaultSchema,
		ID:            *id,
		Title:         *title,
		Status:        "approved",
		UpdatedAt:     now,
		Review: reviewMeta{
			ReviewedBy:   *reviewer,
			ReviewedAt:   now,
			Decision:     "approved",
			DecisionNote: *notes,
		},
	}
	if err := writeJSONAsYAML(metaPath, meta); err != nil {
		return err
	}

	updated := false
	for i := range idx.Entries {
		if idx.Entries[i].ID == *id {
			idx.Entries[i] = indexEntry{
				ID:           *id,
				Type:         *typeValue,
				Domain:       *domain,
				Path:         relContentPath,
				MetadataPath: relMetaPath,
				Status:       "approved",
				UpdatedAt:    now,
				Title:        *title,
			}
			updated = true
			break
		}
	}
	if !updated {
		idx.Entries = append(idx.Entries, indexEntry{
			ID:           *id,
			Type:         *typeValue,
			Domain:       *domain,
			Path:         relContentPath,
			MetadataPath: relMetaPath,
			Status:       "approved",
			UpdatedAt:    now,
			Title:        *title,
		})
	}
	sort.Slice(idx.Entries, func(i, j int) bool { return idx.Entries[i].ID < idx.Entries[j].ID })
	idx.UpdatedAt = now
	if err := writeJSONAsYAML(filepath.Join(*root, "index.yaml"), idx); err != nil {
		return err
	}

	fmt.Printf("wrote entry %s at %s\n", *id, relContentPath)
	return nil
}

func runRetrieve(args []string) error {
	fs := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "", "natural language query")
	domain := fs.String("domain", "", "optional domain filter")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*query) == "" {
		return errors.New("--query is required")
	}

	idx, err := loadIndex(*root)
	if err != nil {
		return err
	}
	if len(idx.Entries) == 0 {
		return errors.New("memory index has no entries")
	}

	candidates, err := loadCandidates(*root, idx.Entries, *domain)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		return errors.New("no candidates found for query/domain")
	}

	q := strings.ToLower(strings.TrimSpace(*query))
	for i := range candidates {
		candidates[i].Score = semanticScore(q, candidates[i])
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].Entry.ID < candidates[j].Entry.ID
		}
		return candidates[i].Score > candidates[j].Score
	})

	top := candidates[0]
	second := 0.0
	if len(candidates) > 1 {
		second = candidates[1].Score
	}

	if isSemanticConfident(top.Score, second) {
		return printResult(retrieveResult{
			SelectedID:    top.Entry.ID,
			SelectionMode: "semantic",
			SourcePath:    top.Entry.Path,
			Confidence:    top.Score,
			Reason:        "semantic confidence gate passed",
		})
	}

	for _, c := range candidates {
		if strings.EqualFold(c.Entry.ID, q) {
			return printResult(retrieveResult{
				SelectedID:    c.Entry.ID,
				SelectionMode: "fallback_exact_key",
				SourcePath:    c.Entry.Path,
				Confidence:    c.Score,
				Reason:        "semantic confidence gate failed; exact-key fallback matched",
			})
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Entry.Path < candidates[j].Entry.Path
	})
	chosen := candidates[0]
	return printResult(retrieveResult{
		SelectedID:    chosen.Entry.ID,
		SelectionMode: "fallback_path_priority",
		SourcePath:    chosen.Entry.Path,
		Confidence:    chosen.Score,
		Reason:        "semantic confidence gate failed; deterministic path-priority fallback used",
	})
}

func printResult(r retrieveResult) error {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func loadCandidates(root string, entries []indexEntry, domain string) ([]candidate, error) {
	candidates := make([]candidate, 0, len(entries))
	for _, e := range entries {
		if domain != "" && e.Domain != domain {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(e.Path)))
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate{Entry: e, Body: string(data)})
	}
	return candidates, nil
}

func semanticScore(query string, c candidate) float64 {
	qTokens := tokenSet(query)
	if len(qTokens) == 0 {
		return 0
	}

	hay := strings.ToLower(strings.Join([]string{c.Entry.Title, c.Entry.ID, c.Entry.Domain, c.Body}, " "))
	hits := 0
	for tok := range qTokens {
		if strings.Contains(hay, tok) {
			hits++
		}
	}
	return float64(hits) / float64(len(qTokens))
}

func tokenSet(s string) map[string]struct{} {
	clean := strings.NewReplacer(".", " ", ",", " ", ":", " ", ";", " ", "/", " ", "-", " ", "_", " ").Replace(strings.ToLower(s))
	parts := strings.Fields(clean)
	out := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		if len(p) > 1 {
			out[p] = struct{}{}
		}
	}
	return out
}

func isSemanticConfident(top, second float64) bool {
	const minConfidence = 0.34
	const minMargin = 0.15
	if top < minConfidence {
		return false
	}
	if top-second < minMargin {
		return false
	}
	return true
}

func enforceWritePolicy(stage, reviewer string, approved bool) error {
	if isTrue(os.Getenv("AUTONOMOUS_RUN")) {
		return errors.New("ERR_MUTATION_NOT_ALLOWED_DURING_AUTONOMOUS_RUN: writes are blocked during autonomous runs")
	}
	allowed := map[string]struct{}{"planning": {}, "architect": {}, "pm": {}}
	if _, ok := allowed[stage]; !ok {
		return errors.New("ERR_MUTATION_STAGE_INVALID: --stage must be planning, architect, or pm")
	}
	if strings.TrimSpace(reviewer) == "" || !approved {
		return errors.New("ERR_MUTATION_REVIEW_REQUIRED: reviewed writes require --reviewer and --approved=true")
	}
	return nil
}

func isTrue(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes"
}

func loadIndex(root string) (indexFile, error) {
	path := filepath.Join(root, "index.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(root, 0o755); err != nil {
				return indexFile{}, err
			}
			return indexFile{SchemaVersion: defaultSchema, UpdatedAt: time.Now().UTC().Format(time.RFC3339), Entries: []indexEntry{}}, nil
		}
		return indexFile{}, err
	}

	var idx indexFile
	if err := json.Unmarshal(data, &idx); err != nil {
		return indexFile{}, fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: cannot parse %s: %w", path, err)
	}
	if strings.TrimSpace(idx.SchemaVersion) == "" {
		return indexFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: index schema_version is required")
	}
	if err := validateSchemaVersion(idx.SchemaVersion); err != nil {
		return indexFile{}, err
	}
	for _, e := range idx.Entries {
		if e.ID == "" || e.Type == "" || e.Domain == "" || e.Path == "" || e.MetadataPath == "" || e.Status == "" || e.UpdatedAt == "" {
			return indexFile{}, errors.New("ERR_SCHEMA_VALIDATION: index entry missing required fields")
		}
	}
	return idx, nil
}

func validateSchemaVersion(version string) error {
	major, _, err := parseMajorMinor(version)
	if err != nil {
		return fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: %w", err)
	}
	if major > supportedMajor {
		return fmt.Errorf("ERR_SCHEMA_MAJOR_UNSUPPORTED: schema version %s is newer than supported major %d", version, supportedMajor)
	}
	return nil
}

func parseMajorMinor(v string) (int, int, error) {
	parts := strings.Split(strings.TrimSpace(v), ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("version must be MAJOR.MINOR")
	}
	var major, minor int
	_, err := fmt.Sscanf(v, "%d.%d", &major, &minor)
	if err != nil {
		return 0, 0, fmt.Errorf("version must contain numeric MAJOR.MINOR")
	}
	return major, minor, nil
}

func writeJSONAsYAML(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
