# Connecting desktops

Once your relay is running, point Flockdeck at it instead of the shared one
at `remote.flockdeck.ai`.

## Enabling remote access against your relay

From any terminal on the desktop:

```
flockdeck remote enable -relay https://relay.example.com
```

or set the `FLOCKDECK_RELAY` environment variable to the same address
before Flockdeck starts, if you'd rather configure it that way across a
fleet of machines than run the command on each one by hand. The flag wins
if both are set.

This registers the desktop with your relay (or, with `-registration invite`,
redeems an invite code an admin issued — see [Admin and
invites](admin-and-invites.html)) and turns remote access on, exactly as
[Remote access](/app/remote.html) describes for the shared relay — pairing a
phone with a QR code, notifications, everything — just against your own
infrastructure instead.

## A second desktop on the same relay

Each desktop that should be reachable through your relay runs its own
`flockdeck remote enable -relay ...` (or has `FLOCKDECK_RELAY` set). There's
no separate "join" step for a second desktop beyond registering it the same
way as the first, subject to whatever `-registration` mode and `-max-hosts`
limit your relay is configured with.

## Switching a desktop back to the shared relay

```
flockdeck remote disable
flockdeck remote enable
```

`disable` removes the desktop from whichever relay it was registered with;
running `enable` again with no `-relay` flag and no `FLOCKDECK_RELAY` set
registers it with the shared one.

## Confirming it's talking to the right relay

Settings → Remote access shows which relay a desktop is currently registered
with. If a desktop that should be on your relay is showing up against the
shared one instead, check that `FLOCKDECK_RELAY` is actually set in the
environment Flockdeck is launched from — a value set in one shell profile
won't reach a copy of Flockdeck started from a desktop shortcut, a different
shell, or a service manager.
