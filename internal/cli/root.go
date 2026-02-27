package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "govm",
	Short: "Go Version Manager",
	Long: `govm is a fast and efficient Go version manager.

It allows you to install, manage, and switch between multiple Go versions easily.
Use 'govm' or 'g' to access the tool.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand, show help
		cmd.Help()
	},
}

var versionFlag bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if versionFlag {
			ui.PrintVersionInfo(config.Version, config.BuildTime)
			os.Exit(0)
		}
	}

	// Add subcommands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.PrintError("%s", err)
		os.Exit(1)
	}
}
