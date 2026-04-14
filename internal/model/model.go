package model

import (
	"log"
	"regexp"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/delsonjabberwo/bible-tui/internal/bible"
	"github.com/delsonjabberwo/bible-tui/internal/buffer"
)

type Model struct {
	Buffer		buffer.Buffer
	ready		bool
	viewport	viewport.Model
	appending	bool
	prepending	bool
	nav			*navModel
	references	[]Reference
}

func (m *Model) Init() tea.Cmd {
	m.references = ReferencesFromVersion(m.Buffer.Version)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	start := time.Now()

	if m.nav != nil {
		cmd := m.nav.Update(msg)
		// check for exit or navigation
		switch msg := msg.(type) {
			
		case CloseNavMsg:
			m.nav = nil
			return &m, nil

		case SelectVerseMsg:
			m.nav = nil
			verseInfo := bible.VerseInfo{
				Book: msg.Ref.BookInd,
				Chapter: msg.Ref.Chapter,
				Verse: msg.Ref.Verse,
			}
			var err error
			m.Buffer, err = buffer.NewBuffer(m.Buffer.LastViewportInfo, 
										strings.ToLower(m.Buffer.Version.Metadata.ShortName), 
										verseInfo.Book)
			if err != nil {
				return nil, nil
			}
			m.viewport.SetContent(m.Buffer.Content)
			m.viewport.SetYOffset(m.Buffer.VerseLocs.Verses[verseInfo])
			return &m, nil
		}

		return &m, cmd
	}

	inSecondBook := m.Buffer.GetBookFromLine(m.viewport.YOffset()) <= m.Buffer.Books[1]
	inFourthBook := m.Buffer.GetBookFromLine(m.viewport.YOffset()) >= m.Buffer.Books[3]
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return &m, tea.Quit
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
			if m.nav == nil {
				newNav := NewNavModel(m.references)
				m.nav = &newNav
			}

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
	return &m, tea.Batch(cmds...)
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

	if m.nav == nil {
		v.SetContent(content)
		return v
	}
	
	dimmedContent := lipgloss.NewStyle().
		Faint(true).
        Render(content)

	navContent := m.nav.View()
	styledNav := lipgloss.NewStyle().
        Width(50).
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Render(navContent)

	finalContent := overlayStrings(
		dimmedContent, 
		styledNav, 
		34,
		len(strings.Split(styledNav, "\n")), 
	)

	v.SetContent(finalContent)

	return v
}

func overlayStrings(base string, overlay string, overlayWidth int, overlayHeight int) string {
    baseLines := strings.Split(base, "\n")
    overlayLines := strings.Split(overlay, "\n")
    
    // Find center position
	visibleWidth := len(stripAnsi(baseLines[0]))
    startY := (len(baseLines) - overlayHeight) / 2
    startX := (visibleWidth - overlayWidth) / 2
    
    for i, line := range overlayLines {
        if startY+i >= 0 && startY+i < len(baseLines) {
            baseLine := baseLines[startY+i]
			visibleBaseWidth := len(stripAnsi(baseLine))

			beforeEnd := visibleToByteIndex(baseLine, startX)
            afterStart := visibleToByteIndex(baseLine, startX+overlayWidth)

            before := ""
            if startX > 0 && startX <= visibleBaseWidth {
                before = baseLine[:beforeEnd]
            }
            after := ""
            if startX+overlayWidth < visibleBaseWidth {
                after = baseLine[afterStart:]
            }
            baseLines[startY+i] = before + line + after
        }
    }
    return strings.Join(baseLines, "\n")
}

func stripAnsi(str string) string {
    re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
    return re.ReplaceAllString(str, "")
}

func visibleToByteIndex(str string, visiblePos int) int {
    ansiRe := regexp.MustCompile(`\x1b\[[0-9;]*m`)
    byteIdx := 0
    visibleCount := 0
    
    for visibleCount < visiblePos && byteIdx < len(str) {
        // Check if current position is start of ANSI sequence
        loc := ansiRe.FindStringIndex(str[byteIdx:])
        if loc != nil && loc[0] == 0 {
            byteIdx += loc[1] // Skip entire ANSI sequence
        } else {
            byteIdx++
            visibleCount++
        }
    }
    
    return byteIdx
}

