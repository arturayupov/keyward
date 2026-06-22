package vault

import (
	"errors"
	"fmt"

	"filippo.io/age"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "keyward"
	keyringAccount = "vault-age-identity"
)

// keystoreHint explains how to recover when the OS keystore is unreachable.
const keystoreHint = "could not reach the OS keystore (macOS Keychain / Windows Credential Manager / Linux libsecret). " +
	"On Linux install a Secret Service (e.g. gnome-keyring) and unlock a session; on macOS run keyward under your normal login (not an overridden $HOME)"

// EnsureIdentity returns the vault's age identity from the OS keystore,
// generating and storing a new one on first use.
func EnsureIdentity() (*age.X25519Identity, error) {
	s, err := keyring.Get(keyringService, keyringAccount)
	if err == nil {
		return age.ParseX25519Identity(s)
	}
	if !errors.Is(err, keyring.ErrNotFound) {
		return nil, fmt.Errorf("%s: %w", keystoreHint, err)
	}
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	if err := keyring.Set(keyringService, keyringAccount, id.String()); err != nil {
		return nil, fmt.Errorf("%s: %w", keystoreHint, err)
	}
	return id, nil
}
