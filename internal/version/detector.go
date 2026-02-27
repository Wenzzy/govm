package version

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// goVersionRegex matches "go X.Y" or "go X.Y.Z" in go.mod/go.work files
	goVersionRegex = regexp.MustCompile(`^go\s+(\d+\.\d+(?:\.\d+)?)`)
)

// DetectVersion detects the Go version from go.mod or go.work in the given directory
// It searches from the given directory up to the root
func DetectVersion(dir string) (string, string, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", "", err
		}
	}

	// Make path absolute
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}

	// Search up the directory tree
	for {
		// Check for go.work first (takes precedence)
		goWorkPath := filepath.Join(dir, "go.work")
		if version, err := parseGoVersionFile(goWorkPath); err == nil && version != "" {
			return version, goWorkPath, nil
		}

		// Then check for go.mod
		goModPath := filepath.Join(dir, "go.mod")
		if version, err := parseGoVersionFile(goModPath); err == nil && version != "" {
			return version, goModPath, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return "", "", fmt.Errorf("no go.mod or go.work found")
}

// DetectVersionInDir detects version only in the specific directory (no parent search)
func DetectVersionInDir(dir string) (string, string, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", "", err
		}
	}

	// Check for go.work first
	goWorkPath := filepath.Join(dir, "go.work")
	if version, err := parseGoVersionFile(goWorkPath); err == nil && version != "" {
		return version, goWorkPath, nil
	}

	// Then check for go.mod
	goModPath := filepath.Join(dir, "go.mod")
	if version, err := parseGoVersionFile(goModPath); err == nil && version != "" {
		return version, goModPath, nil
	}

	return "", "", fmt.Errorf("no go.mod or go.work found in %s", dir)
}

// parseGoVersionFile parses a go.mod or go.work file and extracts the Go version
func parseGoVersionFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Look for "go X.Y" or "go X.Y.Z"
		matches := goVersionRegex.FindStringSubmatch(line)
		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no go version found in %s", path)
}

// HasGoModOrWork checks if a directory has a go.mod or go.work file
func HasGoModOrWork(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	goWorkPath := filepath.Join(dir, "go.work")

	if _, err := os.Stat(goWorkPath); err == nil {
		return true
	}
	if _, err := os.Stat(goModPath); err == nil {
		return true
	}
	return false
}

// NormalizeDetectedVersion converts a version like "1.21" to a full version like "1.21.0"
// by finding the latest patch version available
func NormalizeDetectedVersion(version string) (string, error) {
	// If version already has patch (e.g., "1.21.5"), return as-is
	parts := strings.Split(version, ".")
	if len(parts) >= 3 {
		return version, nil
	}

	// Find the latest patch version for this minor version
	allVersions, err := ListAllVersions()
	if err != nil {
		// If we can't fetch remote versions, try with .0
		return version + ".0", nil
	}

	// Find the latest version that starts with our prefix
	prefix := version + "."
	for _, v := range allVersions {
		if strings.HasPrefix(v, prefix) || v == version {
			return v, nil
		}
	}

	// Fallback to .0
	return version + ".0", nil
}
