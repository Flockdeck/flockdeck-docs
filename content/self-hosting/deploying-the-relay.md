# Deploying the relay

The relay is one Go program, `flockdeck-relay`, run as `flockdeck-relay
serve` with flags (or the matching environment variables — see
[Configuration](configuration.html)).

## As a binary

The relay's source is in a private repository, so there is no public
download or `go install` for it: you build it from a checkout, which an
Enterprise licence gives you access to. It is pure Go, so one command in the
checkout builds it for the machine it runs on:

```
go build -trimpath -o flockdeck-relay ./cmd/flockdeck-relay
```

Building needs a Go toolchain of the version in `go.mod`, and a git
credential that can read the web client's module (`github.com/Flockdeck/flockdeck-remote`,
also private), with `GOPRIVATE` set to it. Run the binary directly, or under
whatever process supervisor your platform already uses (systemd, a container
orchestrator's own restart policy, and so on) — the relay itself doesn't
daemonize.

## As a container

The relay's repository has a Dockerfile. Building it needs a build secret,
because the relay embeds Flockdeck's own web client and that module isn't
public:

```
docker build --secret id=flockdeck_remote_token,env=FLOCKDECK_REMOTE_TOKEN -t flockdeck-relay .
```

An Enterprise licence includes access to a built image, `ghcr.io/flockdeck/flockdeck-relay`
(in a private registry, so pulling it needs the credential issued with the
licence), so you don't need to build it yourself. The image's command is
`serve -data /data -addr :8080`; `-public-url` has no default and is always
given at run time. Running it looks like:

```
docker run \
  -p 127.0.0.1:8080:8080 \
  -v relay-data:/data \
  flockdeck-relay serve \
    -public-url https://relay.example.com \
    -addr :8080 \
    -trust-proxy \
    -data /data
```

`-v relay-data:/data` is the volume that has to survive restarts — see
[Storage and data](storage-and-data.html). The image runs as a non-root user,
and `/data` in it is owned by that user, so a fresh named volume works as is.

## On Kubernetes

The relay's repository also has a Helm chart, `deploy/helm/flockdeck-relay`,
for a licensed company running the relay on its own cluster. It runs exactly
one pod and replaces it on an upgrade rather than rolling it: the relay keeps
every desktop's tunnel in memory, so a second replica would be a second relay
that knows about only some of the desktops. It needs Kubernetes 1.24 or later,
Helm 3, an ingress that passes WebSockets through without buffering, a TLS
certificate for the relay's hostname, and the pull token issued with the
licence. The chart's own README walks through installing it.

## Three ways to terminate TLS

### 1. The relay's own ACME client

If DNS for your relay's hostname points straight at the machine running it,
the relay can get and renew its own certificate:

```
flockdeck-relay serve -acme-domain relay.example.com -acme-email ops@example.com -data /data
```

This is the simplest option when there's nothing else in front of the
relay. It listens on `:443` unless you give `-addr`, and on `:80` for the
challenges and a redirect to https. `-public-url` defaults to `https://` and
the first domain, and a list of domains (comma-separated) gets a certificate
for each; one of them has to be the host of `-public-url`. The certificates
are kept in `acme` under `-data`. It can't be combined with `-tls-cert`, and
it can't get the wildcard certificate `-desk-domain` needs.

### 2. Your own certificate

```
flockdeck-relay serve -public-url https://relay.example.com -addr :443 -tls-cert /etc/relay/fullchain.pem -tls-key /etc/relay/privkey.pem -data /data
```

`-tls-cert` and `-tls-key` go together, and `-public-url` has to be `https://`.
The relay watches these files and reloads them when they change, so a
certificate renewal elsewhere on the machine doesn't need a restart.

### 3. Behind a reverse proxy or load balancer

If something else — an ingress controller, a load balancer, another proxy —
already terminates TLS, run the relay on plain HTTP behind it:

```
flockdeck-relay serve -public-url https://relay.example.com -addr 127.0.0.1:8080 -trust-proxy -data /data
```

`-public-url` is what the relay tells desktops and devices to use, even
though it's listening on plain HTTP internally. `-trust-proxy` tells it to
take the client address from forwarded headers rather than the TCP
connection, which is the proxy's, not the visitor's.

Whichever proxy sits in front, it has to forward WebSocket upgrades without
buffering — see [Requirements](requirements.html).

## Confirm it's up

Two health endpoints:

- `/healthz` answers as soon as the process is up.
- `/readyz` answers only once its storage is ready — with MySQL, that means
  the database connection is live. Point a load balancer or orchestrator's
  readiness probe at this one, not `/healthz`, or it'll route traffic to a
  relay that can't yet serve it.
