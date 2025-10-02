package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/adrg/xdg"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(&cobra.Command{
        Use:   "doctor",
        Short: "Check environment and dependencies",
        RunE: func(cmd *cobra.Command, args []string) error {
            binDir := filepath.Join(xdg.Home, ".local", "bin")
            found := false
            for _, p := range filepath.SplitList(os.Getenv("PATH")) {
                if p == binDir { found = true; break }
            }
            if !found { fmt.Fprintf(cmd.OutOrStdout(), "Hint: add %s to your PATH\n", binDir) }
            fmt.Fprintln(cmd.OutOrStdout(), "OK")
            return nil
        },
    })
}
