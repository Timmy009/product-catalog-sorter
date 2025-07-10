#!/bin/bash

# Product Catalog Sorting - Test Script
# Comprehensive testing automation

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
COVERAGE_THRESHOLD=80
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"

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

# Run unit tests
run_unit_tests() {
    log_info "Running unit tests..."
    go test -v ./test/unit/...
    
    if [ $? -eq 0 ]; then
        log_success "Unit tests passed"
    else
        log_error "Unit tests failed"
        exit 1
    fi
}

# Run integration tests
run_integration_tests() {
    log_info "Running integration tests..."
    go test -v ./test/integration/...
    
    if [ $? -eq 0 ]; then
        log_success "Integration tests passed"
    else
        log_error "Integration tests failed"
        exit 1
    fi
}

# Run all tests
run_all_tests() {
    log_info "Running all tests..."
    go test -v ./...
    
    if [ $? -eq 0 ]; then
        log_success "All tests passed"
    else
        log_error "Some tests failed"
        exit 1
    fi
}

# Run tests with race detection
run_race_tests() {
    log_info "Running tests with race detection..."
    go test -race -v ./...
    
    if [ $? -eq 0 ]; then
        log_success "Race tests passed"
    else
        log_error "Race condition detected"
        exit 1
    fi
}

# Generate coverage report
generate_coverage() {
    log_info "Generating coverage report..."
    
    # Run tests with coverage
    go test -coverprofile=${COVERAGE_FILE} ./...
    
    if [ $? -ne 0 ]; then
        log_error "Coverage generation failed"
        exit 1
    fi
    
    # Generate HTML report
    go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}
    
    # Calculate coverage percentage
    COVERAGE=$(go tool cover -func=${COVERAGE_FILE} | grep total | awk '{print $3}' | sed 's/%//')
    
    log_info "Coverage: ${COVERAGE}%"
    
    # Check coverage threshold
    if (( $(echo "$COVERAGE >= $COVERAGE_THRESHOLD" | bc -l) )); then
        log_success "Coverage meets threshold (${COVERAGE_THRESHOLD}%)"
    else
        log_warning "Coverage below threshold: ${COVERAGE}% < ${COVERAGE_THRESHOLD}%"
    fi
    
    log_info "Coverage report: ${COVERAGE_HTML}"
}

# Run benchmarks
run_benchmarks() {
    log_info "Running benchmarks..."
    go test -bench=. -benchmem ./...
    
    if [ $? -eq 0 ]; then
        log_success "Benchmarks completed"
    else
        log_error "Benchmark execution failed"
        exit 1
    fi
}

# Clean test artifacts
clean_test_artifacts() {
    log_info "Cleaning test artifacts..."
    rm -f ${COVERAGE_FILE} ${COVERAGE_HTML}
    log_success "Test artifacts cleaned"
}

# Show help
show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  unit        Run unit tests only"
    echo "  integration Run integration tests only"
    echo "  all         Run all tests (default)"
    echo "  race        Run tests with race detection"
    echo "  coverage    Generate coverage report"
    echo "  bench       Run benchmarks"
    echo "  clean       Clean test artifacts"
    echo "  ci          Run full CI test suite"
    echo "  help        Show this help message"
}

# CI test suite
run_ci_suite() {
    log_info "Running full CI test suite..."
    
    # Run linting first
    if command -v golangci-lint &> /dev/null; then
        log_info "Running linter..."
        golangci-lint run
    else
        log_warning "golangci-lint not found, skipping linting"
        go vet ./...
        go fmt ./...
    fi
    
    # Run all tests
    run_all_tests
    
    # Run race detection
    run_race_tests
    
    # Generate coverage
    generate_coverage
    
    # Run benchmarks
    run_benchmarks
    
    log_success "CI test suite completed successfully"
}

# Main execution
main() {
    case "${1:-all}" in
        "unit")
            run_unit_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "all")
            run_all_tests
            ;;
        "race")
            run_race_tests
            ;;
        "coverage")
            generate_coverage
            ;;
        "bench")
            run_benchmarks
            ;;
        "clean")
            clean_test_artifacts
            ;;
        "ci")
            run_ci_suite
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
