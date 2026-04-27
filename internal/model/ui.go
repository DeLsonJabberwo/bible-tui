package model

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/delsonjabberwo/bible-tui/internal/buffer"
)

func (m Model) infoBar(width int) string {
	keys := []string{
		"ctrl+k  navigation",
		"ctrl+s  search",
		"ctrl+v  version",
		"q  quit",
	}

	contentWidth := max(min(width-2*buffer.PADDING, buffer.MAX_WIDTH), 0)

	barWidth := contentWidth + 2*buffer.PADDING

	style := lipgloss.NewStyle().
		Width(barWidth).
		Background(lipgloss.Color("#1e2d3d")).
		Foreground(lipgloss.Color("252")).
		Align(lipgloss.Center)

	row1 := style.Render(strings.Join(keys, "  •  "))
	row2 := style.Render("")

	return lipgloss.JoinVertical(0, row1, row2)
}
