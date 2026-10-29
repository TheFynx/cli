package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/TheFynx/cli/internal/disk"
)

func newDiskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "Disk and file size tools",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "usage [PATH]...",
		Short: "Apparent size of files or directories (replaces fs)",
		Long: "Reports apparent bytes per argument, human-formatted like `du -sbh`.\n" +
			"Without arguments it sizes every entry of the working directory,\n" +
			"hidden files first.",
		RunE: func(c *cobra.Command, args []string) error {
			return disk.Usage(os.Stdout, args)
		},
	})
	return cmd
}
