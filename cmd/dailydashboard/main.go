package main

import (
	"DailyDashBoard/internal/tui/app"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(app.InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
