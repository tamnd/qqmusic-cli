package qqmusic

import (
	"context"
	"strings"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func registerOps(app *kit.App) {
	registerCharts(app)
	registerChart(app)
	registerSearch(app)
	registerSong(app)
	registerLyrics(app)
	registerArtist(app)
	registerAlbum(app)
	registerPlaylist(app)
	registerPlaylists(app)
}

func effectiveLimit(n, def int) int {
	if n > 0 {
		return n
	}
	return def
}

// --- charts ---

type chartsIn struct {
	Sess *Session `kit:"inject"`
}

func registerCharts(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "charts",
		Group:   "read",
		List:    true,
		Summary: "List all available QQ Music charts (排行榜)",
	}, func(ctx context.Context, in chartsIn, emit func(ChartInfo) error) error {
		in.Sess.Progressf("fetching chart list")
		items, err := in.Sess.Client.Charts(ctx)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- chart ---

type chartIn struct {
	Chart string   `kit:"arg" help:"chart name (hot/new/rising) or numeric topid"`
	Limit int      `kit:"flag,inherit" help:"max records"`
	Sess  *Session `kit:"inject"`
}

func registerChart(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "chart",
		Group:   "read",
		List:    true,
		Summary: "Songs on a QQ Music chart",
		Args:    []kit.Arg{{Name: "chart", Help: "chart name (hot/new/rising) or topid; default: hot"}},
	}, func(ctx context.Context, in chartIn, emit func(Song) error) error {
		topid, err := chartTopID(in.Chart)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		limit := effectiveLimit(in.Limit, 30)
		in.Sess.Progressf("fetching chart (topid=%d)", topid)
		items, err := in.Sess.Client.Chart(ctx, topid, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- search ---

type searchIn struct {
	Query  []string `kit:"arg,variadic" help:"search terms"`
	Limit  int      `kit:"flag,inherit" help:"max records"`
	Page   int      `kit:"flag" help:"page number (1-based)"`
	Sess   *Session `kit:"inject"`
}

func registerSearch(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "search",
		Group:   "read",
		List:    true,
		Summary: "Search QQ Music for songs",
		Args:    []kit.Arg{{Name: "query", Help: "search terms", Variadic: true}},
	}, func(ctx context.Context, in searchIn, emit func(Song) error) error {
		if len(in.Query) == 0 {
			return errs.Usage("search requires a query")
		}
		q := strings.Join(in.Query, " ")
		limit := effectiveLimit(in.Limit, 20)
		page := in.Page
		if page <= 0 {
			page = 1
		}
		in.Sess.Progressf("searching for %q", q)
		items, err := in.Sess.Client.Search(ctx, q, limit, page)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- song ---

type songIn struct {
	ID   string   `kit:"arg" help:"songmid or song URL"`
	Sess *Session `kit:"inject"`
}

func registerSong(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "song",
		Group:   "read",
		Single:  true,
		Summary: "QQ Music song detail",
		Args:    []kit.Arg{{Name: "mid-or-url", Help: "songmid or song URL"}},
	}, func(ctx context.Context, in songIn, emit func(Song) error) error {
		mid, err := parseSongMid(in.ID)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		in.Sess.Progressf("fetching song %s", mid)
		s, err := in.Sess.Client.Song(ctx, mid)
		if err != nil {
			return MapErr(err)
		}
		return emit(s)
	})
}

// --- lyrics ---

type lyricsIn struct {
	ID   string   `kit:"arg" help:"songmid or song URL"`
	Sess *Session `kit:"inject"`
}

func registerLyrics(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "lyrics",
		Group:   "read",
		Single:  true,
		Summary: "Song lyrics in LRC format",
		Args:    []kit.Arg{{Name: "mid-or-url", Help: "songmid or song URL"}},
	}, func(ctx context.Context, in lyricsIn, emit func(Lyric) error) error {
		mid, err := parseSongMid(in.ID)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		in.Sess.Progressf("fetching lyrics for %s", mid)
		l, err := in.Sess.Client.Lyrics(ctx, mid)
		if err != nil {
			return MapErr(err)
		}
		return emit(l)
	})
}

// --- artist ---

type artistIn struct {
	ID    string   `kit:"arg" help:"singermid or singer URL"`
	Limit int      `kit:"flag,inherit" help:"max records"`
	Sess  *Session `kit:"inject"`
}

func registerArtist(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "artist",
		Group:   "read",
		List:    true,
		Summary: "Artist top tracks",
		Args:    []kit.Arg{{Name: "mid-or-url", Help: "singermid or singer URL"}},
	}, func(ctx context.Context, in artistIn, emit func(Song) error) error {
		mid, err := parseSingerMid(in.ID)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		limit := effectiveLimit(in.Limit, 30)
		in.Sess.Progressf("fetching artist tracks for %s", mid)
		items, err := in.Sess.Client.Artist(ctx, mid, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- album ---

type albumIn struct {
	ID   string   `kit:"arg" help:"albummid or album URL"`
	Sess *Session `kit:"inject"`
}

func registerAlbum(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "album",
		Group:   "read",
		List:    true,
		Summary: "Album track list",
		Args:    []kit.Arg{{Name: "mid-or-url", Help: "albummid or album URL"}},
	}, func(ctx context.Context, in albumIn, emit func(Song) error) error {
		mid, err := parseAlbumMid(in.ID)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		in.Sess.Progressf("fetching album %s", mid)
		items, err := in.Sess.Client.Album(ctx, mid)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- playlist ---

type playlistIn struct {
	ID   string   `kit:"arg" help:"playlist id or URL"`
	Sess *Session `kit:"inject"`
}

func registerPlaylist(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "playlist",
		Group:   "read",
		List:    true,
		Summary: "Songs in a playlist",
		Args:    []kit.Arg{{Name: "id-or-url", Help: "playlist id or URL"}},
	}, func(ctx context.Context, in playlistIn, emit func(Song) error) error {
		id, err := parsePlaylistID(in.ID)
		if err != nil {
			return errs.Usage("%s", err.Error())
		}
		in.Sess.Progressf("fetching playlist %s", id)
		items, err := in.Sess.Client.Playlist(ctx, id)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

// --- playlists ---

type playlistsIn struct {
	Category int      `kit:"flag" help:"category id (default 10000000 = all)"`
	Limit    int      `kit:"flag,inherit" help:"max records"`
	Sess     *Session `kit:"inject"`
}

func registerPlaylists(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "playlists",
		Group:   "read",
		List:    true,
		Summary: "Discover QQ Music playlists by category",
	}, func(ctx context.Context, in playlistsIn, emit func(PlaylistInfo) error) error {
		limit := effectiveLimit(in.Limit, 30)
		in.Sess.Progressf("fetching playlists (category=%d)", in.Category)
		items, err := in.Sess.Client.Playlists(ctx, in.Category, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(items, emit)
	})
}

func emitAll[T any](items []T, emit func(T) error) error {
	for _, it := range items {
		if err := emit(it); err != nil {
			return err
		}
	}
	return nil
}
