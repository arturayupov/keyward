package main

import (
	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/audit"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/config"
	kmcp "github.com/arturayupov/keyward/internal/mcp"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve-mcp",
		Short: "Run the keyward MCP server over stdio",
		RunE: func(cmd *cobra.Command, _ []string) error {
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
				base = approval.TerminalApprover{In: cmd.InOrStdin(), Out: cmd.ErrOrStderr()}
			}
			h := &kmcp.Handlers{
				Store:  store,
				Broker: &broker.Broker{Store: store, Approver: approval.NewSessionCache(base), Audit: audit.New(p.Audit)},
			}
			return server.ServeStdio(kmcp.NewServer(h))
		},
	}
}
