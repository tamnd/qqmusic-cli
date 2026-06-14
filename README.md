# qqmusic

Browse QQ Music charts and songs (QQ音乐)

`qqmusic` is a single pure-Go binary. It reads QQ Music chart data through
the public toplist API, shapes the responses into clean records, and pipes into
the rest of your tools. No API key, nothing to run alongside it.

## Install

```bash
go install github.com/tamnd/qqmusic-cli/cmd/qqmusic@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/qqmusic-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/qqmusic:latest --help
```

## Usage

```bash
# List songs from the hot chart (热歌榜)
qqmusic top

# List songs from the new song chart (新歌榜)
qqmusic top --chart new

# List songs from the rising chart (飙升榜) as JSON
qqmusic top --chart rising -o json

# List top 10 songs
qqmusic top -n 10

# Show available charts
qqmusic charts

# Output as JSONL for piping
qqmusic top -o jsonl | jq .title
```

## Development

```
cmd/qqmusic/   thin main, wires cli.Root into fang
cli/           the cobra command tree
qqmusic/       the library: HTTP client and data models
pkg/render/    output rendering (table, json, jsonl, csv, tsv, url)
docs/          tago documentation site
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
