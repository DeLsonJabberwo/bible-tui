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
	buffer, err := buffer.NewBuffer(viewportInfo, "kjv", 1)
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

