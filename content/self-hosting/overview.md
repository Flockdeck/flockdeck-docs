# Self-hosting the Flockdeck relay

Flockdeck Remote normally goes through the shared relay at `remote.flockdeck.ai`:
your desktop dials out to it, and your phone or laptop reaches your desktop
through it. Flockdeck operates that relay. Only terminal traffic is
end-to-end encrypted through it; everything else, the state of your panes
included, passes through decrypted. See [Security and
privacy](security-and-privacy.html).

Terminal traffic is end-to-end encrypted with keys the relay never holds, so
even the shared relay can't read it, unless a terminal falls back to plaintext
because either side has no registered key, or the relay tampers with the client
it serves or swaps the keys it hands out. Apart from push notification content,
which neither the relay nor the push service can read, whoever operates the
relay can read the rest of what passes through: pane state, the chat view,
diffs, photos and typed API keys (see [Security and
privacy](security-and-privacy.html)). For a company whose rules don't allow a
third party to read those, the relay can instead be run on your own
infrastructure. That is Flockdeck Enterprise: the licensed, supported, self-hosted relay,
single sign-on included (see
[Configuration](configuration.html#single-sign-on)). It's announced on
[flockdeck.ai](https://flockdeck.ai/#enterprise) as coming soon and isn't
generally available yet.

This section documents the relay as it exists and runs today, for a
platform or ops team evaluating what self-hosting involves: how it's
deployed, how it's configured, what it needs, and how a desktop points at it
instead of the shared one. It covers the relay's current behaviour, single
sign-on included. When you're ready to run this in production, get in touch
through [flockdeck.ai](https://flockdeck.ai/#enterprise) for access and
support.

## What the relay does

The relay is a single Go program. A desktop running Flockdeck connects
*out* to it over a WebSocket and keeps that connection open, so nothing has
to be opened on the desktop's own network. A browser (your phone, tablet or
another computer) connects to the relay too, and the relay proxies its
requests down the desktop's tunnel to the copy of Flockdeck running there.
No ports are opened on the desktop's router, and no VPN is involved.

You'll need to be comfortable operating a small stateful Go service: a
binary or a container, a place to keep its data, and a certificate.

## What's on the following pages

- [Requirements](requirements.html): ports, TLS, DNS and storage.
- [Deploying the relay](deploying-the-relay.html): running it as a binary
  or a container, behind a reverse proxy, with your own certificate or with
  the relay's built-in ACME support.
- [Configuration](configuration.html): every flag and environment
  variable, in one table.
- [Storage and data](storage-and-data.html): the embedded database and the
  MySQL alternative.
- [Admin and invites](admin-and-invites.html): the local admin API,
  inviting desktops, and usage stats.
- [Connecting desktops](connecting-desktops.html): pointing Flockdeck at
  your relay instead of the shared one.
- [Security and privacy](security-and-privacy.html): what the relay can
  see, and what it can't.
- [Upgrading](upgrading.html): keeping a self-hosted relay current.
- [Troubleshooting](troubleshooting.html): the usual ways a first deploy
  goes wrong.
