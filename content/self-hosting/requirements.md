# Requirements

What a self-hosted relay needs before you deploy it.

## A domain, and DNS you control

The relay serves everything from the root of one hostname — for example
`relay.example.com` — not from a subpath of a site you already run. Point
that hostname's A/AAAA record at wherever the relay will listen, whether
that's the relay itself (with its built-in ACME support) or a load balancer
or reverse proxy in front of it.

If you want each paired desktop to get its own subdomain (`-desk-domain`,
see [Configuration](configuration.html)), you'll also need a wildcard record
— `*.relay.example.com` — and a wildcard certificate to match. This is
optional; the relay works without it.

## TLS

The relay needs to be reachable over HTTPS. Three ways to get there, covered
in [Deploying the relay](deploying-the-relay.html):

- the relay's own built-in ACME client, if it's the thing your DNS record
  points straight at;
- your own certificate and key, which the relay reloads without a restart
  when the files change; or
- a reverse proxy or load balancer in front of it that terminates TLS
  itself, with the relay behind it on plain HTTP.

Whichever you choose, WebSocket upgrades have to reach the relay unbuffered
— that's how a desktop's connection and a browser's both stay open. A proxy
that buffers or times out long-lived connections will break both.

## Compute

The relay is a single small Go binary with modest resource needs — closer
to a sidecar than a service you'd size a cluster around. It keeps no
per-request state beyond the open tunnels themselves, so it scales with the
number of connected desktops and devices, not with traffic volume.

## Storage

By default the relay keeps its state — paired desktops, devices, sessions —
in an embedded database file on disk, no separate database server needed.
MySQL 8 or later is supported as well, for a deployment that wants its state
in a database it already operates. See [Storage and
data](storage-and-data.html).

Either way, the relay's data directory needs to survive restarts and
redeploys: losing it forgets every paired desktop and device.

## Push notifications (optional)

If you want the relay to send push notifications to paired phones when an
agent is waiting, it needs a VAPID key pair, which the relay can generate
for you. See [Configuration](configuration.html).

## What you don't need

No database migrations to run by hand, no message queue, no separate cache.
The relay is one process and one data directory.
