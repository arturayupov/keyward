package broker

import (
	"errors"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/envfile"
	"github.com/arturayupov/keyward/internal/model"
)

// ErrNotFound is returned when no secret matches the request.
var ErrNotFound = errors.New("secret not found")

// Broker orchestrates request → approval → inject → audit.
type Broker struct {
	Store    *model.Store
	Approver approval.Approver
	Audit    *audit.Logger
}

// Result is intentionally value-free.
type Result struct {
	Status string `json:"status"` // "injected" | "denied"
	Name   string `json:"name"`
	Target string `json:"target"`
}

// Request resolves a secret, seeks approval, and on approval injects it into
// the target env file. It never returns the secret value.
func (b *Broker) Request(r approval.Request) (Result, error) {
	sec, ok := b.Store.Find(r.Name, r.Namespace)
	if !ok {
		return Result{}, ErrNotFound
	}
	d, err := b.Approver.Approve(r)
	if err != nil {
		return Result{}, err
	}
	if d == approval.Deny {
		_ = b.Audit.Record(audit.Entry{Tool: r.Tool, Name: r.Name, Namespace: r.Namespace, Target: r.Target, Decision: "denied"})
		return Result{Status: "denied", Name: r.Name, Target: r.Target}, nil
	}
	if err := envfile.Set(r.Target, sec.Name, sec.Value); err != nil {
		return Result{}, err
	}
	decision := "approved_once"
	if d == approval.ApproveSession {
		decision = "approved_session"
	}
	_ = b.Audit.Record(audit.Entry{Tool: r.Tool, Name: r.Name, Namespace: r.Namespace, Target: r.Target, Decision: decision})
	return Result{Status: "injected", Name: r.Name, Target: r.Target}, nil
}
