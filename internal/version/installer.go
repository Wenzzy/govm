package version

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wenzzy/govm/internal/config"
)

// Installer handles installing and uninstalling Go versions
type Installer struct {
	paths *config.Paths
}

// NewInstaller creates a new installer
func NewInstaller() (*Installer, error) {
	paths, err := config.GetPaths()
	if err != nil {
		return nil, err
	}
	return &Installer{paths: paths}, nil
}

// Install extracts and installs a Go version from an archive
func (i *Installer) Install(archivePath, version string) error {
	// Create version directory
	versionPath := i.paths.VersionPath(version)

	// Remove if exists (fresh install)
	if err := os.RemoveAll(versionPath); err != nil {
		return fmt.Errorf("failed to remove existing version: %w", err)
	}

	if err := os.MkdirAll(versionPath, 0755); err != nil {
		return fmt.Errorf("failed to create version directory: %w", err)
	}

	// Extract archive
	if err := i.extractTarGz(archivePath, versionPath); err != nil {
		// Clean up on error
		os.RemoveAll(versionPath)
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	return nil
}

// Uninstall removes an installed Go version
func (i *Installer) Uninstall(version string) error {
	versionPath := i.paths.VersionPath(version)

	// Check if version exists
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// Check if it's the current version
	current, err := i.GetCurrent()
	if err == nil && current == version {
		// Remove the symlink first
		if err := os.Remove(i.paths.Current); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove current symlink: %w", err)
		}
	}

	// Remove version directory
	if err := os.RemoveAll(versionPath); err != nil {
		return fmt.Errorf("failed to remove version: %w", err)
	}

	return nil
}

// SetCurrent sets the current Go version
func (i *Installer) SetCurrent(version string) error {
	versionPath := i.paths.VersionPath(version)

	// Verify version exists
	goPath := filepath.Join(versionPath, "go")
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// Remove existing symlink
	if err := os.Remove(i.paths.Current); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove current symlink: %w", err)
	}

	// Create new symlink pointing to the go directory inside version
	if err := os.Symlink(goPath, i.paths.Current); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// GetCurrent returns the currently active Go version
func (i *Installer) GetCurrent() (string, error) {
	// Read symlink target
	target, err := os.Readlink(i.paths.Current)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	// Extract version from path
	// target is like /home/user/.govm/versions/1.21.0/go
	// We want "1.21.0"
	dir := filepath.Dir(target) // /home/user/.govm/versions/1.21.0
	version := filepath.Base(dir)
	return version, nil
}

// ListInstalled returns a list of installed Go versions
func (i *Installer) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(i.paths.Versions)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Verify it's a valid Go installation
			goPath := filepath.Join(i.paths.Versions, entry.Name(), "go", "bin", "go")
			if _, err := os.Stat(goPath); err == nil {
				versions = append(versions, entry.Name())
			}
		}
	}

	// Sort versions
	sortVersionsDesc(versions)
	return versions, nil
}

// IsInstalled checks if a version is installed
func (i *Installer) IsInstalled(version string) bool {
	goPath := filepath.Join(i.paths.VersionPath(version), "go", "bin", "go")
	_, err := os.Stat(goPath)
	return err == nil
}

// GetGoBinary returns the path to the Go binary for a version
func (i *Installer) GetGoBinary(version string) (string, error) {
	goPath := filepath.Join(i.paths.VersionPath(version), "go", "bin", "go")
	if _, err := os.Stat(goPath); err != nil {
		return "", fmt.Errorf("go binary not found for version %s", version)
	}
	return goPath, nil
}

// extractTarGz extracts a tar.gz archive to the destination
func (i *Installer) extractTarGz(archivePath, destPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Security: prevent path traversal
		targetPath := filepath.Join(destPath, header.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeSymlink:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}
