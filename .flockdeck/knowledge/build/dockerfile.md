---
type: Build Step
title: Dockerfile
description: Builds the site image by copying the repository into an nginx base image; nothing is compiled.
resource: Dockerfile
tags: [docker, nginx]
timestamp: 2026-10-19T00:00:00Z
code_paths: [Dockerfile, .dockerignore, nginx.conf]
---
```
FROM nginx:1.31-alpine
COPY . /usr/share/nginx/html
RUN mv /usr/share/nginx/html/nginx.conf /etc/nginx/conf.d/default.conf
```

- Base is pinned to nginx's minor version, so rebuilds take patch releases only. Moving to the next minor is a reviewed edit here. The file says it is kept level with the sibling site's pin.
- One `COPY` of the whole directory. What is not site is excluded by [.dockerignore](/build/dockerignore.md), so a new page needs no Dockerfile change.
- `nginx.conf` arrives with the site and is moved to `conf.d/default.conf` so it is not served as a page. See [nginx config](/build/nginx-conf.md).
- Changes to `Dockerfile` count as shipping for [autotag.sh](/cli/autotag-sh.md) even though the file is dockerignored.

Built and checked by the `image` job of [CI](/ci/ci.md). Tags: [image tags](/build/image-tags.md).
