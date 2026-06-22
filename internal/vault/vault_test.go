package vault

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/arturayupov/keyward/internal/model"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "vault.age")
	in := &model.Store{Secrets: []model.Secret{{Name: "K", Namespace: "n", Value: "topsecret"}}}
	if err := Save(path, in, id.Recipient()); err != nil {
		t.Fatalf("save: %v", err)
	}
	out, err := Load(path, id)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got, _ := out.Find("K", "n"); got.Value != "topsecret" {
		t.Fatalf("roundtrip mismatch: %+v", got)
	}
}

func TestFileIsEncrypted(t *testing.T) {
	id, _ := age.GenerateX25519Identity()
	path := filepath.Join(t.TempDir(), "vault.age")
	_ = Save(path, &model.Store{Secrets: []model.Secret{{Name: "K", Value: "topsecret"}}}, id.Recipient())
	raw, _ := os.ReadFile(path)
	if bytes.Contains(raw, []byte("topsecret")) {
		t.Fatal("plaintext secret found in vault file")
	}
}
