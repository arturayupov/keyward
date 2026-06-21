//go:build !darwin && !windows && !linux

package approval

func nativeApprover() (Approver, bool) { return nil, false }
