package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func newDockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Container housekeeping",
	}
	clean := &cobra.Command{
		Use:   "clean [all|fire]",
		Short: "Stop and remove containers and images (replaces dclean)",
		Long: "With no argument: stop and remove every container, then remove dangling\n" +
			"images. `all`: same, plus all images. `fire`: same as `all`, plus a full\n" +
			"`docker system prune -a --volumes`. With --dry-run the commands are only\n" +
			"printed. Docker must be installed.",
		RunE: dockerClean,
	}
	clean.Flags().Bool("dry-run", false, "print the docker commands without running them")
	cmd.AddCommand(clean)
	return cmd
}

func dockerClean(c *cobra.Command, args []string) error {
	scope := ""
	if len(args) > 1 {
		return errors.New("usage: fynx docker clean [all|fire]")
	}
	if len(args) == 1 {
		scope = args[0]
	}
	switch scope {
	case "", "all", "fire":
	default:
		return fmt.Errorf("invalid argument: %s (use no arg, all, or fire)", scope)
	}
	if _, e := exec.LookPath("docker"); e != nil {
		return errors.New("docker is not installed or not in PATH")
	}
	dry, _ := c.Flags().GetBool("dry-run")
	fmt.Println("Stopping and removing all containers and images...")
	if e := dockerStep(dry, "stop"); e != nil {
		return e
	}
	if e := dockerStep(dry, "rm"); e != nil {
		return e
	}
	if scope == "" {
		return dockerStep(dry, "rmi", "--filter", "dangling=true")
	}
	if e := dockerStep(dry, "rmi"); e != nil {
		return e
	}
	if scope == "fire" {
		return dockerStep(dry, "system", "prune", "-a", "--force", "--volumes")
	}
	return nil
}

// dockerStep runs `docker <args...>` prefixed with the matching id list
// from the appropriate query, or prints the command under --dry-run.
// An empty id list skips the step, like the original wrapper.
func dockerStep(dry bool, args ...string) error {
	var ids []string
	switch {
	case args[0] == "stop" || args[0] == "rm":
		ids, _ = dockerIDs("ps", "-a", "-q")
	case args[0] == "rmi" && len(args) > 2:
		ids, _ = dockerIDs("images", "-f", "dangling=true", "-q")
	case args[0] == "rmi":
		ids, _ = dockerIDs("images", "-q")
	}
	if len(ids) == 0 {
		return nil
	}
	full := append(append([]string{args[0]}, ids...), args[1:]...)
	if dry {
		fmt.Println("+ docker", strings.Join(full, " "))
		return nil
	}
	return run("docker", full...)
}
