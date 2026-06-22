package main

import (
	"strings"
	"testing"
)

func TestAddReadsValueFromStdin(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if _, err := runCLI(t, home, "init"); err != nil {
		t.Fatalf("init: %v", err)
	}

	cmd := newRootCmd()
	var out strings.Builder
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader("super-secret-value\n"))
	cmd.SetArgs([]string{"add", "NPM_TOKEN", "--ns", "cli-tools"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}

	lo, err := runCLI(t, home, "ls")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(lo, "NPM_TOKEN") || !strings.Contains(lo, "cli-tools") {
		t.Fatalf("ls missing added key: %s", lo)
	}
	if strings.Contains(lo, "super-secret-value") {
		t.Fatalf("ls leaked the value: %s", lo)
	}
}
