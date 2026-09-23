# Configuration

Every flag `flockdeck-relay serve` takes, with the environment variable that
does the same thing where one exists. A flag wins over its environment
variable if both are set.

## Network and TLS

| Flag | Env | Default | What it does |
| --- | --- | --- | --- |
| `-addr` | | `:8080` (`:443` with `-acme-domain`) | Address to listen on. |
| `-public-url` | | | The address the relay tells desktops and devices to use — required whenever it isn't the same as `-addr`, such as behind a proxy. |
| `-trust-proxy` | | off | Take the client address from forwarded headers instead of the TCP connection. Only turn this on when something trusted sits in front of the relay. |
| `-acme-domain` | | | Domain to get a certificate for automatically. Requires the relay to be reachable on port 80 and 443 directly. |
| `-acme-email` | | | Contact address for the ACME account. |
| `-http-addr` | | `:80` | With `-acme-domain`, the address answering ACME challenges and redirecting to https. Does nothing without `-acme-domain`, and is refused if set without it. |
| `-tls-cert` / `-tls-key` | | | Your own certificate and key. Reloaded automatically when the files change. |
| `-client-dir` | | | Serve the remote client from this directory instead of the one built into the relay. |
| `-desk-domain` | | | Base domain for giving each paired desktop its own subdomain. Needs a wildcard DNS record and certificate. |
| `-admin-addr` | | `127.0.0.1:8081` | Loopback address for the admin endpoints `flockdeck-relay invite`, `stats` and `plan` use — see [Admin and invites](admin-and-invites.html). Empty turns them off. |

## Storage

| Flag | Env | Default | What it does |
| --- | --- | --- | --- |
| `-data` | | `data` | Directory for the embedded database and the admin token file. |
| `-database` | `FLOCKDECK_RELAY_DATABASE_URL` | | MySQL DSN, in place of the embedded database. See [Storage and data](storage-and-data.html). |
| `-database-ca` | `FLOCKDECK_RELAY_DATABASE_CA` | | Path to a CA certificate for the MySQL connection. |

## Registration and limits

| Flag | Default | What it does |
| --- | --- | --- |
| `-registration` | `open` | Who can register a new desktop: `open`, `invite` (needs an invite code from an admin — see [Admin and invites](admin-and-invites.html)), or `closed`. |
| `-max-hosts` | 10 | Maximum number of desktops per account (0: no limit). |
| `-max-devices` | 20 | Maximum number of paired phones/browsers per account (0: no limit). |
| `-max-pairings` | 5 | Unused pairing codes an account can hold at once; a new one replaces the oldest (0: no limit). |
| `-pairing-ttl` | 10m | How long a pairing link or QR code stays valid. |
| `-session-ttl` | 720h (30 days) | How long a device stays paired without being used, before it has to pair again. |
| `-max-streams` | 256 | Connections open to one desktop at once. |
| `-register-per-hour` | 5 | Accounts made, and join codes that don't work, per client address per hour (0: no limit). |
| `-pair-per-minute` | 10 | Pairing attempts per client address per minute (0: no limit). |
| `-connects-per-minute` | 60 | Desktop tunnels opened with a token that doesn't work, per client address per minute (0: no limit). |
| `-tunnels-per-address` | 20 | Maximum desktops connected from one address at once (0: no limit). |
| `-max-tunnels` | 1000 | Maximum desktops connected to the relay at once (0: no limit). |
| `-log-level` | `info` | `debug`, `info`, `warn` or `error`. |
| `-memory-limit` | | A Go `GOMEMLIMIT`-style soft memory cap for the process. |

## Verified registration

| Flag | Env | Default | What it does |
| --- | --- | --- | --- |
| `-require-verified-registration` | `FLOCKDECK_RELAY_REQUIRE_VERIFIED_REGISTRATION` | off | Require a verified email, confirmed with a magic link, before `flockdeck remote enable` may create an account and start a trial. Needs a Postmark token and `-mail-from`. |
| `-mail-from` | | | The relay's own address verification emails are sent from, e.g. `verify@flockdeck.ai`. Needed with `-require-verified-registration`, and refused without it. |
| `-verify-email-per-hour` | | 5 | Verification emails started or resent per client address per hour (0: no limit). |
| `-verification-ttl` | | 30m | How long someone has to click a verification email's link before it expires. |

The Postmark server token itself is read only from `FLOCKDECK_RELAY_POSTMARK_TOKEN`
— there is no flag for it, so it never appears in a process listing.
`-require-verified-registration` refuses to start without both it and
`-mail-from` set.

The desktop app's command-palette dialog can't yet wait on a verification
link, so it fails immediately on a relay set up this way; `flockdeck remote
enable` and `flockdeck remote move` from a terminal open the link and wait
for it instead. See [Remote access](/app/remote.html).

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

## Audit log

| Flag | Default | What it does |
| --- | --- | --- |
| `-audit` | off | Keep the audit log: who reached which desktop, from which device, and when. Never what was typed or shown. |
| `-audit-ip` | off | With `-audit`, also record the address each event came from. Refused without `-audit`. |
| `-audit-retention` | 0 (for good) | With `-audit`, how long events are kept. Refused without `-audit`. |

See [Security and privacy](security-and-privacy.html) for what the audit log
does and doesn't capture.

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
