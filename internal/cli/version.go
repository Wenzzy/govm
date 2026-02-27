package cli

import (
	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print govm version",
	Long:  `Print the version of govm CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintLogo()
		ui.PrintVersionInfo(config.Version, config.BuildTime)
		ui.PrintKeyValue("GitHub", config.RepoURL)
	},
}
