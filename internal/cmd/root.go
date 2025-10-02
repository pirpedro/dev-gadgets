package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagYes    bool
	flagDryRun bool
	version    = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "dev-gadgets",
	Short: "Install and manage dev adjacent tools",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVar(&flagYes, "yes", false, "assume yes to confirmations")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "print plan only, do not execute")
}
