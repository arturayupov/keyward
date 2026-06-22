package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/model"
)

type yes struct{}

func (yes) Approve(approval.Request) (approval.Decision, error) { return approval.ApproveOnce, nil }

func TestListKeysReturnsNoValues(t *testing.T) {
	st := &model.Store{Secrets: []model.Secret{{Name: "API", Namespace: "n", Value: "sekret"}}}
	h := &Handlers{Store: st}
	out, err := h.listKeys(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "API") || strings.Contains(out, "sekret") {
		t.Fatalf("list leaked or missing: %s", out)
	}
}

func TestRequestKeyReturnsConfirmationNotValue(t *testing.T) {
	dir := t.TempDir()
	st := &model.Store{Secrets: []model.Secret{{Name: "API", Namespace: "n", Value: "sekret"}}}
	h := &Handlers{Store: st, Broker: &broker.Broker{Store: st, Approver: yes{}, Audit: audit.New(dir + "/a.jsonl")}}
	out, err := h.requestKey(context.Background(), "API", "n", dir+"/.env", "need it")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "sekret") {
		t.Fatalf("request_key leaked value: %s", out)
	}
	if !strings.Contains(out, "injected") {
		t.Fatalf("expected confirmation, got: %s", out)
	}
}
