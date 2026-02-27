package cli

import (
	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/ui"
	"github.com/wenzzy/govm/internal/version"
)

var forceUninstall bool

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <version>",
	Aliases: []string{"rm", "remove", "delete"},
	Short:   "Uninstall a Go version",
	Long: `Uninstall a specific Go version.

Examples:
  govm uninstall 1.21.0       Uninstall Go 1.21.0
  govm rm 1.20.0              Short form
  g rm 1.19.0 -f              Force uninstall without confirmation`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := version.NewManager()
		if err != nil {
			return err
		}

		ver := args[0]

		// Check if it's the current version
		current, _ := mgr.Current()
		if current == ver && !forceUninstall {
			if !ui.Confirm(ui.Warning.Sprintf("Version %s is currently in use. Uninstall anyway?", ver)) {
				ui.PrintInfo("Aborted")
				return nil
			}
		}

		return mgr.Uninstall(ver)
	},
}

func init() {
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Force uninstall without confirmation")
}
