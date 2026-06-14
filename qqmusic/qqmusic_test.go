package qqmusic_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tamnd/qqmusic-cli/qqmusic"
)

const mockToplistResponse = `{
  "code": 0,
  "date": "2026-06-14",
  "songlist": [
    {
      "Franking_value": "258287",
      "cur_count": "258287",
      "data": {
        "songid": 372025766,
        "songmid": "0030KTZ8141fsb",
        "songname": "NIGHT DANCER",
        "albumname": "NIGHT DANCER",
        "albumid": 29960153,
        "albummid": "002M1v7j4B5Yf7",
        "interval": 210,
        "singer": [{"id": 10178787, "mid": "003o2vm21NeSNX", "name": "imase"}]
      }
    },
    {
      "Franking_value": "123456",
      "cur_count": "123456",
      "data": {
        "songid": 123456789,
        "songmid": "001abcDEF",
        "songname": "Test Song",
        "albumname": "Test Album",
        "albumid": 100,
        "albummid": "abc123",
        "interval": 185,
        "singer": [
          {"id": 1, "mid": "singer1", "name": "Artist A"},
          {"id": 2, "mid": "singer2", "name": "Artist B"}
        ]
      }
    }
  ]
}`

func newTestClient(ts *httptest.Server) *qqmusic.Client {
	cfg := qqmusic.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return qqmusic.NewClient(cfg)
}

func TestTopSetsReferer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Referer")
		if got != "https://y.qq.com/" {
			t.Errorf("Referer = %q, want %q", got, "https://y.qq.com/")
		}
		_, _ = w.Write([]byte(mockToplistResponse))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.Top(context.Background(), "hot", 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTopParsesResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(mockToplistResponse))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	songs, err := c.Top(context.Background(), "hot", 0)
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
		t.Errorf("title = %q, want NIGHT DANCER", s.Title)
	}
	if s.Artists != "imase" {
		t.Errorf("artists = %q, want imase", s.Artists)
	}
	if s.Duration != "3:30" {
		t.Errorf("duration = %q, want 3:30", s.Duration)
	}
	if s.URL != "https://y.qq.com/n/ryqq/songDetail/0030KTZ8141fsb" {
		t.Errorf("url = %q", s.URL)
	}
	if s.PlayCount != "258287" {
		t.Errorf("play_count = %q, want 258287", s.PlayCount)
	}

	s2 := songs[1]
	if s2.Rank != 2 {
		t.Errorf("rank = %d, want 2", s2.Rank)
	}
	if s2.Artists != "Artist A, Artist B" {
		t.Errorf("artists = %q, want 'Artist A, Artist B'", s2.Artists)
	}
	// 185 seconds = 3:05
	if s2.Duration != "3:05" {
		t.Errorf("duration = %q, want 3:05", s2.Duration)
	}
}

func TestTopRespectsLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(mockToplistResponse))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	songs, err := c.Top(context.Background(), "hot", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
}

func TestTopUnknownChart(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(mockToplistResponse))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.Top(context.Background(), "bogus", 10)
	if err == nil {
		t.Fatal("expected error for unknown chart, got nil")
	}
}

func TestTopRetriesOn503(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(mockToplistResponse))
	}))
	defer srv.Close()

	cfg := qqmusic.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	cfg.Retries = 5
	c := qqmusic.NewClient(cfg)

	start := time.Now()
	_, err := c.Top(context.Background(), "hot", 0)
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
