package main

import (
	"os"
	"os/exec"
	"strings"
)

// run streams a command's stdout/stderr and returns its exit status.
func run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// dockerIDs lists ids from a docker query, treating empty result sets and
// query failures as "nothing to do", like the original wrapper's
// silenced substitutions.
func dockerIDs(args ...string) ([]string, error) {
	c := exec.Command("docker", args...)
	c.Stderr = os.Stderr
	out, e := c.Output()
	if e != nil {
		if len(strings.TrimSpace(string(out))) == 0 {
			return nil, nil // docker daemon down or nothing present
		}
		return nil, nil
	}
	var ids []string
	for _, line := range strings.Fields(string(out)) {
		if line != "" {
			ids = append(ids, line)
		}
	}
	return ids, nil
}
