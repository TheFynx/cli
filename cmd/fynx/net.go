package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"fynx.local/cli/internal/netcheck"
)

func newNetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "net",
		Short: "Network probes",
	}
	isup := &cobra.Command{
		Use:   "isup URI",
		Short: "Check reachability and send a desktop notification (replaces isup)",
		Long: "HEADs the URI through redirects with a timeout and reports whether it\n" +
			"serves 2xx/3xx. Prints the verdict, sends a desktop notification when\n" +
			"notify-send is available (critical urgency when down), and exits 1\n" +
			"when the site is down for easy scripting. The old shell version\n" +
			"reported the opposite verdict; that bug is gone here.",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			timeout, e := c.Flags().GetDuration("timeout")
			if e != nil {
				return e
			}
			up, detail, err := netcheck.Check(args[0], timeout)
			if err != nil {
				detail = err.Error()
			}
			state, urgency := "is down", "critical"
			if up {
				state, urgency = "is up", "low"
			}
			fmt.Println(args[0], state, "("+detail+")")
			if _, e := exec.LookPath("notify-send"); e == nil {
				_ = run("notify-send", "--urgency="+urgency, args[0]+" "+state)
			}
			if !up {
				os.Exit(1)
			}
			return nil
		},
	}
	isup.Flags().Duration("timeout", netcheck.DefaultTimeout, "probe timeout")
	cmd.AddCommand(isup)
	return cmd
}
