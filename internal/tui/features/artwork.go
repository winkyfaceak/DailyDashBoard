package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"
)

const (
	artworkTextWidth         = 66
	artworkImageCols         = artworkTextWidth
	artworkImageRows         = 18
	artworkFallbackImageCols = 32
	artworkFallbackImageRows = 10
)

// RandomArtwork shows a random public-domain artwork from the Art Institute
// of Chicago.
func RandomArtwork() *Feature {
	return &Feature{
		MenuLabel: "random artwork",
		Loading:   "Fetching artwork...",
		ErrLabel:  "Artwork",
		BoxWidth:  84,
		Fetch: func() (Session, error) {
			artwork, err := daily.RandomArtwork()
			if err != nil {
				return nil, err
			}

			return &ArtworkSession{
				Artwork:    artwork,
				imageState: newImageState(renderArtworkImage(artwork)),
			}, nil
		},
	}
}

// ArtworkSession renders an artwork with its image.
type ArtworkSession struct {
	imageState
	Artwork daily.Artwork
}

func (s *ArtworkSession) View() string {
	return FormatArtwork(s.Artwork, s.imageText("Rendering high-quality image..."))
}

func (s *ArtworkSession) Links() []LinkAction {
	imageURL := PreferredArtworkImageURL(s.Artwork)
	if imageURL == "" {
		return nil
	}

	return []LinkAction{{
		Key:     "o",
		URL:     imageURL,
		Help:    "o to open image",
		Message: "Opening full artwork image...",
	}}
}

func renderArtworkImage(artwork daily.Artwork) terminalimage.Image {
	var lastErr error

	for _, imageURL := range []string{artwork.MirrorImageURL, artwork.ImageURL} {
		if imageURL == "" {
			continue
		}

		renderedImage, err := terminalimage.RenderFromURLWithReferer(imageURL, artworkImageCols, artworkImageRows, "https://www.artic.edu/")
		if err == nil {
			return renderedImage
		}

		lastErr = err
	}

	renderedImage, err := terminalimage.RenderFromDataURL(
		artwork.ThumbnailDataURL,
		artworkFallbackImageCols,
		artworkFallbackImageRows,
	)
	if err == nil {
		return terminalimage.Image{
			Text: renderedImage + "\n" + subtleTextStyle.Render("Low-res API thumbnail; full IIIF image was blocked."),
		}
	}

	if lastErr != nil {
		return terminalimage.Image{Text: subtleTextStyle.Render(fmt.Sprintf("Image unavailable: %v", lastErr))}
	}

	return terminalimage.Image{Text: subtleTextStyle.Render(fmt.Sprintf("Image unavailable: %v", err))}
}

func PreferredArtworkImageURL(artwork daily.Artwork) string {
	if artwork.MirrorImageURL != "" {
		return artwork.MirrorImageURL
	}

	return artwork.ImageURL
}

func FormatArtwork(artwork daily.Artwork, renderedImage string) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(artworkTextWidth).Render(
		fmt.Sprintf("%s -> %s -> %s", artwork.Artist, artwork.Title, artwork.Date),
	))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n")
	}

	builder.WriteString(labelStyle.Render("Description:"))
	builder.WriteString("\n")
	builder.WriteString(textutil.Wrap(artwork.Description, artworkTextWidth))

	return builder.String()
}
