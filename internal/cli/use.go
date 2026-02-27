package cli

import (
	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/version"
)

var useCmd = &cobra.Command{
	Use:     "use <version|alias>",
	Aliases: []string{"switch", "select"},
	Short:   "Switch to a Go version",
	Long: `Switch to a specific Go version or alias.

If the version is not installed and auto-install is enabled,
it will be downloaded and installed automatically.

Examples:
  govm use 1.22.0             Switch to Go 1.22.0
  govm use stable             Switch to the stable alias
  govm use .                  Use version from go.mod/go.work
  g use 1.21.0                Short form`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := version.NewManager()
		if err != nil {
			return err
		}

		ver := args[0]

		// Special case: "." means use version from current directory
		if ver == "." {
			return mgr.UseFromProject("")
		}

		return mgr.Use(ver)
	},
}
