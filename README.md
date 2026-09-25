# flockdeck-docs

The documentation site for [Flockdeck](https://github.com/jmwri/flockdeck),
served at [docs.flockdeck.ai](https://docs.flockdeck.ai).

It has two parts, built differently on purpose.

## Using Flockdeck (`app/`)

Generated, not written here. `cmd/docgen` in the Flockdeck repository renders
`internal/help`'s own Markdown — the same pages, the same key-table
placeholders, that render behind `F1` in the app — into this section, so a
web copy can't say anything the app itself doesn't. Don't edit anything
under `app/` by hand; the next generation overwrites it. Edit the pages at
their source, `internal/help/pages/*.md` in
[flockdeck](https://github.com/jmwri/flockdeck), instead:

```sh
# from a checkout of github.com/jmwri/flockdeck, beside this one
go run ./cmd/docgen -out ../flockdeck-docs
```

## Self-hosting (`self-hosting/`, from `content/self-hosting/`)

There's no source of truth inside the app to render this from — nothing
about self-hosting a relay lives behind `F1`, since it isn't about the app.
So it's ordinary Markdown, written by hand and kept right here, in
`content/self-hosting/*.md`. The same `docgen` run above reads it and
writes `self-hosting/`; edit the Markdown, not the generated HTML.

The order pages appear in is `selfHostOrder` in `docgen`'s `main.go` — a
page has to be listed there to be picked up, and a file present but
unlisted fails the build rather than going unlinked from the nav.

## Everything else

`index.html`, `docs.css`, the fonts and the icons are generated too, from
`cmd/docgen`'s own embedded assets — the same palette, type and prose rules
as [flockdeck-site](https://github.com/jmwri/flockdeck-site), copied rather
than linked across origins, since the two sites are deployed separately.

## Keeping the docs from going stale

The generated pages come from a specific flockdeck release, recorded in
`.docs-generated-from` (`<tag> <commit>`). Regenerate, and update that stamp,
in one step:

```sh
tools/docdrift/generated.sh regenerate          # from the latest release
tools/docdrift/generated.sh regenerate v0.3.41  # or a named one
```

That needs `go` and `gh`; it clones flockdeck itself. Four guards (the
comparison logic is `tools/docdrift`, a small stdlib-only Go program with its
own tests) keep the docs honest. Each is its own job, so a failure says which:

| Guard | When | What it catches |
| --- | --- | --- |
| **Generated files match** (`docs-guards.yml`) | every PR | A hand edit of `app/*.html` or `self-hosting/*.html`, or `content/` changed without regenerating. Compares against the *stamped* flockdeck, never the newest release, so a flockdeck release cannot turn an unrelated PR red. |
| **Relay flags are documented** | every PR | The relay's `serve -h`, as checked in under `tools/docdrift/snapshot` (see below), against `content/self-hosting/configuration.md`: every flag documented or in `tools/docdrift/allow-relay-flags.txt` with a reason; every documented flag still exists; stated defaults and env vars agree. Fails only on findings the PR *adds* (the base branch's page is checked the same way and diffed). Needs no credential. |
| **Desktop commands are documented** | every PR | The released desktop's `flockdeck -h` subcommands against `app/cli.html`; same base-branch rule, allowlist in `allow-desktop-commands.txt`. |
| **Docs drift** (`docs-drift.yml`) | daily, and on demand | The site behind the latest flockdeck release, or the config/CLI pages disagreeing with the latest relay/desktop: keeps one tracking issue per check, closes it when clean, and does not comment when nothing changed. |

The desktop's commands are read from the linux amd64 release tarball,
verified against `checksums.txt`. None of the guards is a required check by
default: add the job names above to branch protection to make them gating
(leave `Check the site` as it is).

### The relay snapshot

The released relay image, `ghcr.io/flockdeck/flockdeck-relay`, is private on
purpose (the built relay comes with an Enterprise licence), and this
repository is public, so a pull request cannot and must not read it. Pull
requests are checked against `tools/docdrift/snapshot/relay-serve-help.txt`,
the relay's `serve -h` at the version in `snapshot/relay.version`.

Something has to keep that snapshot honest, and it is the `relay` job of
*Docs drift* (daily, and on demand). It is the only thing that reads the image:

- **It is dormant until configured.** With no repository variable
  `RELAY_GHCR_READ_USER` (or an empty one) the job is skipped: no run goes red
  and no issue opens. Until then nothing refreshes the snapshot; pull requests
  still work against the checked-in one. To switch it on:
  1. Create a dedicated GitHub account (a bot) and give it **Read** on the
     private `flockdeck-relay` package. The image stays **private**; this
     grants read access to one account, it does not publish anything.
  2. On that account, create a *classic* personal access token with the
     **`read:packages`** scope and nothing else (fine-grained tokens cannot
     read packages).
  3. Settings → Environments → new **`relay-release`**, with **Deployment
     branches: `main` only**, so a workflow on any other branch cannot read
     the credential. Add the token as the environment secret
     **`RELAY_GHCR_READ_TOKEN`**.
  4. Settings → Secrets and variables → Actions → Variables: add the
     repository variable **`RELAY_GHCR_READ_USER`** with that account's
     username. (Setting this last is what turns the job on.)

  The job also only runs in `Flockdeck/flockdeck-docs`, never in a fork.
- If the credential later stops working (expired, revoked, access lost), the
  job opens one issue, *Relay drift check cannot read the relay image*, ends
  green, and closes the issue itself on the next run that can read the image.
  Any other failure of the job is a bug in the tool and stays red. It never
  falls back to an older image.
- If the snapshot is not the latest release's, it opens or updates one issue,
  *Relay flag snapshot is behind flockdeck-relay vX*, with the flag diff, and
  closes it once they match. It also opens *Docs disagree with flockdeck-relay
  vX* if `configuration.md` disagrees with the latest release itself.

To refresh the snapshot by hand (with `RELAY_GHCR_READ_USER` and `RELAY_GHCR_READ_TOKEN` set in your shell):

```sh
tools/docdrift/refresh-relay-snapshot.sh            # or ... 0.2.51
```

then fix or allowlist what `docdrift relay` reports and commit the snapshot
with the docs. Until it is refreshed, pull requests cannot see a relay change
the snapshot doesn't have, which is why the daily job exists.

An entry in an allowlist is a decision: it needs a one-line reason, and the
check fails when the entry outlives it (the docs now cover it, or the flag is
gone). To run a guard by hand, see the top of each script in `tools/docdrift/`.

## How it is served

Deployed the same way flockdeck-site is, to the wost Kubernetes cluster:

1. CI builds an image from these files on every push. A tag like `v0.1.0`
   publishes `ghcr.io/jmwri/flockdeck-docs:0.1.0` — the bare semver matters,
   because that is what Flux's ImagePolicy matches.
2. Flux, watching this repository from `k8s-infra/flockdeck-docs`, sees the
   new tag, writes it into `deployment.yaml` and reconciles the cluster.
3. `terrawost/svc/flockdeck-docs` owns everything around it: the
   `docs.flockdeck.ai` DNS record, the namespace, and the registry pull
   secrets.

So a release is: regenerate, commit, tag.

## Licence

PolyForm Noncommercial — see [LICENSE](LICENSE), the same licence as the
Flockdeck desktop app this documents. The fonts (Archivo, JetBrains Mono) keep
their own licence, the SIL Open Font License, beside them in `fonts/`.
