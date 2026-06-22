package main

import (
	"fmt"
	"os"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/spf13/cobra"
)

func newInjectCmd() *cobra.Command {
	var ns, into string
	c := &cobra.Command{
		Use:   "inject NAME --ns NS --into PATH",
		Short: "Inject a key into a target env file (prompts for approval)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := config.Default()
			if err != nil {
				return err
			}
			id, err := vault.EnsureIdentity()
			if err != nil {
				return err
			}
			store, err := vault.Load(p.Vault, id)
			if err != nil {
				return err
			}
			var base approval.Approver
			if na, ok := approval.NativeApprover(); ok {
				base = na
			} else {
				base = approval.TerminalApprover{In: os.Stdin, Out: cmd.OutOrStdout()}
			}
			b := &broker.Broker{
				Store:    store,
				Approver: approval.NewSessionCache(base),
				Audit:    audit.New(p.Audit),
			}
			res, err := b.Request(approval.Request{Tool: "cli", Name: args[0], Namespace: ns, Target: into})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s → %s\n", res.Status, res.Name, res.Target)
			return nil
		},
	}
	c.Flags().StringVar(&ns, "ns", "", "namespace")
	c.Flags().StringVar(&into, "into", ".env", "target env file")
	return c
}
