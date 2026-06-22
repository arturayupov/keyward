//go:build manual

// Manual, GUI-driven test. Excluded from normal builds/CI by the `manual` tag.
// Run explicitly: go test -tags manual -run TestLiveApprovalDialog ./internal/approval/ -v
package approval_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/model"
)

const demoValue = "sk_live_DEMO_d4f9_not_a_real_key"

func TestLiveApprovalDialog(t *testing.T) {
	appr, ok := approval.NativeApprover()
	if !ok {
		t.Skip("no native approver on this OS")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, ".env")
	store := &model.Store{Secrets: []model.Secret{
		{Name: "STRIPE_KEY", Namespace: "acme-shop", Value: demoValue},
	}}
	b := &broker.Broker{
		Store:    store,
		Approver: appr,
		Audit:    audit.New(filepath.Join(dir, "audit.jsonl")),
	}

	res, err := b.Request(approval.Request{
		Tool: "claude-code", Name: "STRIPE_KEY", Namespace: "acme-shop",
		Target: target, Reason: "start the dev server",
	})
	if err != nil {
		t.Fatalf("request error: %v", err)
	}

	data, _ := os.ReadFile(target)
	fmt.Printf("\n========== LIVE DIALOG RESULT ==========\n")
	fmt.Printf("broker result: %+v\n", res)
	fmt.Printf("target file contents: %q\n", string(data))
	switch {
	case res.Status == "injected" && strings.Contains(string(data), demoValue):
		fmt.Printf("VERDICT: APPROVED — key injected into target .env ✓\n")
	case res.Status == "denied" && len(data) == 0:
		fmt.Printf("VERDICT: DENIED — nothing injected ✓\n")
	default:
		fmt.Printf("VERDICT: UNEXPECTED — status=%q fileLen=%d\n", res.Status, len(data))
	}
	fmt.Printf("========================================\n")
}
