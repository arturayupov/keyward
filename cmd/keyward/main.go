package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "keyward",
		Short: "A local secret broker for AI agents",
	}
	root.AddCommand(newInitCmd(), newImportCmd(), newLsCmd(), newInjectCmd(), newServeCmd())
	return root
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
