package textutil

import "strings"

func Wrap(value string, width int) string {
	words := strings.Fields(value)
	if len(words) == 0 {
		return ""
	}

	var builder strings.Builder
	lineLength := 0

	for _, word := range words {
		if lineLength == 0 {
			builder.WriteString(word)
			lineLength = len(word)
			continue
		}

		if lineLength+1+len(word) > width {
			builder.WriteString("\n")
			builder.WriteString(word)
			lineLength = len(word)
			continue
		}

		builder.WriteString(" ")
		builder.WriteString(word)
		lineLength += 1 + len(word)
	}

	return builder.String()
}

func IndentWrapped(value string, width int, indent string) string {
	wrapped := Wrap(value, width-len(indent))
	if wrapped == "" {
		return ""
	}

	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}

	return strings.Join(lines, "\n")
}

func Truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}

	trimmed := strings.TrimSpace(string(runes[:limit]))
	trimmed = strings.TrimRight(trimmed, ".,;:")
	if trimmed == "" {
		return string(runes[:limit])
	}

	return trimmed + "..."
}

func LimitStrings(items []string, limit int) []string {
	if len(items) <= limit {
		return items
	}

	return items[:limit]
}
