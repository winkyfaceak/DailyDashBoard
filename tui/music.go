package tui

import (
	"DailyDashBoard/dailys/api"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const (
	musicTextWidth = 66
	musicImageCols = musicTextWidth
	musicImageRows = 18
)

type musicResultMsg struct {
	Pick  api.MusicPick
	Image terminalImage
	Err   error
}

type musicImageReadyMsg struct{}

func fetchMusicCmd() tea.Cmd {
	return func() tea.Msg {
		pick, err := api.RandomMusic()
		if err != nil {
			return musicResultMsg{Err: err}
		}

		return musicResultMsg{
			Pick:  pick,
			Image: renderMusicImage(pick),
		}
	}
}

func renderMusicImage(pick api.MusicPick) terminalImage {
	if pick.ArtworkURL == "" {
		return terminalImage{Text: subtleTextStyle.Render("Album art unavailable.")}
	}

	renderedImage, err := renderTerminalImageFromURL(pick.ArtworkURL, musicImageCols, musicImageRows)
	if err != nil {
		return terminalImage{Text: subtleTextStyle.Render(fmt.Sprintf("Album art unavailable: %v", err))}
	}

	return renderedImage
}

func formatMusic(pick api.MusicPick, renderedImage string) string {
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

func formatMusicMeta(pick api.MusicPick) string {
	var parts []string

	if pick.PrimaryGenre != "" && !strings.EqualFold(pick.PrimaryGenre, pick.Genre) {
		parts = append(parts, pick.PrimaryGenre)
	}

	if !pick.ReleaseDate.IsZero() {
		parts = append(parts, pick.ReleaseDate.Format("2006"))
	}

	return strings.Join(parts, " • ")
}
