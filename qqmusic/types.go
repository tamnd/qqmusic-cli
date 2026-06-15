package qqmusic

import "errors"

// ErrNotFound is returned when the QQ Music API returns HTTP 404 or
// an empty result for a known-valid ID.
var ErrNotFound = errors.New("not found")

// Song is a QQ Music song record, used by chart, search, artist, album,
// and playlist commands.
type Song struct {
	Rank     int    `json:"rank"     table:",right"`
	SongMid  string `json:"songmid"`
	Title    string `json:"title"    table:",truncate"`
	Artists  string `json:"artists"  table:",truncate"`
	Album    string `json:"album"    table:",truncate"`
	AlbumMid string `json:"album_mid"`
	Duration string `json:"duration" table:",right"` // "M:SS"
	Date     string `json:"date"`
	URL      string `json:"url"      kit:"url" table:",truncate"`
}

// ChartInfo is one item from the chart listing.
type ChartInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Cover string `json:"cover" table:",truncate"`
}

// Lyric holds lyrics for one song in LRC format.
type Lyric struct {
	SongMid     string `json:"songmid"`
	LRC         string `json:"lrc"`
	Translation string `json:"translation"`
}

// PlaylistInfo is a lightweight playlist record for the playlists listing.
type PlaylistInfo struct {
	Rank      int    `json:"rank"       table:",right"`
	ID        int64  `json:"id"`
	Name      string `json:"name"       table:",truncate"`
	Creator   string `json:"creator"`
	SongCount int    `json:"song_count" table:",right"`
	Cover     string `json:"cover"      table:",truncate"`
	URL       string `json:"url"        kit:"url" table:",truncate"`
}

// --- wire types (unexported, used for JSON decoding only) ---

type wireChartListResp struct {
	Code int           `json:"code"`
	Data wireChartData `json:"data"`
}

type wireChartData struct {
	TopList []wireChartItem `json:"topList"`
}

type wireChartItem struct {
	ID       int    `json:"id"`
	TopTitle string `json:"topTitle"`
	PicURL   string `json:"picUrl"`
}

type wireChartSongsResp struct {
	Code     int             `json:"code"`
	SongList []wireChartSong `json:"songlist"`
}

type wireChartSong struct {
	CurCount string            `json:"cur_count"`
	Data     wireChartSongData `json:"data"`
}

type wireChartSongData struct {
	SongMid    string       `json:"songmid"`
	SongName   string       `json:"songname"`
	AlbumName  string       `json:"albumname"`
	AlbumMid   string       `json:"albummid"`
	Interval   int          `json:"interval"`
	TimePublic string       `json:"time_public"`
	Singer     []wireSinger `json:"singer"`
}

type wireSinger struct {
	Mid  string `json:"mid"`
	Name string `json:"name"`
}

type wireSearchResp struct {
	Code int            `json:"code"`
	Data wireSearchData `json:"data"`
}

type wireSearchData struct {
	Song wireSearchSongs `json:"song"`
}

type wireSearchSongs struct {
	List []wireSearchSong `json:"list"`
}

type wireSearchSong struct {
	SongMid    string       `json:"songmid"`
	SongName   string       `json:"songname"`
	AlbumName  string       `json:"albumname"`
	AlbumMid   string       `json:"albummid"`
	Interval   int          `json:"interval"`
	TimePublic string       `json:"time_public"` // old endpoint
	PubTime    int64        `json:"pubtime"`     // search_for_qq_cp uses unix timestamp
	Singer     []wireSinger `json:"singer"`
}

type wireSongDetailResp struct {
	Code int              `json:"code"`
	Data []wireSongDetail `json:"data"`
}

type wireSongDetail struct {
	Mid        string          `json:"mid"`         // songmid
	Name       string          `json:"name"`        // song title
	Album      wireSongAlbum   `json:"album"`
	Interval   int             `json:"interval"`
	TimePublic string          `json:"time_public"`
	Singer     []wireSinger    `json:"singer"`
}

type wireSongAlbum struct {
	Mid  string `json:"mid"`
	Name string `json:"name"`
}

type wireLyricResp struct {
	RetCode int    `json:"retcode"`
	Lyric   string `json:"lyric"`
	Trans   string `json:"trans"`
}

type wireArtistResp struct {
	Code int            `json:"code"`
	Data wireArtistData `json:"data"`
}

type wireArtistData struct {
	SingerInfo wireArtistInfo   `json:"singerinfo"`
	List       []wireArtistSong `json:"list"`
}

type wireArtistInfo struct {
	SingerMid  string `json:"singer_mid"`
	SingerName string `json:"singer_name"`
}

type wireArtistSong struct {
	MusicData wireArtistMusicData `json:"musicData"`
}

type wireArtistMusicData struct {
	SongMid    string       `json:"songmid"`
	SongName   string       `json:"songname"`
	AlbumName  string       `json:"albumname"`
	AlbumMid   string       `json:"albummid"`
	Interval   int          `json:"interval"`
	TimePublic string       `json:"time_public"`
	Singer     []wireSinger `json:"singer"`
}

type wireAlbumResp struct {
	Code int           `json:"code"`
	Data wireAlbumData `json:"data"`
}

type wireAlbumData struct {
	Name       string           `json:"name"`
	ADate      string           `json:"aDate"`
	Desc       string           `json:"desc"`
	SingerName string           `json:"singerName"`
	SingerMid  string           `json:"singerMid"`
	AlbumMid   string           `json:"mid"`
	List       []wireAlbumTrack `json:"list"`
}

type wireAlbumTrack struct {
	SongMid  string       `json:"songmid"`
	SongName string       `json:"songname"`
	Interval int          `json:"interval"`
	Singer   []wireSinger `json:"singer"`
}

type wirePlaylistResp struct {
	Code   int              `json:"code"`
	CDList []wirePlaylistCD `json:"cdlist"`
}

type wirePlaylistCD struct {
	DissID   string               `json:"dissid"`
	DissName string               `json:"dissname"`
	Desc     string               `json:"desc"`
	Logo     string               `json:"logo"`
	Creator  wirePlaylistCreator  `json:"creator"`
	SongList []wirePlaylistSong   `json:"songlist"`
}

type wirePlaylistCreator struct {
	Nick string `json:"nick"`
}

type wirePlaylistSong struct {
	SongMid   string       `json:"songmid"`
	SongName  string       `json:"songname"`
	AlbumName string       `json:"albumname"`
	AlbumMid  string       `json:"albummid"`
	Interval  int          `json:"interval"`
	Singer    []wireSinger `json:"singer"`
}

type wirePlaylistsResp struct {
	Code int               `json:"code"`
	Data wirePlaylistsData `json:"data"`
}

type wirePlaylistsData struct {
	List []wirePlaylistItem `json:"list"`
}

type wirePlaylistItem struct {
	DissID    interface{}          `json:"dissid"` // API returns int or string
	DissName  string               `json:"dissname"`
	ImgURL    string               `json:"imgurl"`  // live API field
	Logo      string               `json:"logo"`    // alternate field name
	Creator   wirePlaylistsCreator `json:"creator"`
	SongCnt   int                  `json:"song_cnt"` // not always present
	ListenNum int                  `json:"listennum"`
}

type wirePlaylistsCreator struct {
	Name string `json:"name"`
}

// musicu.fcg POST (u.y.qq.com) - used for artist song list
type wireMusicuArtistResp struct {
	Req0 wireMusicuArtistReq0 `json:"req_0"`
}

type wireMusicuArtistReq0 struct {
	Code int                 `json:"code"`
	Data wireMusicuArtistData `json:"data"`
}

type wireMusicuArtistData struct {
	SingerMid string               `json:"singerMid"`
	TotalNum  int                  `json:"totalNum"`
	SongList  []wireMusicuSongItem `json:"songList"`
}

type wireMusicuSongItem struct {
	SongInfo wireMusicuSongInfo `json:"songInfo"`
}

type wireMusicuSongInfo struct {
	Mid        string       `json:"mid"`
	Name       string       `json:"name"`
	Singer     []wireSinger `json:"singer"`
	Album      wireSongAlbum `json:"album"`
	Interval   int          `json:"interval"`
	TimePublic string       `json:"time_public"`
}
