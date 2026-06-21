package envfile

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	got, err := Parse([]byte("# c\nexport A=1\nB=\"two words\"\n\nC='x'\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"A": "1", "B": "two words", "C": "x"}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("%s: got %q want %q", k, got[k], v)
		}
	}
}

func TestSetCreatesFile0600(t *testing.T) {
	p := filepath.Join(t.TempDir(), ".env")
	if err := Set(p, "API_KEY", "sekret"); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(p)
	// Unix mode bits aren't honored on Windows (files report -rw-rw-rw-);
	// ACL-based hardening on Windows is a v1 item.
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatalf("perm = %v", info.Mode().Perm())
	}
	m, _ := ParseFile(p)
	if m["API_KEY"] != "sekret" {
		t.Fatalf("value not written")
	}
}

func TestSetReplacesExistingKeyNotDuplicate(t *testing.T) {
	p := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(p, []byte("API_KEY=old\nOTHER=keep\n"), 0o600)
	_ = Set(p, "API_KEY", "new")
	m, _ := ParseFile(p)
	if m["API_KEY"] != "new" || m["OTHER"] != "keep" {
		t.Fatalf("got %+v", m)
	}
	data, _ := os.ReadFile(p)
	if strings.Count(string(data), "API_KEY=") != 1 {
		t.Fatalf("API_KEY duplicated:\n%s", data)
	}
}
