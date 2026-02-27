package version

import (
	"fmt"

	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

// Manager coordinates all version management operations
type Manager struct {
	installer  *Installer
	downloader *Downloader
	paths      *config.Paths
}

// NewManager creates a new version manager
func NewManager() (*Manager, error) {
	paths, err := config.GetPaths()
	if err != nil {
		return nil, err
	}

	// Ensure directories exist
	if err := paths.EnsureDirs(); err != nil {
		return nil, err
	}

	installer, err := NewInstaller()
	if err != nil {
		return nil, err
	}

	downloader, err := NewDownloader()
	if err != nil {
		return nil, err
	}

	return &Manager{
		installer:  installer,
		downloader: downloader,
		paths:      paths,
	}, nil
}

// Install downloads and installs a Go version
func (m *Manager) Install(version string, setDefault bool, showProgress bool) error {
	// Normalize version
	version = config.NormalizeVersion(version)

	// Check if already installed
	if m.installer.IsInstalled(version) {
		ui.PrintInfo("Go %s is already installed", version)
		if setDefault {
			return m.Use(version)
		}
		return nil
	}

	// Check if version exists remotely
	available, err := IsVersionAvailable(version)
	if err != nil {
		return fmt.Errorf("failed to check version availability: %w", err)
	}
	if !available {
		return fmt.Errorf("version %s is not available for download", version)
	}

	// Download
	ui.PrintInfo("Downloading Go %s...", version)
	archivePath, err := m.downloader.Download(version, showProgress)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	// Install
	spinner := ui.NewSpinner(fmt.Sprintf("Installing Go %s...", version))
	spinner.Start()

	if err := m.installer.Install(archivePath, version); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to install Go %s", version))
		return fmt.Errorf("failed to install: %w", err)
	}
	spinner.Success(fmt.Sprintf("Installed Go %s", version))

	// Set as current/default if requested or if it's the first version
	if setDefault {
		return m.Use(version)
	}

	// If no version is currently set, use this one
	current, err := m.installer.GetCurrent()
	if err == nil && current == "" {
		ui.PrintHint("Setting as current version (first install)")
		return m.Use(version)
	}

	return nil
}

// Uninstall removes an installed Go version
func (m *Manager) Uninstall(version string) error {
	version = config.NormalizeVersion(version)

	if !m.installer.IsInstalled(version) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// Check if it's the current version
	current, _ := m.installer.GetCurrent()
	if current == version {
		ui.PrintWarning("Removing current version, you may need to switch to another version")
	}

	if err := m.installer.Uninstall(version); err != nil {
		return err
	}

	ui.PrintSuccess("Uninstalled Go %s", version)
	return nil
}

// Use switches to a specific Go version
func (m *Manager) Use(version string) error {
	// Resolve alias if needed
	version = config.ResolveVersion(version)

	if !m.installer.IsInstalled(version) {
		// Check if auto-install is enabled
		cfg := config.Get()
		if cfg.AutoInstall {
			ui.PrintInfo("Version %s not installed, installing...", version)
			if err := m.Install(version, false, true); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("version %s is not installed (auto-install is disabled)", version)
		}
	}

	if err := m.installer.SetCurrent(version); err != nil {
		return err
	}

	ui.PrintSuccess("Now using Go %s", version)
	return nil
}

// UseFromProject detects and uses the Go version from go.mod/go.work
func (m *Manager) UseFromProject(dir string) error {
	version, source, err := DetectVersion(dir)
	if err != nil {
		return err
	}

	// Normalize to full version (e.g., 1.21 -> 1.21.0)
	fullVersion, err := NormalizeDetectedVersion(version)
	if err != nil {
		fullVersion = version
	}

	ui.PrintInfo("Detected Go %s from %s", fullVersion, source)
	return m.Use(fullVersion)
}

// Current returns the current Go version
func (m *Manager) Current() (string, error) {
	return m.installer.GetCurrent()
}

// ListInstalled returns installed versions
func (m *Manager) ListInstalled() ([]string, error) {
	return m.installer.ListInstalled()
}

// IsInstalled checks if a version is installed
func (m *Manager) IsInstalled(version string) bool {
	version = config.NormalizeVersion(version)
	return m.installer.IsInstalled(version)
}

// GetGoBinary returns the path to the Go binary for a version
func (m *Manager) GetGoBinary(version string) (string, error) {
	version = config.ResolveVersion(version)
	return m.installer.GetGoBinary(version)
}

// SetDefault sets the default Go version in config
func (m *Manager) SetDefault(version string) error {
	version = config.NormalizeVersion(version)

	cfg := config.Get()
	cfg.DefaultVersion = version
	if err := config.Save(cfg); err != nil {
		return err
	}

	ui.PrintSuccess("Set Go %s as default", version)
	return nil
}

// GetDefault returns the default Go version from config
func (m *Manager) GetDefault() string {
	cfg := config.Get()
	return cfg.DefaultVersion
}

// QuietUse switches to a version without output (for shell hooks)
func (m *Manager) QuietUse(version string) error {
	version = config.ResolveVersion(version)

	if !m.installer.IsInstalled(version) {
		cfg := config.Get()
		if cfg.AutoInstall {
			// Install quietly
			archivePath, err := m.downloader.Download(version, false)
			if err != nil {
				return err
			}
			if err := m.installer.Install(archivePath, version); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("version %s is not installed", version)
		}
	}

	return m.installer.SetCurrent(version)
}
