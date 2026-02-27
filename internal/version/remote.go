package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
)

const (
	goDevURL     = "https://go.dev/dl/?mode=json&include=all"
	goDownloadURL = "https://go.dev/dl/"
)

// RemoteVersion represents a Go version available for download
type RemoteVersion struct {
	Version string        `json:"version"`
	Stable  bool          `json:"stable"`
	Files   []VersionFile `json:"files"`
}

// VersionFile represents a downloadable file for a Go version
type VersionFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"` // archive, installer, source
}

// FetchRemoteVersions fetches all available Go versions from go.dev
func FetchRemoteVersions() ([]RemoteVersion, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(goDevURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []RemoteVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	return versions, nil
}

// GetLatestStable returns the latest stable version
func GetLatestStable() (string, error) {
	versions, err := FetchRemoteVersions()
	if err != nil {
		return "", err
	}

	for _, v := range versions {
		if v.Stable {
			return normalizeVersionString(v.Version), nil
		}
	}

	return "", fmt.Errorf("no stable version found")
}

// GetVersionInfo returns info about a specific version
func GetVersionInfo(version string) (*RemoteVersion, error) {
	versions, err := FetchRemoteVersions()
	if err != nil {
		return nil, err
	}

	version = normalizeVersionString(version)
	goVersion := "go" + version

	for _, v := range versions {
		if v.Version == goVersion || normalizeVersionString(v.Version) == version {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("version %s not found", version)
}

// GetDownloadURL returns the download URL for a version
func GetDownloadURL(version string) (string, string, error) {
	info, err := GetVersionInfo(version)
	if err != nil {
		return "", "", err
	}

	os := runtime.GOOS
	arch := runtime.GOARCH

	for _, f := range info.Files {
		if f.OS == os && f.Arch == arch && f.Kind == "archive" {
			url := goDownloadURL + f.Filename
			return url, f.SHA256, nil
		}
	}

	return "", "", fmt.Errorf("no archive found for %s/%s", os, arch)
}

// ListStableVersions returns a list of stable versions (sorted newest first)
func ListStableVersions() ([]string, error) {
	versions, err := FetchRemoteVersions()
	if err != nil {
		return nil, err
	}

	var stable []string
	for _, v := range versions {
		if v.Stable {
			stable = append(stable, normalizeVersionString(v.Version))
		}
	}

	sortVersionsDesc(stable)
	return stable, nil
}

// ListAllVersions returns all versions (sorted newest first)
func ListAllVersions() ([]string, error) {
	versions, err := FetchRemoteVersions()
	if err != nil {
		return nil, err
	}

	var all []string
	for _, v := range versions {
		all = append(all, normalizeVersionString(v.Version))
	}

	sortVersionsDesc(all)
	return all, nil
}

// normalizeVersionString removes the "go" prefix from version strings
func normalizeVersionString(version string) string {
	return strings.TrimPrefix(version, "go")
}

// sortVersionsDesc sorts version strings in descending order (newest first)
func sortVersionsDesc(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		vi, err1 := goversion.NewVersion(versions[i])
		vj, err2 := goversion.NewVersion(versions[j])
		if err1 != nil || err2 != nil {
			return versions[i] > versions[j]
		}
		return vi.GreaterThan(vj)
	})
}

// IsVersionAvailable checks if a version is available for download
func IsVersionAvailable(version string) (bool, error) {
	_, err := GetVersionInfo(version)
	if err != nil {
		return false, nil
	}
	return true, nil
}
