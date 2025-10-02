package cmd

import (
	"github.com/pirpedro/dev-gadgets/internal/catalog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List catalog items",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := catalog.Load()
			if err != nil {
				return err
			}
			for _, it := range cfg.Items {
				desc := it.Description
				if desc == "" {
					desc = "(sem descrição)"
				}
				cmd.Printf("%-20s %s\n", it.ID, desc)
			}
			return nil
		},
	})
}
