package qqmusic

import (
	"context"
	"fmt"
	neturl "net/url"
)

// Charts fetches the full list of available QQ Music charts (排行榜).
func (c *Client) Charts(ctx context.Context) ([]ChartInfo, error) {
	u := fmt.Sprintf("%s/v8/fcg-bin/fcg_myqq_toplist.fcg?g_tk=%s&uin=0&needNewCode=1&platform=h5&format=json&inCharset=utf-8&outCharset=utf-8",
		c.cfg.BaseURL, anonGTK)
	var resp wireChartListResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]ChartInfo, 0, len(resp.Data.TopList))
	for _, item := range resp.Data.TopList {
		out = append(out, ChartInfo{
			ID:    item.ID,
			Name:  item.TopTitle,
			Cover: item.PicURL,
		})
	}
	return out, nil
}

// Chart fetches songs from a specific chart by topid. limit=0 returns all.
func (c *Client) Chart(ctx context.Context, topid, limit int) ([]Song, error) {
	u := fmt.Sprintf("%s/v8/fcg-bin/fcg_v8_toplist_cp.fcg?topid=%d&g_tk=%s&uin=0&tpl=3&page=detail&type=top&song_begin=0&format=json&platform=h5&needNewCode=1",
		c.cfg.BaseURL, topid, anonGTK)
	var resp wireChartSongsResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	entries := resp.SongList
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	out := make([]Song, 0, len(entries))
	for i, e := range entries {
		out = append(out, chartSongFrom(e, i+1))
	}
	return out, nil
}

// Search searches QQ Music for songs matching keyword. Page is 1-based.
func (c *Client) Search(ctx context.Context, keyword string, limit, page int) ([]Song, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	u := fmt.Sprintf("%s/soso/fcgi-bin/search_for_qq_cp?format=json&t=0&w=%s&p=%d&n=%d&aggr=1&cr=1&g_tk=%s&loginUin=0&hostUin=0&inCharset=utf8&outCharset=utf-8&platform=yqq&needNewCode=0",
		c.cfg.BaseURL, neturl.QueryEscape(keyword), page, limit, anonGTK)
	var resp wireSearchResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	list := resp.Data.Song.List
	out := make([]Song, 0, len(list))
	rank := (page-1)*limit + 1
	for _, s := range list {
		out = append(out, searchSongFrom(s, rank))
		rank++
	}
	return out, nil
}

// Song fetches detail for one song by songmid.
func (c *Client) Song(ctx context.Context, songmid string) (Song, error) {
	u := fmt.Sprintf("%s/v8/fcg-bin/fcg_play_single_song.fcg?songmid=%s&format=json&tpl=yqq_song_detail",
		c.cfg.BaseURL, neturl.PathEscape(songmid))
	var resp wireSongDetailResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return Song{}, err
	}
	if len(resp.Data) == 0 {
		return Song{}, ErrNotFound
	}
	return songDetailFrom(resp.Data[0]), nil
}

// Lyrics fetches LRC lyrics for a song by songmid.
func (c *Client) Lyrics(ctx context.Context, songmid string) (Lyric, error) {
	u := fmt.Sprintf("%s/lyric/fcgi-bin/fcg_query_lyric_new.fcg?songmid=%s&g_tk=%s&loginUin=0&hostUin=0&format=json&nobase64=1&inCharset=utf8&outCharset=UTF-8&platform=yqq&needNewCode=0",
		c.cfg.BaseURL, neturl.PathEscape(songmid), anonGTK)
	var resp wireLyricResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return Lyric{}, err
	}
	return Lyric{
		SongMid:     songmid,
		LRC:         tryBase64(resp.Lyric),
		Translation: tryBase64(resp.Trans),
	}, nil
}

// Artist fetches a singer's top tracks by singermid via musicu.fcg POST.
// limit=0 returns up to 50.
func (c *Client) Artist(ctx context.Context, singermid string, limit int) ([]Song, error) {
	num := 50
	if limit > 0 && limit < 50 {
		num = limit
	}
	body := map[string]any{
		"comm": map[string]any{"uin": 0, "format": "json", "ct": 24, "cv": 0},
		"req_0": map[string]any{
			"module": "musichall.song_list_server",
			"method": "GetSingerSongList",
			"param":  map[string]any{"singermid": singermid, "begin": 0, "num": num, "order": 1},
		},
	}
	var resp wireMusicuArtistResp
	if err := c.postJSON(ctx, c.cfg.MusicuURL, body, &resp); err != nil {
		return nil, err
	}
	if resp.Req0.Code != 0 {
		return nil, ErrNotFound
	}
	list := resp.Req0.Data.SongList
	out := make([]Song, 0, len(list))
	for i, item := range list {
		if limit > 0 && i >= limit {
			break
		}
		out = append(out, artistSongFromMusicu(item.SongInfo, i+1))
	}
	return out, nil
}

// Album fetches an album's track list by albummid.
func (c *Client) Album(ctx context.Context, albummid string) ([]Song, error) {
	u := fmt.Sprintf("%s/v8/fcg-bin/fcg_v8_album_info_cp.fcg?albummid=%s&format=json&g_tk=%s&platform=h5page",
		c.cfg.BaseURL, neturl.PathEscape(albummid), anonGTK)
	var resp wireAlbumResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	if resp.Data.Name == "" {
		return nil, ErrNotFound
	}
	out := make([]Song, 0, len(resp.Data.List))
	for i, track := range resp.Data.List {
		s := albumTrackFrom(track, resp.Data.Name, i+1)
		out = append(out, s)
	}
	return out, nil
}

// Playlist fetches the song list for a playlist by its disstid.
func (c *Client) Playlist(ctx context.Context, disstid string) ([]Song, error) {
	u := fmt.Sprintf("%s/qzone/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg?type=1&json=1&utf8=1&onlysong=0&disstid=%s&format=json&g_tk=%s&loginUin=0&hostUin=0&platform=yqq&needNewCode=0",
		c.cfg.BaseURL, disstid, anonGTK)
	var resp wirePlaylistResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	if len(resp.CDList) == 0 {
		return nil, ErrNotFound
	}
	cd := resp.CDList[0]
	out := make([]Song, 0, len(cd.SongList))
	for i, s := range cd.SongList {
		out = append(out, playlistSongFrom(s, i+1))
	}
	return out, nil
}

// Playlists discovers playlists by category. categoryID=10000000 means all.
// limit controls how many to return (max ~30 per page; default 30).
func (c *Client) Playlists(ctx context.Context, categoryID, limit int) ([]PlaylistInfo, error) {
	if categoryID <= 0 {
		categoryID = 10000000
	}
	if limit <= 0 {
		limit = 30
	}
	ein := limit - 1
	u := fmt.Sprintf("%s/splcloud/fcgi-bin/fcg_get_diss_by_tag.fcg?categoryId=%d&sortId=5&sin=0&ein=%d&format=json&g_tk=%s&inCharset=utf8&outCharset=utf-8",
		c.cfg.BaseURL, categoryID, ein, anonGTK)
	var resp wirePlaylistsResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]PlaylistInfo, 0, len(resp.Data.List))
	for i, item := range resp.Data.List {
		out = append(out, playlistInfoFrom(item, i+1))
	}
	return out, nil
}
