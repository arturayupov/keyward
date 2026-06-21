package main

import (
	"fmt"

	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/spf13/cobra"
)

func newLsCmd() *cobra.Command {
	var ns string
	c := &cobra.Command{
		Use:   "ls",
		Short: "List stored key names and namespaces (never values)",
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
			for _, m := range store.Meta(ns) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", m.Namespace, m.Name)
			}
			return nil
		},
	}
	c.Flags().StringVar(&ns, "ns", "", "filter by namespace")
	return c
}
