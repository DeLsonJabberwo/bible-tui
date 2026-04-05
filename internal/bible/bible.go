package bible

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

type Verse struct {
	BookName	string		`json:"book_name"`
	Book		int			`json:"book"`
	Chapter		int			`json:"chapter"`
	Verse		int			`json:"verse"`
	Text		string		`json:"text"`
}

type Metadata struct {
	Name		string		`json:"name"`
	ShortName	string		`json:"shortname"`
	Year		string		`json:"year"`
	Publisher	string		`json:"publisher"`
	Owner		string		`json:"owner"`
	Description	string		`json:"description"`
	Lang		string		`json:"lang"`
	LangShort	string		`json:"lang_short"`
	Copyright	int			`json:"copyright"`
	CopyrightStatement	string	`json:"copyright_statement"`
}

type Version struct {
	Metadata	Metadata	`json:"metadata"`
	Verses		[]Verse		`json:"verses"`
}

func LoadVersion(versionCode string) (Version, error) {
	file, err := os.ReadFile(lipgloss.Sprintf("content/%s.json", versionCode))
	if err != nil {
		return Version{}, err
	}

	var version Version
	err = json.Unmarshal(file, &version)
	if err != nil {
		return Version{}, err
	}

	return version, nil
}

func (v *Version) GetBookText(book int) string {
	var text string
	bookStyle := lipgloss.NewStyle().Bold(true).
					Border(lipgloss.ASCIIBorder()).
					Width(20).
					Align(lipgloss.Center)
	chapterStyle := lipgloss.NewStyle().Bold(true).
					Underline(true).
					Foreground(lipgloss.BrightRed)
	verseStyle := lipgloss.NewStyle().Bold(true).
					Foreground(lipgloss.Cyan)
	for _, i := range v.Verses {
		if i.Book < book {
			continue
		}
		if i.Book > book {
			break
		}
		if i.Verse == 1 {
			if i.Chapter == 1{
				text = lipgloss.Sprintf("%s\n\n\n%s", text, bookStyle.Render(i.BookName))
			}
			text = lipgloss.Sprintf("%s\n\n%s", text, chapterStyle.Render(fmt.Sprintf("Chapter %s", strconv.Itoa(i.Chapter))))
		}
		newline := strings.Contains(i.Text, "¶ ")
		verse := strings.Replace(i.Text, "¶ ", "", 1)
		if newline {
			text = lipgloss.Sprintf("%s\n[%s]%s", text, verseStyle.Render(strconv.Itoa(i.Verse)), verse)
		} else {
			text = lipgloss.Sprintf("%s [%s]%s", text, verseStyle.Render(strconv.Itoa(i.Verse)), verse)
		}
	}
	return text
}

