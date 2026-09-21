# Self-hosting the Flockdeck relay

Remote access normally goes through the shared relay at `remote.flockdeck.ai`:
your desktop dials out to it, and your phone or laptop reaches your desktop
through it. That relay is operated by Flockdeck; terminal traffic is
end-to-end encrypted through it, but the state of your panes still passes
through in the clear — see [Security and
privacy](security-and-privacy.html).

For a company whose rules don't allow a third party to decrypt developers'
terminal traffic, the relay can instead be run on your own infrastructure —
single sign-on included; see
[Configuration](configuration.html#single-sign-on). **Flockdeck Enterprise**
is a licensed version of the same self-hosted relay, with support. It's
announced on [flockdeck.ai](https://flockdeck.ai/#enterprise) as coming soon
and isn't generally available yet, but self-hosting the relay — single
sign-on included — doesn't wait on it.

This section documents the relay as it exists and runs today, for a
platform or ops team evaluating what self-hosting will involve: how it's
deployed, how it's configured, what it needs, and how a desktop points at it
instead of the shared one. Everything here is accurate to the relay's
current behaviour, single sign-on included; nothing in it depends on
Enterprise or licensing existing yet. When you're ready to run this in
production, get in touch through
[flockdeck.ai](https://flockdeck.ai/#enterprise) for access and support.

## What the relay does

The relay is a single Go program. A desktop running Flockdeck connects
*out* to it over a WebSocket and keeps that connection open — nothing has to
be opened on the desktop's own network. A browser (your phone, tablet or
another computer) connects to the relay too, and the relay proxies its
requests down the desktop's tunnel to the copy of Flockdeck running there.
No ports are opened on the desktop's router, and no VPN is involved.

Everything else in this section assumes you're comfortable operating a small
stateful Go service: a binary or a container, a place to keep its data, and
a certificate.

## What's on the following pages

- [Requirements](requirements.html) — what the relay needs to run: ports,
  TLS, DNS, storage.
- [Deploying the relay](deploying-the-relay.html) — running it as a binary
  or a container, behind a reverse proxy, with your own certificate or with
  Flockdeck's built-in ACME support.
- [Configuration](configuration.html) — every flag and environment
  variable, in one table.
- [Storage and data](storage-and-data.html) — the embedded database and the
  MySQL alternative.
- [Admin and invites](admin-and-invites.html) — the local admin API,
  inviting desktops, and usage stats.
- [Connecting desktops](connecting-desktops.html) — pointing Flockdeck at
  your relay instead of the shared one.
- [Security and privacy](security-and-privacy.html) — what the relay can
  see, and what it can't.
- [Upgrading](upgrading.html) — keeping a self-hosted relay current.
- [Troubleshooting](troubleshooting.html) — the usual ways a first deploy
  goes wrong.
