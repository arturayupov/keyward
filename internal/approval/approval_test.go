package approval

import "testing"

type stubApprover struct {
	calls int
	ret   Decision
}

func (s *stubApprover) Approve(Request) (Decision, error) { s.calls++; return s.ret, nil }

func TestSessionCacheAvoidsRepeatPrompt(t *testing.T) {
	stub := &stubApprover{ret: ApproveSession}
	c := NewSessionCache(stub)
	r := Request{Tool: "claude", Name: "K", Namespace: "n", Target: "./.env"}
	if d, _ := c.Approve(r); d != ApproveSession {
		t.Fatal("first call should approve")
	}
	if d, _ := c.Approve(r); d != ApproveSession {
		t.Fatal("second call should be cached-approved")
	}
	if stub.calls != 1 {
		t.Fatalf("underlying approver called %d times, want 1", stub.calls)
	}
}

func TestApproveOnceNotCached(t *testing.T) {
	stub := &stubApprover{ret: ApproveOnce}
	c := NewSessionCache(stub)
	r := Request{Tool: "claude", Name: "K", Namespace: "n"}
	_, _ = c.Approve(r)
	_, _ = c.Approve(r)
	if stub.calls != 2 {
		t.Fatalf("ApproveOnce must not cache; calls=%d", stub.calls)
	}
}
