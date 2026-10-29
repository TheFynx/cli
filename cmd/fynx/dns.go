package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/TheFynx/cli/internal/dns"
)

func newDNSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS lookups through the system resolver",
	}
	lookup := &cobra.Command{
		Use:   "lookup NAME",
		Short: "Resolve records for a name (replaces digga)",
		Long: "Resolves A, AAAA, CNAME, MX, TXT, and NS records natively and grouped\n" +
			"by type, with no dig dependency. Filter with --type. TTLs are not\n" +
			"shown; the standard library resolver does not surface them. Fails\n" +
			"with a non-zero exit when a name has none of the requested records.",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			types, e := c.Flags().GetStringSlice("type")
			if e != nil {
				return e
			}
			return dns.Lookup(os.Stdout, args[0], types)
		},
	}
	lookup.Flags().StringSlice("type", nil, "look up only these types (A, AAAA, CNAME, MX, TXT, NS)")
	cmd.AddCommand(lookup)
	return cmd
}
