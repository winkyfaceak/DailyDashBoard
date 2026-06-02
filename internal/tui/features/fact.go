package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const factTextWidth = 70

type FactResultMsg struct {
	Fact daily.Fact
	Err  error
}

func FetchFactCmd() tea.Cmd {
	return func() tea.Msg {
		fact, err := daily.RandomFact("en")

		return FactResultMsg{
			Fact: fact,
			Err:  err,
		}
	}
}

func FormatFact(fact daily.Fact) string {
	var builder strings.Builder

	builder.WriteString(factHeaderStyle.Render("Random Fact"))
	builder.WriteString("\n\n")
	builder.WriteString(textutil.Wrap(fact.Text, factTextWidth))

	meta := formatFactMeta(fact)
	if meta != "" {
		builder.WriteString("\n\n")
		builder.WriteString(subtleTextStyle.Render(meta))
	}

	if fact.Permalink != "" {
		builder.WriteString("\n")
		builder.WriteString(exampleStyle.Underline(true).Hyperlink(fact.Permalink).Render("view fact source"))
	}

	return builder.String()
}

func formatFactMeta(fact daily.Fact) string {
	var parts []string

	if fact.Source != "" {
		parts = append(parts, fact.Source)
	}

	if fact.Language != "" {
		parts = append(parts, strings.ToUpper(fact.Language))
	}

	if fact.SourceURL != "" {
		parts = append(parts, fact.SourceURL)
	}

	return strings.Join(parts, " • ")
}
