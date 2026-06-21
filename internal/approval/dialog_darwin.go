//go:build darwin

package approval

import (
	"fmt"
	"os/exec"
	"strings"
)

func nativeApprover() (Approver, bool) { return osascriptApprover{}, true }

type osascriptApprover struct{}

func (osascriptApprover) Approve(r Request) (Decision, error) {
	msg := fmt.Sprintf("%s\n\nReason: %s", r.Prompt(), r.Reason)
	script := fmt.Sprintf(`display dialog %q with title "keyward" buttons {"Deny","Approve once","Approve for session"} default button "Deny"`, msg)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return Deny, nil // cancel / error → deny (fail closed)
	}
	switch {
	case strings.Contains(string(out), "Approve for session"):
		return ApproveSession, nil
	case strings.Contains(string(out), "Approve once"):
		return ApproveOnce, nil
	default:
		return Deny, nil
	}
}
