package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.compile(); err != nil {
		t.Fatalf("Failed to compile default config: %v", err)
	}

	tests := []struct {
		message   string
		wantMajor bool
		wantMinor bool
		wantPatch bool
	}{
		// Major (breaking changes)
		{"feat!: breaking change", true, false, false},
		{"fix!: breaking fix", true, false, false},
		{"BREAKING CHANGE: something", true, false, false},

		// Minor (features)
		{"feat: new feature", false, true, false},
		{"feat(api): scoped feature", false, true, false},
		{"feature: alias", false, true, false},

		// Patch (fixes)
		{"fix: bug fix", false, false, true},
		{"fix(core): scoped fix", false, false, true},
		{"hotfix: urgent fix", false, false, true},
		{"refactor: code cleanup", false, false, true},
		{"perf: optimization", false, false, true},

		// No match
		{"docs: update readme", false, false, false},
		{"chore: update deps", false, false, false},
		{"random commit message", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			if got := cfg.MatchMajor(tt.message); got != tt.wantMajor {
				t.Errorf("MatchMajor() = %v, want %v", got, tt.wantMajor)
			}
			if got := cfg.MatchMinor(tt.message); got != tt.wantMinor {
				t.Errorf("MatchMinor() = %v, want %v", got, tt.wantMinor)
			}
			if got := cfg.MatchPatch(tt.message); got != tt.wantPatch {
				t.Errorf("MatchPatch() = %v, want %v", got, tt.wantPatch)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temp config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".sem-version.yaml")

	configContent := `
major:
  - '^break:'
minor:
  - '^add:'
patch:
  - '^fix:'
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test custom patterns
	if !cfg.MatchMajor("break: something") {
		t.Error("Expected major match for 'break: something'")
	}
	if !cfg.MatchMinor("add: something") {
		t.Error("Expected minor match for 'add: something'")
	}
	if !cfg.MatchPatch("fix: something") {
		t.Error("Expected patch match for 'fix: something'")
	}

	// Standard patterns should NOT match with custom config
	if cfg.MatchMinor("feat: something") {
		t.Error("Did not expect minor match for 'feat: something' with custom config")
	}
}

func TestLoadDefault(t *testing.T) {
	// Test with no config file (should use defaults)
	dir := t.TempDir()
	cfg, err := LoadDefault(dir)
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	// Default config should match conventional commits
	if !cfg.MatchMinor("feat: something") {
		t.Error("Expected default config to match 'feat: something'")
	}
}
