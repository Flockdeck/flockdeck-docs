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

The state of a pane still passes through the relay in the clear — what each
pane's agent has spent, its usage limits, and which model routing chose for
it — so whoever operates the relay you use can see that.

End-to-end encryption defeats an honestly-run relay. It doesn't yet defend
against a relay that's been actively compromised and tampered with to swap
the keys it hands out at pairing — that needs an out-of-band check that
hasn't been built yet. Running your own relay narrows who mounting that
attack would require compromising to your own infrastructure and the people
who administer it, instead of Flockdeck's — which is one of the reasons
self-hosting exists, alongside not sending that metadata to a third party at
all.

## What the relay routes, and what it stores

The relay proxies a desktop's Flockdeck window — every pane, every agent
conversation, any photo attached from a phone — to a paired device, and
back. None of that is written to the relay's storage; [Storage and
data](storage-and-data.html) is only which desktops and devices exist and
are paired, not anything that passes between them.

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

## Operating a relay securely

Treat it like any other service that terminates or proxies encrypted
traffic for your organisation: TLS between every hop (see
[Requirements](requirements.html)), the admin API kept off the public
network (see [Admin and invites](admin-and-invites.html)), the data
directory or database backed up and access-controlled, and the relay itself
kept current (see [Upgrading](upgrading.html)).
