// Package netcheck probes whether a URL is reachable, replacing the old
// curl-and-grep `isup` workaround (which reported the inverted result)
// with a real HTTP request and a correct verdict.
package netcheck

import (
	"fmt"
	"net/http"
	"time"
)

// DefaultTimeout bounds each probe.
const DefaultTimeout = 5 * time.Second

// Check HEADs uri through redirects. It returns the verdict, the HTTP
// status line, and the error distinguishing network failure from a
// reachable-but-unhappy server.
func Check(uri string, timeout time.Duration) (up bool, detail string, err error) {
	c := &http.Client{Timeout: timeout}
	resp, err := c.Head(uri)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	return code >= 200 && code < 400, fmt.Sprintf("HTTP %d %s", code, http.StatusText(code)), nil
}
