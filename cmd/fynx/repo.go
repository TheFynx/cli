package main

import (
	"os"

	"github.com/spf13/cobra"

	"fynx.local/cli/internal/repo"
)

func newRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Git repository inspection and source export",
	}
	contents := &cobra.Command{
		Use:   "contents DIRECTORY",
		Short: "Concatenate a source tree into one text file (replaces copy_repo_contents)",
		Long: "Writes all .go and .adoc files, .goreleaser.yaml, and the cmd/ and\n" +
			"internal/ trees (with per-file path headers) into a single file for\n" +
			"LLM consumption. The output is rewritten fresh on each run; pass\n" +
			"--out to choose the destination. The old shell version's tar step\n" +
			"silently re-extracted files under /; it is gone here.",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			out, e := c.Flags().GetString("out")
			if e != nil {
				return e
			}
			return repo.Contents(args[0], out, os.Stdout)
		},
	}
	contents.Flags().String("out", "", "output file (default: repo_contents.txt)")
	cmd.AddCommand(
		&cobra.Command{
			Use:   "check BASE_FOLDER",
			Short: "Scan a tree for repos with unpushed or uncommitted work (replaces git-check)",
			Long: "Walks BASE_FOLDER for git work trees and prints each repo's unmerged\n" +
				"branches and local changes, followed by a summary. A discovered repo is\n" +
				"reported but not descended into. Symlinked directories are not followed\n" +
				"(the shell version could loop on them).",
			Args: cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return repo.Check(args[0], os.Stdout)
			},
		},
		contents,
	)
	return cmd
}
