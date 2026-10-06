# Security and privacy

What a relay (yours or the shared one) can and can't see, so you can
weigh that against what your own rules require.

## End-to-end encrypted terminal traffic, with one caveat

Traffic between a desktop and a paired device is encrypted with TLS to the
relay and from the relay onward. On top of that, terminal traffic (what you
type and what an agent prints back) is end-to-end encrypted by the desktop
app and the paired device themselves, using keys the relay hands out but
never holds. The relay has to route it, to know which desktop a browser's
request is for, but it cannot read it unless it actively tampers with the client
it serves (below). This is true of the shared relay at
`remote.flockdeck.ai` and of a relay you run yourself, in exactly the same
way.

Only the terminal is end-to-end encrypted, and only when the desktop and the
paired device have each registered a key. A terminal opened when either has
none is served without it, in plaintext, and a relay that holds back a device's
key can cause that: the manual check below cannot show that a key was withheld.

Everything else passes through the relay decrypted, so whoever operates the
relay you use can see it: the chat view of a pane and what is typed into it,
the state of a pane (what its agent has spent, its usage limits, and which
model routing chose for it), the latest replies, photos attached from a
phone, diffs and commit, push and pull request data, and an API key typed
into a dialog.

End-to-end encryption defeats an honestly-run relay. It doesn't defend
against a relay that's been actively compromised and tampered with to swap
the keys it hands out at pairing. The Verify code beside each device in
the desktop's Flockdeck Remote dialog is the check for that: compare it with
the code the device shows on its own Devices page, and if they differ,
unpair the device. Running your own relay narrows who mounting that attack
would require compromising to your own infrastructure and the people who
administer it, instead of Flockdeck's.

The relay also serves the web client's JavaScript itself: the files are built
into the relay binary (or read from `-client-dir`) and served from its own
origin, under a Content Security Policy that allows only its own scripts. With
`-desk-domain` it also proxies each desktop's own pages on that desktop's
origin. A relay that actively tampers, or an attacker who controls it, can
serve altered client code or altered desktop pages, which can read what the
browser types and shows, and can report matching codes in Verify. So the
guarantee is that the relay can read nothing that passes through it unless it
actively tampers with the client it serves: the web client's JavaScript, and the
desktop pages it proxies on a desk origin. Neither the pairing handshake nor the
manual check can detect that. The remedy is to run your own relay and trust its
operator.

Whoever controls a relay decides which devices are paired and serves the
web page browsers run. A compromised or malicious relay can pair a device of
its own and use any desktop connected to it, and end-to-end encryption does
not prevent that. Run only a relay you trust with the desktops on it.

## Pages and what they may run

The web client's pages take scripts only from the relay (`script-src 'self'`,
no inline script) and open sockets only to it. The two small pages that run a
verified registration (`/auth/verify` and `/verify`) are the exception: their
script is inline, so their policy allows `'unsafe-inline'` for scripts, and
still allows no third-party script, style or connection.

## What the relay routes, and what it stores

The relay proxies a desktop's Flockdeck window (every pane, every agent
conversation, any photo attached from a phone) to a paired device, and
back. None of that is written to the relay's storage; [Storage and
data](storage-and-data.html) lists what is kept, which is account and device
records and never anything that passes between them. A desktop that registered
and never connected is removed after 7 days, and only the desktop: its account
stays. A desktop that did connect is never removed for being quiet. An account
nothing has used for `-account-retention` (90 days by default) is deleted with
its desktops and devices; see [How long accounts are kept](storage-and-data.html#how-long-accounts-are-kept).

## Push notifications are the exception

A push notification's content is encrypted on the desktop, for the specific
device it's going to, before it reaches the relay. The relay delivers it
without being able to read it, and the desktop pads every message to one
length, so even its size tells the push service nothing. This is the one thing a self-hosted relay
handles no differently from the shared one: neither can see what a
notification says.

## What reaches the relay from a desktop, and what doesn't

Turning on remote access is the only thing that makes a desktop talk to a
relay at all. With it off, nothing about a desktop's agents, panes or
projects is sent anywhere related to remote access. With it on, the
desktop's outbound connection to the relay carries exactly what's described
above: proxied window content for paired devices, and nothing about
projects, agents or conversations that aren't being viewed remotely.

## Plans and billing

A relay you run for yourself has no plans. A relay run with `-enforce-plans`
(on by default once `-billing-url` is set, as the shared relay is) holds an
account's plan and when its trial or paid time ends, and nothing about who
pays: the billing service keeps that. When a paired browser asks to subscribe,
the relay makes a link that names the account, is signed with a key the relay
and the billing service share, and works for ten minutes, once; the browser
carries it to the billing service, which opens Paddle. The relay itself calls
nothing: the billing service calls it, to set the plan.

## The audit log

`-audit` (see [Configuration](configuration.html)) turns on a log of who
reached which desktop, from which device, and when. It never records what was typed
or shown, since that's the terminal traffic [above](#end-to-end-encrypted-terminal-traffic-with-one-caveat)
this relay can't read either way. `-audit-ip` adds the address each event
came from, and `-audit-retention` bounds how long events are kept; both are
refused without `-audit` itself on, and the relay won't start.

## Operating a relay securely

Treat it like any other service that terminates or proxies encrypted
traffic for your organisation: TLS between every hop (see
[Requirements](requirements.html)), the admin API kept off the public
network (`-admin-addr` defaults to loopback already; set it to empty to
turn the admin endpoints off entirely, see
[Admin and invites](admin-and-invites.html)), the data directory or
database backed up and access-controlled, and the relay itself kept current
(see [Upgrading](upgrading.html)).
