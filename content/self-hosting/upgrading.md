# Upgrading

The relay is one process with no in-place schema you have to migrate by
hand — restart it on a newer build and it's upgraded.

## Binary installs

```
go install github.com/jmwri/flockdeck-relay/cmd/flockdeck-relay@latest
```

again, then restart the process under whatever supervises it.

## Container deployments

Deploy a newer image tag and let your orchestrator do its usual rolling
restart. Because desktops hold a persistent connection to the relay, a
restart drops every open tunnel; desktops reconnect automatically, so a
brief gap in remote access during the restart is the only visible effect —
nothing an agent is doing locally is interrupted.

## Storage across an upgrade

There's no separate export/import step to run yourself; the relay opens its
existing data directory or database on startup. Back up the data directory
or database before upgrading regardless, the way you would before
upgrading any stateful service, and check that release's own notes for
anything it calls out about storage.

## Checking what's running

`/healthz` and `/readyz` (see [Deploying the relay](deploying-the-relay.html))
don't report a version. Check the binary or image tag you deployed against
the relay's own release notes to know what you're running.
