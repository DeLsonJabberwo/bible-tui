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
	"github.com/delsonjabberwo/bible-tui/internal/bible"
)

type searchModel struct {
	query		textinput.Model
	matches		[]bible.Verse
	all			[]bible.Verse
	cursor		int
}

func NewSearchModel(verses []bible.Verse) searchModel {
	ti := textinput.New()
	ti.Placeholder = "John 1:1"
	ti.Focus()

	return searchModel{
		query: ti,
		all: verses,
	}
}

func (n *searchModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+s", "ctrl+c":
			return closeSearchCmd()

		case "enter":
			if len(n.matches) > 0 {
				return selectVerseCmd(NewReference(n.matches[n.cursor]))
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

func (n searchModel) View() string {
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
		line := fmt.Sprintf("%s %d:%d\t%s", m.BookName, m.Chapter, m.Verse, m.Text)

		if i == n.cursor {
			b.WriteString("> " + lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(line))
		} else {
			lineMax := min(len(line), 30)
			b.WriteString("  " + line[:lineMax] + "...")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (n *searchModel) filter(query string) {
	query = strings.ToLower(query)
	n.matches = nil
	n.cursor = 0

	if query == "" {
		return
	}

	var results []searchResult
	for _, v := range n.all {
		match := fuzzyMatchSearch(query, &v)
		if match != nil {
			results = append(results, *match)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	for _, res := range results {
		n.matches = append(n.matches, *res.verse)
		if len(n.matches) >= 10 {
			break
		}
	}
}

type searchResult struct {
	verse		*bible.Verse
	score		int
}
func fuzzyMatchSearch(pattern string, targetVerse *bible.Verse) *searchResult {
	if len(pattern) == 0 {
		return &searchResult{verse: targetVerse, score: 0}
	}

	pInd, tInd := 0, 0
	matches := []int{}
	score := 0
	prevMatchInd := -1
	target := strings.ToLower(targetVerse.BookName) + " " +
				strconv.Itoa(targetVerse.Chapter) + ":" +
				strconv.Itoa(targetVerse.Verse) + " " +
				strings.ToLower(targetVerse.Text)

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

	return &searchResult{
		verse: targetVerse,
		score: score,
	}

}

type CloseSearchMsg struct {}

func closeSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseSearchMsg{}
	}
}

