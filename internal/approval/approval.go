package approval

import (
	"fmt"
	"sync"
)

// Decision is the outcome of an approval request.
type Decision int

const (
	Deny Decision = iota
	ApproveOnce
	ApproveSession
)

// Request describes who wants which key for what target. Carries no value.
type Request struct {
	Tool      string
	Name      string
	Namespace string
	Target    string
	Reason    string
}

// Prompt renders a one-line human description of the request.
func (r Request) Prompt() string {
	return fmt.Sprintf("%s requests %q (project %q) → %s", r.Tool, r.Name, r.Namespace, r.Target)
}

// Approver decides whether to release a key. Implementations must fail closed
// (return Deny) on error or cancellation.
type Approver interface {
	Approve(Request) (Decision, error)
}

// NativeApprover returns the OS-native dialog approver if one is available.
func NativeApprover() (Approver, bool) { return nativeApprover() }

// SessionCache remembers ApproveSession decisions per {tool,namespace,name}
// for the process lifetime, so repeat requests don't re-prompt.
type SessionCache struct {
	inner Approver
	mu    sync.Mutex
	ok    map[string]bool
}

// NewSessionCache wraps an Approver with session-scoped memoization.
func NewSessionCache(inner Approver) *SessionCache {
	return &SessionCache{inner: inner, ok: map[string]bool{}}
}

func cacheKey(r Request) string { return r.Tool + "\x00" + r.Namespace + "\x00" + r.Name }

// Approve returns a cached ApproveSession without re-prompting; otherwise it
// delegates and caches only ApproveSession outcomes.
func (c *SessionCache) Approve(r Request) (Decision, error) {
	c.mu.Lock()
	cached := c.ok[cacheKey(r)]
	c.mu.Unlock()
	if cached {
		return ApproveSession, nil
	}
	d, err := c.inner.Approve(r)
	if err != nil {
		return Deny, err
	}
	if d == ApproveSession {
		c.mu.Lock()
		c.ok[cacheKey(r)] = true
		c.mu.Unlock()
	}
	return d, nil
}
