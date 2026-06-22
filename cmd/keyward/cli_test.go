package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	keyring.MockInit() // in-memory keyring so tests touch no real keychain
	os.Exit(m.Run())
}

func runCLI(t *testing.T, home string, args ...string) (string, error) {
	t.Helper()
	cmd := newRootCmd()
	var out strings.Builder
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	t.Setenv("HOME", home)
	err := cmd.Execute()
	return out.String(), err
}

func TestLsPrintsNamesNotValues(t *testing.T) {
	home := t.TempDir()
	proj := filepath.Join(home, "demo")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, ".env"), []byte("TOKEN=supersecret\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := runCLI(t, home, "init"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := runCLI(t, home, "import", home); err != nil {
		t.Fatalf("import: %v", err)
	}
	out, err := runCLI(t, home, "ls")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "TOKEN") || !strings.Contains(out, "demo") {
		t.Fatalf("ls missing name/namespace: %s", out)
	}
	if strings.Contains(out, "supersecret") {
		t.Fatalf("ls leaked a value: %s", out)
	}
}
