package bible

import (
	"encoding/json"
	"os"

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

