// Package logs_test (external, black-box) is required here rather than the
// internal "package logs" test file: testservice imports logs, so an
// internal test file importing testservice would be a cyclic import.
package logs_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/scrambledeggs/booky-go-common/logs/testservice"
)

// End-to-end proof that callerService's stack walk correctly identifies a
// real external caller (as opposed to the unit tests in logs_test.go, which
// only exercise the pure string-parsing half of the derivation).
func TestInfo_AttributesServiceFromCaller(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	origStdout := os.Stdout
	os.Stdout = w
	log.SetOutput(w)

	testservice.CallInfo("hello from testservice")

	w.Close()
	os.Stdout = origStdout
	log.SetOutput(os.Stdout)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	// testservice lives inside the booky-go-common repo itself, so the
	// scrambledeggs-segment rule attributes it to the repo — the same
	// outcome booky-ark's own inline (non-upstream-service) handlers get.
	if raw["service"] != "booky-go-common" {
		t.Errorf(`"service" = %v, want %q`, raw["service"], "booky-go-common")
	}
}
