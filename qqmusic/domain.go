package qqmusic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the QQ Music driver for the any-cli/kit framework.
type Domain struct{}

// Info describes the scheme and hostnames matched against pasted URLs.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme:   "qqmusic",
		Hosts:    []string{"y.qq.com", "c.y.qq.com", "u.y.qq.com"},
		Identity: BaseIdentity(),
	}
}

// BaseIdentity is the help and version identity.
func BaseIdentity() kit.Identity {
	return kit.Identity{
		Binary: "qqmusic",
		Short:  "A command line for QQ Music (QQ音乐).",
		Long: `qqmusic reads QQ Music (QQ音乐) data and prints clean, pipeable records.

Browse charts, search songs, fetch lyrics, explore playlists, artists, and albums.
All metadata commands work without an account.

Records come out as table, list, markdown, JSON, JSONL, CSV, TSV, url, or raw.

qqmusic is an independent tool and is not affiliated with Tencent or QQ Music.`,
		Site: "https://y.qq.com",
		Repo: "https://github.com/tamnd/qqmusic-cli",
	}
}

// Defaults seeds the framework baseline from the qqmusic defaults.
func Defaults(c *kit.Config) {
	d := defaultConfig()
	c.Rate = d.Rate
	c.Timeout = d.Timeout
	c.Retries = d.Retries
	c.UserAgent = d.UserAgent
}

// DefaultConfig returns a Config with sensible defaults for the QQ Music API.
func DefaultConfig() Config {
	return Config{
		BaseURL:   defaultBaseURL,
		MusicuURL: defaultMusicuURL,
		UserAgent: defaultUA,
		Rate:      500 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
	}
}

func defaultConfig() Config { return DefaultConfig() }

// Register installs the client factory and all operations onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)
	registerOps(app)
}

// Register is a convenience so callers don't need to name the zero-value Domain.
func Register(app *kit.App) { Domain{}.Register(app) }

// Session is the per-run client kit injects into every operation.
type Session struct {
	Client *Client
	Quiet  bool
}

// Progressf prints a one-line progress note to stderr unless quiet.
func (s *Session) Progressf(format string, args ...any) {
	if s == nil || s.Quiet {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func newClient(_ context.Context, c kit.Config) (any, error) {
	cfg := defaultConfig()
	if c.UserAgent != "" {
		cfg.UserAgent = c.UserAgent
	}
	if c.Rate > 0 {
		cfg.Rate = c.Rate
	}
	if c.Timeout > 0 {
		cfg.Timeout = c.Timeout
	}
	if c.Retries > 0 {
		cfg.Retries = c.Retries
	}
	return &Session{Client: NewClient(cfg), Quiet: c.Quiet}, nil
}

// MapErr converts a library error into the kit error kind with the right exit code.
func MapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrNotFound):
		return errs.NotFound("%s", err.Error())
	default:
		return err
	}
}

// Classify turns an input URL into the canonical (uriType, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty qqmusic reference")
	}
	if strings.Contains(input, "/songDetail/") {
		if mid, e := parseSongMid(input); e == nil {
			return "song", mid, nil
		}
	}
	if strings.Contains(input, "/singer/") {
		if mid, e := parseSingerMid(input); e == nil {
			return "artist", mid, nil
		}
	}
	if strings.Contains(input, "/album/") {
		if mid, e := parseAlbumMid(input); e == nil {
			return "album", mid, nil
		}
	}
	if strings.Contains(input, "/playlist/") {
		if pid, e := parsePlaylistID(input); e == nil {
			return "playlist", pid, nil
		}
	}
	return "", "", errs.Usage("unrecognised qqmusic reference: %q", input)
}

// Locate returns the live https URL for a (uriType, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "song":
		return "https://y.qq.com/n/ryqq/songDetail/" + id, nil
	case "artist":
		return "https://y.qq.com/n/ryqq/singer/" + id, nil
	case "album":
		return "https://y.qq.com/n/ryqq/album/" + id, nil
	case "playlist":
		return "https://y.qq.com/n/ryqq/playlist/" + id, nil
	default:
		return "", errs.Usage("qqmusic has no resource type %q", uriType)
	}
}
