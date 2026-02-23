package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"athenamind/internal/types"
)

const (
	EventSchema  = "1.0"
	TelemetryRel = "telemetry/events.jsonl"
)

func Emit(root, telemetryPath string, ev types.TelemetryEvent) error {
	path := strings.TrimSpace(telemetryPath)
	if path == "" {
		path = filepath.Join(root, filepath.FromSlash(TelemetryRel))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func TelemetryErrorCode(err error) string {
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "ERR_UNKNOWN"
	}
	first := strings.Fields(msg)
	code := strings.TrimSuffix(first[0], ":")
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "ERR_") {
		return code
	}
	return "ERR_COMMAND_FAILED"
}

func NormalizeMemoryType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "procedural", "state", "semantic":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "semantic"
	}
}

func NormalizeOperatorVerdict(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "correct", "partially_correct", "incorrect", "not_scored":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "not_scored"
	}
}

func NormalizeTelemetryValue(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}
