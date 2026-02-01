package upgrade

import (
	"fmt"
	"strings"
)

// Run executes the upgrade command
func Run(currentVersion string, args []string) error {
	// Parse flags
	var force, checkOnly bool
	for _, arg := range args {
		switch arg {
		case "-f", "--force":
			force = true
		case "-c", "--check":
			checkOnly = true
		case "-h", "--help":
			printUsage()
			return nil
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown flag: %s", arg)
			}
		}
	}

	fmt.Printf("Current version: %s\n", formatVersion(currentVersion))

	// Fetch latest release
	fmt.Print("Checking for updates...")
	release, err := FetchLatestRelease()
	if err != nil {
		fmt.Println(" failed")
		return fmt.Errorf("checking for updates: %w", err)
	}
	fmt.Println(" done")

	latestVersion := release.Version()
	fmt.Printf("Latest version:  %s\n", formatVersion(latestVersion))
	fmt.Println()

	// Compare versions
	needsUpgrade, err := compareVersions(currentVersion, latestVersion)
	if err != nil {
		// If version comparison fails (e.g., "dev" version), allow upgrade with --force
		if !force {
			return fmt.Errorf("version comparison failed: %w (use --force to upgrade anyway)", err)
		}
		needsUpgrade = true
	}

	if !needsUpgrade && !force {
		fmt.Println("You're already on the latest version.")
		return nil
	}

	if checkOnly {
		if needsUpgrade {
			fmt.Println("An update is available!")
			fmt.Printf("Run 'skim upgrade' to upgrade to %s\n", formatVersion(latestVersion))
		}
		return nil
	}

	// Find the appropriate asset for this platform
	asset, err := release.FindAssetForPlatform()
	if err != nil {
		return fmt.Errorf("finding download: %w", err)
	}

	// Download
	fmt.Printf("Downloading %s...\n", asset.Name)
	archivePath, err := DownloadAsset(asset, func(downloaded, total int64) {
		percent := float64(downloaded) / float64(total) * 100
		fmt.Printf("\r  Progress: %.1f%% (%d/%d bytes)", percent, downloaded, total)
	})
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer Cleanup(archivePath)
	fmt.Println()

	// Extract
	fmt.Print("Extracting...")
	binaryPath, err := ExtractBinary(archivePath)
	if err != nil {
		return fmt.Errorf("extracting: %w", err)
	}
	fmt.Println(" done")

	// Replace
	fmt.Print("Installing...")
	if err := ReplaceBinary(binaryPath); err != nil {
		return fmt.Errorf("installing: %w", err)
	}
	fmt.Println(" done")

	fmt.Println()
	fmt.Printf("Successfully upgraded to %s\n", formatVersion(latestVersion))

	return nil
}

func printUsage() {
	fmt.Print(`Usage: skim upgrade [flags]

Upgrade skim to the latest version from GitHub releases.

Flags:
  -c, --check    Only check for updates, don't install
  -f, --force    Upgrade even if already on latest version
  -h, --help     Show this help message

Examples:
  skim upgrade           Upgrade to the latest version
  skim upgrade --check   Check if an update is available
  skim upgrade --force   Force upgrade (useful for dev builds)
`)
}

func formatVersion(v string) string {
	if v == "dev" {
		return "dev"
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

// compareVersions compares two semantic versions
// Returns true if latest > current
func compareVersions(current, latest string) (bool, error) {
	// Handle "dev" version
	if current == "dev" {
		return false, fmt.Errorf("cannot compare dev version")
	}

	// Strip 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Parse versions
	currentParts, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("parsing current version: %w", err)
	}

	latestParts, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf("parsing latest version: %w", err)
	}

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true, nil
		}
		if latestParts[i] < currentParts[i] {
			return false, nil
		}
	}

	return false, nil
}

// parseVersion parses a semantic version string into [major, minor, patch]
func parseVersion(v string) ([3]int, error) {
	var parts [3]int

	// Handle versions like "1.0.0-beta.1" by taking only the main part
	if idx := strings.IndexAny(v, "-+"); idx != -1 {
		v = v[:idx]
	}

	// Split by dots
	segments := strings.Split(v, ".")

	for i := 0; i < len(segments) && i < 3; i++ {
		var n int
		_, err := fmt.Sscanf(segments[i], "%d", &n)
		if err != nil {
			return parts, fmt.Errorf("invalid version segment: %s", segments[i])
		}
		parts[i] = n
	}

	return parts, nil
}
