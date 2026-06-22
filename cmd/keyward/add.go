package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/arturayupov/keyward/internal/config"
	"github.com/arturayupov/keyward/internal/model"
	"github.com/arturayupov/keyward/internal/vault"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var ns string
	c := &cobra.Command{
		Use:   "add NAME --ns NS",
		Short: "Add or update a single secret (value read from stdin)",
		Long: "Add or update one secret in the vault. The value is read from stdin so " +
			"it never appears in argv or shell history, e.g.:\n\n" +
			"  printf %s \"$TOKEN\" | keyward add NPM_TOKEN --ns cli-tools",
		Args: cobra.ExactArgs(1),
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
			raw, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return err
			}
			value := strings.TrimRight(string(raw), "\r\n")
			if value == "" {
				return fmt.Errorf("no value on stdin; pipe it, e.g. printf %%s \"$TOKEN\" | keyward add %s --ns %s", args[0], ns)
			}
			store.Upsert(model.Secret{Name: args[0], Namespace: ns, Value: value})
			if err := vault.Save(p.Vault, store, id.Recipient()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "added %s/%s\n", ns, args[0])
			return nil
		},
	}
	c.Flags().StringVar(&ns, "ns", "default", "namespace")
	return c
}
