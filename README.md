# sem-version

A CLI tool to automatically generate semantic versions based on Git commit history using [Conventional Commits](https://www.conventionalcommits.org/) specification.

## Installation

```bash
go install github.com/TheScenery/sem-version@latest
```

Or build from source:

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
```

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
