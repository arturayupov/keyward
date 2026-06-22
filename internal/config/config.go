package config

import (
	"os"
	"path/filepath"
)

// Paths holds the on-disk locations keyward uses, all under ~/.keyward.
type Paths struct {
	Dir   string
	Vault string
	Audit string
	Conf  string
}

// PathsFor builds the path set rooted at the given home directory.
func PathsFor(home string) Paths {
	dir := filepath.Join(home, ".keyward")
	return Paths{
		Dir:   dir,
		Vault: filepath.Join(dir, "vault.age"),
		Audit: filepath.Join(dir, "audit.jsonl"),
		Conf:  filepath.Join(dir, "config.toml"),
	}
}

// Default resolves paths from the user's home directory.
func Default() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, err
	}
	return PathsFor(home), nil
}

// EnsureDir creates ~/.keyward with 0700 perms.
func (p Paths) EnsureDir() error { return os.MkdirAll(p.Dir, 0o700) }
