---
type: Module
title: Desktop subcommand parser and comparer
description: Parses the subcommands out of flockdeck -h and checks that the CLI page documents each one.
resource: tools/docdrift/cli.go
tags: [docs, drift, parser]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/cli.go]
---
## Help side: `parseDesktopHelp`

Starts after the line `Subcommands:` (error if absent) and ends at the first blank line. A subcommand line is two spaces then a lowercase word; deeper-indented lines are descriptions and are skipped. Keys are the top word (`remote`), plus for a nested command `top/second` when the second word is a bare word, plus each bare word that starts an alternative after ` | ` (so `remote pair ... | status | devices` gives `remote/pair`, `remote/status`, `remote/devices`). A key with a space is written with a slash so an allowlist line stays one token. No subcommands found is an error.

## Docs side

The page is HTML. `docText` keeps only the text inside `<code>` elements (tags stripped, entities unescaped), so prose such as "flockdeck is built without a console" is not read as a command. `documentedCommands` matches `flockdeck <cmd> [<sub>]`, with an optional leading `$`, line by line; the second word counts only if the command has subcommands in the help.

`mentionsCommand` is looser: the string `flockdeck <cmd>` anywhere in the page text (prose, heading or code) counts as documenting a command that exists.

## Findings: `compareCLI`

| Finding | Condition |
|---|---|
| not documented | In help, not mentioned, not allowlisted |
| allowlisted but documented | Mentioned or typed, and in the allowlist |
| removed or renamed | Typed in a code element, not in help |
| stale allowlist entry | Allowlisted, not in help |

Allowlist: [allowlist format](/data/allowlist-files.md). Part of [docdrift](/packages/docdrift.md). Input captured by [fetch-desktop-help.sh](/cli/fetch-desktop-help-sh.md).
