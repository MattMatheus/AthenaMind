package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"athenamind/internal/types"
)

const sqliteFileName = "index.db"

type sqliteIndexStore struct{}

func (sqliteIndexStore) Load(root string) (types.IndexFile, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return types.IndexFile{}, err
	}

	dbPath := filepath.Join(root, sqliteFileName)
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		legacy, legacyErr := loadIndexFromYAML(root)
		if legacyErr != nil {
			return types.IndexFile{}, legacyErr
		}
		if err := initSQLite(root); err != nil {
			return types.IndexFile{}, err
		}
		if err := saveIndexToSQLite(root, legacy); err != nil {
			return types.IndexFile{}, err
		}
		return legacy, nil
	}

	if err := initSQLite(root); err != nil {
		return types.IndexFile{}, err
	}

	idx, err := readIndexFromSQLite(root)
	if err != nil {
		return types.IndexFile{}, err
	}
	if err := ValidateSchemaVersion(idx.SchemaVersion); err != nil {
		return types.IndexFile{}, err
	}
	if err := ValidateIndex(idx, root); err != nil {
		return types.IndexFile{}, err
	}
	return idx, nil
}

func (sqliteIndexStore) Save(root string, idx types.IndexFile) error {
	if err := ValidateSchemaVersion(idx.SchemaVersion); err != nil {
		return err
	}
	if err := initSQLite(root); err != nil {
		return err
	}
	return saveIndexToSQLite(root, idx)
}

func initSQLite(root string) error {
	sql := `
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout=5000;
CREATE TABLE IF NOT EXISTS meta (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  schema_version TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS entries (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  domain TEXT NOT NULL,
  path TEXT NOT NULL,
  metadata_path TEXT NOT NULL,
  status TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  title TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS embeddings (
  entry_id TEXT PRIMARY KEY,
  vector_json TEXT NOT NULL,
  model_id TEXT NOT NULL DEFAULT '',
  provider TEXT NOT NULL DEFAULT '',
  dim INTEGER NOT NULL DEFAULT 0,
  content_hash TEXT NOT NULL DEFAULT '',
  commit_sha TEXT NOT NULL DEFAULT '',
  session_id TEXT NOT NULL DEFAULT '',
  generated_at TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
);
INSERT INTO meta (id, schema_version, updated_at)
VALUES (1, '1.0', strftime('%Y-%m-%dT%H:%M:%SZ','now'))
ON CONFLICT(id) DO NOTHING;
`
	if _, err := runSQLite(root, sql, false); err != nil {
		return err
	}
	for _, col := range []string{
		"model_id TEXT NOT NULL DEFAULT ''",
		"provider TEXT NOT NULL DEFAULT ''",
		"dim INTEGER NOT NULL DEFAULT 0",
		"content_hash TEXT NOT NULL DEFAULT ''",
		"commit_sha TEXT NOT NULL DEFAULT ''",
		"session_id TEXT NOT NULL DEFAULT ''",
		"generated_at TEXT NOT NULL DEFAULT ''",
	} {
		stmt := fmt.Sprintf("ALTER TABLE embeddings ADD COLUMN %s;", col)
		if _, err := runSQLite(root, stmt, false); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return err
		}
	}
	return nil
}

func readIndexFromSQLite(root string) (types.IndexFile, error) {
	type row struct {
		SchemaVersion string `json:"schema_version"`
		UpdatedAt     string `json:"updated_at"`
	}
	metaJSON, err := runSQLite(root, "SELECT schema_version, updated_at FROM meta WHERE id=1;", true)
	if err != nil {
		return types.IndexFile{}, err
	}
	metaRows := []row{}
	if strings.TrimSpace(metaJSON) != "" {
		if err := json.Unmarshal([]byte(metaJSON), &metaRows); err != nil {
			return types.IndexFile{}, fmt.Errorf("ERR_SCHEMA_VALIDATION: cannot parse sqlite meta rows: %w", err)
		}
	}

	idx := types.IndexFile{
		SchemaVersion: types.DefaultSchema,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Entries:       []types.IndexEntry{},
	}
	if len(metaRows) > 0 {
		idx.SchemaVersion = metaRows[0].SchemaVersion
		idx.UpdatedAt = metaRows[0].UpdatedAt
	}

	entryJSON, err := runSQLite(root, "SELECT id, type, domain, path, metadata_path, status, updated_at, title FROM entries ORDER BY id ASC;", true)
	if err != nil {
		return types.IndexFile{}, err
	}
	entryRows := []types.IndexEntry{}
	if strings.TrimSpace(entryJSON) != "" {
		if err := json.Unmarshal([]byte(entryJSON), &entryRows); err != nil {
			return types.IndexFile{}, fmt.Errorf("ERR_SCHEMA_VALIDATION: cannot parse sqlite entries: %w", err)
		}
	}
	idx.Entries = entryRows
	return idx, nil
}

func saveIndexToSQLite(root string, idx types.IndexFile) error {
	schemaVersion := sqlQuote(idx.SchemaVersion)
	updatedAt := sqlQuote(idx.UpdatedAt)
	if updatedAt == "''" {
		updatedAt = sqlQuote(time.Now().UTC().Format(time.RFC3339))
	}
	entries := append([]types.IndexEntry(nil), idx.Entries...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	stmts := []string{
		fmt.Sprintf("UPDATE meta SET schema_version=%s, updated_at=%s WHERE id=1;", schemaVersion, updatedAt),
		"DELETE FROM entries;",
	}
	for _, e := range entries {
		stmts = append(stmts, fmt.Sprintf(
			"INSERT INTO entries (id, type, domain, path, metadata_path, status, updated_at, title) VALUES (%s,%s,%s,%s,%s,%s,%s,%s);",
			sqlQuote(e.ID),
			sqlQuote(e.Type),
			sqlQuote(e.Domain),
			sqlQuote(e.Path),
			sqlQuote(e.MetadataPath),
			sqlQuote(e.Status),
			sqlQuote(e.UpdatedAt),
			sqlQuote(e.Title),
		))
	}
	_, err := runSQLiteTx(root, stmts)
	return err
}

func runSQLite(root, sql string, jsonMode bool) (string, error) {
	dbPath := filepath.Join(root, sqliteFileName)
	args := []string{}
	if !jsonMode {
		args = append(args, "-cmd", "PRAGMA busy_timeout=5000;")
	}
	if jsonMode {
		args = append(args, "-json")
	}
	args = append(args, dbPath, sql)
	cmd := exec.Command("sqlite3", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sqlite command failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func runSQLiteTx(root string, statements []string) (string, error) {
	if len(statements) == 0 {
		return "", nil
	}
	var script strings.Builder
	script.WriteString("BEGIN IMMEDIATE;\n")
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		script.WriteString(trimmed)
		if !strings.HasSuffix(trimmed, ";") {
			script.WriteString(";")
		}
		script.WriteString("\n")
	}
	script.WriteString("COMMIT;")
	return runSQLite(root, script.String(), false)
}

func sqlQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

func UpsertEmbedding(root, entryID string, vector []float64) error {
	return UpsertEmbeddingRecord(root, types.EmbeddingRecord{
		EntryID:     entryID,
		Vector:      vector,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

func UpsertEmbeddingRecord(root string, record types.EmbeddingRecord) error {
	entryID := strings.TrimSpace(record.EntryID)
	vector := record.Vector
	if strings.TrimSpace(entryID) == "" {
		return errors.New("entry id is required for embedding upsert")
	}
	if len(vector) == 0 {
		return errors.New("embedding vector cannot be empty")
	}
	if err := initSQLite(root); err != nil {
		return err
	}
	raw, err := json.Marshal(vector)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	dim := record.Dim
	if dim <= 0 {
		dim = len(vector)
	}
	generatedAt := strings.TrimSpace(record.GeneratedAt)
	if generatedAt == "" {
		generatedAt = now
	}
	stmt := fmt.Sprintf(
		"INSERT INTO embeddings (entry_id, vector_json, model_id, provider, dim, content_hash, commit_sha, session_id, generated_at, updated_at) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) ON CONFLICT(entry_id) DO UPDATE SET vector_json=excluded.vector_json, model_id=excluded.model_id, provider=excluded.provider, dim=excluded.dim, content_hash=excluded.content_hash, commit_sha=excluded.commit_sha, session_id=excluded.session_id, generated_at=excluded.generated_at, updated_at=excluded.updated_at;",
		sqlQuote(entryID),
		sqlQuote(string(raw)),
		sqlQuote(strings.TrimSpace(record.ModelID)),
		sqlQuote(strings.TrimSpace(record.Provider)),
		strconv.Itoa(dim),
		sqlQuote(strings.TrimSpace(record.ContentHash)),
		sqlQuote(strings.TrimSpace(record.CommitSHA)),
		sqlQuote(strings.TrimSpace(record.SessionID)),
		sqlQuote(generatedAt),
		sqlQuote(now),
	)
	_, err = runSQLiteTx(root, []string{stmt})
	return err
}

func GetEmbeddings(root string, ids []string) (map[string][]float64, error) {
	records, err := GetEmbeddingRecords(root, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]float64, len(records))
	for id, rec := range records {
		out[id] = rec.Vector
	}
	return out, nil
}

func GetEmbeddingRecords(root string, ids []string) (map[string]types.EmbeddingRecord, error) {
	out := map[string]types.EmbeddingRecord{}
	if err := initSQLite(root); err != nil {
		return nil, err
	}
	query := "SELECT entry_id, vector_json, model_id, provider, dim, content_hash, commit_sha, session_id, generated_at, updated_at FROM embeddings;"
	if len(ids) > 0 {
		values := make([]string, 0, len(ids))
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			values = append(values, sqlQuote(id))
		}
		if len(values) == 0 {
			return out, nil
		}
		query = fmt.Sprintf("SELECT entry_id, vector_json, model_id, provider, dim, content_hash, commit_sha, session_id, generated_at, updated_at FROM embeddings WHERE entry_id IN (%s);", strings.Join(values, ","))
	}
	raw, err := runSQLite(root, query, true)
	if err != nil {
		return nil, err
	}
	type row struct {
		EntryID     string `json:"entry_id"`
		VectorJSON  string `json:"vector_json"`
		ModelID     string `json:"model_id"`
		Provider    string `json:"provider"`
		Dim         int    `json:"dim"`
		ContentHash string `json:"content_hash"`
		CommitSHA   string `json:"commit_sha"`
		SessionID   string `json:"session_id"`
		GeneratedAt string `json:"generated_at"`
		UpdatedAt   string `json:"updated_at"`
	}
	rows := []row{}
	if strings.TrimSpace(raw) == "" {
		return out, nil
	}
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	for _, r := range rows {
		vec := []float64{}
		if err := json.Unmarshal([]byte(r.VectorJSON), &vec); err != nil {
			continue
		}
		out[r.EntryID] = types.EmbeddingRecord{
			EntryID:     r.EntryID,
			Vector:      vec,
			ModelID:     r.ModelID,
			Provider:    r.Provider,
			Dim:         r.Dim,
			ContentHash: r.ContentHash,
			CommitSHA:   r.CommitSHA,
			SessionID:   r.SessionID,
			GeneratedAt: r.GeneratedAt,
			LastUpdated: r.UpdatedAt,
		}
	}
	return out, nil
}
