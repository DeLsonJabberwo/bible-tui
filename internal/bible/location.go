package bible

type VerseInfo struct {
	Book	int
	Chapter	int
	Verse	int
}

type VerseLocs struct {
	/*
	Key: verse info
	Value: line number
	*/
	Verses		map[VerseInfo]int
	LineCount	int
}

func (v *VerseLocs) GetVerseFromLine(line int) VerseInfo {
	var verse VerseInfo
	var closest int
	for currVerse, currLine := range v.Verses {
		if currLine <= line && currLine > closest {
			closest = currLine
			verse = currVerse
		}
	}
	return verse
}

type ChapterInfo struct {
	Book	int
	Chapter	int
}

