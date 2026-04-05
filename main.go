package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/delsonjabberwo/bible-tui/internal/buffer"
)

func main() {
	if os.Getenv("DEBUG") == "1" {
		f, err := tea.LogToFile("tmp/debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		f, err := tea.LogToFile("/dev/null", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	viewportInfo := buffer.NewViewportInfo(0)
	buffer, err := buffer.NewBuffer(viewportInfo, "kjv", 1)
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		model{buffer: buffer},
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	buffer   buffer.Buffer
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	//start := time.Now()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		padding := buffer.PADDING
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height))
			m.viewport.YPosition = 0
			m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
			m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
			m.viewport.Style = m.viewport.Style.
				Margin(0, padding)
			m.buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset())
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height)
			viewportInfo := buffer.NewViewportInfo(m.viewport.Width())
			//log.Printf("Old Offset: %d\n", m.viewport.YOffset())
			m.viewport.SetYOffset(m.buffer.UpdateBuffer(viewportInfo, m.viewport.YOffset()))
			//log.Printf("New Offset: -> %d\n", m.viewport.YOffset())
		}
		m.viewport.SetContent(m.buffer.Content)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	//duration := time.Since(start)
	//log.Printf("Update Time: %s\n", duration)
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
