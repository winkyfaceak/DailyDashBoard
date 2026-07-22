package features

// appendFiltered appends the runes of text that pass keep to current,
// capping the result at maxRunes. It is the shared building block for the
// hand-typed inputs (math answers, music trivia guesses).
func appendFiltered(current string, text string, maxRunes int, keep func(rune) bool) string {
	for _, char := range text {
		if !keep(char) {
			continue
		}

		if len([]rune(current)) >= maxRunes {
			break
		}

		current += string(char)
	}

	return current
}

// trimLastRune removes the final rune from value.
func trimLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}

	return string(runes[:len(runes)-1])
}
