package model

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

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
	ti.SetWidth(20)

	return navModel{
		query: ti,
		all: references,
	}
}

func (n *navModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+k", "ctrl+c":
			return closeNavCmd()

		case "enter":
			if len(n.matches) > 0 {
				return selectVerseCmd(n.matches[n.cursor])
			}

		case "up", "ctrl+n":
			if n.cursor > 0 {
				n.cursor--
			} else if n.cursor == 0 {
				n.cursor = len(n.matches) - 1
			}
			return nil

		case "down", "ctrl+p":
			if n.cursor < len(n.matches) - 1 {
				n.cursor++
			} else if n.cursor == len(n.matches) - 1 {
				n.cursor = 0
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
	
	if n.query.Value() == "" {
		b.WriteString("\n\n\n\n\n\n\n\n\n\n\n")
		return b.String()
	}

	if len(n.matches) == 0 && n.query.Value() != "" {
		b.WriteString("\n\n\tNo matches\n\n\n\n\n\n\n\n")
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

func (n *navModel) filter(query string) {
	query = strings.ToLower(query)
	n.matches = nil
	n.cursor = 0

	if query == "" {
		return
	}

	var results []result
	for _, v := range n.all {
		match := fuzzyMatch(query, &v)
		if match != nil {
			results = append(results, *match)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	for _, res := range results {
		n.matches = append(n.matches, *res.reference)
		if len(n.matches) >= 10 {
			break
		}
	}
}

type result struct {
	reference	*Reference
	score		int
}
func fuzzyMatch(pattern string, targetRef *Reference) *result {
	if len(pattern) == 0 {
		return &result{reference: targetRef, score: 0}
	}

	pInd, tInd := 0, 0
	matches := []int{}
	score := 0
	prevMatchInd := -1
	target := strings.ToLower(targetRef.Book) + " " +
				strconv.Itoa(targetRef.Chapter) + ":" +
				strconv.Itoa(targetRef.Verse)

	for tInd < len(target) && pInd < len(pattern) {
		pChar := pattern[pInd]
		tChar := target[tInd]

		if byte(unicode.ToLower(rune(pChar))) == byte(unicode.ToLower(rune(tChar))) ||
			isSeparator(pChar) && isSeparator(tChar) {
			matches = append(matches, tInd)
			score += 10

			if prevMatchInd != -1 && tInd == prevMatchInd+1 {
				score += 15
			}

			if tInd == 0 || isSeparator(target[tInd-1]) {
				score += 20
			}

			if tInd > 0 && unicode.IsLower(rune(target[tInd-1])) && unicode.IsUpper(rune(tChar)) {
				score += 20
			}
			prevMatchInd = tInd
			pInd++
		} else if pInd > 0 {
			score -= 3
		}
		tInd++
	}

	if pInd < len(pattern)-1 {
		return nil
	}

	if len(matches) > 0 {
		score -= matches[0] * 2
		tail := len(target) - matches[len(matches)-1] - 1
		score -= tail
	}

	return &result{
		reference: targetRef,
		score: score,
	}

}

func isSeparator(c byte) bool {
	return c == ' ' || c == '-' || c == '_' || c == '/' || c == '\\' || c == ':'
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

