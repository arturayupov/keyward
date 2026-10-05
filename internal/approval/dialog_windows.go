//go:build windows

package approval

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

func nativeApprover() (Approver, bool) { return msgBoxApprover{}, true }

type msgBoxApprover struct{}

// Approve shows the native MessageBox via PowerShell.
//
// The message is passed base64-encoded: embedding it directly breaks the
// PowerShell command, because Request.Prompt() always contains double quotes
// (around the key name) which terminate the interpolated PS string early.
// That made every approval request fail to even render the dialog and fail
// closed (instant Deny) for all Windows users.
func (msgBoxApprover) Approve(r Request) (Decision, error) {
	msg := fmt.Sprintf("%s\n\nReason: %s\n\nYes = approve for session, No = approve once, Cancel = deny", r.Prompt(), r.Reason)
	enc := base64.StdEncoding.EncodeToString([]byte(msg))
	ps := fmt.Sprintf(
		`Add-Type -AssemblyName PresentationFramework; $m=[Text.Encoding]::UTF8.GetString([Convert]::FromBase64String('%s')); [System.Windows.MessageBox]::Show($m,'keyward','YesNoCancel','Question')`,
		enc)
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
