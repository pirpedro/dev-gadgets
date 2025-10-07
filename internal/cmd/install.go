package cmd

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pirpedro/dev-gadgets/internal/catalog"
	"github.com/pirpedro/dev-gadgets/internal/install"
	"github.com/pirpedro/dev-gadgets/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	flagAll         bool
	flagInteractive bool
	flagOnly        string
)

func init() {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install curated tools and add-ons",
		RunE:  runInstall,
	}
	cmd.Flags().BoolVar(&flagAll, "all", false, "install curated defaults")
	cmd.Flags().BoolVar(&flagInteractive, "interactive", false, "interactive TUI selection")
	cmd.Flags().StringVar(&flagOnly, "only", "", "comma-separated subset of item IDs")
	rootCmd.AddCommand(cmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	cfg, err := catalog.Load()
	if err != nil {
		return err
	}

	var toInstall []catalog.Item
	switch {
	case flagAll:
		toInstall = cfg.Curated()
	case flagOnly != "":
		ids := strings.Split(flagOnly, ",")
		toInstall = cfg.ByIDs(ids)
	case flagInteractive:
		// Integração TUI: seleção interativa
		importUI := func() ([]catalog.Item, error) {
			// Constrói lista de SelectItem
			var items []ui.SelectItem
			for _, it := range cfg.Items {
				items = append(items, ui.SelectItem{
					ID:        it.ID,
					Name:      it.Name,
					Desc:      it.Description,
					Installed: ui.IsInstalled(it),
				})
			}
			model := ui.NewSelectItemsModel(items)
			p := tea.NewProgram(model, tea.WithAltScreen())
			finalModel, err := p.Run()
			if err != nil {
				return nil, err
			}
			selected := finalModel.(ui.SelectItemsModel).SelectedIDs()
			var result []catalog.Item
			for _, id := range selected {
				for _, it := range cfg.Items {
					if it.ID == id {
						result = append(result, it)
					}
				}
			}
			return result, nil
		}
		selected, err := importUI()
		if err != nil {
			return err
		}
		toInstall = selected
	default:
		toInstall = cfg.Curated()
	}

	if flagDryRun {
		for _, it := range toInstall {
			fmt.Fprintf(cmd.OutOrStdout(), "PLAN: %s\n", it.ID)
		}
		return nil
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, it := range toInstall {
		it := it
		g.Go(func() error { return install.Install(ctx, it, install.Options{AssumeYes: flagYes}) })
	}
	return g.Wait()
}
