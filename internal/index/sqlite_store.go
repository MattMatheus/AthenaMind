package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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
  updated_at TEXT NOT NULL
);
INSERT INTO meta (id, schema_version, updated_at)
VALUES (1, '1.0', strftime('%Y-%m-%dT%H:%M:%SZ','now'))
ON CONFLICT(id) DO NOTHING;
`
	_, err := runSQLite(root, sql, false)
	return err
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
	if _, err := runSQLite(root, fmt.Sprintf("UPDATE meta SET schema_version=%s, updated_at=%s WHERE id=1;", schemaVersion, updatedAt), false); err != nil {
		return err
	}
	if _, err := runSQLite(root, "DELETE FROM entries;", false); err != nil {
		return err
	}

	entries := append([]types.IndexEntry(nil), idx.Entries...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	for _, e := range entries {
		stmt := fmt.Sprintf(
			"INSERT INTO entries (id, type, domain, path, metadata_path, status, updated_at, title) VALUES (%s,%s,%s,%s,%s,%s,%s,%s);",
			sqlQuote(e.ID),
			sqlQuote(e.Type),
			sqlQuote(e.Domain),
			sqlQuote(e.Path),
			sqlQuote(e.MetadataPath),
			sqlQuote(e.Status),
			sqlQuote(e.UpdatedAt),
			sqlQuote(e.Title),
		)
		if _, err := runSQLite(root, stmt, false); err != nil {
			return err
		}
	}
	return nil
}

func runSQLite(root, sql string, jsonMode bool) (string, error) {
	dbPath := filepath.Join(root, sqliteFileName)
	args := []string{}
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

func sqlQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

func UpsertEmbedding(root, entryID string, vector []float64) error {
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
	stmt := fmt.Sprintf(
		"INSERT INTO embeddings (entry_id, vector_json, updated_at) VALUES (%s,%s,%s) ON CONFLICT(entry_id) DO UPDATE SET vector_json=excluded.vector_json, updated_at=excluded.updated_at;",
		sqlQuote(strings.TrimSpace(entryID)),
		sqlQuote(string(raw)),
		sqlQuote(now),
	)
	_, err = runSQLite(root, stmt, false)
	return err
}

func GetEmbeddings(root string, ids []string) (map[string][]float64, error) {
	out := map[string][]float64{}
	if len(ids) == 0 {
		return out, nil
	}
	if err := initSQLite(root); err != nil {
		return nil, err
	}
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
	query := fmt.Sprintf("SELECT entry_id, vector_json FROM embeddings WHERE entry_id IN (%s);", strings.Join(values, ","))
	raw, err := runSQLite(root, query, true)
	if err != nil {
		return nil, err
	}
	type row struct {
		EntryID    string `json:"entry_id"`
		VectorJSON string `json:"vector_json"`
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
			return nil, err
		}
		out[r.EntryID] = vec
	}
	return out, nil
}
