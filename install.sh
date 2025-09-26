#!/bin/sh
#
# Terra Installation Script
#
# This script installs the latest terra binary from GitHub releases.
# It automatically detects your operating system and architecture.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh
#   wget -qO- https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh
#
# Or download and run locally:
#   wget https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh
#   chmod +x install.sh
#   ./install.sh
#
# Environment variables:
#   TERRA_INSTALL_DIR - installation directory (default: ~/.local/bin)
#   TERRA_VERSION     - specific version to install (default: latest)
#   TERRA_FORCE       - force installation even if already installed (default: false)
#   TERRA_DRY_RUN     - show what would be done without installing (default: false)
#
# Command line options:
#   --help            - show this help message
#   --version VER     - install specific version
#   --force           - force installation
#   --dry-run         - show what would be done without installing
#   --install-dir DIR - custom installation directory
#

set -e

# Script constants
TERRA_REPO_OWNER="rios0rios0"
TERRA_REPO_NAME="terra"
GITHUB_API_BASE="https://api.github.com"
GITHUB_RELEASE_BASE="https://github.com"

# Default values
DEFAULT_INSTALL_DIR="$HOME/.local/bin"
TERRA_INSTALL_DIR="${TERRA_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
TERRA_VERSION="${TERRA_VERSION:-latest}"
TERRA_FORCE="${TERRA_FORCE:-false}"
TERRA_DRY_RUN="${TERRA_DRY_RUN:-false}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Utility functions
info() {
    printf "${BLUE}INFO:${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}WARN:${NC} %s\n" "$1"
}

error() {
    printf "${RED}ERROR:${NC} %s\n" "$1" >&2
}

success() {
    printf "${GREEN}SUCCESS:${NC} %s\n" "$1"
}

# Help function
show_help() {
    cat << EOF
Terra Installation Script

This script installs the latest terra binary from GitHub releases.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --help              Show this help message
    --version VERSION   Install specific version (e.g. v1.0.0 or 1.0.0)
    --force             Force installation even if already installed
    --dry-run           Show what would be done without installing
    --install-dir DIR   Custom installation directory (default: ~/.local/bin)

ENVIRONMENT VARIABLES:
    TERRA_INSTALL_DIR   Installation directory (default: ~/.local/bin)
    TERRA_VERSION       Specific version to install (default: latest)
    TERRA_FORCE         Force installation (true/false, default: false)
    TERRA_DRY_RUN       Dry run mode (true/false, default: false)

EXAMPLES:
    # Install latest version
    $0

    # Install specific version
    $0 --version v1.0.0

    # Install to custom directory
    $0 --install-dir /usr/local/bin

    # Dry run to see what would happen
    $0 --dry-run

    # Force reinstallation
    $0 --force

EOF
}

# Parse command line arguments
parse_args() {
    while [ $# -gt 0 ]; do
        case $1 in
            --help|-h)
                show_help
                exit 0
                ;;
            --version)
                if [ -z "$2" ]; then
                    error "Version argument required"
                    exit 1
                fi
                TERRA_VERSION="$2"
                shift 2
                ;;
            --force)
                TERRA_FORCE="true"
                shift
                ;;
            --dry-run)
                TERRA_DRY_RUN="true"
                shift
                ;;
            --install-dir)
                if [ -z "$2" ]; then
                    error "Install directory argument required"
                    exit 1
                fi
                TERRA_INSTALL_DIR="$2"
                shift 2
                ;;
            *)
                error "Unknown option: $1"
                error "Use --help to see available options"
                exit 1
                ;;
        esac
    done
}

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        CYGWIN*|MINGW*|MSYS*)
                    echo "windows" ;;
        *)          error "Unsupported operating system: $(uname -s)"
                    exit 1 ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        i386|i686)      echo "386" ;;
        arm64|aarch64)  echo "arm64" ;;
        armv7l)         echo "arm" ;;
        armv6l)         echo "arm" ;;
        *)              error "Unsupported architecture: $(uname -m)"
                        exit 1 ;;
    esac
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if curl or wget is available
check_download_tools() {
    if command_exists curl; then
        DOWNLOAD_CMD="curl"
    elif command_exists wget; then
        DOWNLOAD_CMD="wget"
    else
        error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download file using available tool
download_file() {
    local url="$1"
    local output="$2"
    
    if [ "$DOWNLOAD_CMD" = "curl" ]; then
        curl -fsSL -o "$output" "$url"
    else
        wget -q -O "$output" "$url"
    fi
}

# Get latest release information from GitHub API
get_latest_release() {
    local api_url="${GITHUB_API_BASE}/repos/${TERRA_REPO_OWNER}/${TERRA_REPO_NAME}/releases/latest"
    local temp_file
    
    temp_file=$(mktemp)
    
    if ! download_file "$api_url" "$temp_file"; then
        rm -f "$temp_file"
        error "Failed to fetch release information from GitHub API"
        exit 1
    fi
    
    # Extract tag_name from JSON (simple parsing without jq)
    local tag_name
    tag_name=$(grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' "$temp_file" | cut -d'"' -f4)
    
    if [ -z "$tag_name" ]; then
        rm -f "$temp_file"
        error "Could not parse release information"
        exit 1
    fi
    
    rm -f "$temp_file"
    echo "$tag_name"
}

# Get download URL for specific version and platform
get_download_url() {
    local version="$1"
    local os="$2"
    local arch="$3"
    
    # Remove 'v' prefix if present for consistency
    version=$(echo "$version" | sed 's/^v//')
    
    # Construct expected asset name (matches terra's self-update logic)
    local asset_name="terra_${os}_${arch}"
    
    # For latest version, use the releases/latest API endpoint
    if [ "$version" = "latest" ]; then
        local api_url="${GITHUB_API_BASE}/repos/${TERRA_REPO_OWNER}/${TERRA_REPO_NAME}/releases/latest"
    else
        local api_url="${GITHUB_API_BASE}/repos/${TERRA_REPO_OWNER}/${TERRA_REPO_NAME}/releases/tags/v${version}"
    fi
    
    local temp_file
    temp_file=$(mktemp)
    
    if ! download_file "$api_url" "$temp_file"; then
        rm -f "$temp_file"
        error "Failed to fetch release information for version $version"
        exit 1
    fi
    
    # Extract download URL for the specific asset
    local download_url
    download_url=$(grep -A 1 "\"name\"[[:space:]]*:[[:space:]]*\"${asset_name}\"" "$temp_file" | \
                   grep "browser_download_url" | \
                   cut -d'"' -f4)
    
    if [ -z "$download_url" ]; then
        rm -f "$temp_file"
        error "No binary found for platform ${os}_${arch} in version $version"
        error "Available assets:"
        grep '"name"' "$temp_file" | cut -d'"' -f4 | sed 's/^/  /'
        exit 1
    fi
    
    rm -f "$temp_file"
    echo "$download_url"
}

# Check if terra is already installed
check_existing_installation() {
    if [ -f "$TERRA_INSTALL_DIR/terra" ]; then
        if [ "$TERRA_FORCE" = "false" ]; then
            local current_version
            current_version=$("$TERRA_INSTALL_DIR/terra" --version 2>/dev/null | head -n1 | cut -d' ' -f3 || echo "unknown")
            warn "terra is already installed at $TERRA_INSTALL_DIR/terra (version: $current_version)"
            warn "Use --force to reinstall"
            return 1
        else
            info "Forcing reinstallation (--force specified)"
        fi
    fi
    return 0
}

# Create installation directory
create_install_dir() {
    if [ ! -d "$TERRA_INSTALL_DIR" ]; then
        if [ "$TERRA_DRY_RUN" = "true" ]; then
            info "[DRY RUN] Would create directory: $TERRA_INSTALL_DIR"
        else
            info "Creating installation directory: $TERRA_INSTALL_DIR"
            mkdir -p "$TERRA_INSTALL_DIR"
        fi
    fi
}

# Install terra binary
install_terra() {
    local download_url="$1"
    local os="$2"
    local arch="$3"
    local version="$4"
    
    info "Installing terra for ${os}/${arch} (version: $version)"
    info "Installation directory: $TERRA_INSTALL_DIR"
    
    if [ "$TERRA_DRY_RUN" = "true" ]; then
        info "[DRY RUN] Would download from: $download_url"
        info "[DRY RUN] Would install to: $TERRA_INSTALL_DIR/terra"
        return 0
    fi
    
    # Create temporary file for download
    local temp_file
    temp_file=$(mktemp)
    
    info "Downloading terra binary..."
    if ! download_file "$download_url" "$temp_file"; then
        rm -f "$temp_file"
        error "Failed to download terra binary"
        exit 1
    fi
    
    # Verify download
    if [ ! -s "$temp_file" ]; then
        rm -f "$temp_file"
        error "Downloaded file is empty"
        exit 1
    fi
    
    # Move to installation directory and make executable
    info "Installing binary to $TERRA_INSTALL_DIR/terra"
    mv "$temp_file" "$TERRA_INSTALL_DIR/terra"
    chmod +x "$TERRA_INSTALL_DIR/terra"
    
    success "terra has been successfully installed!"
}

# Verify installation
verify_installation() {
    if [ "$TERRA_DRY_RUN" = "true" ]; then
        info "[DRY RUN] Would verify installation at: $TERRA_INSTALL_DIR/terra"
        return 0
    fi
    
    if [ -x "$TERRA_INSTALL_DIR/terra" ]; then
        local installed_version
        installed_version=$("$TERRA_INSTALL_DIR/terra" --version 2>/dev/null | head -n1 || echo "unknown")
        success "Installation verified: $installed_version"
        
        # Check if install directory is in PATH
        case ":$PATH:" in
            *":$TERRA_INSTALL_DIR:"*) ;;
            *)
                warn "Installation directory $TERRA_INSTALL_DIR is not in your PATH"
                info "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
                info "  export PATH=\"\$PATH:$TERRA_INSTALL_DIR\""
                info ""
                info "Or run terra with the full path:"
                info "  $TERRA_INSTALL_DIR/terra --help"
                ;;
        esac
    else
        error "Installation verification failed"
        exit 1
    fi
}

# Main installation function
main() {
    info "Terra Installation Script"
    info "========================="
    
    # Parse command line arguments
    parse_args "$@"
    
    # Check prerequisites
    check_download_tools
    
    # Detect platform
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)
    info "Detected platform: ${os}/${arch}"
    
    # Get version to install
    local version="$TERRA_VERSION"
    local download_url=""
    
    if [ "$TERRA_DRY_RUN" = "true" ]; then
        # In dry run mode, use mock data
        if [ "$version" = "latest" ]; then
            version="v1.0.0"
            info "[DRY RUN] Would fetch latest release information..."
            info "[DRY RUN] Mock latest version: $version"
        else
            info "[DRY RUN] Would install version: $version"
        fi
        download_url="https://github.com/${TERRA_REPO_OWNER}/${TERRA_REPO_NAME}/releases/download/${version}/terra_${os}_${arch}"
        info "[DRY RUN] Mock download URL: $download_url"
    else
        if [ "$version" = "latest" ]; then
            info "Fetching latest release information..."
            version=$(get_latest_release)
            info "Latest version: $version"
        else
            info "Installing version: $version"
        fi
        
        # Get download URL
        info "Getting download URL..."
        download_url=$(get_download_url "$version" "$os" "$arch")
    fi
    
    # Check existing installation
    if ! check_existing_installation; then
        exit 0
    fi
    
    # Create installation directory
    create_install_dir
    
    # Install terra
    install_terra "$download_url" "$os" "$arch" "$version"
    
    # Verify installation
    verify_installation
    
    info ""
    success "Installation complete!"
    if [ "$TERRA_DRY_RUN" = "false" ]; then
        info "Run 'terra --help' to get started"
    fi
}

# Run the main function with all arguments
main "$@"