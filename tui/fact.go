package tui

import (
	"DailyDashBoard/dailys/api"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const factTextWidth = 70

type factResultMsg struct {
	Fact api.Fact
	Err  error
}

func fetchFactCmd() tea.Cmd {
	return func() tea.Msg {
		fact, err := api.RandomFact("en")

		return factResultMsg{
			Fact: fact,
			Err:  err,
		}
	}
}

func formatFact(fact api.Fact) string {
	var builder strings.Builder

	builder.WriteString(factHeaderStyle.Render("Random Fact"))
	builder.WriteString("\n\n")
	builder.WriteString(wrapText(fact.Text, factTextWidth))

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

func formatFactMeta(fact api.Fact) string {
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
