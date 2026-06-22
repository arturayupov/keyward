package vault

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/arturayupov/keyward/internal/model"
)

// Save serializes the store to JSON and writes it age-encrypted to path (0600).
// The write is atomic (temp file + rename) so an interrupted write can never
// truncate or corrupt an existing vault.
func Save(path string, store *model.Store, recipient age.Recipient) error {
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".vault-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op once renamed; cleans up on any error path
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// Load decrypts and parses the vault at path.
func Load(path string, id age.Identity) (*model.Store, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := age.Decrypt(f, id)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var store model.Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}
