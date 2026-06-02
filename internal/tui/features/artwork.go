package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"DailyDashBoard/internal/tui/textutil"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const (
	artworkTextWidth         = 66
	artworkImageCols         = artworkTextWidth
	artworkImageRows         = 18
	artworkFallbackImageCols = 32
	artworkFallbackImageRows = 10
)

type ArtworkResultMsg struct {
	Artwork daily.Artwork
	Image   terminalimage.Image
	Err     error
}

func FetchArtworkCmd() tea.Cmd {
	return func() tea.Msg {
		artwork, err := daily.RandomArtwork()
		if err != nil {
			return ArtworkResultMsg{Err: err}
		}

		renderedImage := renderArtworkImage(artwork)

		return ArtworkResultMsg{
			Artwork: artwork,
			Image:   renderedImage,
		}
	}
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

	builder.WriteString(artworkHeaderStyle.Render(
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
