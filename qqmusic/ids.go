package qqmusic

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSongMid extracts a QQ Music songmid from a bare mid or a song URL.
// Accepted forms:
//   - "002Zkt5709ysSQ"                                (bare mid)
//   - "https://y.qq.com/n/ryqq/songDetail/002Zkt5709ysSQ"
func ParseSongMid(s string) (string, error) {
	mid, err := parseSongMid(s)
	if err != nil {
		return "", fmt.Errorf("invalid song mid or URL %q: %w", s, err)
	}
	return mid, nil
}

func parseSongMid(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty input")
	}
	if strings.Contains(s, "songDetail/") {
		parts := strings.Split(s, "songDetail/")
		mid := strings.Trim(parts[len(parts)-1], "/")
		if mid != "" {
			return mid, nil
		}
	}
	if !strings.HasPrefix(s, "http") {
		return s, nil
	}
	return "", fmt.Errorf("unrecognised URL format")
}

// ParseSingerMid extracts a QQ Music singermid from a bare mid or a singer URL.
// Accepted forms:
//   - "003o2vm21NeSNX"                               (bare mid)
//   - "https://y.qq.com/n/ryqq/singer/003o2vm21NeSNX"
func ParseSingerMid(s string) (string, error) {
	mid, err := parseSingerMid(s)
	if err != nil {
		return "", fmt.Errorf("invalid singer mid or URL %q: %w", s, err)
	}
	return mid, nil
}

func parseSingerMid(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty input")
	}
	if strings.Contains(s, "/singer/") {
		parts := strings.Split(s, "/singer/")
		mid := strings.Trim(parts[len(parts)-1], "/")
		if mid != "" {
			return mid, nil
		}
	}
	if !strings.HasPrefix(s, "http") {
		return s, nil
	}
	return "", fmt.Errorf("unrecognised URL format")
}

// ParseAlbumMid extracts a QQ Music albummid from a bare mid or an album URL.
// Accepted forms:
//   - "002M1v7j4B5Yf7"                               (bare mid)
//   - "https://y.qq.com/n/ryqq/album/002M1v7j4B5Yf7"
func ParseAlbumMid(s string) (string, error) {
	mid, err := parseAlbumMid(s)
	if err != nil {
		return "", fmt.Errorf("invalid album mid or URL %q: %w", s, err)
	}
	return mid, nil
}

func parseAlbumMid(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty input")
	}
	if strings.Contains(s, "/album/") {
		parts := strings.Split(s, "/album/")
		mid := strings.Trim(parts[len(parts)-1], "/")
		if mid != "" {
			return mid, nil
		}
	}
	if !strings.HasPrefix(s, "http") {
		return s, nil
	}
	return "", fmt.Errorf("unrecognised URL format")
}

// ParsePlaylistID extracts a QQ Music playlist disstid (integer) from a bare
// number or a playlist URL.
// Accepted forms:
//   - "8888"                                         (bare integer)
//   - "https://y.qq.com/n/ryqq/playlist/8888"
func ParsePlaylistID(s string) (string, error) {
	id, err := parsePlaylistID(s)
	if err != nil {
		return "", fmt.Errorf("invalid playlist id or URL %q: %w", s, err)
	}
	return id, nil
}

func parsePlaylistID(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty input")
	}
	if strings.Contains(s, "/playlist/") {
		parts := strings.Split(s, "/playlist/")
		id := strings.Trim(parts[len(parts)-1], "/")
		if _, err := strconv.ParseInt(id, 10, 64); err == nil {
			return id, nil
		}
	}
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return s, nil
	}
	if strings.HasPrefix(s, "http") {
		return "", fmt.Errorf("unrecognised URL format")
	}
	return "", fmt.Errorf("not a numeric playlist id")
}

// chartTopID resolves a chart name or numeric topid string to an integer.
// Known names: hot=4, new=26, rising=27.
func chartTopID(nameOrID string) (int, error) {
	switch strings.ToLower(strings.TrimSpace(nameOrID)) {
	case "hot", "热歌榜", "":
		return 4, nil
	case "new", "新歌榜":
		return 26, nil
	case "rising", "飙升榜":
		return 27, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(nameOrID))
	if err == nil && n > 0 {
		return n, nil
	}
	return 0, fmt.Errorf("unknown chart %q (use a name like hot/new/rising or a numeric topid)", nameOrID)
}
