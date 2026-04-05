package buffer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/delsonjabberwo/bible-tui/internal/bible"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

type Buffer struct {
	Version			bible.Version
	Books 			[]int
	Content 		string
	VerseLocs		bible.VerseLocs
	ChapterLocs		map[bible.ChapterInfo]int
	BookLocs		map[int]int

}

func NewBuffer(viewportInfo ViewportInfo, versionCode string, book int) (Buffer, error) {
	var buffer Buffer
	var err error
	buffer.Version, err = bible.LoadVersion(versionCode)
	if err != nil {
		return Buffer{}, err
	}
	buffer.Books = []int{ book }
	buffer.VerseLocs = bible.VerseLocs{
		Verses: make(map[bible.VerseInfo]int),
		LineCount: 0,
	}
	buffer.ChapterLocs = make(map[bible.ChapterInfo]int)
	buffer.BookLocs = make(map[int]int)
	buffer.AppendBook(viewportInfo, book)

	return buffer, nil
}

func (b *Buffer) UpdateBuffer(viewportInfo ViewportInfo, yOffset int) int {
	verse := b.VerseLocs.GetVerseFromLine(yOffset)
	b.Content = ""
	for _, i := range b.Books {
		b.AppendBook(viewportInfo, i)
	}
	return b.VerseLocs.Verses[verse]
}

func (b *Buffer) AppendBook(viewportInfo ViewportInfo, bookNum int) error {
	widthLimit := viewportInfo.WordWidthLimit()
	bookStyle := lipgloss.NewStyle().Bold(true).
					Border(lipgloss.ASCIIBorder()).
					Width(20).
					Align(lipgloss.Center)
	chapterStyle := lipgloss.NewStyle().Bold(true).
					Underline(true).
					Foreground(lipgloss.BrightRed)
	verseStyle := lipgloss.NewStyle().Bold(true).
					Foreground(lipgloss.Cyan)

	var sb strings.Builder
	for _, i := range b.Version.Verses {
		if i.Book < bookNum {
			continue
		}
		if i.Book > bookNum {
			break
		}
		if i.Verse == 1 {
			if i.Chapter == 1{
				_, err := sb.WriteString(lipgloss.Sprintf("\n\n\n%s", bookStyle.Render(i.BookName)))
				if err != nil {
					return err
				}
			}
			_, err := sb.WriteString(lipgloss.Sprintf("\n\n%s", chapterStyle.Render(fmt.Sprintf("Chapter %s", strconv.Itoa(i.Chapter)))))
			if err != nil {
				return err
			}
		}
		verseInfo := bible.VerseInfo{
			Book: i.Book,
			Chapter: i.Chapter,
			Verse: i.Verse,
		}
		b.VerseLocs.Verses[verseInfo] = b.VerseLocs.LineCount
		newline := strings.Contains(i.Text, "¶ ")
		if newline {
			verse := strings.Replace(i.Text, "¶ ", "", 1)
			_, err := sb.WriteString(lipgloss.Sprintf("\n[%s]%s", verseStyle.Render(strconv.Itoa(i.Verse)), verse))
			if err != nil {
				return err
			}
		} else {
			_, err := sb.WriteString(lipgloss.Sprintf(" [%s]%s", verseStyle.Render(strconv.Itoa(i.Verse)), i.Text))
			if err != nil {
				return err
			}
		}
	}

	b.Content = wrap.String(wordwrap.String(sb.String(), widthLimit), viewportInfo.MaxWidth())

	plainContent := ansi.Strip(b.Content)

	lines := strings.Split(plainContent, "\n")
	b.VerseLocs.LineCount = len(lines)

	b.VerseLocs.Verses = make(map[bible.VerseInfo]int)
	re := regexp.MustCompile(`\[(\d+)\]`)
	var ch, vs int
	for num, line := range lines {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				if verseNum, err := strconv.Atoi(match[1]); err == nil {
					switch verseNum {
						case 1:
							ch++
							vs = 1
						case vs + 1:
							vs++
						default:
							return fmt.Errorf("error: content failure")
					}
					verseInfo := bible.VerseInfo{
						Book: bookNum,
						Chapter: ch,
						Verse: vs,
					}
					b.VerseLocs.Verses[verseInfo] = num
				}
			}
		}
	}

	return nil
}
