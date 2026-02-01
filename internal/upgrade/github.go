package upgrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHub repository information
	repoOwner = "Ayushlm10"
	repoName  = "skim"

	// API endpoint
	githubAPIURL = "https://api.github.com/repos/%s/%s/releases/latest"

	// Request timeout
	httpTimeout = 30 * time.Second
)

// Release represents a GitHub release
type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Assets  []Asset `json:"assets"`
	HTMLURL string  `json:"html_url"`
}

// Asset represents a release asset (downloadable file)
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// FetchLatestRelease fetches the latest release information from GitHub
func FetchLatestRelease() (*Release, error) {
	url := fmt.Sprintf(githubAPIURL, repoOwner, repoName)

	client := &http.Client{
		Timeout: httpTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers for GitHub API
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "skim-upgrade")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found for %s/%s", repoOwner, repoName)
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API rate limit exceeded, try again later")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &release, nil
}

// FindAssetForPlatform finds the appropriate asset for the current platform
func (r *Release) FindAssetForPlatform() (*Asset, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Build expected filename pattern: skim_<version>_<os>_<arch>.<ext>
	// Version from tag might have 'v' prefix, strip it
	version := strings.TrimPrefix(r.TagName, "v")

	var ext string
	if goos == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}

	expectedName := fmt.Sprintf("skim_%s_%s_%s.%s", version, goos, goarch, ext)

	for _, asset := range r.Assets {
		if asset.Name == expectedName {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no asset found for %s/%s (looking for %s)", goos, goarch, expectedName)
}

// Version returns the version string without the 'v' prefix
func (r *Release) Version() string {
	return strings.TrimPrefix(r.TagName, "v")
}
