package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/athakur/local-md/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse CLI arguments - default to current directory
	rootPath := "."
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Verify the path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Path does not exist: %s\n", absPath)
		} else {
			fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", err)
		}
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Path is not a directory: %s\n", absPath)
		os.Exit(1)
	}

	// Create and run the Bubble Tea program
	model := app.New(absPath)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
