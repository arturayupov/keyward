package vault

import (
	"errors"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestEnsureIdentityWrapsKeystoreError(t *testing.T) {
	keyring.MockInitWithError(errors.New("boom"))
	t.Cleanup(func() { keyring.MockInit() })
	_, err := EnsureIdentity()
	if err == nil {
		t.Fatal("expected error when keystore is unreachable")
	}
	if !strings.Contains(err.Error(), "OS keystore") {
		t.Fatalf("error lacks guidance: %v", err)
	}
}

func TestEnsureIdentityRoundTrip(t *testing.T) {
	keyring.MockInit() // in-memory keyring for tests
	id1, err := EnsureIdentity()
	if err != nil {
		t.Fatalf("ensure: %v", err)
	}
	id2, err := EnsureIdentity() // second call returns the same stored identity
	if err != nil {
		t.Fatalf("ensure2: %v", err)
	}
	if id1.String() != id2.String() {
		t.Fatalf("identity not stable across calls")
	}
}
