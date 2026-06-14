# qqmusic

Browse QQ Music charts and songs (QQ音乐)

`qqmusic` is a single pure-Go binary. It speaks to qqmusic-cli over plain
HTTPS, shapes the responses into clean records, and pipes into the rest of your
tools. No API key, nothing to run alongside it.

## Install

```bash
go install github.com/tamnd/qqmusic-cli-cli/cmd/qqmusic@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/qqmusic-cli-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/qqmusic:latest --help
```

## Usage

```bash
qqmusic --help
qqmusic version
```

This is a fresh scaffold. The command tree starts with `version`; build out the
real commands in `cli/` on top of the `qqmusic-cli` library package.

## Development

```
cmd/qqmusic/   thin main, wires cli.Root into fang
cli/                 the cobra command tree
qqmusic-cli/                the library: HTTP client and data models
docs/                tago documentation site
```

```bash
make build      # ./bin/qqmusic
make test       # go test ./...
make vet        # go vet ./...
```

## Releasing

Push a version tag and GitHub Actions runs GoReleaser, which builds the
archives, Linux packages, the multi-arch GHCR image, checksums, SBOMs, and a
cosign signature:

```bash
git tag v0.1.0
git push --tags
```

The Homebrew and Scoop steps self-disable until their tokens exist, so the first
release works with no extra secrets.

## License

Apache-2.0. See [LICENSE](LICENSE).
