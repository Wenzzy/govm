package version

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/wenzzy/govm/internal/config"
)

// Downloader handles downloading Go versions
type Downloader struct {
	paths *config.Paths
}

// NewDownloader creates a new downloader
func NewDownloader() (*Downloader, error) {
	paths, err := config.GetPaths()
	if err != nil {
		return nil, err
	}
	return &Downloader{paths: paths}, nil
}

// Download downloads a Go version archive
// Returns the path to the downloaded file
func (d *Downloader) Download(version string, showProgress bool) (string, error) {
	url, expectedHash, err := GetDownloadURL(version)
	if err != nil {
		return "", fmt.Errorf("failed to get download URL: %w", err)
	}

	// Determine filename from URL
	filename := "go" + version + ".tar.gz"
	destPath := d.paths.CachePath(filename)

	// Check if already cached with correct hash
	if d.isValidCache(destPath, expectedHash) {
		return destPath, nil
	}

	// Create cache directory
	if err := os.MkdirAll(d.paths.Cache, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Download the file
	client := &http.Client{
		Timeout: 30 * time.Minute, // Large file, long timeout
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp(d.paths.Cache, "download-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // Clean up temp file on error

	// Setup progress bar if requested
	var writer io.Writer = tmpFile
	if showProgress {
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionSetDescription("Downloading"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionThrottle(100*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stderr, "\n")
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
		writer = io.MultiWriter(tmpFile, bar)
	}

	// Calculate hash while downloading
	hash := sha256.New()
	multiWriter := io.MultiWriter(writer, hash)

	_, err = io.Copy(multiWriter, resp.Body)
	tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Verify hash
	actualHash := hex.EncodeToString(hash.Sum(nil))
	if expectedHash != "" && actualHash != expectedHash {
		return "", fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	// Move to final location
	if err := os.Rename(tmpPath, destPath); err != nil {
		// If rename fails (cross-device), try copy
		if err := copyFile(tmpPath, destPath); err != nil {
			return "", fmt.Errorf("failed to save file: %w", err)
		}
	}

	return destPath, nil
}

// isValidCache checks if a cached file exists and has the correct hash
func (d *Downloader) isValidCache(path, expectedHash string) bool {
	if expectedHash == "" {
		return false
	}

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return false
	}

	actualHash := hex.EncodeToString(hash.Sum(nil))
	return actualHash == expectedHash
}

// CleanCache removes all cached files
func (d *Downloader) CleanCache() error {
	return os.RemoveAll(d.paths.Cache)
}

// CacheSize returns the total size of cached files
func (d *Downloader) CacheSize() (int64, error) {
	var size int64
	entries, err := os.ReadDir(d.paths.Cache)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		size += info.Size()
	}
	return size, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
