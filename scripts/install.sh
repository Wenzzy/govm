#!/usr/bin/env bash
#
# govm installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/wenzzy/govm/main/scripts/install.sh | bash
#

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
GOVM_ROOT="${GOVM_ROOT:-$HOME/.govm}"
GITHUB_REPO="wenzzy/govm"
BINARY_NAME="govm"

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)  os="linux" ;;
        Darwin*) os="darwin" ;;
        *)       echo -e "${RED}Unsupported OS: $(uname -s)${NC}"; exit 1 ;;
    esac

    case "$(uname -m)" in
        x86_64)  arch="amd64" ;;
        amd64)   arch="amd64" ;;
        arm64)   arch="arm64" ;;
        aarch64) arch="arm64" ;;
        *)       echo -e "${RED}Unsupported architecture: $(uname -m)${NC}"; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# Get latest release version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" 2>/dev/null \
        | grep -o '"tag_name":"[^"]*"' | head -1 | cut -d'"' -f4
}

# Download and install govm
install_govm() {
    local platform version download_url temp_file

    echo -e "${CYAN}${BOLD}"
    echo "   __ _  _____   ___ __ ___"
    echo "  / _\` |/ _ \\ \\ / / '_ \` _ \\"
    echo " | (_| | (_) \\ V /| | | | | |"
    echo "  \\__, |\\___/ \\_/ |_| |_| |_|"
    echo "   __/ |"
    echo "  |___/   Go Version Manager"
    echo -e "${NC}"
    echo

    # Detect platform
    echo -e "${CYAN}Detecting platform...${NC}"
    platform=$(detect_platform)
    echo -e "  Platform: ${GREEN}${platform}${NC}"

    # Get latest version
    echo -e "${CYAN}Fetching latest version...${NC}"
    version=$(get_latest_version)

    if [[ -z "$version" ]]; then
        echo -e "${YELLOW}Could not determine latest version, building from source...${NC}"
        install_from_source
        return
    fi

    echo -e "  Version: ${GREEN}${version}${NC}"

    # Create directories
    echo -e "${CYAN}Creating directories...${NC}"
    mkdir -p "$GOVM_ROOT/bin"
    mkdir -p "$GOVM_ROOT/versions"
    mkdir -p "$GOVM_ROOT/cache"

    # Download binary
    download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${BINARY_NAME}-${platform}"
    echo -e "${CYAN}Downloading govm...${NC}"
    echo -e "  URL: ${download_url}"

    temp_file=$(mktemp)
    if ! curl -fsSL "$download_url" -o "$temp_file" 2>/dev/null; then
        echo -e "${YELLOW}Pre-built binary not available, building from source...${NC}"
        rm -f "$temp_file"
        install_from_source
        return
    fi

    # Install binary
    echo -e "${CYAN}Installing...${NC}"
    chmod +x "$temp_file"
    mv "$temp_file" "$GOVM_ROOT/bin/$BINARY_NAME"

    # Create 'g' symlink
    ln -sf "$GOVM_ROOT/bin/$BINARY_NAME" "$GOVM_ROOT/bin/g"

    echo -e "${GREEN}${BOLD}Installation complete!${NC}"
    echo
    setup_shell
}

# Build from source
install_from_source() {
    echo -e "${CYAN}Building from source...${NC}"

    # Check for Go
    if ! command -v go &>/dev/null; then
        echo -e "${RED}Go is required to build from source.${NC}"
        echo -e "Please install Go first: https://go.dev/dl/"
        exit 1
    fi

    # Clone and build
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    echo -e "  Cloning repository..."
    git clone --depth 1 "https://github.com/${GITHUB_REPO}.git" govm

    cd govm
    echo -e "  Building..."
    go build -o "$GOVM_ROOT/bin/$BINARY_NAME" ./cmd/govm

    # Create 'g' symlink
    ln -sf "$GOVM_ROOT/bin/$BINARY_NAME" "$GOVM_ROOT/bin/g"

    # Cleanup
    cd /
    rm -rf "$temp_dir"

    echo -e "${GREEN}${BOLD}Build complete!${NC}"
    echo
    setup_shell
}

# Setup shell integration
setup_shell() {
    echo -e "${CYAN}${BOLD}Shell Setup${NC}"
    echo
    echo "Add the following to your shell configuration file:"
    echo

    local shell_name=$(basename "$SHELL")

    case "$shell_name" in
        bash)
            echo -e "${YELLOW}# Add to ~/.bashrc or ~/.bash_profile:${NC}"
            echo -e "${GREEN}export GOVM_ROOT=\"\$HOME/.govm\"${NC}"
            echo -e "${GREEN}export PATH=\"\$GOVM_ROOT/bin:\$GOVM_ROOT/current/bin:\$PATH\"${NC}"
            echo -e "${GREEN}eval \"\$(govm init bash)\"${NC}"
            ;;
        zsh)
            echo -e "${YELLOW}# Add to ~/.zshrc:${NC}"
            echo -e "${GREEN}export GOVM_ROOT=\"\$HOME/.govm\"${NC}"
            echo -e "${GREEN}export PATH=\"\$GOVM_ROOT/bin:\$GOVM_ROOT/current/bin:\$PATH\"${NC}"
            echo -e "${GREEN}eval \"\$(govm init zsh)\"${NC}"
            ;;
        *)
            echo -e "${YELLOW}# Add to your shell config:${NC}"
            echo -e "${GREEN}export GOVM_ROOT=\"\$HOME/.govm\"${NC}"
            echo -e "${GREEN}export PATH=\"\$GOVM_ROOT/bin:\$GOVM_ROOT/current/bin:\$PATH\"${NC}"
            echo -e "${GREEN}eval \"\$(govm init bash)\"  # or zsh${NC}"
            ;;
    esac

    echo
    echo -e "${CYAN}Then restart your terminal or run:${NC}"
    echo -e "  source ~/.${shell_name}rc"
    echo
    echo -e "${CYAN}Quick start:${NC}"
    echo -e "  govm install latest    # Install latest Go version"
    echo -e "  govm list              # List installed versions"
    echo -e "  govm use 1.22.0        # Switch to a version"
    echo -e "  g --help               # Short alias"
    echo
    echo -e "${GREEN}${BOLD}Happy coding!${NC}"
}

# Uninstall govm
uninstall_govm() {
    echo -e "${CYAN}Uninstalling govm...${NC}"

    if [[ -d "$GOVM_ROOT" ]]; then
        read -p "This will remove $GOVM_ROOT and all installed Go versions. Continue? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$GOVM_ROOT"
            echo -e "${GREEN}govm has been uninstalled.${NC}"
            echo -e "Don't forget to remove the shell integration from your config file."
        else
            echo "Aborted."
        fi
    else
        echo "govm is not installed."
    fi
}

# Main
main() {
    case "${1:-}" in
        --uninstall|-u)
            uninstall_govm
            ;;
        --help|-h)
            echo "govm installer"
            echo
            echo "Usage:"
            echo "  install.sh           Install govm"
            echo "  install.sh --uninstall   Uninstall govm"
            echo "  install.sh --help        Show this help"
            ;;
        *)
            install_govm
            ;;
    esac
}

main "$@"
