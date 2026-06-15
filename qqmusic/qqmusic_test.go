package qqmusic_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tamnd/qqmusic-cli/qqmusic"
)

// newTestClient points the client at a local httptest server with no rate limit.
func newTestClient(ts *httptest.Server) *qqmusic.Client {
	cfg := qqmusic.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.MusicuURL = ts.URL + "/musicu"
	cfg.Rate = 0
	cfg.Retries = 0
	return qqmusic.NewClient(cfg)
}

// --- fixtures ---

const fixtureChartList = `{
  "code": 0,
  "data": {
    "topList": [
      {"id": 4, "topTitle": "热歌榜", "picUrl": "https://y.gtimg.cn/hot.jpg"},
      {"id": 26, "topTitle": "新歌榜", "picUrl": "https://y.gtimg.cn/new.jpg"}
    ]
  }
}`

const fixtureChartSongs = `{
  "code": 0,
  "date": "2026-06-15",
  "songlist": [
    {
      "cur_count": "258287",
      "data": {
        "songmid": "0030KTZ8141fsb",
        "songname": "NIGHT DANCER",
        "albumname": "NIGHT DANCER",
        "albummid": "002M1v7j4B5Yf7",
        "interval": 210,
        "time_public": "2023-09-06",
        "singer": [{"mid": "003o2vm21NeSNX", "name": "imase"}]
      }
    },
    {
      "cur_count": "123456",
      "data": {
        "songmid": "001abcDEFxyz",
        "songname": "Test Song",
        "albumname": "Test Album",
        "albummid": "abc123xyz",
        "interval": 185,
        "time_public": "2024-01-01",
        "singer": [
          {"mid": "s1", "name": "Artist A"},
          {"mid": "s2", "name": "Artist B"}
        ]
      }
    }
  ]
}`

const fixtureSearch = `{
  "code": 0,
  "data": {
    "song": {
      "curnum": 1,
      "curpage": 1,
      "list": [
        {
          "songmid": "002Zkt5709ysSQ",
          "songname": "晴天",
          "albumname": "叶惠美",
          "albummid": "001IPgQR3UHOcs",
          "interval": 269,
          "time_public": "2003-07-31",
          "singer": [{"mid": "000Sp0Bz4JXH0E", "name": "周杰伦"}]
        }
      ]
    }
  }
}`

const fixtureSongDetail = `{
  "code": 0,
  "data": [
    {
      "mid": "002Zkt5709ysSQ",
      "name": "晴天",
      "album": {"mid": "001IPgQR3UHOcs", "name": "叶惠美"},
      "interval": 269,
      "time_public": "2003-07-31",
      "singer": [{"mid": "000Sp0Bz4JXH0E", "name": "周杰伦"}]
    }
  ]
}`

const fixtureLyrics = `{
  "retcode": 0,
  "lyric": "[00:00.00]晴天\n[00:05.00]故事的小黄花",
  "trans": ""
}`

const fixtureArtist = `{
  "req_0": {
    "code": 0,
    "data": {
      "singerMid": "0025NhlN2yWrP4",
      "totalNum": 1013,
      "songList": [
        {
          "songInfo": {
            "mid": "002Zkt5709ysSQ",
            "name": "晴天",
            "singer": [{"mid": "0025NhlN2yWrP4", "name": "周杰伦"}],
            "album": {"mid": "001IPgQR3UHOcs", "name": "叶惠美"},
            "interval": 269,
            "time_public": "2003-07-31"
          }
        }
      ]
    }
  }
}`

const fixtureAlbum = `{
  "code": 0,
  "data": {
    "name": "叶惠美",
    "aDate": "2003-07-31",
    "desc": "周杰伦第四张专辑",
    "singerName": "周杰伦",
    "singerMid": "000Sp0Bz4JXH0E",
    "mid": "001IPgQR3UHOcs",
    "list": [
      {
        "songmid": "002Zkt5709ysSQ",
        "songname": "晴天",
        "interval": 269,
        "singer": [{"mid": "000Sp0Bz4JXH0E", "name": "周杰伦"}]
      }
    ]
  }
}`

const fixturePlaylist = `{
  "code": 0,
  "cdlist": [
    {
      "dissid": "8888",
      "dissname": "我喜欢的歌",
      "desc": "test playlist",
      "logo": "https://y.gtimg.cn/logo.jpg",
      "creator": {"nick": "user123"},
      "songlist": [
        {
          "songmid": "002Zkt5709ysSQ",
          "songname": "晴天",
          "albumname": "叶惠美",
          "albummid": "001IPgQR3UHOcs",
          "interval": 269,
          "singer": [{"mid": "000Sp0Bz4JXH0E", "name": "周杰伦"}]
        }
      ]
    }
  ]
}`

const fixturePlaylists = `{
  "code": 0,
  "data": {
    "list": [
      {
        "dissid": "12345678",
        "dissname": "夜晚安静时刻",
        "logo": "https://y.gtimg.cn/logo.jpg",
        "creator": {"name": "creator_user"},
        "song_cnt": 30
      }
    ]
  }
}`

// --- tests ---

func TestChartsReferer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Referer"); got != "https://y.qq.com/" {
			t.Errorf("Referer = %q, want https://y.qq.com/", got)
		}
		_, _ = w.Write([]byte(fixtureChartList))
	}))
	defer srv.Close()
	_, err := newTestClient(srv).Charts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestChartsParsesItems(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureChartList))
	}))
	defer srv.Close()
	items, err := newTestClient(srv).Charts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d charts, want 2", len(items))
	}
	if items[0].ID != 4 || items[0].Name != "热歌榜" {
		t.Errorf("item[0] = %+v", items[0])
	}
}

func TestChartParsesSongs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureChartSongs))
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Chart(context.Background(), 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 2 {
		t.Fatalf("got %d songs, want 2", len(songs))
	}
	s := songs[0]
	if s.Rank != 1 {
		t.Errorf("rank = %d, want 1", s.Rank)
	}
	if s.Title != "NIGHT DANCER" {
		t.Errorf("title = %q", s.Title)
	}
	if s.Artists != "imase" {
		t.Errorf("artists = %q", s.Artists)
	}
	if s.Duration != "3:30" {
		t.Errorf("duration = %q, want 3:30", s.Duration)
	}
	if s.URL != "https://y.qq.com/n/ryqq/songDetail/0030KTZ8141fsb" {
		t.Errorf("url = %q", s.URL)
	}
	s2 := songs[1]
	if s2.Artists != "Artist A, Artist B" {
		t.Errorf("artists = %q", s2.Artists)
	}
	if s2.Duration != "3:05" {
		t.Errorf("duration = %q, want 3:05", s2.Duration)
	}
}

func TestChartLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureChartSongs))
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Chart(context.Background(), 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
}

func TestSearchParsesResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureSearch))
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Search(context.Background(), "晴天", 20, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
	if songs[0].Title != "晴天" {
		t.Errorf("title = %q", songs[0].Title)
	}
}

func TestSongDetailParses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureSongDetail))
	}))
	defer srv.Close()
	s, err := newTestClient(srv).Song(context.Background(), "002Zkt5709ysSQ")
	if err != nil {
		t.Fatal(err)
	}
	if s.SongMid != "002Zkt5709ysSQ" {
		t.Errorf("songmid = %q", s.SongMid)
	}
	if s.Title != "晴天" {
		t.Errorf("title = %q", s.Title)
	}
}

func TestSongDetailEmptyIsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"code":0,"data":[]}`))
	}))
	defer srv.Close()
	_, err := newTestClient(srv).Song(context.Background(), "bogus")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLyricsParsesLRC(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureLyrics))
	}))
	defer srv.Close()
	l, err := newTestClient(srv).Lyrics(context.Background(), "002Zkt5709ysSQ")
	if err != nil {
		t.Fatal(err)
	}
	if l.SongMid != "002Zkt5709ysSQ" {
		t.Errorf("songmid = %q", l.SongMid)
	}
	if l.LRC == "" {
		t.Error("LRC is empty")
	}
}

func TestArtistParsesTrackList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/musicu" {
			_, _ = w.Write([]byte(fixtureArtist))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Artist(context.Background(), "0025NhlN2yWrP4", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
	if songs[0].Title != "晴天" {
		t.Errorf("title = %q", songs[0].Title)
	}
}

func TestAlbumParsesTrackList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixtureAlbum))
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Album(context.Background(), "001IPgQR3UHOcs")
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
	if songs[0].Album != "叶惠美" {
		t.Errorf("album = %q", songs[0].Album)
	}
	if songs[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", songs[0].Rank)
	}
}

func TestPlaylistParsesSongs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixturePlaylist))
	}))
	defer srv.Close()
	songs, err := newTestClient(srv).Playlist(context.Background(), "8888")
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
}

func TestPlaylistsDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fixturePlaylists))
	}))
	defer srv.Close()
	items, err := newTestClient(srv).Playlists(context.Background(), 10000000, 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d playlists, want 1", len(items))
	}
	if items[0].Name != "夜晚安静时刻" {
		t.Errorf("name = %q", items[0].Name)
	}
	if items[0].SongCount != 30 {
		t.Errorf("song_count = %d, want 30", items[0].SongCount)
	}
}

func TestRetriesOn503(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(fixtureChartSongs))
	}))
	defer srv.Close()

	cfg := qqmusic.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	cfg.Retries = 5
	c := qqmusic.NewClient(cfg)

	start := time.Now()
	_, err := c.Chart(context.Background(), 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
	if time.Since(start) < 500*time.Millisecond {
		t.Error("retries did not back off")
	}
}
