package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendRedactsValue(t *testing.T) {
	p := filepath.Join(t.TempDir(), "audit.jsonl")
	a := New(p)
	if err := a.Record(Entry{Tool: "claude", Name: "API_KEY", Namespace: "n", Target: "./.env", Decision: "approved_once"}); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if !strings.Contains(string(b), `"name":"API_KEY"`) {
		t.Fatalf("missing name in audit: %s", b)
	}
	if strings.Contains(string(b), "value") {
		t.Fatalf("audit entry must have no value field: %s", b)
	}
}
