# qqmusic

Browse QQ Music charts, search songs, fetch lyrics, and explore playlists (QQ音乐)

`qqmusic` is a single pure-Go binary. It reads QQ Music public data without an API
key or account, shapes the responses into clean records, and pipes into the rest of
your tools.

qqmusic is not affiliated with Tencent or QQ Music.

## Install

```bash
go install github.com/tamnd/qqmusic-cli/cmd/qqmusic@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/qqmusic-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/qqmusic:latest --help
```

## Commands

```
charts                  list all available charts
chart <name|id>         songs from a chart  (hot / new / rising / <topid>)
search <query...>       search songs
song <songmid|url>      details for one song
lyrics <songmid|url>    LRC lyrics for a song
artist <singermid|url>  top tracks for an artist
album <albummid|url>    track list for an album
playlist <id|url>       songs in a playlist
playlists               discover playlists by category
```

## Usage

```bash
# List all charts
qqmusic charts

# Hot chart top 10
qqmusic chart hot -n 10

# New songs chart
qqmusic chart new

# Rising chart as JSON
qqmusic chart rising -o json

# Search
qqmusic search 周杰伦 -n 5

# Song detail
qqmusic song 0030KTZ8141fsb

# Lyrics
qqmusic lyrics 0030KTZ8141fsb

# Artist top tracks (use singermid or y.qq.com artist URL)
qqmusic artist 0025NhlN2yWrP4 -n 10

# Album tracks
qqmusic album 002M1v7j4B5Yf7

# Playlist songs
qqmusic playlist 7707261125

# Discover playlists
qqmusic playlists -n 20

# Pipe JSONL
qqmusic chart hot -o jsonl | jq .title
```

Every command accepts `-o table|list|json|jsonl|csv|tsv|url|markdown` (default: jsonl).
Use `-n <N>` to limit results.

## Development

```
cmd/qqmusic/   thin main
cli/           app wiring (any-cli/kit)
qqmusic/       library: HTTP client, types, API methods
pkg/render/    output rendering
docs/          documentation site
```

```bash
go build ./...
go test ./...
go vet ./...
```

## Releasing

Push a version tag and GitHub Actions runs GoReleaser:

```bash
git tag v0.1.0
git push --tags
```

## License

Apache-2.0. See [LICENSE](LICENSE).
