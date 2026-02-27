# govm

Go version manager. Install, switch, and manage multiple Go versions.

## Why govm over gvm?

- **Single binary, zero dependencies** — no bash framework, no rvm, just one Go binary
- **Symlink-based switching** — versions switch by updating one symlink (`~/.govm/current`), not by rewriting `PATH` or shell variables. Works correctly in any context: subshells, scripts, AI agents
- **`cd` is not broken** — shell hook wraps `cd` cleanly, so tools that use `builtin cd` or spawn subprocesses still see the correct Go version
- **Auto-switch from `go.mod` / `go.work`** — detects and switches version when you enter a project directory
- **`govm exec`** — run a command with a specific Go version without switching globally
- **Aliases** — map names like `stable` or `dev` to versions, use them anywhere

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/wenzzy/govm/main/scripts/install.sh | bash
```

Or build from source:

```bash
git clone https://github.com/wenzzy/govm.git
cd govm
make install
```

### Shell setup

Add to `~/.bashrc` or `~/.zshrc`:

```bash
export GOVM_ROOT="$HOME/.govm"
export PATH="$GOVM_ROOT/bin:$GOVM_ROOT/current/bin:$PATH"
eval "$(govm init bash)"  # or: eval "$(govm init zsh)"
```

## Usage

```bash
govm install 1.22.0          # Install a version
govm install latest           # Install latest stable
govm use 1.22.0               # Switch version
govm use .                    # Use version from go.mod
govm list                     # List installed versions
govm list remote              # List available versions
govm alias dev 1.23.0         # Create alias
govm exec 1.21.0 go test ./.. # Run with specific version
govm current                  # Show active version
```

## Commands

| Command | Aliases | Description |
| --- | --- | --- |
| `govm install <version>` | `i`, `add` | Install a Go version |
| `govm uninstall <version>` | `rm`, `remove` | Remove a Go version |
| `govm use <version>` | `switch`, `select` | Switch active version |
| `govm list` | `ls` | List versions |
| `govm alias [name] [version]` | | Manage aliases |
| `govm exec <ver> <cmd>` | | Run command with version |
| `govm current` | `now` | Show current version |
| `govm config [get\|set]` | | Manage configuration |
| `govm upgrade` | | Upgrade govm |

## Configuration

`~/.govm/config.toml`:

```toml
default_version = "1.22.0"
auto_install = true
inherit_version = false

[aliases]
stable = "1.22.0"
dev = "1.23.0"
```

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| `default_version` | string | `""` | Go version used when no project-specific version is detected |
| `auto_install` | bool | `true` | Automatically install a missing version when `govm use` or auto-switch requires it |
| `inherit_version` | bool | `false` | Search parent directories for `go.mod`/`go.work`. When `false`, only the current directory is checked |

### Aliases

The `[aliases]` section maps short names to specific versions. Use them anywhere a version is expected:

```bash
govm alias stable 1.22.0   # create
govm use stable             # use
govm alias rm stable        # remove
```

`stable` and `latest` are reserved — they auto-resolve to the newest stable version when set to `""`.

## Uninstall

```bash
rm -rf ~/.govm
# Remove shell integration lines from ~/.bashrc or ~/.zshrc
```

## License

[AGPL-3.0](LICENSE)
