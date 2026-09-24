package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// loadAllow reads an allowlist: one entry per line, the name then a reason,
// separated by whitespace. Blank lines and #-comments are skipped. An entry
// with no reason is an error — the file exists to record *why* something is
// deliberately undocumented, so a bare name defeats it.
func loadAllow(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	allow := map[string]string{}
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, reason, _ := strings.Cut(line, " ")
		reason = strings.TrimSpace(reason)
		if reason == "" {
			return nil, fmt.Errorf("%s:%d: %q has no reason; every allowlist entry needs a one-line why", path, n, name)
		}
		if _, dup := allow[name]; dup {
			return nil, fmt.Errorf("%s:%d: %q is listed twice", path, n, name)
		}
		allow[name] = reason
	}
	return allow, sc.Err()
}
