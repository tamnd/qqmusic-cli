package qqmusic

import "testing"

func TestParseSongMid(t *testing.T) {
	cases := []struct{ in, want string; ok bool }{
		{"002Zkt5709ysSQ", "002Zkt5709ysSQ", true},
		{"https://y.qq.com/n/ryqq/songDetail/002Zkt5709ysSQ", "002Zkt5709ysSQ", true},
		{"", "", false},
	}
	for _, tc := range cases {
		got, err := ParseSongMid(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("ParseSongMid(%q) = (%q, %v)", tc.in, got, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("ParseSongMid(%q) succeeded, want error", tc.in)
		}
	}
}

func TestParseSingerMid(t *testing.T) {
	cases := []struct{ in, want string; ok bool }{
		{"003o2vm21NeSNX", "003o2vm21NeSNX", true},
		{"https://y.qq.com/n/ryqq/singer/003o2vm21NeSNX", "003o2vm21NeSNX", true},
		{"", "", false},
	}
	for _, tc := range cases {
		got, err := ParseSingerMid(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("ParseSingerMid(%q) = (%q, %v)", tc.in, got, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("ParseSingerMid(%q) succeeded, want error", tc.in)
		}
	}
}

func TestParseAlbumMid(t *testing.T) {
	cases := []struct{ in, want string; ok bool }{
		{"002M1v7j4B5Yf7", "002M1v7j4B5Yf7", true},
		{"https://y.qq.com/n/ryqq/album/002M1v7j4B5Yf7", "002M1v7j4B5Yf7", true},
		{"", "", false},
	}
	for _, tc := range cases {
		got, err := ParseAlbumMid(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("ParseAlbumMid(%q) = (%q, %v)", tc.in, got, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("ParseAlbumMid(%q) succeeded, want error", tc.in)
		}
	}
}

func TestParsePlaylistID(t *testing.T) {
	cases := []struct{ in, want string; ok bool }{
		{"8888", "8888", true},
		{"https://y.qq.com/n/ryqq/playlist/8888", "8888", true},
		{"not-a-number", "", false},
		{"", "", false},
	}
	for _, tc := range cases {
		got, err := ParsePlaylistID(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("ParsePlaylistID(%q) = (%q, %v)", tc.in, got, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("ParsePlaylistID(%q) succeeded, want error", tc.in)
		}
	}
}

func TestChartTopID(t *testing.T) {
	cases := []struct {
		in   string
		want int
		ok   bool
	}{
		{"hot", 4, true},
		{"new", 26, true},
		{"rising", 27, true},
		{"4", 4, true},
		{"", 4, true}, // empty defaults to hot
		{"bogus", 0, false},
	}
	for _, tc := range cases {
		got, err := chartTopID(tc.in)
		if tc.ok && (err != nil || got != tc.want) {
			t.Errorf("chartTopID(%q) = (%d, %v)", tc.in, got, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("chartTopID(%q) succeeded, want error", tc.in)
		}
	}
}
