# Connecting desktops

Once your relay is running, point Flockdeck at it instead of the shared one
at `remote.flockdeck.ai`.

## Enabling remote access against your relay

From any terminal on the desktop:

```
flockdeck remote enable -relay https://relay.example.com
```

or set the `FLOCKDECK_RELAY` environment variable to the same address
before Flockdeck starts, which is easier across a fleet of machines than
running the command on each one. The flag wins if both are set. `-invite`, `-join` and `-name` have environment variables
too (`FLOCKDECK_REMOTE_INVITE`, `FLOCKDECK_REMOTE_JOIN` and
`FLOCKDECK_REMOTE_NAME`), for enrolling with no one at the keyboard.

This registers the desktop with your relay and turns remote access on. With
`-registration invite`, the relay needs an invite code an admin issued.
See [Admin and invites](admin-and-invites.html), where `flockdeck-relay
invite` prints the exact command to hand over:

```
flockdeck remote enable -relay https://relay.example.com -invite <code>
```

Everything else works as
[Flockdeck Remote](/app/remote.html) describes for the shared relay, including
pairing a phone with a QR code and notifications, but against your own
infrastructure.

## A second desktop on the same relay

A plain `flockdeck remote enable` creates a separate account. To put a second
desktop on the same account, run `flockdeck remote pair -desktop` on a
desktop already on the relay, and on the new one run `flockdeck remote enable
-relay https://relay.example.com -join <code>`. Joining never needs an
invite and works even with `-registration closed`; `-max-hosts` is per
account, so a join fails once the account has that many desktops.

## Joining by email

On a relay run with `-require-verified-registration`, `flockdeck remote enable`
asks for an email address and sends a link to it. The desktop prints a match
code, and the page the link opens asks you to type it. If the address already
has an account, the desktop joins it: same plan, same paired devices, no second
account. The page says how many desktops the account has before you go back to
the terminal. The account's paired phones and browsers are told a desktop joined,
and the address is emailed the same, so a desktop that someone else enrolled
does not go unnoticed. If you did not enrol anything, remove the desktop with
`flockdeck remote remove` from another machine.

## A desktop that is already on another relay

```
flockdeck remote move https://relay.example.com
```

enrols the desktop with your relay first, and leaves the one it was on only
once yours answers. Add `-invite <code>` or `-join <code>` as for `enable`,
and `-yes` to skip the question it asks. Every device paired with the desktop
has to pair again afterwards, because a device's pairing belongs to the relay
it was made on.

## Switching a desktop back to the shared relay

```
flockdeck remote move https://remote.flockdeck.ai
```

or, step by step:

```
flockdeck remote remove
flockdeck remote enable
```

`remove` takes the desktop off whichever relay it was registered with;
running `enable` again with no `-relay` flag and no `FLOCKDECK_RELAY` set
registers it with the shared one. Either way the devices pair again.

## Confirming it's talking to the right relay

Settings → Flockdeck Remote shows which relay a desktop is currently registered
with. If a desktop that should be on your relay is showing up against the
shared one instead, check that `FLOCKDECK_RELAY` is set in the
environment Flockdeck is launched from. A value set in one shell profile
won't reach a copy of Flockdeck started from a desktop shortcut, a different
shell, or a service manager.
