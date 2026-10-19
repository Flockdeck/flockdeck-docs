---
type: File Format
title: Drift allowlists
description: Two text files listing relay flags and desktop subcommands the docs leave out on purpose, each with a reason.
resource: tools/docdrift/allow-relay-flags.txt
tags: [docs, drift]
timestamp: 2026-10-19T00:00:00Z
code_paths: [tools/docdrift/allow-relay-flags.txt, tools/docdrift/allow-desktop-commands.txt, tools/docdrift/allow.go]
---
| File | Names | Checked against |
|---|---|---|
| `allow-relay-flags.txt` | Relay flags such as `-name` | `content/self-hosting/configuration.md` |
| `allow-desktop-commands.txt` | Subcommands; a nested one is `remote/enable` | `app/cli.html` |

Format, enforced by [the loader](/packages/docdrift-allowlist.md): one entry per line, name then whitespace then a one-line reason; `#` comments and blank lines allowed; no reason or a duplicate is an error.

Three conditions fail the check, so an entry cannot outlive its decision: a name with no reason, a name the docs now mention, a name the binary no longer has.

The relay file's header says its seed entries are flags tied to the hosted service's paid plans, which a self-hosted relay never uses. The desktop file's entries are not described further here.

Consumed by the [docdrift command](/cli/docdrift.md).
