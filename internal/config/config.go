package config

import (
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for commit parsing
type Config struct {
	// Major version bump patterns (e.g., breaking changes)
	Major []string `yaml:"major"`
	// Minor version bump patterns (e.g., new features)
	Minor []string `yaml:"minor"`
	// Patch version bump patterns (e.g., bug fixes)
	Patch []string `yaml:"patch"`

	// Compiled regexes (internal use)
	majorRegexes []*regexp.Regexp
	minorRegexes []*regexp.Regexp
	patchRegexes []*regexp.Regexp
}

// DefaultConfig returns the default configuration based on Conventional Commits
func DefaultConfig() *Config {
	return &Config{
		Major: []string{
			`^.+!:`,            // type!: breaking change
			`BREAKING CHANGE:`, // in commit body
			`BREAKING-CHANGE:`, // alternative format
		},
		Minor: []string{
			`^feat(\(.+\))?:`,    // feat: or feat(scope):
			`^feature(\(.+\))?:`, // feature: alias
		},
		Patch: []string{
			`^fix(\(.+\))?:`,      // fix:
			`^bugfix(\(.+\))?:`,   // bugfix: alias
			`^hotfix(\(.+\))?:`,   // hotfix:
			`^refactor(\(.+\))?:`, // refactor:
			`^perf(\(.+\))?:`,     // perf:
		},
	}
}

// DefaultConfigYAML returns the default config as YAML string
func DefaultConfigYAML() string {
	return `# sem-version configuration
# Each section contains regex patterns to match commit messages

# Major version bump (breaking changes)
major:
  - '^.+!:'                # type!: breaking change
  - 'BREAKING CHANGE:'     # in commit body
  - 'BREAKING-CHANGE:'     # alternative format

# Minor version bump (new features)
minor:
  - '^feat(\(.+\))?:'      # feat: or feat(scope):
  - '^feature(\(.+\))?:'   # feature: alias

# Patch version bump (bug fixes, refactoring)
patch:
  - '^fix(\(.+\))?:'       # fix:
  - '^bugfix(\(.+\))?:'    # bugfix: alias
  - '^hotfix(\(.+\))?:'    # hotfix:
  - '^refactor(\(.+\))?:'  # refactor:
  - '^perf(\(.+\))?:'      # perf:
`
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.compile(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadDefault looks for .sem-version.yaml or .sem-version.yml in the given directory
// Returns default config if no config file is found
func LoadDefault(dir string) (*Config, error) {
	configNames := []string{".sem-version.yaml", ".sem-version.yml"}

	for _, name := range configNames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return Load(path)
		}
	}

	// No config file found, use default
	cfg := DefaultConfig()
	if err := cfg.compile(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// compile compiles all regex patterns
func (c *Config) compile() error {
	var err error

	c.majorRegexes, err = compilePatterns(c.Major)
	if err != nil {
		return err
	}

	c.minorRegexes, err = compilePatterns(c.Minor)
	if err != nil {
		return err
	}

	c.patchRegexes, err = compilePatterns(c.Patch)
	if err != nil {
		return err
	}

	return nil
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	regexes := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		regexes = append(regexes, re)
	}
	return regexes, nil
}

// MatchMajor returns true if the message matches any major bump pattern
func (c *Config) MatchMajor(message string) bool {
	for _, re := range c.majorRegexes {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}

// MatchMinor returns true if the message matches any minor bump pattern
func (c *Config) MatchMinor(message string) bool {
	for _, re := range c.minorRegexes {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}

// MatchPatch returns true if the message matches any patch bump pattern
func (c *Config) MatchPatch(message string) bool {
	for _, re := range c.patchRegexes {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}
