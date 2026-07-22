package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"strings"
)

const factTextWidth = 70

// RandomFact shows a random fact from the useless facts API.
func RandomFact() *Feature {
	return &Feature{
		MenuLabel: "random fact",
		Loading:   "Fetching a random fact...",
		ErrLabel:  "Fact",
		BoxWidth:  84,
		Fetch: func() (Session, error) {
			fact, err := daily.RandomFact("en")
			if err != nil {
				return nil, err
			}

			return &FactSession{Fact: fact}, nil
		},
	}
}

// FactSession renders a fetched fact.
type FactSession struct {
	Fact daily.Fact
}

func (s *FactSession) View() string {
	return FormatFact(s.Fact)
}

func FormatFact(fact daily.Fact) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(factTextWidth).Render("Random Fact"))
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
