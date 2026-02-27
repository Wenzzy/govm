package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wenzzy/govm/internal/config"
	"github.com/wenzzy/govm/internal/version"
)

var execCmd = &cobra.Command{
	Use:   "exec <version> <command> [args...]",
	Short: "Run a command with a specific Go version",
	Long: `Execute a command using a specific Go version without switching globally.

The specified Go version must be installed. If not, install it first with 'govm install'.

Examples:
  govm exec 1.21.0 go version        Run 'go version' with Go 1.21.0
  govm exec 1.22.0 go build ./...    Build project with Go 1.22.0
  govm exec 1.20.0 go test -v        Run tests with Go 1.20.0
  g exec 1.21.0 go run main.go       Short form`,
	Args:               cobra.MinimumNArgs(2),
	DisableFlagParsing: true, // Allow flags to be passed to the subcommand
	RunE: func(cmd *cobra.Command, args []string) error {
		ver := config.ResolveVersion(args[0])
		command := args[1:]

		mgr, err := version.NewManager()
		if err != nil {
			return err
		}

		// Check if version is installed
		if !mgr.IsInstalled(ver) {
			return fmt.Errorf("version %s is not installed. Run 'govm install %s' first", ver, ver)
		}

		// Get the Go binary path for this version
		goBinary, err := mgr.GetGoBinary(ver)
		if err != nil {
			return err
		}

		// Get the bin directory for this version
		binDir := filepath.Dir(goBinary)

		// Prepare environment
		env := os.Environ()
		env = updateEnv(env, "GOROOT", filepath.Dir(binDir))
		env = prependPath(env, binDir)

		// If the command is "go", use the specific binary
		cmdName := command[0]
		cmdArgs := command[1:]

		if cmdName == "go" {
			cmdName = goBinary
		} else {
			// Check if the command exists in the version's bin directory
			versionBinCmd := filepath.Join(binDir, cmdName)
			if _, err := os.Stat(versionBinCmd); err == nil {
				cmdName = versionBinCmd
			}
		}

		// Execute the command
		return executeCommand(cmdName, cmdArgs, env)
	},
}

// executeCommand runs a command with the given environment
func executeCommand(name string, args []string, env []string) error {
	// Try to use syscall.Exec for a cleaner process replacement
	// This replaces the current process with the new command
	binary, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("command not found: %s", name)
	}

	// Prepare args (first element must be the command name)
	argv := append([]string{name}, args...)

	// Execute and replace current process
	return syscall.Exec(binary, argv, env)
}

// updateEnv updates or adds an environment variable
func updateEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if len(e) > len(prefix) && e[:len(prefix)] == prefix {
			env[i] = key + "=" + value
			return env
		}
	}
	return append(env, key+"="+value)
}

// prependPath prepends a directory to the PATH environment variable
func prependPath(env []string, dir string) []string {
	for i, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			env[i] = "PATH=" + dir + string(os.PathListSeparator) + e[5:]
			return env
		}
	}
	return append(env, "PATH="+dir)
}
