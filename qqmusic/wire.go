package qqmusic

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// joinArtists joins singer names into a comma-separated string.
func joinArtists(singers []wireSinger) string {
	names := make([]string, 0, len(singers))
	for _, s := range singers {
		if s.Name != "" {
			names = append(names, s.Name)
		}
	}
	return strings.Join(names, ", ")
}

// fmtDuration converts seconds to "M:SS".
func fmtDuration(secs int) string {
	if secs <= 0 {
		return ""
	}
	return fmt.Sprintf("%d:%02d", secs/60, secs%60)
}

// songURL builds the canonical song page URL from a songmid.
func songURL(mid string) string {
	return "https://y.qq.com/n/ryqq/songDetail/" + mid
}

// tryBase64 decodes s if it looks like base64, otherwise returns s as-is.
// QQ Music lyrics endpoints sometimes return base64-encoded LRC even with nobase64=1.
func tryBase64(s string) string {
	if s == "" {
		return ""
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}
	return string(b)
}

func chartSongFrom(w wireChartSong, rank int) Song {
	d := w.Data
	return Song{
		Rank:     rank,
		SongMid:  d.SongMid,
		Title:    d.SongName,
		Artists:  joinArtists(d.Singer),
		Album:    d.AlbumName,
		AlbumMid: d.AlbumMid,
		Duration: fmtDuration(d.Interval),
		Date:     d.TimePublic,
		URL:      songURL(d.SongMid),
	}
}

func searchSongFrom(w wireSearchSong, rank int) Song {
	date := w.TimePublic
	if date == "" && w.PubTime > 0 {
		date = time.Unix(w.PubTime, 0).UTC().Format("2006-01-02")
	}
	return Song{
		Rank:     rank,
		SongMid:  w.SongMid,
		Title:    w.SongName,
		Artists:  joinArtists(w.Singer),
		Album:    w.AlbumName,
		AlbumMid: w.AlbumMid,
		Duration: fmtDuration(w.Interval),
		Date:     date,
		URL:      songURL(w.SongMid),
	}
}

func songDetailFrom(w wireSongDetail) Song {
	return Song{
		Rank:     1,
		SongMid:  w.Mid,
		Title:    w.Name,
		Artists:  joinArtists(w.Singer),
		Album:    w.Album.Name,
		AlbumMid: w.Album.Mid,
		Duration: fmtDuration(w.Interval),
		Date:     w.TimePublic,
		URL:      songURL(w.Mid),
	}
}

func artistSongFrom(w wireArtistMusicData, rank int) Song {
	return Song{
		Rank:     rank,
		SongMid:  w.SongMid,
		Title:    w.SongName,
		Artists:  joinArtists(w.Singer),
		Album:    w.AlbumName,
		AlbumMid: w.AlbumMid,
		Duration: fmtDuration(w.Interval),
		Date:     w.TimePublic,
		URL:      songURL(w.SongMid),
	}
}

func artistSongFromMusicu(w wireMusicuSongInfo, rank int) Song {
	return Song{
		Rank:     rank,
		SongMid:  w.Mid,
		Title:    w.Name,
		Artists:  joinArtists(w.Singer),
		Album:    w.Album.Name,
		AlbumMid: w.Album.Mid,
		Duration: fmtDuration(w.Interval),
		Date:     w.TimePublic,
		URL:      songURL(w.Mid),
	}
}

func albumTrackFrom(w wireAlbumTrack, albumName string, rank int) Song {
	artists := joinArtists(w.Singer)
	return Song{
		Rank:     rank,
		SongMid:  w.SongMid,
		Title:    w.SongName,
		Artists:  artists,
		Album:    albumName,
		Duration: fmtDuration(w.Interval),
		URL:      songURL(w.SongMid),
	}
}

func playlistSongFrom(w wirePlaylistSong, rank int) Song {
	return Song{
		Rank:     rank,
		SongMid:  w.SongMid,
		Title:    w.SongName,
		Artists:  joinArtists(w.Singer),
		Album:    w.AlbumName,
		AlbumMid: w.AlbumMid,
		Duration: fmtDuration(w.Interval),
		URL:      songURL(w.SongMid),
	}
}

// parseDissID converts the API's dissid (int or string) to int64.
func parseDissID(v interface{}) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	}
	return 0
}

func playlistInfoFrom(w wirePlaylistItem, rank int) PlaylistInfo {
	id := parseDissID(w.DissID)
	cover := w.ImgURL
	if cover == "" {
		cover = w.Logo
	}
	return PlaylistInfo{
		Rank:      rank,
		ID:        id,
		Name:      w.DissName,
		Creator:   w.Creator.Name,
		SongCount: w.SongCnt,
		Cover:     cover,
		URL:       fmt.Sprintf("https://y.qq.com/n/ryqq/playlist/%d", id),
	}
}
