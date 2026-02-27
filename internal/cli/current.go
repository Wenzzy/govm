package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
	"github.com/wenzzy/govm/internal/version"
)

var currentPath bool

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"now", "active"},
	Short:   "Show current Go version",
	Long: `Show the currently active Go version.

Examples:
  govm current                Show current version
  govm current --path         Show path to current Go binary
  g current                   Short form`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := version.NewManager()
		if err != nil {
			return err
		}

		current, err := mgr.Current()
		if err != nil {
			return err
		}

		if current == "" {
			ui.PrintWarning("No Go version is currently active")
			ui.PrintHint("Run 'govm use <version>' to activate a version")
			return nil
		}

		if currentPath {
			goBinary, err := mgr.GetGoBinary(current)
			if err != nil {
				return err
			}
			fmt.Println(goBinary)
			return nil
		}

		// Get actual Go version string
		goBinary, _ := mgr.GetGoBinary(current)
		goVersion := getGoVersionString(goBinary)

		ui.PrintKeyValue("Current", ui.GreenBold.Sprint(current))
		if goVersion != "" && goVersion != current {
			ui.PrintKeyValue("Go version", goVersion)
		}

		// Show default if different
		defaultVer := mgr.GetDefault()
		if defaultVer != "" && defaultVer != current {
			ui.PrintKeyValue("Default", defaultVer)
		}

		// Show aliases pointing to current version
		aliases := findAliasesForVersion(current)
		if len(aliases) > 0 {
			ui.PrintKeyValue("Aliases", strings.Join(aliases, ", "))
		}

		// Check for project version
		projectVer, source, err := version.DetectVersion("")
		if err == nil && projectVer != "" {
			normalizedProject, _ := version.NormalizeDetectedVersion(projectVer)
			if normalizedProject != current {
				ui.PrintWarning("Project requires Go %s (from %s)", normalizedProject, source)
				ui.PrintHint("Run 'govm use .' to switch to project version")
			} else {
				ui.PrintInfo("Matches project requirement from %s", ui.Dim.Sprint(source))
			}
		}

		return nil
	},
}

// getGoVersionString runs 'go version' and extracts the version string
func getGoVersionString(goBinary string) string {
	if goBinary == "" {
		return ""
	}

	cmd := exec.Command(goBinary, "version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Output is like "go version go1.21.0 darwin/amd64"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go")
	}
	return ""
}

func init() {
	currentCmd.Flags().BoolVarP(&currentPath, "path", "p", false, "Show path to Go binary")
}

// findAliasesForVersion is defined in list.go but we need it here too
func init() {
	// Already defined in list.go
	_ = config.ListAliases
}
