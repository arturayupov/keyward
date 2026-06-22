package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportWalksAndNamespaces(t *testing.T) {
	root := t.TempDir()
	mk := func(rel, body string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(body), 0o600)
	}
	mk("shop-app/.env", "SHOPIFY_TOKEN=abc\n")
	mk("alice/.env.local", "OPENAI_API_KEY=xyz\n")   // .env.local must be imported
	mk("alice/node_modules/pkg/.env", "IGNORED=1\n") // must be skipped
	mk("alice/.env.example", "OPENAI_API_KEY=\n")    // must be skipped

	secs, err := Import(root)
	if err != nil {
		t.Fatal(err)
	}
	byNS := map[string]string{}
	for _, s := range secs {
		byNS[s.Namespace+"/"+s.Name] = s.Value
	}
	if byNS["shop-app/SHOPIFY_TOKEN"] != "abc" || byNS["alice/OPENAI_API_KEY"] != "xyz" {
		t.Fatalf("unexpected import: %+v", byNS)
	}
	if len(secs) != 2 {
		t.Fatalf("want 2 secrets (examples/node_modules skipped), got %d", len(secs))
	}
}
