// Command fynx is a personal toolkit CLI: asset upscaling, native theme
// generation, and the housekeeping utilities retired from ~/.functions.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/TheFynx/cli/internal/config"
)

const version = "0.2.0"

func main() {
	if e := newRoot().Execute(); e != nil {
		os.Exit(1)
	}
}

func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "fynx",
		Short: "Personal toolkit: upscaling, theme assets, and housekeeping",
		Long: "fynx is the single maintained binary behind the retired ~/~.functions\n" +
			"dotfile. Each former shell function now has a real subcommand; keep\n" +
			"short aliases with `fynx aliases <shell>`.",
		SilenceUsage:       true,
		DisableSuggestions: true,
		Version:            version,
		CompletionOptions:  cobra.CompletionOptions{HiddenDefaultCmd: true},
	}
	root.SetVersionTemplate("fynx {{.Version}}\n")
	root.AddCommand(
		newAssetsCmd(),
		newThemeCmd(),
		newDockerCmd(),
		newDiskCmd(),
		newDNSCmd(),
		newNetCmd(),
		newManCmd(),
		newRepoCmd(),
		newAliasesCmd(),
	)
	return root
}

// loadConfig populates the shared configuration lazily; commands that do
// not read settings never touch the file.
func loadConfig() error {
	c, e := config.Load()
	if e != nil {
		return fmt.Errorf("config: %w", e)
	}
	config.C = c
	return nil
}
