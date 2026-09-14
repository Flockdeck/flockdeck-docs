# Configuration

Every flag `flockdeck-relay serve` takes, with the environment variable that
does the same thing where one exists. A flag wins over its environment
variable if both are set.

## Network and TLS

| Flag | Env | Default | What it does |
| --- | --- | --- | --- |
| `-addr` | | `:443` | Address to listen on. |
| `-public-url` | | | The address the relay tells desktops and devices to use — required whenever it isn't the same as `-addr`, such as behind a proxy. |
| `-trust-proxy` | | off | Take the client address from forwarded headers instead of the TCP connection. Only turn this on when something trusted sits in front of the relay. |
| `-acme-domain` | | | Domain to get a certificate for automatically. Requires the relay to be reachable on port 80 and 443 directly. |
| `-acme-email` | | | Contact address for the ACME account. |
| `-tls-cert` / `-tls-key` | | | Your own certificate and key. Reloaded automatically when the files change. |
| `-desk-domain` | | | Base domain for giving each paired desktop its own subdomain. Needs a wildcard DNS record and certificate. |

## Storage

| Flag | Env | Default | What it does |
| --- | --- | --- | --- |
| `-data` | | | Directory for the embedded database and the admin token file. |
| `-database` | `FLOCKDECK_RELAY_DATABASE_URL` | | MySQL DSN, in place of the embedded database. See [Storage and data](storage-and-data.html). |
| `-database-ca` | `FLOCKDECK_RELAY_DATABASE_CA` | | Path to a CA certificate for the MySQL connection. |

## Registration and limits

| Flag | Default | What it does |
| --- | --- | --- |
| `-registration` | `open` | Who can register a new desktop: `open`, `invite` (needs an invite code from an admin — see [Admin and invites](admin-and-invites.html)), or `closed`. |
| `-max-hosts` | 10 | Maximum number of registered desktops. |
| `-max-devices` | 20 | Maximum number of paired phones/browsers per desktop. |
| `-pairing-ttl` | 10m | How long a pairing link or QR code stays valid. |
| `-session-ttl` | 720h (30 days) | How long a paired device stays paired before it has to pair again. |
| `-register-per-hour` | 5 | Rate limit on new desktop registrations, per IP. |
| `-pair-per-minute` | 10 | Rate limit on pairing attempts, per IP. |
| `-connects-per-minute` | 60 | Rate limit on new connections, per IP. |
| `-tunnels-per-address` | 20 | Maximum concurrent tunnels from one address. |
| `-max-tunnels` | 1000 | Maximum concurrent tunnels, relay-wide. |
| `-memory-limit` | | A Go `GOMEMLIMIT`-style soft memory cap for the process. |

## Push notifications

| Flag | Env | What it does |
| --- | --- | --- |
| `-vapid-private-key` | `FLOCKDECK_RELAY_VAPID_PRIVATE_KEY` | Private key for signing web push messages. |
| `-vapid-public-key` | `FLOCKDECK_RELAY_VAPID_PUBLIC_KEY` | The matching public key, which paired devices are given. |

Generate a pair with:

```
flockdeck-relay vapid-keys
```

Without these set, the relay runs fine; paired devices just don't get push
notifications.

## Single sign-on

Not available in the relay's main branch yet — it's part of what
Enterprise will add. See [Overview](overview.html).

## An example

Behind a reverse proxy, with MySQL, registration by invite only, and push
notifications on:

```
flockdeck-relay serve \
  -public-url https://relay.example.com \
  -addr 127.0.0.1:8080 \
  -trust-proxy \
  -database "flockdeck:$MYSQL_PASSWORD@tcp(db.internal:3306)/flockdeck_relay" \
  -registration invite \
  -vapid-private-key "$VAPID_PRIVATE_KEY" \
  -vapid-public-key "$VAPID_PUBLIC_KEY"
```
