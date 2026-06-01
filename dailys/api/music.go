package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	itunesSearchURL   = "https://itunes.apple.com/search"
	itunesCountryCode = "US"
	musicSearchLimit  = 50
)

type MusicGenre struct {
	Name       string
	SearchTerm string
}

type MusicPick struct {
	Genre        string
	TrackName    string
	ArtistName   string
	AlbumName    string
	PrimaryGenre string
	ReleaseDate  time.Time
	PreviewURL   string
	ArtworkURL   string
	TrackURL     string
}

type itunesSearchResponse struct {
	ResultCount int           `json:"resultCount"`
	Results     []itunesTrack `json:"results"`
}

type itunesTrack struct {
	Kind         string `json:"kind"`
	TrackName    string `json:"trackName"`
	ArtistName   string `json:"artistName"`
	Collection   string `json:"collectionName"`
	PrimaryGenre string `json:"primaryGenreName"`
	ReleaseDate  string `json:"releaseDate"`
	PreviewURL   string `json:"previewUrl"`
	ArtworkURL   string `json:"artworkUrl100"`
	TrackURL     string `json:"trackViewUrl"`
}

func DefaultMusicGenres() []MusicGenre {
	return []MusicGenre{
		{Name: "Ambient", SearchTerm: "ambient"},
		{Name: "Blues", SearchTerm: "blues"},
		{Name: "Classical", SearchTerm: "classical"},
		{Name: "Disco", SearchTerm: "disco"},
		{Name: "Electronic", SearchTerm: "electronic"},
		{Name: "Folk", SearchTerm: "folk"},
		{Name: "Funk", SearchTerm: "funk"},
		{Name: "Hip-Hop", SearchTerm: "hip hop"},
		{Name: "Indie", SearchTerm: "indie"},
		{Name: "Jazz", SearchTerm: "jazz"},
		{Name: "Latin", SearchTerm: "latin"},
		{Name: "Metal", SearchTerm: "metal"},
		{Name: "New Wave", SearchTerm: "new wave"},
		{Name: "Punk", SearchTerm: "punk"},
		{Name: "R&B", SearchTerm: "r&b"},
		{Name: "Reggae", SearchTerm: "reggae"},
		{Name: "Rock", SearchTerm: "rock"},
		{Name: "Soul", SearchTerm: "soul"},
		{Name: "Synthpop", SearchTerm: "synthpop"},
		{Name: "World", SearchTerm: "world music"},
	}
}

func RandomMusic() (MusicPick, error) {
	return RandomMusicFromGenres(DefaultMusicGenres())
}

func RandomMusicFromGenres(genres []MusicGenre) (MusicPick, error) {
	if len(genres) == 0 {
		return MusicPick{}, fmt.Errorf("no music genres configured")
	}

	genreIndex, err := randomInt(len(genres))
	if err != nil {
		return MusicPick{}, err
	}

	genre := genres[genreIndex]
	tracks, err := searchITunesSongs(genre.SearchTerm)
	if err != nil {
		return MusicPick{}, err
	}

	if len(tracks) == 0 {
		return MusicPick{}, fmt.Errorf("no songs found for %s", genre.Name)
	}

	trackIndex, err := randomInt(len(tracks))
	if err != nil {
		return MusicPick{}, err
	}

	track := tracks[trackIndex]
	return musicPickFromTrack(genre, track), nil
}

func searchITunesSongs(term string) ([]itunesTrack, error) {
	endpoint, err := url.Parse(itunesSearchURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing music search URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("term", term)
	query.Set("country", itunesCountryCode)
	query.Set("media", "music")
	query.Set("entity", "song")
	query.Set("limit", fmt.Sprint(musicSearchLimit))
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating music request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending music request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("music API error: %s", resp.Status)
	}

	var response itunesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding music response: %w", err)
	}

	return usableITunesTracks(response.Results), nil
}

func usableITunesTracks(tracks []itunesTrack) []itunesTrack {
	filtered := make([]itunesTrack, 0, len(tracks))

	for _, track := range tracks {
		if track.Kind != "song" {
			continue
		}

		if strings.TrimSpace(track.TrackName) == "" || strings.TrimSpace(track.ArtistName) == "" {
			continue
		}

		filtered = append(filtered, track)
	}

	return filtered
}

func musicPickFromTrack(genre MusicGenre, track itunesTrack) MusicPick {
	return MusicPick{
		Genre:        genre.Name,
		TrackName:    strings.TrimSpace(track.TrackName),
		ArtistName:   strings.TrimSpace(track.ArtistName),
		AlbumName:    strings.TrimSpace(track.Collection),
		PrimaryGenre: strings.TrimSpace(track.PrimaryGenre),
		ReleaseDate:  parseITunesReleaseDate(track.ReleaseDate),
		PreviewURL:   strings.TrimSpace(track.PreviewURL),
		ArtworkURL:   highResolutionArtworkURL(strings.TrimSpace(track.ArtworkURL)),
		TrackURL:     strings.TrimSpace(track.TrackURL),
	}
}

func parseITunesReleaseDate(value string) time.Time {
	releasedAt, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}

	return releasedAt
}

func highResolutionArtworkURL(imageURL string) string {
	lastSlash := strings.LastIndex(imageURL, "/")
	if imageURL == "" || lastSlash == -1 || lastSlash == len(imageURL)-1 {
		return imageURL
	}

	sizeSegment := imageURL[lastSlash+1:]
	if !strings.Contains(sizeSegment, "x") {
		return imageURL
	}

	return imageURL[:lastSlash+1] + "600x600bb.jpg"
}
