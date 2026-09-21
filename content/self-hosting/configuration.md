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

| Flag | Env | What it does |
| --- | --- | --- |
| `-oidc-issuer` | `FLOCKDECK_RELAY_OIDC_ISSUER` | Have everyone pairing a device sign in with this OpenID Connect issuer first, e.g. `https://example.okta.com`. |
| `-oidc-client-id` | `FLOCKDECK_RELAY_OIDC_CLIENT_ID` | The relay's client id at the identity provider. |
| `-oidc-client-secret-file` | `FLOCKDECK_RELAY_OIDC_CLIENT_SECRET` | File holding the relay's client secret at the identity provider — or, in the environment variable, the secret itself. There's no `-oidc-client-secret` flag, so it never has to appear in a process listing. |
| `-oidc-name` | `FLOCKDECK_RELAY_OIDC_NAME` | What the sign-in button calls the identity provider. The provider's own host by default. |
| `-oidc-email-domains` | `FLOCKDECK_RELAY_OIDC_EMAIL_DOMAINS` | Comma-separated email domains whose people may sign in. Anyone the provider signs in, by default. |

`-oidc-issuer` needs a client id and a secret to go with it; either given
without the other, or without `-oidc-issuer`, is a startup error. This is
part of the relay itself — it doesn't need Flockdeck Enterprise or a
licence. See [Overview](overview.html).

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
