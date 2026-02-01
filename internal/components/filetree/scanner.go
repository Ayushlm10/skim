package filetree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ScanOptions configures directory scanning behavior
type ScanOptions struct {
	// ShowHidden includes hidden files/directories (starting with .)
	ShowHidden bool

	// MarkdownOnly only shows .md files (directories always shown if they contain md files)
	MarkdownOnly bool

	// MaxDepth limits recursion depth (-1 = unlimited)
	MaxDepth int

	// IgnoreDirs is a list of directory names to skip (e.g., node_modules, vendor)
	IgnoreDirs []string

	// ShowIgnored when true shows directories that would normally be ignored
	ShowIgnored bool
}

// DefaultIgnoreDirs is the list of directories to ignore by default
// These are common development noise directories that often contain
// many files but rarely have meaningful markdown documentation
var DefaultIgnoreDirs = []string{
	"node_modules",     // JavaScript/Node.js dependencies
	"vendor",           // Go modules, PHP Composer
	"__pycache__",      // Python bytecode cache
	".venv",            // Python virtual environments
	"venv",             // Python virtual environments (alternative)
	"dist",             // Build output directories
	"build",            // Build output directories
	"target",           // Rust/Java build output
	".cache",           // Generic cache directories
	".next",            // Next.js build cache
	".nuxt",            // Nuxt.js build cache
	"coverage",         // Test coverage reports
	".terraform",       // Terraform state/cache
	".serverless",      // Serverless framework
	"bower_components", // Bower dependencies (legacy)
}

// DefaultScanOptions returns sensible defaults
func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		ShowHidden:   false,
		MarkdownOnly: true,
		MaxDepth:     -1,
		IgnoreDirs:   DefaultIgnoreDirs,
		ShowIgnored:  false,
	}
}

// ScanDirectory scans a directory and returns tree items
// Only returns the root level items; children are loaded on expand
func ScanDirectory(rootPath string, opts ScanOptions) ([]*Item, error) {
	// Ensure path is absolute
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	// Verify directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrInvalid
	}

	return scanLevel(absPath, 0, opts)
}

// isIgnoredDir checks if a directory name is in the ignore list
func isIgnoredDir(name string, ignoreDirs []string) bool {
	for _, ignored := range ignoreDirs {
		if name == ignored {
			return true
		}
	}
	return false
}

// scanLevel scans a single directory level
func scanLevel(dirPath string, depth int, opts ScanOptions) ([]*Item, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var items []*Item

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless enabled
		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)
		isDir := entry.IsDir()

		// Skip ignored directories unless ShowIgnored is true
		if isDir && !opts.ShowIgnored && isIgnoredDir(name, opts.IgnoreDirs) {
			continue
		}

		// For files, check if markdown (when MarkdownOnly is true)
		if !isDir && opts.MarkdownOnly {
			if !isMarkdownFile(name) {
				continue
			}
		}

		// For directories, check if they contain any markdown files
		if isDir && opts.MarkdownOnly {
			hasMarkdown, _ := dirContainsMarkdown(fullPath, opts)
			if !hasMarkdown {
				continue
			}
		}

		item := NewItem(fullPath, isDir, depth)
		items = append(items, item)
	}

	// Sort: directories first, then alphabetically
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir // dirs come first
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

// ScanChildren loads children for a directory item
func ScanChildren(item *Item, opts ScanOptions) error {
	if !item.IsDir {
		return nil
	}

	// Check max depth
	if opts.MaxDepth >= 0 && item.Depth >= opts.MaxDepth {
		return nil
	}

	children, err := scanLevel(item.Path, item.Depth+1, opts)
	if err != nil {
		return err
	}

	// Set parent reference
	for _, child := range children {
		child.Parent = item
	}

	item.Children = children
	return nil
}

// isMarkdownFile checks if a filename is a markdown file
func isMarkdownFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown"
}

// dirContainsMarkdown recursively checks if a directory contains markdown files
func dirContainsMarkdown(dirPath string, opts ScanOptions) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden
		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			// Skip ignored directories unless ShowIgnored is true
			if !opts.ShowIgnored && isIgnoredDir(name, opts.IgnoreDirs) {
				continue
			}

			// Recursively check subdirectories
			fullPath := filepath.Join(dirPath, name)
			has, _ := dirContainsMarkdown(fullPath, opts)
			if has {
				return true, nil
			}
		} else if isMarkdownFile(name) {
			return true, nil
		}
	}

	return false, nil
}

// CountMarkdownFiles counts total markdown files in a directory tree
func CountMarkdownFiles(rootPath string, opts ScanOptions) (int, error) {
	count := 0

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		name := d.Name()

		// Skip hidden
		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip ignored directories unless ShowIgnored is true
		if d.IsDir() && !opts.ShowIgnored && isIgnoredDir(name, opts.IgnoreDirs) {
			return filepath.SkipDir
		}

		if !d.IsDir() && isMarkdownFile(name) {
			count++
		}

		return nil
	})

	return count, err
}
