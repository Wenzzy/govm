package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/ui"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage govm configuration",
	Long: `View and modify govm configuration settings.

Available settings:
  auto_install     - Automatically install missing versions (true/false)
  inherit_version  - Search parent directories for go.mod/go.work (true/false)
  default_version  - Default Go version to use

Examples:
  govm config                           Show all settings
  govm config get auto_install          Get a specific setting
  govm config set inherit_version true  Enable version inheritance
  g config set auto_install false       Disable auto-install`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showConfig()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getConfig(args[0])
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setConfig(args[0], args[1])
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.GetPaths()
		if err != nil {
			return err
		}
		fmt.Println(paths.Config)
		return nil
	},
}

func showConfig() error {
	cfg := config.Get()

	ui.PrintHeader("Configuration")
	ui.PrintKeyValue("auto_install", formatBool(cfg.AutoInstall))
	ui.PrintKeyValue("inherit_version", formatBool(cfg.InheritVersion))
	ui.PrintKeyValue("default_version", formatString(cfg.DefaultVersion))

	paths, _ := config.GetPaths()
	ui.Println()
	ui.PrintHint(fmt.Sprintf("Config file: %s", paths.Config))

	return nil
}

func getConfig(key string) error {
	cfg := config.Get()

	switch strings.ToLower(key) {
	case "auto_install", "autoinstall":
		fmt.Println(cfg.AutoInstall)
	case "inherit_version", "inheritversion", "inherit":
		fmt.Println(cfg.InheritVersion)
	case "default_version", "defaultversion", "default":
		fmt.Println(cfg.DefaultVersion)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return nil
}

func setConfig(key, value string) error {
	cfg := config.Get()

	switch strings.ToLower(key) {
	case "auto_install", "autoinstall":
		b, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value for auto_install: %s (use true/false)", value)
		}
		cfg.AutoInstall = b
		ui.PrintSuccess("Set auto_install = %v", b)

	case "inherit_version", "inheritversion", "inherit":
		b, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value for inherit_version: %s (use true/false)", value)
		}
		cfg.InheritVersion = b
		ui.PrintSuccess("Set inherit_version = %v", b)
		if b {
			ui.PrintHint("govm will now search parent directories for go.mod/go.work")
		} else {
			ui.PrintHint("govm will only check the current directory for go.mod/go.work")
		}

	case "default_version", "defaultversion", "default":
		cfg.DefaultVersion = config.NormalizeVersion(value)
		ui.PrintSuccess("Set default_version = %s", cfg.DefaultVersion)

	default:
		return fmt.Errorf("unknown config key: %s\n\nAvailable keys: auto_install, inherit_version, default_version", key)
	}

	return config.Save(cfg)
}

func parseBool(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return strconv.ParseBool(s)
	}
}

func formatBool(b bool) string {
	if b {
		return ui.Green.Sprint("true")
	}
	return ui.Dim.Sprint("false")
}

func formatString(s string) string {
	if s == "" {
		return ui.Dim.Sprint("(not set)")
	}
	return s
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
}
