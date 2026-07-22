package features

import (
	tea "charm.land/bubbletea/v2"
)

func keyPress(value string) tea.KeyPressMsg {
	if value == "enter" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	}

	if value == "backspace" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyBackspace})
	}

	if value == "esc" {
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape})
	}

	return tea.KeyPressMsg(tea.Key{Text: value, Code: []rune(value)[0]})
}
