// Package qqmusic is the library behind the qqmusic command: the HTTP client,
// request shaping, and the typed data models for QQ Music (QQ音乐).
//
// The client fetches chart data from the public QQ Music toplist endpoint at
// https://c.y.qq.com. No authentication is required. It sets a real
// User-Agent and Referer, paces requests, and retries transient 429/5xx
// errors with exponential backoff.
package qqmusic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// charts maps chart name to topid parameter.
var charts = map[string]int{
	"hot":    4,  // 热歌榜
	"new":    26, // 新歌榜
	"rising": 27, // 飙升榜
}

// ChartLabels maps chart name to its Chinese label.
var ChartLabels = map[string]string{
	"hot":    "热歌榜",
	"new":    "新歌榜",
	"rising": "飙升榜",
}

// Config holds constructor parameters for the Client.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://c.y.qq.com",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Rate:      500 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the QQ Music toplist API.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

// Top fetches the chart specified by chart ("hot", "new", or "rising") and
// returns up to limit songs in ranked order.
func (c *Client) Top(ctx context.Context, chart string, limit int) ([]Song, error) {
	topid, ok := charts[chart]
	if !ok {
		return nil, fmt.Errorf("unknown chart %q", chart)
	}

	url := fmt.Sprintf(
		"%s/v8/fcg-bin/fcg_v8_toplist_opt.fcg?tpl=3&page=detail&type=top&topid=%d&needNewCode=1&uin=0&format=json&inCharset=utf-8&outCharset=utf-8&notice=0&platform=h5",
		c.cfg.BaseURL, topid,
	)

	raw, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp wireResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("decode toplist: %w", err)
	}

	entries := resp.SongList
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	out := make([]Song, 0, len(entries))
	for i, e := range entries {
		out = append(out, wireToSong(e, i+1))
	}
	return out, nil
}

// get performs an HTTP GET with the Referer and User-Agent headers set,
// retrying on transient failures.
func (c *Client) get(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, url)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", url, lastErr)
}

func (c *Client) do(ctx context.Context, url string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Referer", "https://y.qq.com/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
