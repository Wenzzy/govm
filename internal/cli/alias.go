package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

var aliasCmd = &cobra.Command{
	Use:   "alias [name] [version]",
	Short: "Manage version aliases",
	Long: `Create, list, or remove version aliases.

Aliases allow you to assign memorable names to specific Go versions.
Built-in aliases: stable, latest (auto-updated)

Examples:
  govm alias                   List all aliases
  govm alias dev 1.22.0        Create alias 'dev' for version 1.22.0
  govm alias rm dev            Remove alias 'dev'
  g alias lts 1.21.0           Short form`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			// List aliases
			return listAliases()
		case 1:
			// Show specific alias
			return showAlias(args[0])
		case 2:
			// Create or remove alias
			if args[0] == "rm" || args[0] == "remove" || args[0] == "delete" {
				return removeAlias(args[1])
			}
			return createAlias(args[0], args[1])
		default:
			return fmt.Errorf("too many arguments")
		}
	},
}

func listAliases() error {
	aliases := config.ListAliases()

	if len(aliases) == 0 {
		ui.PrintInfo("No aliases defined")
		ui.PrintHint("Create an alias with: govm alias <name> <version>")
		return nil
	}

	ui.PrintHeader("Aliases")

	table := ui.NewTable("Alias", "Version")
	for name, version := range aliases {
		displayVersion := version
		if version == "" {
			displayVersion = ui.Dim.Sprint("(not set)")
		}

		// Mark reserved aliases
		for _, reserved := range config.ReservedAliases {
			if name == reserved {
				name = ui.Cyan.Sprint(name) + ui.Dim.Sprint(" [builtin]")
				break
			}
		}

		table.AddRow(name, displayVersion)
	}
	table.Render()

	return nil
}

func showAlias(name string) error {
	version, ok := config.GetAlias(name)
	if !ok {
		return fmt.Errorf("alias '%s' not found", name)
	}

	if version == "" {
		ui.PrintInfo("Alias '%s' is not set", name)
	} else {
		ui.PrintKeyValue(name, version)
	}
	return nil
}

func createAlias(name, version string) error {
	if err := config.ValidateAliasName(name); err != nil {
		return err
	}

	// Normalize version
	version = config.NormalizeVersion(version)

	if err := config.SetAlias(name, version); err != nil {
		return err
	}

	ui.PrintSuccess("Created alias '%s' -> %s", name, version)
	return nil
}

func removeAlias(name string) error {
	// Check if alias exists
	if _, ok := config.GetAlias(name); !ok {
		return fmt.Errorf("alias '%s' not found", name)
	}

	// Warn about removing reserved aliases
	for _, reserved := range config.ReservedAliases {
		if name == reserved {
			ui.PrintWarning("Removing built-in alias '%s'", name)
			break
		}
	}

	if err := config.RemoveAlias(name); err != nil {
		return err
	}

	ui.PrintSuccess("Removed alias '%s'", name)
	return nil
}

// Subcommand for removing aliases
var aliasRmCmd = &cobra.Command{
	Use:     "rm <name>",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove an alias",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAlias(args[0])
	},
}

func init() {
	aliasCmd.AddCommand(aliasRmCmd)
}
