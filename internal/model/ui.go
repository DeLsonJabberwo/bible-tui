package model

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (m Model) infoBar() string {
	keys := []string{
		"↑/k     up",
        "↓/j     down",
		"ctrl+k  navigation",
        "?       toggle more info",
        "q       quit",
	}

	barText := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(strings.Join(keys, "  •  "))

	return lipgloss.NewStyle().
		Width(m.viewport.Width()).
		Background(lipgloss.Color("236")).
        Foreground(lipgloss.Color("252")).
        Padding(0, 1).
        Align(lipgloss.Center).
        Render(barText)
}
