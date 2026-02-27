package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

var upgradeCheck bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade govm to the latest version",
	Long: `Upgrade govm CLI to the latest version.

Examples:
  govm upgrade                 Upgrade to latest version
  govm upgrade --check         Check for updates without installing
  g upgrade                    Short form`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if upgradeCheck {
			return checkForUpdates()
		}
		return performUpgrade()
	},
}

// GitHub release info
type releaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func checkForUpdates() error {
	spinner := ui.NewSpinner("Checking for updates...")
	spinner.Start()

	latest, err := getLatestRelease()
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	currentVersion := config.Version
	latestVersion := latest.TagName

	if latestVersion == currentVersion || latestVersion == "v"+currentVersion {
		ui.PrintSuccess("You are running the latest version (%s)", currentVersion)
		return nil
	}

	ui.PrintInfo("Current version: %s", currentVersion)
	ui.PrintInfo("Latest version:  %s", ui.GreenBold.Sprint(latestVersion))
	ui.PrintHint("Run 'govm upgrade' to update")

	return nil
}

func performUpgrade() error {
	spinner := ui.NewSpinner("Checking for updates...")
	spinner.Start()

	latest, err := getLatestRelease()
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	currentVersion := config.Version
	latestVersion := latest.TagName

	if latestVersion == currentVersion || latestVersion == "v"+currentVersion {
		ui.PrintSuccess("You are already running the latest version (%s)", currentVersion)
		return nil
	}

	ui.PrintInfo("Upgrading from %s to %s", currentVersion, latestVersion)

	// Find the appropriate asset for this OS/arch
	assetName := fmt.Sprintf("govm-%s-%s", runtime.GOOS, runtime.GOARCH)
	var downloadURL string

	for _, asset := range latest.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		ui.PrintWarning("No pre-built binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
		ui.PrintHint("Please download manually from: %s/releases", config.RepoURL)
		return nil
	}

	// Download new binary
	spinner = ui.NewSpinner("Downloading...")
	spinner.Start()

	tempFile, err := downloadFile(downloadURL)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tempFile)

	// Make executable
	if err := os.Chmod(tempFile, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Get current binary path
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}

	// Resolve symlinks
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Backup current binary
	backupPath := currentBinary + ".bak"
	if err := os.Rename(currentBinary, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary in place
	if err := copyFile(tempFile, currentBinary); err != nil {
		// Restore backup
		os.Rename(backupPath, currentBinary)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Make executable again (just in case)
	os.Chmod(currentBinary, 0755)

	// Remove backup
	os.Remove(backupPath)

	ui.PrintSuccess("Upgraded to %s", latestVersion)

	// Verify
	verifyUpgrade()

	return nil
}

func getLatestRelease() (*releaseInfo, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		config.GitHubOwner,
		config.GitHubRepo)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var release releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadFile(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "govm-upgrade-*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

func verifyUpgrade() {
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command(exe, "--version")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	ui.PrintInfo("Verified: %s", string(output))
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

func init() {
	upgradeCmd.Flags().BoolVarP(&upgradeCheck, "check", "c", false, "Check for updates without installing")
}
