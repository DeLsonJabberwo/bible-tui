package buffer

import (
	"regexp"

	"charm.land/lipgloss/v2"
)


func (b *Buffer) versionStyling(verse string) string {
	redLetter := lipgloss.NewStyle().
		Foreground(lipgloss.BrightRed)
	italics := lipgloss.NewStyle().
		Faint(true)
	switch b.Version.Metadata.ShortName {
	case "KJV":
		ital := regexp.MustCompile(`\[([^\]]+)\]`)
		verse = ital.ReplaceAllStringFunc(verse, func(match string) string {
			groups := ital.FindStringSubmatch(match)
			if len(groups) > 1 {
				return italics.Render(groups[1])
			}
			return match
		})

		red := regexp.MustCompile(`‹([^›]+)›`)
		verse = red.ReplaceAllStringFunc(verse, func(match string) string {
			groups := red.FindStringSubmatch(match)
			if len(groups) > 1 {
				return redLetter.Render(groups[1])
			}
			return match
		})
	}
	return verse
}
