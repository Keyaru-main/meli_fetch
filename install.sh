#!/bin/bash
# Melifetch Installation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_header() {
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "$1"
    echo "═══════════════════════════════════════════════════════════"
    echo ""
}

# Check prerequisites
check_prerequisites() {
    log_header "🔍 Checking Prerequisites"
    
    # Check Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        log_info "Please install Go 1.25 or higher from https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go $GO_VERSION found"
    
    # Check Git
    if ! command -v git &> /dev/null; then
        log_error "Git is not installed"
        log_info "Please install Git from https://git-scm.com/"
        exit 1
    fi
    
    log_success "Git found"
}

# Install
install_melifetch() {
    log_header "📦 Installing MeliFetch"
    
    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    log_info "Cloning repository..."
    git clone https://github.com/Keyaru-main/melli_pet.git
    cd melli_pet
    
    log_info "Building binary..."
    go build -o melifetch cmd/melifetch/main.go
    
    # Install binary
    INSTALL_DIR="/usr/local/bin"
    
    if [ -w "$INSTALL_DIR" ]; then
        mv melifetch "$INSTALL_DIR/"
        log_success "Installed to $INSTALL_DIR/melifetch"
    else
        log_warning "No write permission to $INSTALL_DIR"
        log_info "Installing with sudo..."
        sudo mv melifetch "$INSTALL_DIR/"
        log_success "Installed to $INSTALL_DIR/melifetch"
    fi
    
    # Cleanup
    cd ~
    rm -rf "$TEMP_DIR"
    
    log_success "Installation complete!"
}

# Configure
configure_melifetch() {
    log_header "⚙️  Configuration"
    
    log_info "You need to configure melifetch with your GitHub token and repository"
    echo ""
    echo "Steps:"
    echo "1. Create a GitHub Personal Access Token:"
    echo "   https://github.com/settings/tokens/new"
    echo "   Required scopes: repo, workflow"
    echo ""
    echo "2. Create a repository for downloads (e.g., 'melli-downloads')"
    echo ""
    echo "3. Copy the contents of 'repo/' directory to your repository"
    echo ""
    echo "4. Run the following commands:"
    echo "   melifetch config --token YOUR_TOKEN"
    echo "   melifetch config --repo username/repo-name"
    echo ""
    
    read -p "Do you want to configure now? (y/n) " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter your GitHub token: " token
        read -p "Enter your repository (username/repo): " repo
        
        melifetch config --token "$token"
        melifetch config --repo "$repo"
        
        log_success "Configuration saved!"
    else
        log_info "You can configure later with: melifetch config"
    fi
}

# Test installation
test_installation() {
    log_header "🧪 Testing Installation"
    
    if command -v melifetch &> /dev/null; then
        log_success "melifetch is available"
        
        VERSION=$(melifetch --version 2>&1 || echo "unknown")
        log_info "Version: $VERSION"
        
        return 0
    else
        log_error "melifetch is not available"
        return 1
    fi
}

# Main
main() {
    log_header "🚀 MeliFetch Installation"
    
    check_prerequisites
    install_melifetch
    
    if test_installation; then
        configure_melifetch
        
        log_header "✅ Installation Successful!"
        echo "Next steps:"
        echo "1. Configure melifetch: melifetch config"
        echo "2. Try your first download: melifetch fetch https://httpbin.org/json"
        echo "3. Read the documentation: https://github.com/Keyaru-main/melli_pet"
        echo ""
        log_success "Happy downloading! 🎉"
    else
        log_error "Installation failed"
        exit 1
    fi
}

# Run
main "$@"
