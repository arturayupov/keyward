package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/vault"
)

type approveAll struct{}

func (approveAll) Approve(approval.Request) (approval.Decision, error) {
	return approval.ApproveOnce, nil
}

// TestSentinelNeverLeaks is the cross-package security invariant: a known
// sentinel value must appear ONLY in the encrypted vault (as ciphertext) and
// the injected target — never in stdout, the audit log, or the broker result.
func TestSentinelNeverLeaks(t *testing.T) {
	const sentinel = "SENTINEL-d4f9-do-not-leak"
	home := t.TempDir()
	proj := filepath.Join(home, "demoproj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, ".env"), []byte("DEMO_TOKEN="+sentinel+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)

	// init → import → ls via the real CLI; the sentinel must never be printed.
	for _, args := range [][]string{{"init"}, {"import", home}, {"ls"}} {
		out, err := runCLI(t, home, args...)
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if strings.Contains(out, sentinel) {
			t.Fatalf("sentinel leaked to stdout of %v:\n%s", args, out)
		}
	}

	// inject via the broker directly (CLI inject would raise a native dialog).
	p, err := config.Default()
	if err != nil {
		t.Fatal(err)
	}
	id, err := vault.EnsureIdentity()
	if err != nil {
		t.Fatal(err)
	}
	store, err := vault.Load(p.Vault, id)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(home, "out", ".env")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	b := &broker.Broker{Store: store, Approver: approveAll{}, Audit: audit.New(p.Audit)}
	res, err := b.Request(approval.Request{Tool: "test", Name: "DEMO_TOKEN", Namespace: "demoproj", Target: target})
	if err != nil || res.Status != "injected" {
		t.Fatalf("inject failed: %+v err=%v", res, err)
	}

	// The vault holds the sentinel ONLY as ciphertext.
	if vb, _ := os.ReadFile(p.Vault); bytes.Contains(vb, []byte(sentinel)) {
		t.Fatal("sentinel found as plaintext in vault.age")
	}
	// The audit log must never contain the value.
	if ab, _ := os.ReadFile(p.Audit); bytes.Contains(ab, []byte(sentinel)) {
		t.Fatal("sentinel found in audit.jsonl")
	}
	// The broker result must never contain the value.
	if rb, _ := json.Marshal(res); strings.Contains(string(rb), sentinel) {
		t.Fatalf("sentinel found in broker result: %s", rb)
	}
	// The target IS where the value is supposed to land.
	if tb, _ := os.ReadFile(target); !bytes.Contains(tb, []byte(sentinel)) {
		t.Fatal("sentinel was not injected into the target")
	}
}
