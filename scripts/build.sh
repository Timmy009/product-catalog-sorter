#!/bin/bash

# Product Catalog Sorting - Build Script
# Professional build automation with error handling

set -e  # Exit on any error
set -u  # Exit on undefined variables

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="catalog-sorter"
BUILD_DIR="bin"
VERSION=${VERSION:-"dev"}
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Main build function
build_app() {
    log_info "Starting build process..."
    
    # Create build directory
    mkdir -p ${BUILD_DIR}
    
    # Build flags
    LDFLAGS="-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME}"
    
    # Build for current platform
    log_info "Building ${APP_NAME} for current platform..."
    go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME} cmd/main.go
    
    if [ $? -eq 0 ]; then
        log_success "Build completed successfully"
        log_info "Binary location: ${BUILD_DIR}/${APP_NAME}"
        
        # Show binary info
        ls -lh ${BUILD_DIR}/${APP_NAME}
    else
        log_error "Build failed"
        exit 1
    fi
}

# Cross-compilation function
build_cross_platform() {
    log_info "Building for multiple platforms..."
    
    platforms=(
        "linux/amd64"
        "linux/arm64"
        "darwin/amd64"
        "darwin/arm64"
        "windows/amd64"
    )
    
    for platform in "${platforms[@]}"; do
        IFS='/' read -r GOOS GOARCH <<< "$platform"
        
        output_name="${APP_NAME}-${GOOS}-${GOARCH}"
        if [ "$GOOS" = "windows" ]; then
            output_name="${output_name}.exe"
        fi
        
        log_info "Building for ${GOOS}/${GOARCH}..."
        
        env GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags "${LDFLAGS}" \
            -o ${BUILD_DIR}/${output_name} \
            cmd/main.go
        
        if [ $? -eq 0 ]; then
            log_success "Built ${output_name}"
        else
            log_error "Failed to build for ${GOOS}/${GOARCH}"
        fi
    done
}

# Clean function
clean() {
    log_info "Cleaning build artifacts..."
    rm -rf ${BUILD_DIR}
    log_success "Clean completed"
}

# Help function
show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  build       Build for current platform (default)"
    echo "  cross       Build for multiple platforms"
    echo "  clean       Clean build artifacts"
    echo "  help        Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  VERSION     Set build version (default: dev)"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  VERSION=1.0.0 $0 build"
    echo "  $0 cross"
}

# Main execution
main() {
    case "${1:-build}" in
        "build")
            build_app
            ;;
        "cross")
            build_cross_platform
            ;;
        "clean")
            clean
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
