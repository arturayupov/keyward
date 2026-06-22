//go:build windows

package approval

import (
	"fmt"
	"os/exec"
	"strings"
)

func nativeApprover() (Approver, bool) { return msgBoxApprover{}, true }

type msgBoxApprover struct{}

func (msgBoxApprover) Approve(r Request) (Decision, error) {
	msg := fmt.Sprintf("%s`n`nReason: %s`n`nYes = approve for session, No = approve once, Cancel = deny", r.Prompt(), r.Reason)
	ps := fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework; [System.Windows.MessageBox]::Show(%q,'keyward','YesNoCancel','Question')`, msg)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", ps).Output()
	if err != nil {
		return Deny, nil
	}
	switch strings.TrimSpace(string(out)) {
	case "Yes":
		return ApproveSession, nil
	case "No":
		return ApproveOnce, nil
	default:
		return Deny, nil
	}
}
