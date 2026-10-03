# Storage and data

The relay keeps account and device records: accounts, desktops (name and
when last seen), devices (name, browser user agent and session expiry),
pairing codes and invites not yet used (with the admin's note on an
invite), push subscriptions, and the end-to-end public keys of desktops and
devices. Some settings add more: with single sign-on, each person's
identity and email; with `-audit`, the audit log (emails and names, and
addresses with `-audit-ip`); with verified registration, the verified email.
There's no per-message or per-keystroke history to store: traffic is
proxied through, not recorded, and nothing typed or shown is kept.

## The embedded database

By default, pointing `-data` at a directory is all the storage configuration
the relay needs. It keeps a single database file (`relay.db`) there, using an
embedded key-value store, with no separate database server and no migrations
to run by hand.

This is the right choice for most self-hosted deployments: one relay
process, one volume, nothing else to operate.

## MySQL

For a deployment that wants the relay's state in a database it already
runs, set `-database` (or `FLOCKDECK_RELAY_DATABASE_URL`) to a MySQL 8 or
later DSN instead of the embedded database. `-data` is still used, for the
admin token and ACME certificates, just not for state. Use `-database-ca`
if the connection needs a CA certificate verified. The relay creates and
updates its own tables on start, so its database user needs to be able to
create and alter tables.

`/readyz` doesn't answer until the database connection is live, so a load
balancer or orchestrator waiting on readiness won't route traffic to a
relay that can't reach its database yet.

## Backing it up

Back it up like any other small stateful service, whichever storage you
use: a volume snapshot for the embedded database, or your
existing MySQL backup process. Losing it means every desktop and device has
to register and pair again. Nobody's agents or conversations are affected,
since none of that passes through or is kept by the relay.

## Moving between them

There's no built-in migration between the embedded database and MySQL.
Moving from one to the other means every desktop and device pairs again.
