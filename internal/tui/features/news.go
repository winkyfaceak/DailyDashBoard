package features

import (
	"DailyDashBoard/internal/daily"
	"DailyDashBoard/internal/tui/textutil"
	"DailyDashBoard/internal/tui/theme"
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

const (
	newsItemLimit = 8
	newsTextWidth = 78
)

// DailyNews shows top stories from the configured RSS feeds.
func DailyNews() *Feature {
	return &Feature{
		MenuLabel:       "daily news",
		Loading:         "Fetching daily news...",
		ErrLabel:        "News",
		BoxWidth:        92,
		Accent:          lipgloss.Color("#6C90C4"),
		BackgroundWords: []string{"headline", "breaking", "bulletin", "wire", "§"},
		Fetch: func() (Session, error) {
			items, err := daily.DailyNews(daily.DefaultNewsFeeds(), newsItemLimit)
			if err != nil {
				return nil, err
			}

			if len(items) == 0 {
				return nil, Notice("No news stories found.")
			}

			return &NewsSession{Items: items}, nil
		},
	}
}

// NewsSession renders a numbered list of news stories.
type NewsSession struct {
	Items []daily.NewsItem
}

func (s *NewsSession) View() string {
	return FormatNewsItems(s.Items)
}

func (s *NewsSession) StoryLinks() []string {
	links := make([]string, len(s.Items))
	for i, item := range s.Items {
		links[i] = item.Link
	}

	return links
}

func FormatNewsItems(items []daily.NewsItem) string {
	var builder strings.Builder

	builder.WriteString(theme.Header(newsTextWidth).Render("Daily News"))
	builder.WriteString("\n\n")

	for i, item := range items {
		if i > 0 {
			builder.WriteString("\n\n")
		}

		number := strconv.Itoa(i + 1)
		builder.WriteString(formatNewsHeadline(number, item.Title))

		meta := formatNewsMeta(item)
		if meta != "" {
			builder.WriteString("\n")
			builder.WriteString(subtleTextStyle.Render("   " + meta))
		}

		if item.Summary != "" {
			builder.WriteString("\n")
			builder.WriteString(exampleStyle.Render(textutil.IndentWrapped(textutil.Truncate(item.Summary, 180), newsTextWidth, "   ")))
		}

		if item.Link != "" {
			builder.WriteString("\n")
			builder.WriteString(exampleStyle.Underline(true).Hyperlink(item.Link).Render("   press " + number + " to open story"))
		}
	}

	return builder.String()
}

// NewsLinkIndex maps the keys 1-9 to a story index.
func NewsLinkIndex(key string) (int, bool) {
	if len(key) != 1 || key[0] < '1' || key[0] > '9' {
		return 0, false
	}

	return int(key[0] - '1'), true
}

func formatNewsHeadline(number string, title string) string {
	lines := strings.Split(textutil.Wrap(title, newsTextWidth-3), "\n")
	if len(lines) == 0 {
		return ""
	}

	for i, line := range lines {
		if i == 0 {
			lines[i] = fmt.Sprintf("%s. %s", number, line)
			continue
		}

		lines[i] = "   " + line
	}

	return strongTextStyle.Render(strings.Join(lines, "\n"))
}

func formatNewsMeta(item daily.NewsItem) string {
	var parts []string

	if item.Source != "" {
		parts = append(parts, item.Source)
	}

	if !item.PublishedAt.IsZero() {
		parts = append(parts, formatNewsTime(item.PublishedAt))
	}

	return strings.Join(parts, " • ")
}

func formatNewsTime(publishedAt time.Time) string {
	localTime := publishedAt.Local()
	now := time.Now()

	if sameDay(localTime, now) {
		return localTime.Format("15:04")
	}

	return localTime.Format("2 Jan 15:04")
}

func sameDay(left time.Time, right time.Time) bool {
	leftYear, leftMonth, leftDay := left.Date()
	rightYear, rightMonth, rightDay := right.Date()

	return leftYear == rightYear && leftMonth == rightMonth && leftDay == rightDay
}
