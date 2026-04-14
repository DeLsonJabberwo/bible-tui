package main

import (
	"embed"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/delsonjabberwo/bible-tui/internal/bible"
	"github.com/delsonjabberwo/bible-tui/internal/buffer"
	"github.com/delsonjabberwo/bible-tui/internal/model"
)

//go:embed content/*.json
var contentFS embed.FS

func main() {
	if os.Getenv("DEBUG") == "1" {
		f, err := tea.LogToFile("tmp/debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		f, err := tea.LogToFile("/dev/null", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Connect embedded content filesystem
	bible.ContentFS = contentFS

	viewportInfo := buffer.NewViewportInfo(0)

	var versionCode string
	if len(os.Args) > 1 && versions[os.Args[1]] {
		versionCode = os.Args[1]
	} else {
		versionCode = "kjv"
	}

	buffer, err := buffer.NewBuffer(viewportInfo, versionCode, 1)
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		&model.Model{Buffer: buffer},
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

var versions map[string]bool = map[string]bool{
	"kjv": true,
	"geneva": true,
	"asv": true,
	"web": true,
}

