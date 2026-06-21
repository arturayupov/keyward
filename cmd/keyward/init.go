package main

import (
	"fmt"

	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/model"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create the encrypted vault and master key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := config.Default()
			if err != nil {
				return err
			}
			if err := p.EnsureDir(); err != nil {
				return err
			}
			id, err := vault.EnsureIdentity()
			if err != nil {
				return err
			}
			if err := vault.Save(p.Vault, &model.Store{}, id.Recipient()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Initialized vault at %s\n", p.Vault)
			return nil
		},
	}
}
