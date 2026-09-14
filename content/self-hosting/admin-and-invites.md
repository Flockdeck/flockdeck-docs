# Admin and invites

The relay has a small admin surface, separate from the public one desktops
and devices talk to.

## The admin API

The admin API listens on `127.0.0.1:8081` — local only, never exposed to
the network the relay itself is reachable from. Reach it from the same
machine, or over SSH port forwarding, or from a sidecar in the same pod;
don't put it behind a public reverse proxy.

It's authenticated with a token the relay writes to `admin.json` in the
data directory (`-data`) the first time it starts. Keep that file as
tightly held as you would a root credential — anyone with it has full admin
access.

## Inviting a desktop

If `-registration` is set to `invite` (see [Configuration](configuration.html)),
a new desktop can't register itself; an admin has to issue an invite code
first:

```
flockdeck-relay invite -data /data -note "jane's laptop"
```

The `-note` is just for your own records in the stats output below — it
isn't shown to the desktop that redeems the code.

## Usage stats

```
flockdeck-relay stats -data /data
```

prints a summary of what's registered: how many desktops, how many paired
devices, and how close the deployment is to the limits set by `-max-hosts`
and `-max-devices`.

## Revoking access

There's no separate admin command for this yet — turning off remote access
from a desktop's own Settings, or running `flockdeck remote disable` there,
removes it from the relay it was registered with.
