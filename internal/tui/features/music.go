package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	musicTextWidth = 66
	musicImageCols = musicTextWidth
	musicImageRows = 18
)

// RandomMusic picks a random genre and song from iTunes.
func RandomMusic() *Feature {
	return &Feature{
		MenuLabel:       "random music",
		Loading:         "Finding a genre and song...",
		ErrLabel:        "Music",
		BoxWidth:        84,
		Accent:          lipgloss.Color("#3FB68C"),
		BackgroundWords: []string{"melody", "tempo", "chord", "♪", "♫"},
		Fetch: func() (Session, error) {
			pick, err := daily.RandomMusic()
			if err != nil {
				return nil, err
			}

			return &MusicSession{
				Pick:       pick,
				imageState: newImageState(renderMusicImage(pick)),
			}, nil
		},
	}
}

// MusicSession renders a music pick with album art and preview audio.
type MusicSession struct {
	imageState
	Pick daily.MusicPick
}

func (s *MusicSession) View() string {
	return FormatMusic(s.Pick, s.imageText("Rendering album art..."))
}

func (s *MusicSession) Audio() Audio {
	return Audio{
		URL:           s.Pick.PreviewURL,
		Label:         "preview",
		AllowDownload: true,
	}
}

func (s *MusicSession) Links() []LinkAction {
	var links []LinkAction

	if s.Pick.ArtworkURL != "" {
		links = append(links, LinkAction{
			Key:     "o",
			URL:     s.Pick.ArtworkURL,
			Help:    "o to open album art",
			Message: "Opening album art...",
		})
	}

	if s.Pick.TrackURL != "" {
		links = append(links, LinkAction{
			Key:     "s",
			URL:     s.Pick.TrackURL,
			Help:    "s to open song",
			Message: "Opening song...",
		})
	}

	return links
}

func renderMusicImage(pick daily.MusicPick) terminalimage.Image {
	if pick.ArtworkURL == "" {
		return terminalimage.Image{Text: subtleTextStyle.Render("Album art unavailable.")}
	}

	renderedImage, err := terminalimage.RenderFromURL(pick.ArtworkURL, musicImageCols, musicImageRows)
	if err != nil {
		return terminalimage.Image{Text: subtleTextStyle.Render(fmt.Sprintf("Album art unavailable: %v", err))}
	}

	return renderedImage
}

func FormatMusic(pick daily.MusicPick, renderedImage string) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(musicTextWidth).Render(fmt.Sprintf("Random Music -> %s", pick.Genre)))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n")
	}

	builder.WriteString(labelStyle.Render("Song: "))
	builder.WriteString(strongTextStyle.Render(pick.TrackName))
	builder.WriteString("\n")

	builder.WriteString(labelStyle.Render("Artist: "))
	builder.WriteString(pick.ArtistName)
	builder.WriteString("\n")

	if pick.AlbumName != "" {
		builder.WriteString(labelStyle.Render("Album: "))
		builder.WriteString(pick.AlbumName)
		builder.WriteString("\n")
	}

	meta := formatMusicMeta(pick)
	if meta != "" {
		builder.WriteString(subtleTextStyle.Render(meta))
		builder.WriteString("\n")
	}

	if pick.PreviewURL != "" {
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(pick.PreviewURL).Render("press p to play/pause preview"))
		builder.WriteString("\n")
	}

	if pick.TrackURL != "" {
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(pick.TrackURL).Render("press s to open song"))
		builder.WriteString("\n")
	}

	builder.WriteString(subtleTextStyle.Render("Preview and artwork provided courtesy of iTunes."))

	return builder.String()
}

func formatMusicMeta(pick daily.MusicPick) string {
	var parts []string

	if pick.PrimaryGenre != "" && !strings.EqualFold(pick.PrimaryGenre, pick.Genre) {
		parts = append(parts, pick.PrimaryGenre)
	}

	if !pick.ReleaseDate.IsZero() {
		parts = append(parts, pick.ReleaseDate.Format("2006"))
	}

	return strings.Join(parts, " • ")
}
