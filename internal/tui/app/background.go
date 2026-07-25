package app

import (
	"hash/fnv"
	"image/color"
	"math/rand"
	"strings"

	"charm.land/lipgloss/v2"
)

// backgroundWordDensity controls how sparse the scattered background is:
// roughly one word placed per this many cells of terminal area.
const backgroundWordDensity = 70

const (
	minBackgroundWords = 4
	maxBackgroundWords = 60
	placementAttempts  = 6
)

// composeBackground centers content within a width x height canvas and
// layers it on top of a faint, feature-themed scatter of words, so the
// background only shows through the margins around the (opaque) content.
func composeBackground(width, height int, content string, words []string, accent color.Color) string {
	bg := buildBackground(width, height, words, accent)
	if bg == "" {
		return content
	}

	x := (width - lipgloss.Width(content)) / 2
	if x < 0 {
		x = 0
	}

	y := (height - lipgloss.Height(content)) / 2
	if y < 0 {
		y = 0
	}

	return lipgloss.NewCompositor(
		lipgloss.NewLayer(bg).Z(0),
		lipgloss.NewLayer(content).X(x).Y(y).Z(1),
	).Render()
}

// buildBackground scatters words/glyphs across a width x height canvas,
// styled as a faint watermark derived from accent. Placement is seeded from
// the size and word bank, so the same feature at the same terminal size
// always renders the same layout instead of jittering between frames.
func buildBackground(width, height int, words []string, accent color.Color) string {
	if width <= 0 || height <= 0 || len(words) == 0 {
		return ""
	}

	grid := make([][]rune, height)
	for row := range grid {
		grid[row] = make([]rune, width)
		for col := range grid[row] {
			grid[row][col] = ' '
		}
	}

	rng := rand.New(rand.NewSource(backgroundSeed(width, height, words)))

	count := (width * height) / backgroundWordDensity
	if count < minBackgroundWords {
		count = minBackgroundWords
	}
	if count > maxBackgroundWords {
		count = maxBackgroundWords
	}

	for i := 0; i < count; i++ {
		word := []rune(words[rng.Intn(len(words))])
		if len(word) > width-2 {
			continue
		}

		for attempt := 0; attempt < placementAttempts; attempt++ {
			row := rng.Intn(height)
			col := rng.Intn(width - len(word) + 1)

			if !rowRangeBlank(grid[row], col, len(word)) {
				continue
			}

			copy(grid[row][col:col+len(word)], word)
			break
		}
	}

	lines := make([]string, height)
	for row := range grid {
		lines[row] = string(grid[row])
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Darken(accent, 0.55)).Faint(true)

	return style.Render(strings.Join(lines, "\n"))
}

func rowRangeBlank(row []rune, col, length int) bool {
	for i := col; i < col+length; i++ {
		if row[i] != ' ' {
			return false
		}
	}

	return true
}

// backgroundSeed derives a deterministic seed from the canvas size and word
// bank so identical (size, words) pairs always scatter the same way.
func backgroundSeed(width, height int, words []string) int64 {
	h := fnv.New64a()
	h.Write([]byte(strings.Join(words, ",")))

	return int64(h.Sum64()) + int64(width)*1000003 + int64(height)
}
