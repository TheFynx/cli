package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newManCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "man [SECTION] PAGE...",
		Short: "Colored manual pages (drop-in for the colored man function)",
		Long: "Executes man with the LESS_TERMCAP color variables the old ~/.functions\n" +
			"function set. Replace that function with `alias man='fynx man'`.",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("usage: fynx man [SECTION] PAGE...")
			}
			c := exec.Command("man", args...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Env = append(os.Environ(),
				"LESS_TERMCAP_mb=\x1b[1;31m",
				"LESS_TERMCAP_md=\x1b[1;31m",
				"LESS_TERMCAP_me=\x1b[0m",
				"LESS_TERMCAP_se=\x1b[0m",
				"LESS_TERMCAP_so=\x1b[1;44;33m",
				"LESS_TERMCAP_ue=\x1b[0m",
				"LESS_TERMCAP_us=\x1b[1;32m",
			)
			return c.Run()
		},
	}
}
