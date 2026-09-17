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
