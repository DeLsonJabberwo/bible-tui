package model

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type navModel struct {
	query		textinput.Model
	matches		[]Reference
	all			[]Reference
	cursor		int
}

func NewNavModel(references []Reference) navModel {
	ti := textinput.New()
	ti.Placeholder = "John 1:1"
	ti.Focus()

	return navModel{
		query: ti,
		all: references,
	}
}

func (n *navModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+k":
			return closeNavCmd()

		case "enter":
			if len(n.matches) > 0 {
				return selectVerseCmd(n.matches[n.cursor])
			}

		case "up", "ctrl+n":
			if n.cursor > 0 {
				n.cursor--
			}
			return nil

		case "down", "ctrl+p":
			if n.cursor < len(n.matches) - 1 {
				n.cursor++
			}
			return nil

		default:
			var cmd tea.Cmd
			n.query, cmd = n.query.Update(msg)
			n.filter(n.query.Value())
			return cmd

		}
	}
	return nil
}

func (n navModel) View() string {
	var b strings.Builder

	b.WriteString(n.query.View())
	b.WriteString("\n\n")

	if len(n.matches) == 0 && n.query.Value() != "" {
		b.WriteString("No matches")
		return b.String()
	}

	for i, m := range n.matches {
		line := fmt.Sprintf("%s %d:%d", m.Book, m.Chapter, m.Verse)

		if i == n.cursor {
			b.WriteString("> " + lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(line))
		} else {
			b.WriteString("  " + line)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (n navModel) filter(query string) {
	query = strings.ToLower(query)
	n.matches = nil
	n.cursor = 0

	if query == "" {
		return
	}

	for _, v := range n.all {
		match := v.Book + " " +
				 strconv.Itoa(v.Chapter) + ":" +
				 strconv.Itoa(v.Verse)

		if strings.HasPrefix(match, query) ||
			strings.Contains(match, query) {
			n.matches = append(n.matches, v)
		}

		if len(n.matches) >= 10 {
			break
		}
	}
}

type CloseNavMsg struct {}
type SelectVerseMsg struct {
	Ref Reference
}

func closeNavCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseNavMsg{}
	}
}

func selectVerseCmd(ref Reference) tea.Cmd {
	return func() tea.Msg {
		return SelectVerseMsg{Ref: ref}
	}
}

