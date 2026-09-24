package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// relayFlag is one flag from `flockdeck-relay serve -h`.
type relayFlag struct {
	Name    string
	Default string // as printed; "" when the relay prints none (a zero value)
	Usage   string
	Envs    []string // environment variables the usage text says do the same thing
}

var (
	flagLine    = regexp.MustCompile(`^  -([A-Za-z0-9][A-Za-z0-9_-]*)(?: .*)?$`)
	defaultTail = regexp.MustCompile(`\(default ("[^"]*"|[^)"]*)\)\s*$`)
	orGroup     = regexp.MustCompile(`\(or\b[^)]*\)`)
	envName     = regexp.MustCompile(`FLOCKDECK_RELAY_[A-Z0-9_]+`)
)

// parseRelayHelp reads the flags out of `flockdeck-relay serve -h`, in the
// format Go's flag package prints: a two-space-indented "-name type" line,
// then the usage on tab-indented lines. Only the section after "Flags for
// serve:" is read; without that heading it is an error rather than an empty
// result, so a changed help format fails loudly instead of passing.
func parseRelayHelp(text string) (map[string]relayFlag, error) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "Flags for serve:" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("relay help has no %q heading; did the help format change?", "Flags for serve:")
	}
	flags := map[string]relayFlag{}
	var cur *relayFlag
	flush := func() {
		if cur == nil {
			return
		}
		cur.Usage = strings.TrimSpace(cur.Usage)
		if m := defaultTail.FindStringSubmatch(cur.Usage); m != nil {
			cur.Default = strings.Trim(m[1], `"`)
		}
		for _, g := range orGroup.FindAllString(cur.Usage, -1) {
			cur.Envs = append(cur.Envs, envName.FindAllString(g, -1)...)
		}
		flags[cur.Name] = *cur
		cur = nil
	}
	for _, l := range lines[start:] {
		if m := flagLine.FindStringSubmatch(l); m != nil {
			flush()
			cur = &relayFlag{Name: "-" + m[1]}
		} else if cur != nil && strings.HasPrefix(strings.TrimLeft(l, " "), "\t") { // flag prints "    \t"
			cur.Usage += " " + strings.TrimSpace(l)
		} else if strings.TrimSpace(l) == "" {
			flush()
		}
	}
	flush()
	if len(flags) == 0 {
		return nil, fmt.Errorf("relay help listed no flags")
	}
	return flags, nil
}

// docRow is one flag's row in a configuration table.
type docRow struct {
	Default    string // the Default cell, "" when the table has none or it is blank
	HasDefault bool   // the table has a Default column at all
	Envs       []string
	HasEnv     bool
}

// docFlags is what a configuration page says about flags.
type docFlags struct {
	Rows     map[string]docRow // flags in a table's first column
	Mentions map[string]bool   // every flag the page names, table or prose or example
}

var (
	codeSpan   = regexp.MustCompile("`([^`]+)`")
	flagToken  = regexp.MustCompile(`^-[a-z][a-z0-9-]*$`)
	flagInCode = regexp.MustCompile(`(?:^|\s)(-[a-z][a-z0-9-]*)`)
)

// proseFlags returns the flags text names in a code span of their own.
// "There's no `-oidc-client-secret` flag" says the flag is absent, which is
// not a claim that it exists, so a span right after "no " is skipped.
func proseFlags(text string) []string {
	var out []string
	for _, m := range codeSpan.FindAllStringSubmatchIndex(text, -1) {
		span := text[m[2]:m[3]]
		if flagToken.MatchString(span) && !strings.HasSuffix(strings.ToLower(text[:m[0]]), "no ") {
			out = append(out, span)
		}
	}
	return out
}

func splitRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	cells := strings.Split(line, "|")
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}
	return cells
}

// parseConfigDoc reads flag tables and mentions out of the configuration page.
//
// A table is any run of "|" lines whose header's first cell is "Flag". A
// mention outside a table counts only when a code span is exactly one flag
// (`-audit`), or inside a fenced example that runs `flockdeck-relay serve`;
// other spans (`flockdeck-relay invite`, a flag of another command) are not
// this page's claims about serve.
func parseConfigDoc(md string) docFlags {
	d := docFlags{Rows: map[string]docRow{}, Mentions: map[string]bool{}}
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	var header []string
	inFence, serveFence := false, false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "```") {
			inFence = !inFence
			serveFence = false
			header = nil
			continue
		}
		if inFence {
			if strings.Contains(t, "flockdeck-relay serve") {
				serveFence = true
			}
			if serveFence {
				for _, m := range flagInCode.FindAllStringSubmatch(t, -1) {
					d.Mentions[m[1]] = true
				}
			}
			continue
		}
		if !strings.HasPrefix(t, "|") {
			header = nil
			for _, f := range proseFlags(t) {
				d.Mentions[f] = true
			}
			continue
		}
		cells := splitRow(t)
		if header == nil {
			header = cells
			continue
		}
		if strings.Trim(cells[0], "-: ") == "" { // the |---|---| rule
			continue
		}
		if !strings.EqualFold(header[0], "Flag") {
			continue
		}
		col := func(name string) (string, bool) {
			for i, h := range header {
				if strings.EqualFold(h, name) && i < len(cells) {
					return cells[i], true
				}
			}
			return "", false
		}
		def, hasDef := col("Default")
		env, hasEnv := col("Env")
		var envs []string
		for _, m := range codeSpan.FindAllStringSubmatch(env, -1) {
			envs = append(envs, m[1])
		}
		for _, m := range codeSpan.FindAllStringSubmatch(cells[0], -1) {
			if !flagToken.MatchString(m[1]) {
				continue
			}
			d.Mentions[m[1]] = true
			d.Rows[m[1]] = docRow{Default: def, HasDefault: hasDef, Envs: envs, HasEnv: hasEnv}
		}
		// Prose in the other cells can name flags too.
		for _, c := range cells[1:] {
			for _, f := range proseFlags(c) {
				d.Mentions[f] = true
			}
		}
	}
	return d
}

// canonDefault reduces a default, however written, to something comparable:
// "720h (30 days)" and "720h0m0s" are the same duration; "off", "0", "0s"
// and nothing are all a zero value. ok is false when there is nothing to
// compare (the docs state no default).
func canonDefault(s string) (string, bool) {
	s = strings.TrimSpace(strings.ReplaceAll(s, "`", ""))
	if i := strings.Index(s, " ("); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	s = strings.Trim(s, `"`)
	switch strings.ToLower(s) {
	case "":
		return "", false
	case "off", "none", "false", "0", "0s":
		return "0", true
	}
	if d, err := time.ParseDuration(s); err == nil {
		return "d:" + d.String(), true
	}
	return s, true
}

// compareRelay returns one finding per disagreement between the relay's real
// flags and the configuration page. Findings are stable, sorted strings so a
// baseline run can be diffed against a head run.
func compareRelay(flags map[string]relayFlag, doc docFlags, allow map[string]string) []string {
	var out []string
	for _, name := range sortedKeys(flags) {
		f := flags[name]
		_, allowed := allow[name]
		if !doc.Mentions[name] && !allowed {
			out = append(out, fmt.Sprintf("relay flag %s is undocumented (document it, or allowlist it with a reason)", name))
			continue
		}
		if doc.Mentions[name] && allowed {
			out = append(out, fmt.Sprintf("relay flag %s is allowlisted as intentionally undocumented, but the docs now mention it; remove it from the allowlist", name))
		}
		row, ok := doc.Rows[name]
		if !ok {
			continue
		}
		if row.HasDefault {
			if want, stated := canonDefault(row.Default); stated {
				got, _ := canonDefault(f.Default)
				if got == "" {
					got = "0"
				}
				if want != got {
					out = append(out, fmt.Sprintf("relay flag %s: docs give default %q, the relay's is %s", name, strings.TrimSpace(row.Default), displayDefault(f.Default)))
				}
			}
		}
		if row.HasEnv {
			for _, e := range row.Envs {
				if !strings.Contains(f.Usage, e) {
					out = append(out, fmt.Sprintf("relay flag %s: docs name env var %s, which the relay's help does not", name, e))
				}
			}
			for _, e := range f.Envs {
				if !contains(row.Envs, e) {
					out = append(out, fmt.Sprintf("relay flag %s: relay reads env var %s, which the docs' Env column omits", name, e))
				}
			}
		}
	}
	for _, name := range sortedKeys(doc.Mentions) {
		if _, ok := flags[name]; !ok {
			out = append(out, fmt.Sprintf("docs mention %s, which the relay no longer has (renamed or removed)", name))
		}
	}
	for _, name := range sortedKeys(allow) {
		if _, ok := flags[name]; !ok {
			out = append(out, fmt.Sprintf("allowlist entry %s names a flag the relay no longer has; remove it", name))
		}
	}
	sort.Strings(out)
	return out
}

func displayDefault(s string) string {
	if s == "" {
		return "none"
	}
	return fmt.Sprintf("%q", s)
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
