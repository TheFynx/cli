package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/TheFynx/cli/internal/theme"
)

func newThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Native asset generation for installed omarchy themes",
	}
	cmd.AddCommand(&cobra.Command{
		Use:       "build minecraft|midgar",
		Short:     "Regenerate native theme assets for an installed omarchy theme",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"minecraft", "midgar"},
		Long: "Builds SVG icons, GTK/Qt assets, and screensaver text directly from the\n" +
			"installed theme colors.toml with no Python, Pillow, or network. Generated\n" +
			"assets in the theme directory are overwritten; reapply with\n" +
			"`omarchy theme set <slug>` afterwards.",
		RunE: func(_ *cobra.Command, args []string) error {
			home, e := os.UserHomeDir()
			if e != nil {
				return e
			}
			return theme.Build(home, args[0])
		},
	})
	return cmd
}
