package buffer

import (
	"fmt"
	"regexp"
	"slices"
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
	LastViewportInfo ViewportInfo

}

func NewBuffer(viewportInfo ViewportInfo, versionCode string, book int) (Buffer, error) {
	var buffer Buffer
	var err error
	buffer.Version, err = bible.LoadVersion(versionCode)
	if err != nil {
		return Buffer{}, err
	}
	buffer.Books = make([]int, 5)
	var curr int
	if book - 2 > 0 {
		buffer.Books[curr] = book - 2
		curr++
	}
	if book - 1 > 0 {
		buffer.Books[curr] = book - 1
		curr++
	}
	for curr < 5 {
		buffer.Books[curr] = book
		book++
		curr++
	}
	buffer.VerseLocs = bible.VerseLocs{
		Verses: make(map[bible.VerseInfo]int),
		LineCount: 0,
	}
	buffer.ChapterLocs = make(map[bible.ChapterInfo]int)
	buffer.BookLocs = make(map[int]int)
	buffer.LastViewportInfo = viewportInfo
	buffer.RenderBooks(viewportInfo)

	return buffer, nil
}

func (b *Buffer) UpdateBuffer(viewportInfo ViewportInfo, yOffset int) int {
	verse := b.VerseLocs.GetVerseFromLine(yOffset)
	b.Content = ""
	b.RenderBooks(viewportInfo)
	return b.VerseLocs.Verses[verse]
}

func (b *Buffer) RenderBooks(viewportInfo ViewportInfo) error {
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
		if !slices.Contains(b.Books, i.Book) {
			continue
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
	b.VerseLocs.LineCount = 0
	b.VerseLocs.Verses = make(map[bible.VerseInfo]int)
	b.BookLocs = make(map[int]int)

	plainContent := ansi.Strip(b.Content)

	books := strings.Split(plainContent, "\n\n\n")
	books = books[1:]
	currentLineOffset := 0
	for bookInd, book := range books {
		b.BookLocs[b.Books[bookInd]] = currentLineOffset
		lines := strings.Split(book, "\n")
		b.VerseLocs.LineCount += len(lines)

		re := regexp.MustCompile(`\[(\d+)\]`)
		bookNum := b.Books[bookInd]
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
						b.VerseLocs.Verses[verseInfo] = num + currentLineOffset
					}
				}
			}
		}
		currentLineOffset += len(lines)
	}

	return nil
}

func (b *Buffer) ShiftBooksNext() {
	b.Books = append(b.Books[1:], b.Books[len(b.Books) - 1] + 1)
}

func (b *Buffer) ShiftBooksPrev() {
	curr := b.Books[:len(b.Books) - 1]
	b.Books = append([]int{b.Books[0] - 1}, curr...)
}

func (b *Buffer) GetBookFromLine(line int) int {
	var book int
	var closest int
	for _, currBook := range b.Books {
		currLine := b.BookLocs[currBook]
		if currLine <= line && currLine >= closest {
			closest = currLine
			book = currBook
		}
	}
	return book
}
