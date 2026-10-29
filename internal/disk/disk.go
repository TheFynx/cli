// Package disk reports apparent file sizes, replacing the `fs` du wrapper
// with a native implementation that needs no subprocess.
package disk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Usage prints one line per argument, or every entry of the working
// directory (hidden first, then visible, like the old glob order).
func Usage(w io.Writer, paths []string) error {
	if len(paths) == 0 {
		entries, e := os.ReadDir(".")
		if e != nil {
			return e
		}
		var hidden, visible []string
		for _, en := range entries {
			n := en.Name()
			if strings.HasPrefix(n, ".") && n != "." && n != ".." {
				hidden = append(hidden, n)
			} else {
				visible = append(visible, n)
			}
		}
		sort.Strings(hidden)
		sort.Strings(visible)
		paths = append(paths, hidden...)
		paths = append(paths, visible...)
	}
	for _, p := range paths {
		n, e := size(p)
		if e != nil {
			return fmt.Errorf("%s: %w", p, e)
		}
		fmt.Fprintf(w, "%s\t%s\n", humanize(n), p)
	}
	return nil
}

// size sums apparent bytes beneath path, like `du -b`.
func size(path string) (int64, error) {
	st, e := os.Lstat(path)
	if e != nil {
		return 0, e
	}
	if !st.IsDir() {
		return st.Size(), nil
	}
	var sum int64
	e = filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		sum += info.Size()
		return nil
	})
	return sum, e
}

// humanize formats bytes GNU-style: plain under 1024, one decimal below
// ten units, otherwise an integer ("4.0K", "12K", "1.5G").
func humanize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d", n)
	}
	units := []struct {
		symbol string
		size   int64
	}{{"K", 1024}, {"M", 1024 * 1024}, {"G", 1024 * 1024 * 1024}, {"T", 1024 * 1024 * 1024 * 1024}}
	suffix := ""
	var v float64
	for _, u := range units {
		if n < u.size*1024 {
			v = float64(n) / float64(u.size)
			suffix = u.symbol
			break
		}
	}
	if suffix == "" {
		// >= 1 PiB: cap at TiB display.
		v = float64(n) / float64(units[3].size)
		suffix = "T"
	}
	if v < 10 {
		return fmt.Sprintf("%.1f%s", v, suffix)
	}
	return fmt.Sprintf("%.0f%s", v, suffix)
}
