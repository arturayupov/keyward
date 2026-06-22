//go:build linux

package approval

import (
	"fmt"
	"os/exec"
)

func nativeApprover() (Approver, bool) {
	if _, err := exec.LookPath("zenity"); err == nil {
		return zenityApprover{}, true
	}
	return nil, false // caller falls back to terminal
}

type zenityApprover struct{}

func (zenityApprover) Approve(r Request) (Decision, error) {
	text := fmt.Sprintf("%s\n\nReason: %s", r.Prompt(), r.Reason)
	// OK = approve once, extra button = approve for session, close/cancel = deny.
	cmd := exec.Command("zenity", "--question", "--title=keyward",
		"--text="+text, "--ok-label=Approve once", "--cancel-label=Deny",
		"--extra-button=Approve for session")
	out, err := cmd.Output()
	if string(out) == "Approve for session\n" {
		return ApproveSession, nil
	}
	if err == nil {
		return ApproveOnce, nil
	}
	return Deny, nil
}
