package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ayushlm10/skim/internal/app"
	"github.com/Ayushlm10/skim/internal/upgrade"
	tea "github.com/charmbracelet/bubbletea"
)

// version is set via ldflags at build time
// Default to "dev" for local development builds
var version = "dev"

func main() {
	// Handle subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("skim %s\n", version)
			return
		case "upgrade":
			if err := upgrade.Run(version, os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	// Parse CLI arguments - default to current directory
	rootPath := "."
	if len(os.Args) > 1 {
		// Skip if first arg looks like a flag (already handled above for known flags)
		if os.Args[1][0] != '-' {
			rootPath = os.Args[1]
		}
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

func printHelp() {
	fmt.Printf(`skim - A terminal markdown viewer

Usage:
  skim [path]          Open skim in the specified directory (default: current directory)
  skim version         Print version information
  skim upgrade         Upgrade skim to the latest version
  skim help            Show this help message

Flags:
  -h, --help           Show this help message
  -v, --version        Print version information

Navigation:
  ↑/k, ↓/j             Move selection up/down
  Enter                Open file or toggle directory
  Tab                  Switch focus between panels
  /                    Filter files (file tree) or search (preview)
  n/N                  Next/previous search match
  i                    Toggle ignored directories
  ?                    Show help overlay
  q, Ctrl+C            Quit

Version: %s
`, version)
}
