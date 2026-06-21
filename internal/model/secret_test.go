package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUpsertAndFind(t *testing.T) {
	var s Store
	s.Upsert(Secret{Name: "API", Namespace: "proj", Value: "v1"})
	s.Upsert(Secret{Name: "API", Namespace: "proj", Value: "v2"}) // same key → replace
	if len(s.Secrets) != 1 {
		t.Fatalf("want 1 secret, got %d", len(s.Secrets))
	}
	got, ok := s.Find("API", "proj")
	if !ok || got.Value != "v2" {
		t.Fatalf("want v2, got %+v ok=%v", got, ok)
	}
}

func TestMetaHasNoValue(t *testing.T) {
	var s Store
	s.Upsert(Secret{Name: "API", Namespace: "proj", Value: "supersecret"})
	metas := s.Meta("")
	if len(metas) != 1 || metas[0].Name != "API" {
		t.Fatalf("unexpected metas: %+v", metas)
	}
	b, _ := json.Marshal(metas)
	if strings.Contains(string(b), "supersecret") {
		t.Fatalf("SecretMeta JSON leaked a value: %s", b)
	}
}
