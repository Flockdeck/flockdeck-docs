# Troubleshooting

The ways a first self-hosted deployment usually goes wrong, and where to
look.

## Desktops or devices can't connect at all

Almost always a proxy swallowing WebSocket upgrades. If the relay is behind
a reverse proxy or load balancer, confirm it's configured to forward
WebSocket connections through unbuffered — see
[Requirements](requirements.html). A proxy that works for ordinary HTTP
traffic can still silently break this.

## `/readyz` never returns healthy

With MySQL storage, this means the relay can't reach the database: check
`-database` (or `FLOCKDECK_RELAY_DATABASE_URL`), that the database exists,
and `-database-ca` if the connection needs a CA certificate. With the
embedded database, check that `-data` points at a directory the relay's
process can actually write to.

## Certificate errors

- Using `-acme-domain`: the relay needs to be reachable on ports 80 and 443
  directly from the internet for the ACME challenge to complete. Behind a
  proxy or load balancer, use `-tls-cert`/`-tls-key` or terminate TLS at the
  proxy instead — see [Deploying the relay](deploying-the-relay.html).
- Using `-tls-cert`/`-tls-key`: the relay reloads these when they change, so
  check the files were actually replaced (not just touched) and that the
  process can read them.
- Using `-desk-domain`: this needs a *wildcard* certificate covering
  `*.relay.example.com`, not just the base domain.

## A desktop registers against the shared relay instead of yours

`FLOCKDECK_RELAY` has to be set in the environment the copy of Flockdeck
that's actually running was launched from — see [Connecting
desktops](connecting-desktops.html). A variable set in one shell profile
doesn't reach Flockdeck launched from a desktop shortcut, a different
shell, or a service manager. Passing `-relay` on `flockdeck remote enable`
directly sidesteps this.

## New desktops can't register

If `-registration` is `invite`, every new desktop needs a code from
`flockdeck-relay invite` first — see [Admin and invites](admin-and-invites.html).
If it's `closed`, no new desktop can register regardless. Check
`-max-hosts` too: once it's reached, registration fails until an existing
desktop is removed.

## Rate limiting

`-register-per-hour`, `-pair-per-minute` and `-connects-per-minute` (see
[Configuration](configuration.html)) are all per source address. Behind a
proxy, if `-trust-proxy` isn't set, the relay sees every request as coming
from the proxy's own address and rate-limits everyone behind it together —
turn it on, and make sure the proxy is actually setting the forwarded
headers it relies on.

## Still stuck

Confirm the relay itself considers itself healthy — `/healthz` and
`/readyz` — before chasing anything else; a relay that isn't up yet
produces symptoms that look like every problem above at once.
