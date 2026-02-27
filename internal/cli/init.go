package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/shell"
	"github.com/wenzzy/govm/internal/ui"
)

var initCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Initialize shell integration",
	Long: `Output shell integration code for automatic version switching.

Supported shells: bash, zsh

Add this to your shell configuration file:

For Bash (~/.bashrc):
  eval "$(govm init bash)"

For Zsh (~/.zshrc):
  eval "$(govm init zsh)"

Examples:
  govm init bash               Show bash integration
  govm init zsh                Show zsh integration
  g init zsh                   Short form`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		shellName := args[0]

		var code string
		var err error

		switch shellName {
		case "bash":
			code, err = shell.BashInit()
		case "zsh":
			code, err = shell.ZshInit()
		default:
			return fmt.Errorf("unsupported shell: %s (supported: bash, zsh)", shellName)
		}

		if err != nil {
			return err
		}

		// Export config-based environment variables
		cfg := config.Get()
		if cfg.InheritVersion {
			fmt.Println(`export GOVM_INHERIT_VERSION="true"`)
		} else {
			fmt.Println(`export GOVM_INHERIT_VERSION="false"`)
		}

		// Output the shell code (will be eval'd)
		fmt.Print(code)
		return nil
	},
}

// Also add a setup command for guided setup
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive shell setup guide",
	Long: `Interactive guide to set up shell integration.

This command will help you configure your shell for automatic
Go version switching when entering directories with go.mod files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintHeader("Shell Integration Setup")

		ui.Println()
		ui.PrintInfo("Add the following to your shell configuration file:")
		ui.Println()

		ui.PrintBullet("For Bash (~/.bashrc or ~/.bash_profile):")
		ui.PrintCommand(`export GOVM_ROOT="$HOME/.govm"`)
		ui.PrintCommand(`export PATH="$GOVM_ROOT/bin:$GOVM_ROOT/current/bin:$PATH"`)
		ui.PrintCommand(`eval "$(govm init bash)"`)
		ui.Println()

		ui.PrintBullet("For Zsh (~/.zshrc):")
		ui.PrintCommand(`export GOVM_ROOT="$HOME/.govm"`)
		ui.PrintCommand(`export PATH="$GOVM_ROOT/bin:$GOVM_ROOT/current/bin:$PATH"`)
		ui.PrintCommand(`eval "$(govm init zsh)"`)
		ui.Println()

		ui.PrintHint("After adding, restart your terminal or run: source ~/.bashrc (or ~/.zshrc)")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
