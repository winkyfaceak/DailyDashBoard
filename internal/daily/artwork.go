package daily

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

const (
	articBaseURL         = "https://api.artic.edu/api/v1"
	articImageMirrorURL  = "https://lakeimagesweb.artic.edu/iiif/2"
	articSearchLimit     = 100
	articMaxSearchResult = 1000
)

type Artwork struct {
	ID               int
	Artist           string
	Title            string
	Date             string
	Description      string
	ImageURL         string
	MirrorImageURL   string
	ThumbnailDataURL string
}

type artworkSearchResponse struct {
	Pagination artworkPagination `json:"pagination"`
	Data       []artworkData     `json:"data"`
	Config     artworkConfig     `json:"config"`
}

type artworkDetailResponse struct {
	Data   artworkData   `json:"data"`
	Config artworkConfig `json:"config"`
}

type artworkPagination struct {
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type artworkConfig struct {
	IIIFURL string `json:"iiif_url"`
}

type artworkData struct {
	ID            int              `json:"id"`
	Title         string           `json:"title"`
	ArtistTitle   string           `json:"artist_title"`
	ArtistDisplay string           `json:"artist_display"`
	DateDisplay   string           `json:"date_display"`
	Description   string           `json:"description"`
	ImageID       string           `json:"image_id"`
	Thumbnail     artworkThumbnail `json:"thumbnail"`
}

type artworkThumbnail struct {
	LQIP    string `json:"lqip"`
	AltText string `json:"alt_text"`
}

func RandomArtwork() (Artwork, error) {
	totalPages, err := publicDomainArtworkPages()
	if err != nil {
		return Artwork{}, err
	}

	if totalPages == 0 {
		return Artwork{}, fmt.Errorf("no public domain artworks found")
	}

	var lastErr error

	for attempt := 1; attempt <= 5; attempt++ {
		page, err := randomInt(totalPages)
		if err != nil {
			return Artwork{}, err
		}

		search, err := searchPublicDomainArtworks(page + 1)
		if err != nil {
			lastErr = err
			continue
		}

		candidates := artworksWithImages(search.Data)
		if len(candidates) == 0 {
			lastErr = fmt.Errorf("page %d did not contain artworks with images", page+1)
			continue
		}

		index, err := randomInt(len(candidates))
		if err != nil {
			return Artwork{}, err
		}

		artwork, err := ArtworkByID(candidates[index].ID)
		if err != nil {
			lastErr = err
			continue
		}

		if artwork.ImageURL == "" {
			lastErr = fmt.Errorf("%q did not include an image URL", artwork.Title)
			continue
		}

		return artwork, nil
	}

	return Artwork{}, fmt.Errorf("could not find a renderable artwork: %w", lastErr)
}

func ArtworkByID(id int) (Artwork, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/artworks/%d", articBaseURL, id))
	if err != nil {
		return Artwork{}, fmt.Errorf("error parsing artwork URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("fields", "id,title,artist_title,artist_display,date_display,description,image_id,thumbnail")
	endpoint.RawQuery = query.Encode()

	var response artworkDetailResponse
	if err := getArtworkJSON(endpoint.String(), &response); err != nil {
		return Artwork{}, err
	}

	return artworkFromData(response.Data, response.Config), nil
}

func publicDomainArtworkPages() (int, error) {
	endpoint, err := publicDomainSearchURL(1, 1)
	if err != nil {
		return 0, err
	}

	var response artworkSearchResponse
	if err := getArtworkJSON(endpoint, &response); err != nil {
		return 0, err
	}

	total := response.Pagination.Total
	if total > articMaxSearchResult {
		total = articMaxSearchResult
	}

	return int(math.Ceil(float64(total) / float64(articSearchLimit))), nil
}

func searchPublicDomainArtworks(page int) (artworkSearchResponse, error) {
	endpoint, err := publicDomainSearchURL(page, articSearchLimit)
	if err != nil {
		return artworkSearchResponse{}, err
	}

	var response artworkSearchResponse
	if err := getArtworkJSON(endpoint, &response); err != nil {
		return artworkSearchResponse{}, err
	}

	return response, nil
}

func publicDomainSearchURL(page int, limit int) (string, error) {
	endpoint, err := url.Parse(articBaseURL + "/artworks/search")
	if err != nil {
		return "", fmt.Errorf("error parsing artwork search URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("query[term][is_public_domain]", "true")
	query.Set("fields", "id,title,image_id")
	query.Set("page", fmt.Sprint(page))
	query.Set("limit", fmt.Sprint(limit))
	endpoint.RawQuery = query.Encode()

	return endpoint.String(), nil
}

func getArtworkJSON(endpoint string, target any) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating artwork request: %w", err)
	}

	req.Header.Set("AIC-User-Agent", "DailyDashBoard")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", apiUserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending artwork request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("artwork API error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding artwork response: %w", err)
	}

	return nil
}

func artworkFromData(data artworkData, config artworkConfig) Artwork {
	description := cleanArtworkText(data.Description)
	if description == "" {
		description = cleanArtworkText(data.Thumbnail.AltText)
	}

	if description == "" {
		description = "No description is available for this piece."
	}

	return Artwork{
		ID:               data.ID,
		Artist:           artworkArtist(data),
		Title:            fallbackString(cleanArtworkText(data.Title), "Untitled"),
		Date:             fallbackString(cleanArtworkText(data.DateDisplay), "date unknown"),
		Description:      description,
		ImageURL:         artworkImageURL(config.IIIFURL, data.ImageID),
		MirrorImageURL:   artworkImageURL(articImageMirrorURL, data.ImageID),
		ThumbnailDataURL: data.Thumbnail.LQIP,
	}
}

func artworkArtist(data artworkData) string {
	artist := cleanArtworkText(data.ArtistTitle)
	if artist != "" {
		return artist
	}

	display := strings.TrimSpace(data.ArtistDisplay)
	if display == "" {
		return "Artist unknown"
	}

	lines := strings.Split(display, "\n")
	return fallbackString(cleanArtworkText(lines[0]), "Artist unknown")
}

func artworkImageURL(iiifURL string, imageID string) string {
	if iiifURL == "" || imageID == "" {
		return ""
	}

	return fmt.Sprintf(
		"%s/%s/full/843,/0/default.jpg",
		strings.TrimRight(iiifURL, "/"),
		url.PathEscape(imageID),
	)
}

func artworksWithImages(artworks []artworkData) []artworkData {
	var filtered []artworkData

	for _, artwork := range artworks {
		if artwork.ID != 0 && artwork.ImageID != "" {
			filtered = append(filtered, artwork)
		}
	}

	return filtered
}

func cleanArtworkText(value string) string {
	value = stripHTML(value)
	value = html.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}

func stripHTML(value string) string {
	var builder strings.Builder
	inTag := false

	for _, char := range value {
		switch char {
		case '<':
			builder.WriteRune(' ')
			inTag = true
		case '>':
			builder.WriteRune(' ')
			inTag = false
		default:
			if !inTag {
				builder.WriteRune(char)
			}
		}
	}

	return builder.String()
}

func fallbackString(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func randomInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("random max must be greater than zero")
	}

	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("error generating random number: %w", err)
	}

	return int(value.Int64()), nil
}
