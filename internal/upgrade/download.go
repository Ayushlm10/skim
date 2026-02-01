package upgrade

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DownloadAsset downloads an asset to a temporary directory and returns the path
func DownloadAsset(asset *Asset, progressFn func(downloaded, total int64)) (string, error) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "skim-upgrade-*")
	if err != nil {
		return "", fmt.Errorf("creating temp directory: %w", err)
	}

	destPath := filepath.Join(tempDir, asset.Name)

	// Download the file
	client := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := client.Get(asset.BrowserDownloadURL)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("downloading asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	// Copy with progress tracking
	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("writing to file: %w", writeErr)
			}
			downloaded += int64(n)
			if progressFn != nil {
				progressFn(downloaded, asset.Size)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("reading response: %w", err)
		}
	}

	return destPath, nil
}

// ExtractBinary extracts the skim binary from the downloaded archive
func ExtractBinary(archivePath string) (string, error) {
	tempDir := filepath.Dir(archivePath)
	binaryName := "skim"
	if runtime.GOOS == "windows" {
		binaryName = "skim.exe"
	}

	extractedPath := filepath.Join(tempDir, binaryName)

	if strings.HasSuffix(archivePath, ".zip") {
		if err := extractFromZip(archivePath, binaryName, extractedPath); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		if err := extractFromTarGz(archivePath, binaryName, extractedPath); err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	return extractedPath, nil
}

func extractFromTarGz(archivePath, binaryName, destPath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("opening archive: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}

		// Look for the binary (might be at root or in a subdirectory)
		if filepath.Base(header.Name) == binaryName && header.Typeflag == tar.TypeReg {
			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("creating binary file: %w", err)
			}
			defer out.Close()

			if _, err := io.Copy(out, tr); err != nil {
				return fmt.Errorf("extracting binary: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("binary %s not found in archive", binaryName)
}

func extractFromZip(archivePath, binaryName, destPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		// Look for the binary (might be at root or in a subdirectory)
		if filepath.Base(f.Name) == binaryName && !f.FileInfo().IsDir() {
			src, err := f.Open()
			if err != nil {
				return fmt.Errorf("opening file in zip: %w", err)
			}
			defer src.Close()

			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("creating binary file: %w", err)
			}
			defer out.Close()

			if _, err := io.Copy(out, src); err != nil {
				return fmt.Errorf("extracting binary: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("binary %s not found in archive", binaryName)
}

// ReplaceBinary replaces the current binary with the new one
func ReplaceBinary(newBinaryPath string) error {
	// Get the path to the currently running binary
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	// Resolve any symlinks
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return fmt.Errorf("resolving symlinks: %w", err)
	}

	// Check if we can write to the current binary location
	currentDir := filepath.Dir(currentPath)
	if err := checkWritePermission(currentDir); err != nil {
		return fmt.Errorf("no write permission to %s: %w (try running with sudo)", currentDir, err)
	}

	// On Windows, we can't replace a running binary directly
	// We need to rename the current one first
	if runtime.GOOS == "windows" {
		oldPath := currentPath + ".old"
		// Remove any existing .old file
		os.Remove(oldPath)
		if err := os.Rename(currentPath, oldPath); err != nil {
			return fmt.Errorf("renaming current binary: %w", err)
		}
	}

	// Copy the new binary to the current location
	// We use copy instead of rename to handle cross-device moves
	src, err := os.Open(newBinaryPath)
	if err != nil {
		return fmt.Errorf("opening new binary: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(currentPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("creating destination binary: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copying binary: %w", err)
	}

	return nil
}

func checkWritePermission(dir string) error {
	testFile := filepath.Join(dir, ".skim-upgrade-test")
	f, err := os.Create(testFile)
	if err != nil {
		return err
	}
	f.Close()
	os.Remove(testFile)
	return nil
}

// Cleanup removes the temporary directory
func Cleanup(archivePath string) {
	tempDir := filepath.Dir(archivePath)
	os.RemoveAll(tempDir)
}
