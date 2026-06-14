package qqmusic

import (
	"fmt"
	"strings"
)

// wireResponse is the top-level JSON returned by the toplist endpoint.
type wireResponse struct {
	Code     int         `json:"code"`
	Date     string      `json:"date"`
	SongList []wireEntry `json:"songlist"`
}

// wireEntry is one entry in the songlist array.
type wireEntry struct {
	FrankingValue string   `json:"Franking_value"`
	CurCount      string   `json:"cur_count"`
	Data          wireData `json:"data"`
}

// wireData holds the song metadata embedded in each entry.
type wireData struct {
	SongID    int64        `json:"songid"`
	SongMid   string       `json:"songmid"`
	SongName  string       `json:"songname"`
	AlbumName string       `json:"albumname"`
	AlbumID   int64        `json:"albumid"`
	Interval  int          `json:"interval"` // seconds
	Singer    []wireSinger `json:"singer"`
}

// wireSinger is one artist in a song's singer array.
type wireSinger struct {
	ID   int64  `json:"id"`
	Mid  string `json:"mid"`
	Name string `json:"name"`
}

// Song is the public record emitted for every chart song.
type Song struct {
	Rank      int    `json:"rank"`
	Title     string `json:"title"`
	Artists   string `json:"artists"`
	Album     string `json:"album"`
	Duration  string `json:"duration"`
	PlayCount string `json:"play_count"`
	URL       string `json:"url"`
}

// wireToSong converts a wireEntry to a Song.
func wireToSong(e wireEntry, rank int) Song {
	names := make([]string, 0, len(e.Data.Singer))
	for _, s := range e.Data.Singer {
		names = append(names, s.Name)
	}
	secs := e.Data.Interval
	dur := fmt.Sprintf("%d:%02d", secs/60, secs%60)
	url := "https://y.qq.com/n/ryqq/songDetail/" + e.Data.SongMid
	return Song{
		Rank:      rank,
		Title:     e.Data.SongName,
		Artists:   strings.Join(names, ", "),
		Album:     e.Data.AlbumName,
		Duration:  dur,
		PlayCount: e.CurCount,
		URL:       url,
	}
}
