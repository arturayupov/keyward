package config

import (
	"path/filepath"
	"testing"
)

func TestPaths(t *testing.T) {
	p := PathsFor("/home/u")
	if p.Dir != filepath.Join("/home/u", ".keyward") {
		t.Fatalf("dir: %s", p.Dir)
	}
	if filepath.Base(p.Vault) != "vault.age" {
		t.Fatalf("vault: %s", p.Vault)
	}
	if filepath.Base(p.Audit) != "audit.jsonl" {
		t.Fatalf("audit: %s", p.Audit)
	}
}
