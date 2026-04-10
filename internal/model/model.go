package model

import (
	"log"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/delsonjabberwo/bible-tui/internal/buffer"
)

type Model struct {
	Buffer     buffer.Buffer
	ready      bool
	viewport   viewport.Model
	appending  bool
	prepending bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	start := time.Now()

	inSecondBook := m.Buffer.GetBookFromLine(m.viewport.YOffset()) <= m.Buffer.Books[1]
	inFourthBook := m.Buffer.GetBookFromLine(m.viewport.YOffset()) >= m.Buffer.Books[3]
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}

		switch msg.String() {
		case "up", "k", "pgup":
			if !m.prepending && inSecondBook && m.Buffer.Books[0] != 1 {
				m.prepending = true
				log.Printf("Prepending book:\n")
				log.Printf("\tOld Books: %v\n", m.Buffer.Books)
				m.Buffer.ShiftBooksPrev()
				log.Printf("\tNew Books: %v\n", m.Buffer.Books)
				m.viewport.SetYOffset(m.Buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset()))
			}
		case "down", "j", "pgdown":
			if !m.appending && inFourthBook && m.Buffer.Books[len(m.Buffer.Books)-1] != 66 {
				m.appending = true
				log.Printf("Appending book:\n")
				log.Printf("\tOld Books: %v\n", m.Buffer.Books)
				m.Buffer.ShiftBooksNext()
				log.Printf("\tNew Books: %v\n", m.Buffer.Books)
				m.viewport.SetYOffset(m.Buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset()))
			}
		case "ctrl+k", ":":
			// TODO: add location input box
		}

		m.viewport.SetContent(m.Buffer.Content)
	case tea.MouseWheelMsg:
		switch msg.Mouse().Button {
		case tea.MouseWheelUp:
			if !m.prepending && inSecondBook && m.Buffer.Books[0] != 1 {
				m.prepending = true
				log.Printf("Prepending book:\n")
				log.Printf("\tOld Books: %v\n", m.Buffer.Books)
				m.Buffer.ShiftBooksPrev()
				log.Printf("\tNew Books: %v\n", m.Buffer.Books)
				m.viewport.SetYOffset(m.Buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset()))
			}
		case tea.MouseWheelDown:
			if !m.appending && inFourthBook && m.Buffer.Books[len(m.Buffer.Books)-1] != 66 {
				m.appending = true
				log.Printf("Appending book:\n")
				log.Printf("\tOld Books: %v\n", m.Buffer.Books)
				m.Buffer.ShiftBooksNext()
				log.Printf("\tNew Books: %v\n", m.Buffer.Books)
				m.viewport.SetYOffset(m.Buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset()))
			}
		}

		m.viewport.SetContent(m.Buffer.Content)
	case tea.WindowSizeMsg:
		padding := buffer.PADDING
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height))
			m.viewport.YPosition = 0
			m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
			m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
			m.viewport.Style = m.viewport.Style.Margin(0, padding, 2)

			m.Buffer.UpdateBuffer(buffer.NewViewportInfo(m.viewport.Width()), m.viewport.YOffset())
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height)
			viewportInfo := buffer.NewViewportInfo(m.viewport.Width())
			//log.Printf("Old Offset: %d\n", m.viewport.YOffset())
			m.viewport.SetYOffset(m.Buffer.UpdateBuffer(viewportInfo, m.viewport.YOffset()))
			//log.Printf("New Offset: -> %d\n", m.viewport.YOffset())
		}
		m.viewport.SetContent(m.Buffer.Content)
	}

	if m.appending {
		m.appending = false
	}
	if m.prepending {
		m.prepending = false
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	duration := time.Since(start)
	log.Printf("Update Time: %s\n", duration)
	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	if !m.ready {
		v.SetContent("\n  Initializing...")
		return v
	}

	content := m.viewport.View()
	//infoBar := m.infoBar()
	fullContent := content
	/*
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		infoBar,
	)
	*/

	v.SetContent(fullContent)
	return v
}
