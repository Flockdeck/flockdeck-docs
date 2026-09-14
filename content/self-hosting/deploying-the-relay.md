# Deploying the relay

The relay is one Go program, `flockdeck-relay`, run as `flockdeck-relay
serve` with flags (or the matching environment variables — see
[Configuration](configuration.html)).

## As a binary

```
go install github.com/jmwri/flockdeck-relay/cmd/flockdeck-relay@latest
```

installs the current release. Run it directly, or under whatever process
supervisor your platform already uses (systemd, a container orchestrator's
own restart policy, and so on) — the relay itself doesn't daemonize.

## As a container

The relay's repository has a Dockerfile. Building it needs a build secret,
because the relay embeds Flockdeck's own web client and that module isn't
public:

```
docker build --secret id=flockdeck_remote_token,env=FLOCKDECK_REMOTE_TOKEN -t flockdeck-relay .
```

An Enterprise licence includes access to a built image, so you don't need to
build it yourself. Once you have an image, running it looks like:

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
[Storage and data](storage-and-data.html).

## Three ways to terminate TLS

### 1. The relay's own ACME client

If DNS for your relay's hostname points straight at the machine running it,
the relay can get and renew its own certificate:

```
flockdeck-relay serve -acme-domain relay.example.com -acme-email ops@example.com -data /data
```

This is the simplest option when there's nothing else in front of the
relay.

### 2. Your own certificate

```
flockdeck-relay serve -tls-cert /etc/relay/fullchain.pem -tls-key /etc/relay/privkey.pem -data /data
```

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
