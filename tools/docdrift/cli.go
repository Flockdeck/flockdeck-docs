package main

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
)

var (
	cmdLine    = regexp.MustCompile(`^  ([a-z][a-z0-9-]*)\b(.*)$`)
	bareWord   = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	codeBlock  = regexp.MustCompile(`(?s)<code[^>]*>(.*?)</code>`)
	htmlTag    = regexp.MustCompile(`<[^>]*>`)
	docCommand = regexp.MustCompile(`(?m)^[ \t]*(?:\$[ \t]*)?flockdeck[ \t]+([a-z][a-z0-9-]*)(?:[ \t]+([a-z][a-z0-9-]*))?`)
)

// key names a command in findings and the allowlist. A subcommand of a
// subcommand is "remote/enable", so an allowlist line stays one token.
func cmdKey(parts ...string) string { return strings.Join(parts, "/") }
func cmdText(key string) string     { return strings.ReplaceAll(key, "/", " ") }

// parseDesktopHelp lists the subcommands `flockdeck -h` prints under
// "Subcommands:". Each is an indented line (two spaces) with its
// description on deeper-indented lines below. "remote enable ..." and
// "remote pair ... | status | devices" contribute "remote/enable",
// "remote/pair", "remote/status", "remote/devices": the second word when it
// is a bare word, and any bare word that starts an alternative after " | ".
func parseDesktopHelp(text string) (map[string]bool, error) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "Subcommands:" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("desktop help has no %q heading; did the help format change?", "Subcommands:")
	}
	cmds := map[string]bool{}
	for _, l := range lines[start:] {
		if strings.TrimSpace(l) == "" {
			break // the section ends at the blank line before "Environment:"
		}
		m := cmdLine.FindStringSubmatch(l)
		if m == nil {
			continue // a description line
		}
		top := m[1]
		cmds[cmdKey(top)] = true
		alts := strings.Split(strings.TrimSpace(m[0]), " | ")
		for i, alt := range alts {
			f := strings.Fields(alt)
			if i == 0 {
				if len(f) > 1 && bareWord.MatchString(f[1]) {
					cmds[cmdKey(top, f[1])] = true
				}
			} else if len(f) > 0 && bareWord.MatchString(f[0]) && f[0] != top {
				cmds[cmdKey(top, f[0])] = true
			}
		}
	}
	if len(cmds) == 0 {
		return nil, fmt.Errorf("desktop help listed no subcommands")
	}
	return cmds, nil
}

// docText is the text of a page's code elements — inline and fenced — which
// is where a command is written as something to type. Prose is left out on
// purpose: "flockdeck is built without a console" is not a subcommand.
func docText(page string) string {
	var b strings.Builder
	for _, m := range codeBlock.FindAllStringSubmatch(page, -1) {
		b.WriteString(html.UnescapeString(htmlTag.ReplaceAllString(m[1], "")))
		b.WriteString("\n")
	}
	return b.String()
}

// documentedCommands returns the commands the page's code elements type as
// `flockdeck <cmd> ...`, and, for a command that has subcommands in the
// help, `flockdeck <cmd> <sub>`.
func documentedCommands(page string, known map[string]bool) map[string]bool {
	out := map[string]bool{}
	for _, m := range docCommand.FindAllStringSubmatch(docText(page), -1) {
		out[cmdKey(m[1])] = true
		if m[2] != "" && hasSub(known, m[1]) {
			out[cmdKey(m[1], m[2])] = true
		}
	}
	return out
}

func hasSub(known map[string]bool, top string) bool {
	for k := range known {
		if strings.HasPrefix(k, top+"/") {
			return true
		}
	}
	return false
}

// mentionsCommand reports whether the page's text says `flockdeck <cmd>`
// somewhere — a prose sentence, a heading or a code block all count as
// documenting a command that exists.
func mentionsCommand(page, key string) bool {
	text := html.UnescapeString(htmlTag.ReplaceAllString(page, ""))
	re := regexp.MustCompile(`\bflockdeck\s+` + strings.ReplaceAll(regexp.QuoteMeta(cmdText(key)), `\ `, `\s+`) + `\b`)
	return re.MatchString(text)
}

// compareCLI returns one finding per disagreement between the desktop's real
// subcommands and the CLI page.
func compareCLI(cmds map[string]bool, page string, allow map[string]string) []string {
	var out []string
	typed := documentedCommands(page, cmds)
	for _, k := range sortedKeys(cmds) {
		_, allowed := allow[k]
		documented := mentionsCommand(page, k) || typed[k]
		if !documented && !allowed {
			out = append(out, fmt.Sprintf("desktop subcommand `%s` is not documented (document it, or allowlist it with a reason)", cmdText(k)))
		}
		if documented && allowed {
			out = append(out, fmt.Sprintf("desktop subcommand `%s` is allowlisted as intentionally undocumented, but the page now documents it; remove it from the allowlist", cmdText(k)))
		}
	}
	for _, k := range sortedKeys(typed) {
		if !cmds[k] {
			out = append(out, fmt.Sprintf("the page shows `flockdeck %s`, which the desktop no longer has (renamed or removed)", cmdText(k)))
		}
	}
	for _, k := range sortedKeys(allow) {
		if !cmds[k] {
			out = append(out, fmt.Sprintf("allowlist entry %s names a subcommand the desktop no longer has; remove it", k))
		}
	}
	sort.Strings(out)
	return out
}
