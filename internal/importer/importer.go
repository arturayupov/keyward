package importer

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/arturayupov/keyward/internal/envfile"
	"github.com/arturayupov/keyward/internal/model"
)

var skipDirs = map[string]bool{
	"node_modules": true, ".git": true, "Library": true,
	".vscode": true, ".gemini": true, ".antigravity": true,
	"venv": true, ".venv": true, "site-packages": true,
}

// Import walks root for .env files and returns secrets, deriving the namespace
// from the path segment directly under root.
func Import(root string) ([]model.Secret, error) {
	var out []model.Secret
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // tolerate unreadable dirs
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if !isEnvFile(d.Name()) {
			return nil
		}
		ns := namespaceFor(root, path)
		kv, err := envfile.ParseFile(path)
		if err != nil {
			return nil
		}
		for k, v := range kv {
			if v == "" {
				continue
			}
			out = append(out, model.Secret{Name: k, Value: v, Namespace: ns})
		}
		return nil
	})
	return out, err
}

func isEnvFile(name string) bool {
	if strings.Contains(name, ".example") || strings.Contains(name, ".template") {
		return false
	}
	return name == ".env" || strings.HasSuffix(name, ".env")
}

func namespaceFor(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "default"
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 2 {
		return parts[0]
	}
	return "default"
}
