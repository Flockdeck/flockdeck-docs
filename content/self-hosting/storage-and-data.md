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

Each account also has one stored time, when anything on it was last used. The
relay uses it to decide when to delete an account nobody uses.

## How long accounts are kept

An account outlives its desktops. Removing the last desktop, with `flockdeck
remote remove` or from a paired browser, leaves the account with its devices
and verified email, and a desktop can join it again by enrolling with the same
email address, or with a join code from a paired browser. Removing a desktop
never deletes the account.

The relay deletes an account, with its desktops, devices and codes, in three
ways:

- **Retention.** Once an hour, in the ten-minute sweep, on every relay, the
  relay deletes each account that has not been used for `-account-retention`.
  The default is 90 days. `0` keeps accounts for good, and a value above 0 and
  under 24 hours is refused when the relay starts. Use is a desktop connecting,
  a paired device making a request, a registration, a join, or a desktop being
  removed. An account with a subscription is not deleted. The relay logs how
  many accounts, desktops and devices it deleted, never an id or an address.
  With `-audit` it also records an `account.deleted` event for each.
- **Its owner.** `flockdeck remote delete-account` on a desktop deletes the
  account, every machine in it and every paired device. The relay refuses while
  a subscription is running, because it does not cancel subscriptions.
- **You, by hand.** `flockdeck-relay account delete`, described in [Admin and
  invites](admin-and-invites.html#deleting-accounts).

A desktop that registered and never connected is removed after 7 days, and only
the desktop. There is no sweep of quiet desktops: one that has connected is
never removed for being quiet.

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
