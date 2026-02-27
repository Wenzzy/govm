package shell

// BashInit returns bash initialization code for shell integration
func BashInit() (string, error) {
	return `# govm shell integration for bash
# Add this to your ~/.bashrc

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
                    echo -e "\033[0;36mgovm:\033[0m switched to Go $target_version (from $go_file)"
                elif command -v govm &>/dev/null; then
                    # Version not installed, try to install if auto_install is enabled
                    echo -e "\033[0;33mgovm:\033[0m Go $version required (from $go_file), installing..."
                    govm install "$version" --default
                fi
            fi
        fi
    fi
}

# Hook into cd command
_govm_cd() {
    builtin cd "$@" && _govm_auto_switch
}

# Create alias for cd
alias cd='_govm_cd'

# Run on shell startup for current directory
_govm_auto_switch

# Completions
_govm_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local cmd="${COMP_WORDS[1]}"

    case "$cmd" in
        install|i|add)
            # Complete with remote versions (cached)
            COMPREPLY=()
            ;;
        uninstall|rm|remove|delete|use|switch|select|exec)
            # Complete with installed versions
            if [[ -d "$GOVM_ROOT/versions" ]]; then
                COMPREPLY=($(compgen -W "$(ls "$GOVM_ROOT/versions" 2>/dev/null)" -- "$cur"))
            fi
            ;;
        alias)
            if [[ "${COMP_WORDS[2]}" == "rm" ]]; then
                # Complete with existing aliases
                COMPREPLY=()
            fi
            ;;
        init)
            COMPREPLY=($(compgen -W "bash zsh" -- "$cur"))
            ;;
        *)
            COMPREPLY=($(compgen -W "install uninstall use list alias exec current init upgrade version setup" -- "$cur"))
            ;;
    esac
}

complete -F _govm_completions govm
complete -F _govm_completions g
`, nil
}
