package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
	"github.com/wenzzy/govm/internal/version"
)

var (
	listAll    bool
	listRemote bool
	listLimit  int
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List Go versions",
	Long: `List installed or available Go versions.

Examples:
  govm list                   List installed versions
  govm list remote            List available remote versions
  govm list remote --all      List all remote versions (including RCs, betas)
  govm ls -r -n 20            List last 20 remote versions
  g ls                        Short form`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for "remote" subcommand/argument
		if len(args) > 0 && (args[0] == "remote" || args[0] == "r") {
			listRemote = true
		}

		if listRemote {
			return listRemoteVersions()
		}
		return listInstalledVersions()
	},
}

func listInstalledVersions() error {
	mgr, err := version.NewManager()
	if err != nil {
		return err
	}

	installed, err := mgr.ListInstalled()
	if err != nil {
		return err
	}

	if len(installed) == 0 {
		ui.PrintInfo("No Go versions installed")
		ui.PrintHint("Run 'govm install <version>' to install a version")
		return nil
	}

	current, _ := mgr.Current()
	defaultVer := mgr.GetDefault()

	ui.PrintHeader("Installed Go Versions")

	for _, ver := range installed {
		isCurrent := ver == current

		prefix := "  "
		suffix := ""

		if isCurrent {
			prefix = ui.Green.Sprint(ui.SymbolArrow) + " "
			ver = ui.GreenBold.Sprint(ver)
		} else {
			ver = ui.White.Sprint(ver)
		}

		if current != "" && isCurrent {
			suffix += ui.Dim.Sprint(" (current)")
		}
		if defaultVer != "" && ver == defaultVer {
			suffix += ui.Cyan.Sprint(" [default]")
		}

		// Check for aliases pointing to this version
		aliases := findAliasesForVersion(ver)
		if len(aliases) > 0 {
			for _, a := range aliases {
				suffix += ui.Magenta.Sprintf(" @%s", a)
			}
		}

		fmt.Printf("%s%s%s\n", prefix, ver, suffix)
	}

	return nil
}

func listRemoteVersions() error {
	spinner := ui.NewSpinner("Fetching available versions...")
	spinner.Start()

	var versions []string
	var err error

	if listAll {
		versions, err = version.ListAllVersions()
	} else {
		versions, err = version.ListStableVersions()
	}

	spinner.Stop()

	if err != nil {
		return err
	}

	if len(versions) == 0 {
		ui.PrintInfo("No versions available")
		return nil
	}

	// Apply limit
	if listLimit > 0 && len(versions) > listLimit {
		versions = versions[:listLimit]
	}

	// Get installed versions for marking
	mgr, _ := version.NewManager()
	var installedMap map[string]bool
	var current string
	if mgr != nil {
		installed, _ := mgr.ListInstalled()
		installedMap = make(map[string]bool)
		for _, v := range installed {
			installedMap[v] = true
		}
		current, _ = mgr.Current()
	}

	title := "Available Go Versions (stable)"
	if listAll {
		title = "Available Go Versions (all)"
	}
	ui.PrintHeader(title)

	for _, ver := range versions {
		isCurrent := ver == current
		isInstalled := installedMap[ver]
		ui.PrintVersion(ver, isCurrent, isInstalled)
	}

	if listLimit > 0 && listLimit < len(versions) {
		ui.PrintHint(fmt.Sprintf("Showing first %d versions. Use -n 0 to show all.", listLimit))
	}

	return nil
}

func findAliasesForVersion(ver string) []string {
	aliases := config.ListAliases()
	var result []string
	for name, version := range aliases {
		if version == ver {
			result = append(result, name)
		}
	}
	return result
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all versions including RCs and betas")
	listCmd.Flags().BoolVarP(&listRemote, "remote", "r", false, "List remote (available) versions")
	listCmd.Flags().IntVarP(&listLimit, "number", "n", 30, "Limit number of versions shown (0 for all)")
}
