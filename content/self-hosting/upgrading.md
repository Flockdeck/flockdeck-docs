# Upgrading

The relay is one process, and you don't migrate its database by hand: a newer
build updates its own tables when it starts. Changes to the schema only ever
add tables and columns, so the release before runs on the database a newer one
left behind, unless a release's own notes say otherwise.

## Before an upgrade

Read the release notes of every release between the one you run and the one you
are moving to; each says whether it changes the database and anything you have
to do. Back up the state first: a snapshot of the volume, or a copy of `relay.db`
taken while the relay is stopped; with MySQL, a dump of the database. Check the
relay is healthy now (`/readyz`, and `flockdeck-relay stats -data /data`, noting
how many desktops are connected), so anything wrong afterwards is known to be
new. The relay's repository has an upgrade guide, `docs/UPGRADING.md`, with the
order of operations for the Helm chart and for rolling back; it comes with the
kit a licensed company is given.

## Binary installs

Build the new version from a checkout of the relay's repository, as in [Deploying
the relay](deploying-the-relay.html#as-a-binary), stop the old process, replace
the binary and start it again under whatever supervises it. Never run two
relays against the same data.

## Container deployments

Pin the image tag, so a relay is upgraded when you say so. Deploy a newer one, stopping the old relay before starting the new one
(a recreate, not a rolling restart): never run two relays against the same
data. Because desktops hold a persistent connection to the relay, a
restart drops every open tunnel; desktops reconnect automatically, so a
brief gap in remote access during the restart is the only visible effect —
nothing an agent is doing locally is interrupted.

## Storage across an upgrade

There's no separate export/import step to run yourself; the relay opens its
existing data directory or database on startup and updates its tables. To go
back, put the previous image or binary back: it ignores what it doesn't know. A
release whose notes say otherwise is the exception, and says what to do
instead, restoring the backup among it.

## Checking what's running

`/healthz` and `/readyz` (see [Deploying the relay](deploying-the-relay.html))
don't report a version. Three things do: `flockdeck-relay version`, the first
line of `flockdeck-relay stats`, and the `version=` field in the relay's
"relay listening" log line at startup. Check that against the relay's own
release notes to know what you're running.
