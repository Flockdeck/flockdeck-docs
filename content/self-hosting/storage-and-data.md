# Storage and data

The relay keeps one kind of state: which desktops are registered, which
devices are paired to them, and their active sessions. There's no
per-message or per-keystroke history to store — traffic is proxied through,
not recorded.

## The embedded database

By default, pointing `-data` at a directory is all storage configuration the
relay needs. It keeps a single database file (`relay.db`) there, using an
embedded key-value store — no separate database server, no migrations to
run by hand.

This is the right choice for most self-hosted deployments: one relay
process, one volume, nothing else to operate.

## MySQL

For a deployment that wants the relay's state in a database it already
runs, set `-database` (or `FLOCKDECK_RELAY_DATABASE_URL`) to a MySQL 8 or
later DSN instead of `-data`. Use `-database-ca` if the connection needs a
CA certificate verified.

`/readyz` doesn't answer until the database connection is live, so a load
balancer or orchestrator waiting on readiness won't route traffic to a
relay that can't reach its database yet.

## Backing it up

Whichever storage you use, back it up like you would any other small
stateful service: a volume snapshot for the embedded database, or your
existing MySQL backup process. Losing it means every desktop and device has
to register and pair again — nobody's agents or conversations are affected,
since none of that passes through or is kept by the relay.

## Moving between them

There's no built-in migration between the embedded database and MySQL.
Moving from one to the other means every desktop and device pairs again.
