package cli

import (
	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/ui"
	"github.com/wenzzy/govm/internal/version"
)

var (
	installDefault bool
)

var installCmd = &cobra.Command{
	Use:     "install <version>",
	Aliases: []string{"i", "add"},
	Short:   "Install a Go version",
	Long: `Install a specific Go version.

Examples:
  govm install 1.22.0         Install Go 1.22.0
  govm install 1.22.0 -d      Install and set as default
  govm install latest         Install the latest stable version
  g i 1.21.0                  Short form`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := version.NewManager()
		if err != nil {
			return err
		}

		ver := args[0]

		// Handle special aliases
		if ver == "latest" || ver == "stable" {
			latestVer, err := version.GetLatestStable()
			if err != nil {
				return err
			}
			ui.PrintInfo("Latest stable version: %s", latestVer)
			ver = latestVer
		}

		return mgr.Install(ver, installDefault, true)
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installDefault, "default", "d", false, "Set as default version after install")
}
