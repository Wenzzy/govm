package config

import (
	"os"
	"path/filepath"
)

var (
	// Version is set at build time
	Version = "dev"
	// BuildTime is set at build time
	BuildTime = "unknown"
	// RepoURL is the URL of the GitHub repository
	RepoURL = "https://github.com/wenzzy/govm"
	// GitHubOwner is the GitHub repository owner
	GitHubOwner = "wenzzy"
	// GitHubRepo is the GitHub repository name
	GitHubRepo = "govm"
)

// Paths holds all the paths used by govm
type Paths struct {
	Root       string // ~/.govm
	Versions   string // ~/.govm/versions
	Current    string // ~/.govm/current (symlink)
	Cache      string // ~/.govm/cache
	Config     string // ~/.govm/config.toml
	Bin        string // ~/.govm/bin
}

// GetPaths returns the paths for govm
func GetPaths() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(home, ".govm")

	// Check for GOVM_ROOT env override
	if envRoot := os.Getenv("GOVM_ROOT"); envRoot != "" {
		root = envRoot
	}

	return &Paths{
		Root:     root,
		Versions: filepath.Join(root, "versions"),
		Current:  filepath.Join(root, "current"),
		Cache:    filepath.Join(root, "cache"),
		Config:   filepath.Join(root, "config.toml"),
		Bin:      filepath.Join(root, "bin"),
	}, nil
}

// EnsureDirs creates all necessary directories
func (p *Paths) EnsureDirs() error {
	dirs := []string{p.Root, p.Versions, p.Cache, p.Bin}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// VersionPath returns the path to a specific version
func (p *Paths) VersionPath(version string) string {
	return filepath.Join(p.Versions, version)
}

// VersionBinPath returns the path to a version's bin directory
func (p *Paths) VersionBinPath(version string) string {
	return filepath.Join(p.Versions, version, "go", "bin")
}

// CachePath returns the path for a cached archive
func (p *Paths) CachePath(filename string) string {
	return filepath.Join(p.Cache, filename)
}
