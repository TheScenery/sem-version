#!/bin/bash
set -e

# Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
U_RED='\033[4;31m'
NC='\033[0m'

print_header() {
    echo -e "\n${CYAN}=== $1 ===${NC}"
}

print_step() {
    echo -e "\n${YELLOW}> $1${NC}"
}

check_command() {
    echo -e "${GREEN}$@${NC}"
    "$@"
}

# 1. Build the tool
print_header "Building sem-version..."
go build -o sem-version
if [ ! -f "sem-version" ]; then
    echo -e "${U_RED}Failed to build sem-version${NC}"
    exit 1
fi

# 2. Setup demo environment
DEMO_DIR="demo_repo"
if [ -d "$DEMO_DIR" ]; then
    rm -rf "$DEMO_DIR"
fi
mkdir "$DEMO_DIR"
cd "$DEMO_DIR"

print_header "Setting up demo repository..."
check_command git init

# 3. Scenario 1: Initial Feature
print_header "Scenario 1: Initial Feature (Minor Bump)"
print_step "Commit: 'feat: initial project structure'"
echo "package main" > main.go
check_command git add .
check_command git commit -m "feat: initial project structure"

print_step "Running sem-version..."
../sem-version --verbose

# 4. Scenario 2: Bug Fix
print_header "Scenario 2: Bug Fix (Patch Bump)"
# Tag the previous version so we have a baseline (simulate release)
check_command git tag v0.1.0
print_step "Current Tag: v0.1.0"

print_step "Commit: 'fix: resolve startup crash'"
check_command git commit --allow-empty -m "fix: resolve startup crash"

print_step "Running sem-version..."
../sem-version --verbose

# 5. Scenario 3: New Feature
print_header "Scenario 3: New Feature (Minor Bump)"
# Simulate releasing v0.1.1
check_command git tag v0.1.1
print_step "Current Tag: v0.1.1"

print_step "Commit: 'feat: add user login'"
check_command git commit --allow-empty -m "feat: add user login"

print_step "Running sem-version..."
../sem-version --verbose

# 6. Scenario 4: Breaking Change
print_header "Scenario 4: Breaking Change (Major Bump)"
# Simulate releasing v0.2.0
check_command git tag v0.2.0
print_step "Current Tag: v0.2.0"

print_step "Commit: 'feat!: redesign api structure'"
check_command git commit --allow-empty -m "feat!: redesign api structure"

print_step "Running sem-version..."
../sem-version --verbose

# Cleanup
cd ..
# rm -rf "$DEMO_DIR"
print_header "Demo Complete!"
echo "You can inspect the '$DEMO_DIR' directory."
