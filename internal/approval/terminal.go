package approval

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TerminalApprover prompts on a reader/writer pair (fallback when no GUI).
type TerminalApprover struct {
	In  io.Reader
	Out io.Writer
}

// Approve prints the request and reads a one-line decision. Fails closed.
func (t TerminalApprover) Approve(r Request) (Decision, error) {
	fmt.Fprintf(t.Out, "\n[keyward] %s\nReason: %s\nApprove? [o]nce / [s]ession / [N]o: ", r.Prompt(), r.Reason)
	line, _ := bufio.NewReader(t.In).ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "o", "once":
		return ApproveOnce, nil
	case "s", "session":
		return ApproveSession, nil
	default:
		return Deny, nil
	}
}
