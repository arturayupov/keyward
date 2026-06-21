package vault

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"filippo.io/age"
	"github.com/arturayupov/keyward/internal/model"
)

// Save serializes the store to JSON and writes it age-encrypted to path (0600).
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
	return os.WriteFile(path, buf.Bytes(), 0o600)
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
