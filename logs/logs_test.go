package logs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"sync"
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

// serviceFromFuncName is pure string parsing — exercise both module-path
// shapes upstream services actually use (confirmed by checking their go.mod
// files directly): the github.com/scrambledeggs/X convention, and a bare
// module name with no scrambledeggs/ prefix at all.
func TestServiceFromFuncName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"scrambledeggs prefix", "github.com/scrambledeggs/booky-athena/functions/GetCollectionByIDAdminV1/handler.Handler", "booky-athena"},
		{"bare module", "booky-freyja/functions/NearbyVouchersV1/handler.Handler", "booky-freyja"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := serviceFromFuncName(c.in); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

// Regression for the data race a reviewer caught on the previous approach
// (a package-level Service variable written by the caller before invoking a
// handler): concurrent callers must not share any mutable state. This alone
// doesn't prove correct attribution across two different real callers (that's
// covered by TestInfo_AttributesServiceFromCaller in service_test.go), but it
// is the exact failure mode `go test -race` caught — run with -race to verify.
func TestLogIt_ConcurrentCallsDoNotRace(t *testing.T) {
	const n = 50

	lines := captureOutput(t, func() {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Info("concurrent note")
			}()
		}
		wg.Wait()
	})

	if len(lines) != n {
		t.Fatalf("expected %d lines, got %d", n, len(lines))
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
