package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// aliasPairs map the retired ~/.functions names to their fynx commands.
var aliasPairs = []struct{ fn, cmd string }{
	{"dclean", "fynx docker clean"},
	{"fs", "fynx disk usage"},
	{"digga", "fynx dns lookup"},
	{"man", "fynx man"},
	{"isup", "fynx net isup"},
	{"git-check", "fynx repo check"},
	{"copy_repo_contents", "fynx repo contents"},
}

// newAliasesCmd prints the shell alias block that retires ~/.functions.
func newAliasesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "aliases [bash|zsh|fish]",
		Short: "Print the alias block that replaces ~/.functions",
		Long: "Emits the short shell aliases that point the old function names at\n" +
			"their fynx replacements. Paste the printed block into your shell rc\n" +
			"(bashrc/zshrc/config.fish) and retire ~/.functions.",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		RunE: func(_ *cobra.Command, args []string) error {
			shell := "bash"
			if len(args) == 1 {
				shell = args[0]
			}
			switch shell {
			case "fish":
				for _, p := range aliasPairs {
					fmt.Printf("alias %s '%s'\n", p.fn, p.cmd)
				}
			case "bash", "zsh":
				for _, p := range aliasPairs {
					fmt.Printf("alias %s='%s'\n", p.fn, p.cmd)
				}
			default:
				return fmt.Errorf("unknown shell %q (use bash, zsh, or fish)", shell)
			}
			fmt.Println("# paste into your shell rc, then retire ~/.functions")
			return nil
		},
	}
}
