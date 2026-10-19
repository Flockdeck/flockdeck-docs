# Build and release steps

* [Dockerfile](dockerfile.md) - Builds the site image by copying the repository into an nginx base image; nothing is compiled.
* [.dockerignore](dockerignore.md) - Decides which repository files stay out of the image, and doubles as the rule autotag uses for "ships".
* [Line-ending policy](gitattributes.md) - .gitattributes forces LF for every text file and marks known binary types.
* [Root go.mod retraction](go-mod-retraction.md) - A go.mod with no code that exists only to retract two MIT-licensed versions.
* [Image tags](image-tags.md) - Which container tags CI produces for pull requests, main, and version tags.
* [nginx.conf](nginx-conf.md) - The site's single nginx server block: compression, security headers, no cache, no access log.
* [Page generation by docgen](regenerate-pages.md) - How app/ and self-hosting/ are produced from the flockdeck repository's cmd/docgen and this repo's content directory.
* [Release by tag](release-by-tag.md) - End-to-end path from a merge to main to a deployed site, which only happens through a v* tag.
