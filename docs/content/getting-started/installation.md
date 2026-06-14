---
title: "Installation"
description: "Install qqmusic from a release, with go install, or from source."
weight: 20
---

## Prebuilt binaries

Every [release](https://github.com/tamnd/qqmusic-cli-cli/releases) carries archives for Linux, macOS,
and Windows on amd64 and arm64, plus deb, rpm, and apk packages for Linux.
Download, unpack, put `qqmusic` on your `PATH`, done. The `checksums.txt`
on each release is signed with keyless [cosign](https://docs.sigstore.dev/) if
you want to verify before running.

## With Go

```bash
go install github.com/tamnd/qqmusic-cli-cli/cmd/qqmusic@latest
```

That puts `qqmusic` in `$(go env GOPATH)/bin`, which is `~/go/bin` unless
you moved it. Make sure that directory is on your `PATH`.

## From source

```bash
git clone https://github.com/tamnd/qqmusic-cli-cli
cd qqmusic-cli-cli
make build        # produces ./bin/qqmusic
./bin/qqmusic version
```

## Container image

```bash
docker run --rm ghcr.io/tamnd/qqmusic:latest --help
```

## Checking the install

```bash
qqmusic version
```

prints the version and exits.
