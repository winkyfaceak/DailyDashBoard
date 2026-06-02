package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/terminalimage"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const (
	musicTextWidth = 66
	musicImageCols = musicTextWidth
	musicImageRows = 18
)

type MusicResultMsg struct {
	Pick  daily.MusicPick
	Image terminalimage.Image
	Err   error
}

func FetchMusicCmd() tea.Cmd {
	return func() tea.Msg {
		pick, err := daily.RandomMusic()
		if err != nil {
			return MusicResultMsg{Err: err}
		}

		return MusicResultMsg{
			Pick:  pick,
			Image: renderMusicImage(pick),
		}
	}
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

	builder.WriteString(musicHeaderStyle.Render(fmt.Sprintf("Random Music -> %s", pick.Genre)))
	builder.WriteString("\n\n")

	if renderedImage != "" {
		builder.WriteString(renderedImage)
		builder.WriteString("\n")
	}

	builder.WriteString(labelStyle.Render("Song: "))
	builder.WriteString(musicTrackStyle.Render(pick.TrackName))
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
