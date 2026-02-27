package shell

import (
	"os"
	"path/filepath"

	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/version"
)

// AutoSwitch performs automatic version switching based on go.mod/go.work
// This is called by shell hooks when changing directories
func AutoSwitch(dir string) (string, bool, error) {
	// Detect version from go.mod/go.work
	ver, _, err := version.DetectVersionInDir(dir)
	if err != nil {
		return "", false, nil // No go.mod/go.work, not an error
	}

	// Normalize version
	fullVersion, err := version.NormalizeDetectedVersion(ver)
	if err != nil {
		fullVersion = ver
	}

	// Get current version
	mgr, err := version.NewManager()
	if err != nil {
		return "", false, err
	}

	current, _ := mgr.Current()
	if current == fullVersion {
		return fullVersion, false, nil // Already using correct version
	}

	// Switch to the version
	if err := mgr.QuietUse(fullVersion); err != nil {
		return fullVersion, false, err
	}

	return fullVersion, true, nil
}

// GetGovmBin returns the path to the govm binary directory
func GetGovmBin() string {
	paths, err := config.GetPaths()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".govm", "bin")
	}
	return paths.Bin
}

// GetCurrentBin returns the path to the current Go binary directory
func GetCurrentBin() string {
	paths, err := config.GetPaths()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".govm", "current", "bin")
	}
	return filepath.Join(paths.Current, "bin")
}
