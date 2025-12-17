# sem-version

[![Release](https://img.shields.io/github/v/release/TheScenery/sem-version?style=flat-square)](https://github.com/TheScenery/sem-version/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/TheScenery/sem-version?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/github/license/TheScenery/sem-version?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/TheScenery/sem-version?style=flat-square)](https://goreportcard.com/report/github.com/TheScenery/sem-version)

A CLI tool to automatically generate semantic versions based on Git commit history using [Conventional Commits](https://www.conventionalcommits.org/) specification.

## Installation

### Go Install

```bash
go install github.com/TheScenery/sem-version@latest
```

### One-liner Script

**macOS (Apple Silicon)**
```bash
curl -sSL https://github.com/TheScenery/sem-version/releases/latest/download/sem-version-darwin-arm64 -o /usr/local/bin/sem-version && chmod +x /usr/local/bin/sem-version
```

**macOS (Intel)**
```bash
curl -sSL https://github.com/TheScenery/sem-version/releases/latest/download/sem-version-darwin-amd64 -o /usr/local/bin/sem-version && chmod +x /usr/local/bin/sem-version
```

**Linux (amd64)**
```bash
curl -sSL https://github.com/TheScenery/sem-version/releases/latest/download/sem-version-linux-amd64 -o /usr/local/bin/sem-version && chmod +x /usr/local/bin/sem-version
```

**Linux (arm64)**
```bash
curl -sSL https://github.com/TheScenery/sem-version/releases/latest/download/sem-version-linux-arm64 -o /usr/local/bin/sem-version && chmod +x /usr/local/bin/sem-version
```

**Windows (PowerShell)**
```powershell
Invoke-WebRequest -Uri https://github.com/TheScenery/sem-version/releases/latest/download/sem-version-windows-amd64.exe -OutFile $env:USERPROFILE\bin\sem-version.exe
```

### Build from Source

```bash
git clone https://github.com/TheScenery/sem-version.git
cd sem-version
go build -o sem-version
```

## Usage

```bash
# Generate next version for current repository
sem-version

# With verbose output
sem-version --verbose

# Output without 'v' prefix
sem-version --no-prefix

# Specify repository path
sem-version --path /path/to/repo

# Custom prefix
sem-version --prefix "ver"

# Generate default config file
sem-version --init

# Use custom config file
sem-version --config /path/to/.sem-version.yaml
```

## Configuration

Generate a default config file with `sem-version --init`, which creates `.sem-version.yaml`:

```yaml
# Major version bump (breaking changes)
major:
  - '^.+!:'                # type!: breaking change
  - 'BREAKING CHANGE:'     # in commit body

# Minor version bump (new features)
minor:
  - '^feat(\(.+\))?:'      # feat: or feat(scope):

# Patch version bump (bug fixes)
patch:
  - '^fix(\(.+\))?:'       # fix:
  - '^hotfix(\(.+\))?:'    # hotfix:
  - '^refactor(\(.+\))?:'  # refactor:
  - '^perf(\(.+\))?:'      # perf:
```

Each section contains **regex patterns** to match commit messages. Customize patterns to fit your workflow.

## Conventional Commits

This tool follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `feat` | Minor | `feat: add user login` |
| `fix` | Patch | `fix: resolve null pointer` |
| `refactor` | Patch | `refactor: clean up code` |
| `perf` | Patch | `perf: optimize query` |
| `feat!` | Major | `feat!: breaking API change` |
| `BREAKING CHANGE` | Major | (in commit body) |

### Examples

```bash
# Starting from no tags
git commit -m "feat: initial implementation"
sem-version
# Output: v0.1.0

# After tagging v0.1.0
git commit -m "fix: bug fix"
sem-version
# Output: v0.1.1

git commit -m "feat: new feature"
sem-version
# Output: v0.2.0

git commit -m "feat!: breaking change"
sem-version
# Output: v1.0.0
```

## How It Works

1. Finds the latest semantic version tag (e.g., `v1.2.3`)
2. Collects all commits since that tag
3. Parses each commit message using Conventional Commits format
4. Calculates the next version based on commit types:
   - **BREAKING CHANGE** → Major bump (reset minor and patch)
   - **feat** → Minor bump (reset patch)
   - **fix/refactor/perf** → Patch bump

## License

MIT
