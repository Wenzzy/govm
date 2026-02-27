package shell

// ZshInit returns zsh initialization code for shell integration
func ZshInit() (string, error) {
	return `# govm shell integration for zsh
# Add this to your ~/.zshrc

export GOVM_ROOT="${GOVM_ROOT:-$HOME/.govm}"

# Ensure govm paths are in PATH
[[ ":$PATH:" != *":$GOVM_ROOT/bin:"* ]] && export PATH="$GOVM_ROOT/bin:$PATH"
[[ ":$PATH:" != *":$GOVM_ROOT/current/bin:"* ]] && export PATH="$GOVM_ROOT/current/bin:$PATH"

# Set GOROOT to current version
export GOROOT="$GOVM_ROOT/current"

# Add GOPATH/bin to PATH for go install packages
export GOPATH="${GOPATH:-$HOME/go}"
[[ ":$PATH:" != *":$GOPATH/bin:"* ]] && export PATH="$GOPATH/bin:$PATH"

# Auto-switch Go version based on go.mod/go.work
_govm_auto_switch() {
    local go_file=""
    local search_dir="$PWD"

    # Check current directory first
    if [[ -f "go.work" ]]; then
        go_file="$PWD/go.work"
    elif [[ -f "go.mod" ]]; then
        go_file="$PWD/go.mod"
    fi

    # If inherit_version is enabled, search parent directories
    if [[ -z "$go_file" && "$GOVM_INHERIT_VERSION" == "true" ]]; then
        while [[ "$search_dir" != "/" ]]; do
            search_dir="$(dirname "$search_dir")"
            if [[ -f "$search_dir/go.work" ]]; then
                go_file="$search_dir/go.work"
                break
            elif [[ -f "$search_dir/go.mod" ]]; then
                go_file="$search_dir/go.mod"
                break
            fi
        done
    fi

    if [[ -n "$go_file" ]]; then
        # Extract version from file
        local version=$(grep -E '^go [0-9]+\.[0-9]+' "$go_file" | head -1 | awk '{print $2}')

        if [[ -n "$version" ]]; then
            # Get current version
            local current=""
            if [[ -L "$GOVM_ROOT/current" ]]; then
                current=$(basename "$(dirname "$(readlink "$GOVM_ROOT/current")")")
            fi

            # Check if we need to switch (handle both X.Y and X.Y.Z formats)
            if [[ "$current" != "$version"* ]]; then
                # Try to find exact or compatible version
                local target_version=""
                if [[ -d "$GOVM_ROOT/versions/$version" ]]; then
                    target_version="$version"
                else
                    # Find latest patch version
                    target_version=$(ls -1 "$GOVM_ROOT/versions" 2>/dev/null | grep "^${version}" | sort -V | tail -1)
                fi

                if [[ -n "$target_version" && -d "$GOVM_ROOT/versions/$target_version" ]]; then
                    # Switch version silently
                    ln -sfn "$GOVM_ROOT/versions/$target_version/go" "$GOVM_ROOT/current" 2>/dev/null
                    print -P "%F{cyan}govm:%f switched to Go $target_version (from $go_file)"
                elif (( $+commands[govm] )); then
                    # Version not installed, try to install if auto_install is enabled
                    print -P "%F{yellow}govm:%f Go $version required (from $go_file), installing..."
                    govm install "$version" --default
                fi
            fi
        fi
    fi
}

# Hook into directory change
autoload -U add-zsh-hook
add-zsh-hook chpwd _govm_auto_switch

# Run on shell startup for current directory
_govm_auto_switch

# Completions
_govm() {
    local -a commands
    commands=(
        'install:Install a Go version'
        'uninstall:Uninstall a Go version'
        'use:Switch to a Go version'
        'list:List Go versions'
        'alias:Manage version aliases'
        'exec:Run command with specific Go version'
        'current:Show current Go version'
        'init:Initialize shell integration'
        'upgrade:Upgrade govm'
        'version:Print govm version'
        'setup:Interactive shell setup guide'
    )

    local -a installed_versions
    if [[ -d "$GOVM_ROOT/versions" ]]; then
        installed_versions=(${(f)"$(ls "$GOVM_ROOT/versions" 2>/dev/null)"})
    fi

    _arguments -C \
        '1: :->command' \
        '*: :->args'

    case $state in
        command)
            _describe -t commands 'govm commands' commands
            ;;
        args)
            case $words[2] in
                install|i|add)
                    _message 'version to install'
                    ;;
                uninstall|rm|remove|delete)
                    _describe -t versions 'installed versions' installed_versions
                    ;;
                use|switch|select)
                    _describe -t versions 'installed versions' installed_versions
                    ;;
                exec)
                    if (( CURRENT == 3 )); then
                        _describe -t versions 'installed versions' installed_versions
                    else
                        _command_names -e
                    fi
                    ;;
                alias)
                    if (( CURRENT == 3 )); then
                        local -a subcmds
                        subcmds=('rm:Remove an alias')
                        _describe -t subcmds 'alias commands' subcmds
                    fi
                    ;;
                init)
                    local -a shells
                    shells=('bash' 'zsh')
                    _describe -t shells 'shells' shells
                    ;;
            esac
            ;;
    esac
}

compdef _govm govm
compdef _govm g
`, nil
}
