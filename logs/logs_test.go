package logs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// captureOutput redirects both os.Stdout and the log package's output to a
// pipe for the duration of fn, then returns the captured output split into
// individual lines (log.Print always appends its own trailing newline).
func captureOutput(t *testing.T, fn func()) []string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	origStdout := os.Stdout
	os.Stdout = w
	log.SetOutput(w)

	fn()

	w.Close()
	os.Stdout = origStdout
	log.SetOutput(os.Stdout)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	trimmed := strings.TrimRight(buf.String(), "\n")
	if trimmed == "" {
		return nil
	}

	return strings.Split(trimmed, "\n")
}

// Railway (and any line-oriented log collector) requires one JSON object per
// line. Debug used to be pretty-printed via json.MarshalIndent, which split a
// single log entry across many lines and broke level-based filtering.
func TestLogIt_EmitsSingleLineJSON(t *testing.T) {
	t.Setenv("APP_ENV", PRODUCTION_ENV)

	for _, level := range []struct {
		name string
		log  func(note string, data ...any)
	}{
		{"Debug", Debug},
		{"Info", Info},
		{"Warn", Warn},
		{"Error", Error},
	} {
		t.Run(level.name, func(t *testing.T) {
			lines := captureOutput(t, func() {
				level.log("note", map[string]any{"nested": map[string]any{"foo": "bar"}})
			})

			if len(lines) != 1 {
				t.Fatalf("expected exactly 1 line of output, got %d: %v", len(lines), lines)
			}

			var entry logEntry
			if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
				t.Fatalf("output is not valid single-line JSON: %v\nline: %s", err, lines[0])
			}
		})
	}
}

// Railway's structured-log parser requires a top-level "message" field.
func TestLogIt_SetsMessage(t *testing.T) {
	lines := captureOutput(t, func() {
		Info("info note")
	})

	if len(lines) != 1 {
		t.Fatalf("expected exactly 1 line of output, got %d: %v", len(lines), lines)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if raw["message"] != "info note" {
		t.Errorf(`"message" = %v, want %q`, raw["message"], "info note")
	}
}

// Service lets a caller (e.g. booky-ark's apigwadapter) attribute a log line
// to the upstream microservice that produced it, since many microservices'
// handlers run inside one shared process.
func TestLogIt_SetsService(t *testing.T) {
	t.Cleanup(func() { Service = "" })
	Service = "booky-athena"

	lines := captureOutput(t, func() {
		Info("info note")
	})

	if len(lines) != 1 {
		t.Fatalf("expected exactly 1 line of output, got %d: %v", len(lines), lines)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if raw["service"] != "booky-athena" {
		t.Errorf(`"service" = %v, want %q`, raw["service"], "booky-athena")
	}
}

// Print never runs in production, so it never reaches Railway/CloudWatch —
// pretty-printing it for local terminal readability is safe.
func TestPrint_PrettyPrintsLocallyOnly(t *testing.T) {
	t.Run("pretty-prints outside production", func(t *testing.T) {
		t.Setenv("APP_ENV", "development")

		lines := captureOutput(t, func() {
			Print("print note", map[string]any{"foo": "bar"})
		})

		if len(lines) <= 1 {
			t.Fatalf("expected multi-line pretty-printed output, got %d line(s): %v", len(lines), lines)
		}
	})

	t.Run("suppressed in production", func(t *testing.T) {
		t.Setenv("APP_ENV", PRODUCTION_ENV)

		lines := captureOutput(t, func() {
			Print("print note")
		})

		if len(lines) != 0 {
			t.Fatalf("expected no output in production, got: %v", lines)
		}
	})
}
