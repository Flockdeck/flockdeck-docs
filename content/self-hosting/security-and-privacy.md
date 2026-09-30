# Security and privacy

What a relay — yours or the shared one — can and can't see, so you can
weigh that against what your own rules require.

## End-to-end encrypted terminal traffic, with one caveat

Traffic between a desktop and a paired device is encrypted with TLS to the
relay and from the relay onward. On top of that, terminal traffic — what you
type and what an agent prints back — is end-to-end encrypted by the desktop
app and the paired device themselves, using keys the relay hands out but
never holds: the relay has to route it, to know which desktop a browser's
request is for, but it cannot read it. This is true of the shared relay at
`remote.flockdeck.ai` and of a relay you run yourself in exactly the same
way.

Only the terminal is end-to-end encrypted, and only when the desktop and the
paired device have each registered a key. A terminal opened when either has
none is served without it, and a relay that holds back a device's key can
cause that.

Everything else passes through the relay decrypted, so whoever operates the
relay you use can see it: the chat view of a pane and what is typed into it,
the state of a pane (what its agent has spent, its usage limits, and which
model routing chose for it), the latest replies, photos attached from a
phone, diffs and commit, push and pull request data, and an API key typed
into a dialog.

End-to-end encryption defeats an honestly-run relay. It doesn't defend
against a relay that's been actively compromised and tampered with to swap
the keys it hands out at pairing. The **Verify** code beside each device in
the desktop's Remote access dialog is the check for that: compare it with
the code the device shows on its own Devices page, and if they differ,
unpair the device. Running your own relay narrows who mounting that attack
would require compromising to your own infrastructure and the people who
administer it, instead of Flockdeck's.

Whoever controls a relay decides which devices are paired and serves the
web page browsers run. A compromised or malicious relay can pair a device of
its own and use any desktop connected to it, and end-to-end encryption does
not prevent that. Run only a relay you trust with the desktops on it.

## What the relay routes, and what it stores

The relay proxies a desktop's Flockdeck window — every pane, every agent
conversation, any photo attached from a phone — to a paired device, and
back. None of that is written to the relay's storage; [Storage and
data](storage-and-data.html) lists what is kept, which is account and device
records and never anything that passes between them.

## Push notifications are the exception

A push notification's content is encrypted on the desktop, for the specific
device it's going to, before it reaches the relay — the relay delivers it
without being able to read it. This is the one thing a self-hosted relay
handles no differently from the shared one: neither can see what a
notification says.

## What reaches the relay from a desktop, and what doesn't

Turning on remote access is the only thing that makes a desktop talk to a
relay at all. With it off, nothing about a desktop's agents, panes or
projects is sent anywhere related to remote access. With it on, the
desktop's outbound connection to the relay carries exactly what's described
above — proxied window content for paired devices, and nothing about
projects, agents or conversations that aren't being viewed remotely.

## The audit log

`-audit` (see [Configuration](configuration.html)) turns on a log of who
reached which desktop, from which device, and when — never what was typed
or shown, since that's the terminal traffic [above](#end-to-end-encrypted-terminal-traffic-with-one-caveat)
this relay can't read either way. `-audit-ip` adds the address each event
came from, and `-audit-retention` bounds how long events are kept; both are
refused without `-audit` itself on, and the relay won't start.

## Operating a relay securely

Treat it like any other service that terminates or proxies encrypted
traffic for your organisation: TLS between every hop (see
[Requirements](requirements.html)), the admin API kept off the public
network — `-admin-addr` defaults to loopback already; set it to empty to
turn the admin endpoints off entirely (see
[Admin and invites](admin-and-invites.html)) — the data directory or
database backed up and access-controlled, and the relay itself kept current
(see [Upgrading](upgrading.html)).
