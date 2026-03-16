package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/delsonjabberwo/bible-tui/internal/bible"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

func main() {
	version, err := bible.LoadVersion("kjv")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}
	content := version.GetBookText(1)

	p := tea.NewProgram(
		model{content: string(content)},
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	content		string
	ready		bool
	viewport	viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		maxWidth := 100
		padding := 4
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height))
			m.viewport.YPosition = 0
			m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
			m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
			m.viewport.Style = m.viewport.Style.
								Margin(0, padding)
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height)
		}
		var wordWidthLimit int
		if m.viewport.Width() < maxWidth + padding * 2 {
			wordWidthLimit = (m.viewport.Width() - padding * 2)
		} else {
			wordWidthLimit = maxWidth
		}
		m.viewport.SetContent(wrap.String(wordwrap.String(m.content, wordWidthLimit), m.viewport.Width() - padding * 2))
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	if !m.ready {
		v.SetContent("\n  Initializing...")
	} else {
		v.SetContent(fmt.Sprintf("\n%s\n", m.viewport.View()))
	}
	return v
}

