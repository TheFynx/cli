// Package repo inspects git repositories and exports source trees,
// replacing the git-check and copy_repo_contents shell functions. The
// shell version of repo contents had a broken tar pipe that re-extracted
// files under /; that section now inlines real file contents.
package repo

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	cyan   = "\033[0;36m"
	yellow = "\033[1;33m"
	reset  = "\033[0m"
)

// Check walks base for git work trees (immediate discovery only; a repo
// is found, then not descended into) and reports unmerged branches and
// local changes for each, followed by a summary.
func Check(base string, w io.Writer) error {
	st, e := os.Stat(base)
	if e != nil {
		return fmt.Errorf("'%s' is not a valid folder: %w", base, e)
	}
	if !st.IsDir() {
		return fmt.Errorf("'%s' is not a valid folder", base)
	}
	abs, e := filepath.Abs(base)
	if e != nil {
		return e
	}
	var repos, dirty, unmerged int
	e = filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || p == abs {
			return nil
		}
		gitEntry, e := os.Stat(filepath.Join(p, ".git"))
		if e != nil || !gitEntry.IsDir() {
			return nil // descend further; base itself is never reported
		}
		repos++
		changed, unpushed, e := checkRepo(p, w)
		if e != nil {
			fmt.Fprintf(w, "  %s%s%s\n", yellow, e, reset)
		}
		if unpushed {
			unmerged++
		}
		if changed {
			dirty++
		}
		fmt.Fprintln(w)
		return filepath.SkipDir
	})
	if e != nil {
		return e
	}
	if repos == 0 {
		return fmt.Errorf("no git repositories found under %s", abs)
	}
	fmt.Fprintf(w, "%s%d repositories: %d with unmerged branches, %d with local changes%s\n",
		cyan, repos, unmerged, dirty, reset)
	return nil
}

func checkRepo(dir string, w io.Writer) (dirty, unpushed bool, err error) {
	fmt.Fprintf(w, "%sRepository: %s%s\n", cyan, dir, reset)
	fmt.Fprintf(w, "%sBranches not pushed:%s\n", yellow, reset)
	if b, e := exec.Command("git", "-C", dir, "branch", "--no-merged").Output(); e == nil {
		w.Write(b)
		if len(strings.TrimSpace(string(b))) > 0 {
			unpushed = true
		}
	} else {
		fmt.Fprintf(w, "  git branch failed: %v\n", e)
	}
	fmt.Fprintf(w, "%sLocal changes:%s\n", yellow, reset)
	b, e := exec.Command("git", "-C", dir, "status", "-uno").Output()
	if e != nil {
		return dirty, unpushed, nil
	}
	w.Write(b)
	dirty = !strings.Contains(string(b), "nothing to commit") &&
		!strings.Contains(string(b), "nothing added to commit")
	return dirty, unpushed, nil
}

// Contents concatenates a source tree's code into a single text file for
// LLM consumption: all .go and .adoc files, .goreleaser.yaml, and the
// cmd/ and internal/ directories with per-file path headers. The output
// file is rewritten fresh on every run.
func Contents(dir, out string, w io.Writer) error {
	st, e := os.Stat(dir)
	if e != nil {
		return fmt.Errorf("'%s' is not a valid folder: %w", dir, e)
	}
	if !st.IsDir() {
		return fmt.Errorf("'%s' is not a valid folder", dir)
	}
	if out == "" {
		out = "repo_contents.txt"
	}
	abs, e := filepath.Abs(dir)
	if e != nil {
		return e
	}
	if e := os.MkdirAll(filepath.Dir(absOut(out)), 0o755); e != nil {
		return e
	}
	f, e := os.Create(absOut(out))
	if e != nil {
		return e
	}
	defer f.Close()

	// Apparent original order: all .go, then all .adoc, then goreleaser
	// config, then command and internal trees.
	for _, ext := range []string{".go", ".adoc"} {
		matches := expandGlobDir(abs, ext)
		for _, m := range matches {
			b, e := os.ReadFile(m)
			if e != nil {
				return e
			}
			f.Write(b)
		}
		if len(matches) > 0 {
			f.WriteString("\n---\n\n")
		}
	}
	if b, e := os.ReadFile(filepath.Join(abs, ".goreleaser.yaml")); e == nil {
		f.Write(b)
		f.WriteString("\n---\n\n")
	}
	for _, sub := range []string{"cmd", "internal"} {
		root := filepath.Join(abs, sub)
		if _, e := os.Stat(root); e != nil {
			continue
		}
		fmt.Fprintf(f, "# %s\n", sub)
		e := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !d.Type().IsRegular() {
				return nil
			}
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			rel, _ := filepath.Rel(abs, p)
			fmt.Fprintf(f, "# %s\n", rel)
			f.Write(b)
			f.WriteString("\n")
			return nil
		})
		if e != nil {
			return e
		}
		f.WriteString("\n---\n\n")
	}
	fmt.Fprintf(w, "Directory contents saved to: %s\n", absOut(out))
	return nil
}

func absOut(out string) string {
	if filepath.IsAbs(out) {
		return out
	}
	a, e := filepath.Abs(out)
	if e != nil {
		return out
	}
	return a
}

// expandGlobDir lists regular files ending in ext anywhere beneath dir,
// in stable lexical order (the Go-native stand-in for find -name).
func expandGlobDir(dir, ext string) []string {
	var out []string
	filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && strings.HasSuffix(p, ext) {
			out = append(out, p)
		}
		return nil
	})
	return out
}
