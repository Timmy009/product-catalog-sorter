#!/bin/bash

# Product Catalog Sorting - Lint Script
# Code quality and formatting automation

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# Go formatting
run_gofmt() {
    log_info "Running gofmt..."
    
    # Check if files need formatting
    UNFORMATTED=$(gofmt -l .)
    
    if [ -n "$UNFORMATTED" ]; then
        log_warning "The following files need formatting:"
        echo "$UNFORMATTED"
        
        # Auto-format if requested
        if [ "${1:-}" = "fix" ]; then
            log_info "Auto-formatting files..."
            gofmt -w .
            log_success "Files formatted"
        else
            log_error "Files need formatting. Run with 'fix' to auto-format."
            exit 1
        fi
    else
        log_success "All files are properly formatted"
    fi
}

# Go imports
run_goimports() {
    if ! command -v goimports &> /dev/null; then
        log_warning "goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"
        return
    fi
    
    log_info "Running goimports..."
    
    # Check imports
    UNFORMATTED=$(goimports -l .)
    
    if [ -n "$UNFORMATTED" ]; then
        log_warning "The following files have import issues:"
        echo "$UNFORMATTED"
        
        if [ "${1:-}" = "fix" ]; then
            log_info "Fixing imports..."
            goimports -w .
            log_success "Imports fixed"
        else
            log_error "Import issues found. Run with 'fix' to auto-fix."
            exit 1
        fi
    else
        log_success "All imports are properly organized"
    fi
}

# Go vet
run_govet() {
    log_info "Running go vet..."
    
    go vet ./...
    
    if [ $? -eq 0 ]; then
        log_success "go vet passed"
    else
        log_error "go vet found issues"
        exit 1
    fi
}

# golangci-lint
run_golangci_lint() {
    if ! command -v golangci-lint &> /dev/null; then
        log_warning "golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"
        return
    fi
    
    log_info "Running golangci-lint..."
    
    golangci-lint run
    
    if [ $? -eq 0 ]; then
        log_success "golangci-lint passed"
    else
        log_error "golangci-lint found issues"
        exit 1
    fi
}

# Check for common issues
run_custom_checks() {
    log_info "Running custom checks..."
    
    # Check for TODO/FIXME comments
    TODOS=$(grep -r "TODO\|FIXME" --include="*.go" . || true)
    if [ -n "$TODOS" ]; then
        log_warning "Found TODO/FIXME comments:"
        echo "$TODOS"
    fi
    
    # Check for debug prints
    DEBUG_PRINTS=$(grep -r "fmt.Print\|log.Print" --include="*.go" . | grep -v "_test.go" || true)
    if [ -n "$DEBUG_PRINTS" ]; then
        log_warning "Found potential debug prints:"
        echo "$DEBUG_PRINTS"
    fi
    
    # Check for hardcoded strings that should be constants
    MAGIC_NUMBERS=$(grep -r "magic\|hardcoded" --include="*.go" . || true)
    if [ -n "$MAGIC_NUMBERS" ]; then
        log_info "Consider reviewing hardcoded values"
    fi
    
    log_success "Custom checks completed"
}

# Run all linting
run_all_linting() {
    log_info "Running complete linting suite..."
    
    run_gofmt "$1"
    run_goimports "$1"
    run_govet
    run_golangci_lint
    run_custom_checks
    
    log_success "All linting checks passed"
}

# Show help
show_help() {
    echo "Usage: $0 [OPTION] [fix]"
    echo ""
    echo "Options:"
    echo "  fmt         Run gofmt"
    echo "  imports     Run goimports"
    echo "  vet         Run go vet"
    echo "  golangci    Run golangci-lint"
    echo "  custom      Run custom checks"
    echo "  all         Run all linting (default)"
    echo "  help        Show this help message"
    echo ""
    echo "Modifiers:"
    echo "  fix         Auto-fix issues where possible"
    echo ""
    echo "Examples:"
    echo "  $0 all"
    echo "  $0 fmt fix"
    echo "  $0 all fix"
}

# Main execution
main() {
    case "${1:-all}" in
        "fmt")
            run_gofmt "$2"
            ;;
        "imports")
            run_goimports "$2"
            ;;
        "vet")
            run_govet
            ;;
        "golangci")
            run_golangci_lint
            ;;
        "custom")
            run_custom_checks
            ;;
        "all")
            run_all_linting "$2"
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

main "$@"
