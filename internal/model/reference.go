package model

import (
	"github.com/delsonjabberwo/bible-tui/internal/bible"
)

type Reference struct {
	Book	string
	BookInd int
	Chapter	int
	Verse	int
}

func ReferencesFromVersion(version bible.Version) []Reference {
	var references []Reference
	for _, v := range version.Verses {
		references = append(references, Reference{
			Book: v.BookName,
			BookInd: v.Book,
			Chapter: v.Chapter,
			Verse: v.Verse,
		})
	}
	return references
}

func NewReference(verse bible.Verse) Reference {
	return Reference{
		Book: verse.BookName,
		BookInd: verse.Book,
		Chapter: verse.Chapter,
		Verse: verse.Verse,
	}
}

