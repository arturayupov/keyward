package broker

import (
	"path/filepath"
	"testing"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/envfile"
	"github.com/arturayupov/keyward/internal/model"
)

type fixedApprover struct{ d approval.Decision }

func (f fixedApprover) Approve(approval.Request) (approval.Decision, error) { return f.d, nil }

func newBroker(t *testing.T, d approval.Decision) (*Broker, string) {
	dir := t.TempDir()
	store := &model.Store{Secrets: []model.Secret{{Name: "API", Namespace: "proj", Value: "sekret"}}}
	b := &Broker{
		Store:    store,
		Approver: fixedApprover{d: d},
		Audit:    audit.New(filepath.Join(dir, "audit.jsonl")),
	}
	return b, dir
}

func TestRequestApprovedInjects(t *testing.T) {
	b, dir := newBroker(t, approval.ApproveOnce)
	target := filepath.Join(dir, ".env")
	res, err := b.Request(approval.Request{Tool: "claude", Name: "API", Namespace: "proj", Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "injected" {
		t.Fatalf("status=%s", res.Status)
	}
	m, _ := envfile.ParseFile(target)
	if m["API"] != "sekret" {
		t.Fatalf("not injected: %+v", m)
	}
}

func TestRequestDeniedInjectsNothing(t *testing.T) {
	b, dir := newBroker(t, approval.Deny)
	target := filepath.Join(dir, ".env")
	res, _ := b.Request(approval.Request{Tool: "claude", Name: "API", Namespace: "proj", Target: target})
	if res.Status != "denied" {
		t.Fatalf("status=%s", res.Status)
	}
	if _, err := envfile.ParseFile(target); err == nil {
		t.Fatal("target file should not exist on deny")
	}
}

func TestResultHasNoValueField(t *testing.T) {
	b, dir := newBroker(t, approval.ApproveOnce)
	res, _ := b.Request(approval.Request{Tool: "claude", Name: "API", Namespace: "proj", Target: filepath.Join(dir, ".env")})
	if res.Name != "API" || res.Status == "" {
		t.Fatalf("unexpected result %+v", res)
	}
}
