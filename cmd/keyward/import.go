package main

import (
	"fmt"

	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/importer"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import [root]",
		Short: "Import secrets from .env files under root into the vault",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
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
			found, err := importer.Import(root)
			if err != nil {
				return err
			}
			for _, s := range found {
				store.Upsert(s)
			}
			if err := vault.Save(p.Vault, store, id.Recipient()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Imported %d secrets (names only):\n", len(found))
			for _, m := range store.Meta("") {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s/%s\n", m.Namespace, m.Name)
			}
			return nil
		},
	}
}
