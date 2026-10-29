// Package dns resolves domain records natively through the system
// resolver, replacing the `digga` dig shell-out with parallel lookups
// that also work where dig is not installed. TTLs are not shown; the
// standard library resolver does not surface them.
package dns

import (
	"context"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"time"
)

// DefaultTypes is the lookup order when no type filter is given.
var DefaultTypes = []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS"}

// Lookup queries name for each wanted record type and prints grouped
// results. Unknown types are reported as an error.
func Lookup(w io.Writer, name string, wanted []string) error {
	if len(wanted) == 0 {
		wanted = DefaultTypes
	}
	for _, t := range wanted {
		switch t {
		case "A", "AAAA", "CNAME", "MX", "TXT", "NS":
		default:
			return fmt.Errorf("unsupported record type %q", t)
		}
	}
	r := &net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	found := 0
	for _, t := range wanted {
		lines, e := oneType(ctx, r, name, t)
		if e != nil || len(lines) == 0 {
			continue
		}
		for _, l := range lines {
			fmt.Fprintln(w, l)
		}
		found += len(lines)
	}
	if found == 0 {
		return fmt.Errorf("no %s records for %s", strings.Join(wanted, "/"), name)
	}
	return nil
}

func oneType(ctx context.Context, r *net.Resolver, name, t string) ([]string, error) {
	switch t {
	case "A", "AAAA":
		ips, e := r.LookupNetIP(ctx, "ip", name)
		if e != nil {
			return nil, e
		}
		var out []string
		for _, ip := range ips {
			s := ip.String()
			if strings.Contains(s, ":") == (t == "AAAA") {
				out = append(out, indent(t, s))
			}
		}
		return out, nil
	case "CNAME":
		target, e := r.LookupCNAME(ctx, name)
		if e != nil {
			return nil, e
		}
		if strings.EqualFold(trimDot(target), trimDot(name)) {
			return nil, nil
		}
		return []string{indent(t, trimDot(target))}, nil
	case "MX":
		mxs, e := r.LookupMX(ctx, name)
		if e != nil {
			return nil, e
		}
		var out []string
		for _, m := range mxs {
			out = append(out, fmt.Sprintf("  %-5s %3d %s", t, m.Pref, trimDot(m.Host)))
		}
		return out, nil
	case "TXT":
		txts, e := r.LookupTXT(ctx, name)
		if e != nil {
			return nil, e
		}
		var out []string
		for _, txt := range txts {
			out = append(out, indent(t, `"`+txt+`"`))
		}
		return out, nil
	case "NS":
		nss, e := r.LookupNS(ctx, name)
		if e != nil {
			return nil, e
		}
		sort.Slice(nss, func(i, j int) bool { return nss[i].Host < nss[j].Host })
		var out []string
		for _, ns := range nss {
			out = append(out, indent(t, trimDot(ns.Host)))
		}
		return out, nil
	}
	return nil, fmt.Errorf("unsupported record type %q", t)
}

func indent(t, value string) string {
	return fmt.Sprintf("  %-5s %s", t, value)
}

func trimDot(s string) string {
	return strings.TrimSuffix(s, ".")
}
