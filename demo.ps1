<#
.SYNOPSIS
    Demo script for sem-version tool.
    Run this script to see sem-version in action!
#>

$ErrorActionPreference = "Stop"

# Colors
$Cyan = "Cyan"
$Green = "Green"
$Yellow = "Yellow"

function Print-Header {
    param($Text)
    Write-Host "`n=== $Text ===" -ForegroundColor $Cyan
}

function Print-Step {
    param($Text)
    Write-Host "`n> $Text" -ForegroundColor $Yellow
}

function Check-Command {
    param($Cmd, $ArgsArr)
    Write-Host "$Cmd $ArgsArr" -ForegroundColor $Green
    & $Cmd $ArgsArr
}

# 1. Build the tool
Print-Header "Building sem-version..."
go build -o sem-version.exe
if (-not (Test-Path "sem-version.exe")) {
    Write-Error "Failed to build sem-version.exe"
}

# 2. Setup demo environment
$DemoDir = "demo_repo"
if (Test-Path $DemoDir) {
    Remove-Item -Path $DemoDir -Recurse -Force
}
New-Item -ItemType Directory -Path $DemoDir | Out-Null
Set-Location $DemoDir

Print-Header "Setting up demo repository..."
Check-Command git "init"

# 3. Scenario 1: Initial Feature
Print-Header "Scenario 1: Initial Feature (Minor Bump)"
Print-Step "Commit: 'feat: initial project structure'"
New-Item -Path "main.go" -Value "package main" | Out-Null
Check-Command git "add ."
Check-Command git "commit", "-m", "feat: initial project structure"

Print-Step "Running sem-version..."
..\sem-version.exe --verbose

# 4. Scenario 2: Bug Fix
Print-Header "Scenario 2: Bug Fix (Patch Bump)"
# Tag the previous version so we have a baseline (simulate release)
Check-Command git "tag", "v0.1.0"
Print-Step "Current Tag: v0.1.0"

Print-Step "Commit: 'fix: resolve startup crash'"
Check-Command git "commit", "--allow-empty", "-m", "fix: resolve startup crash"

Print-Step "Running sem-version..."
..\sem-version.exe --verbose

# 5. Scenario 3: New Feature
Print-Header "Scenario 3: New Feature (Minor Bump)"
# Simulate releasing v0.1.1
Check-Command git "tag", "v0.1.1"
Print-Step "Current Tag: v0.1.1"

Print-Step "Commit: 'feat: add user login'"
Check-Command git "commit", "--allow-empty", "-m", "feat: add user login"

Print-Step "Running sem-version..."
..\sem-version.exe --verbose

# 6. Scenario 4: Breaking Change
Print-Header "Scenario 4: Breaking Change (Major Bump)"
# Simulate releasing v0.2.0
Check-Command git "tag", "v0.2.0"
Print-Step "Current Tag: v0.2.0"

Print-Step "Commit: 'feat!: redesign api structure'"
Check-Command git "commit", "--allow-empty", "-m", "feat!: redesign api structure"

Print-Step "Running sem-version..."
..\sem-version.exe --verbose

# Cleanup
Set-Location ..
# Remove-Item -Path $DemoDir -Recurse -Force
Print-Header "Demo Complete!"
Write-Host "You can inspect the '$DemoDir' directory."
