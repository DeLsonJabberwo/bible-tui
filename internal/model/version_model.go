package model

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type verModel struct {
	list		[]string
	cursor		int
}

func NewVerModel() verModel {
	var model verModel
	for _, v := range versions {
		model.list = append(model.list, v)
	}
	return model
}

func (v *verModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+v", "ctrl+c":
			return closeVerCmd()

		case "enter":
			return selectVersionCmd(v.list[v.cursor])

		case "k", "up":
			if v.cursor > 0 {
				v.cursor--
			}

		case "j", "down":
			if v.cursor < len(v.list) - 1 {
				v.cursor++
			}

		}
	}
	return nil
}

func (v verModel) View() string {
	var b strings.Builder

	if len(v.list) == 0 {
		b.WriteString("\n\n\tNo versions available.\n\n\n\n\n\n\n\n")
		return b.String()
	}

	for i, m := range v.list {
		line := m

		if i == v.cursor {
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

type CloseVerMsg struct {}
type SelectVersionMsg struct {
	Code string
}

func closeVerCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseVerMsg{}
	}
}

func selectVersionCmd(code string) tea.Cmd {
	return func() tea.Msg {
		return SelectVersionMsg{Code: code}
	}
}

var versions []string = []string{
	"kjv",
	"geneva",
	"asv",
	"web",
}

