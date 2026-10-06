# Admin and invites

The relay has a small admin surface, separate from the public one desktops
and devices talk to.

## The admin API

The admin API listens on `-admin-addr`, which defaults to loopback
(`127.0.0.1:8081`); keep it there, off the network the relay itself is
reachable from. Reach it from the same machine, or over SSH port
forwarding, or from a sidecar in the same pod; don't put it behind a public
reverse proxy. The relay warns if it is listening beyond loopback.

It's authenticated with a token the relay writes to `admin.json` in the
data directory (`-data`). The relay writes a fresh token on every start and
removes the file when it stops, so run the admin commands (`invite`,
`stats`, `plan`, `user revoke`, `account` and `audit export`) where they can read
`-data`, for example `docker exec <container> /flockdeck-relay invite -data
/data`. Keep that file as tightly held as you would a root credential:
anyone with it has full admin access.

## Inviting a desktop

If `-registration` is set to `invite` (see [Configuration](configuration.html)),
a new desktop can't register itself. An admin has to issue an invite code
first:

```
flockdeck-relay invite -data /data -note "jane's laptop"
```

The command prints the code, and the exact `flockdeck remote enable -relay
… -invite <code>` command to hand over. The `-note` is kept with the invite
for your own records; `stats` doesn't print it, and it isn't shown to the
desktop that redeems the code.

## Usage stats

```
flockdeck-relay stats -data /data
```

prints the relay's version and uptime, and a summary of what it holds: how
many accounts, desktops (and how many are connected), devices (and how many
are notified by push), pairing codes waiting and invites unused, the plans
if the relay enforces them, and how many requests each limit has refused.
The refused counts show a proxy that makes every client look like one
address. `-max-hosts` and `-max-devices` are per account, so `stats` doesn't
show headroom against them.

## Revoking access

To revoke a person, run `flockdeck-relay user revoke -data /data
sam@example.com`. It revokes every device that person paired, and needs a
relay run with `-oidc-issuer`, which is how a device is recorded against an
email; disable them at the identity provider too, or they could pair a
device again.

To take a desktop off its account, run `flockdeck remote remove` on it, or
remove it from a paired browser. That removes only the desktop, the last one
included: the account keeps its devices and verified email, and can have no
desktops. `flockdeck remote disable` only turns Flockdeck Remote off on that
desktop, which stays enrolled and shows as offline.

## Deleting accounts

Nothing deletes an account when its last desktop goes. The relay deletes one
once nothing on it has been used for `-account-retention` (see [Storage and
data](storage-and-data.html#how-long-accounts-are-kept)). To delete one
yourself:

```
flockdeck-relay account delete -data /data -email sam@example.com
flockdeck-relay account delete -data /data -email sam@example.com -apply -export sam.json
```

Without `-apply` it is a dry run that says what would go. `-apply` needs
`-export`, which writes the account's rows out first; the account is deleted only
after that. A subscribed account is refused unless you give `-force`, since
deleting does not cancel a subscription.

`account duplicates` lists the emails that have more than one account, and
`account dedupe` (a dry run without `-apply`) deletes the unpaid duplicates, with
the same `-export` rule. Both are permanent and the export is the only copy, so
take a database snapshot first.

When the command runs inside a container with `docker exec` or `kubectl exec`,
a file named with `-export FILE` is written inside the container, where you
cannot easily copy it out and where a restart loses it. Give `-export -
-confirm-exported` instead, which writes the export to standard output, and
redirect that to a file on your own machine:

```
kubectl exec deploy/relay -- /flockdeck-relay account dedupe -data /data -apply -export - -confirm-exported > rows.json
```

Standard output is not synced, so check that `rows.json` is not empty and parses
as JSON before you rely on it.

## Exporting the audit log

To export the audit log, run `flockdeck-relay audit export -data /data
-format csv -o audit.csv`. `-format` is `csv` (the default) or `json`; `-since`
and `-until` bound it by a date (`2026-09-01`) or a time
(`2026-09-01T09:00:00Z`); without `-o` it writes to the standard output. The
relay does not have to be running with `-audit`: without it, the command says
the relay is not keeping the log now, and writes what was kept before.
