package parser

import (
	"regexp"
	"strings"
)

// CommitType represents the type of a conventional commit
type CommitType string

const (
	TypeFeat     CommitType = "feat"
	TypeFix      CommitType = "fix"
	TypeDocs     CommitType = "docs"
	TypeStyle    CommitType = "style"
	TypeRefactor CommitType = "refactor"
	TypePerf     CommitType = "perf"
	TypeTest     CommitType = "test"
	TypeChore    CommitType = "chore"
	TypeBuild    CommitType = "build"
	TypeCI       CommitType = "ci"
	TypeUnknown  CommitType = "unknown"
)

// ParsedCommit represents a parsed conventional commit
type ParsedCommit struct {
	Type           CommitType
	Scope          string
	Description    string
	IsBreaking     bool
	BreakingChange string
	RawMessage     string
}

// conventionalCommitRegex matches conventional commit format
// Pattern: type(scope)!: description
var conventionalCommitRegex = regexp.MustCompile(`^(\w+)(?:\(([^)]*)\))?(!)?:\s*(.*)$`)

// ParseCommit parses a commit message according to Conventional Commits spec
func ParseCommit(message string) ParsedCommit {
	result := ParsedCommit{
		Type:       TypeUnknown,
		RawMessage: message,
	}

	// Split message into subject and body
	lines := strings.SplitN(message, "\n", 2)
	subject := strings.TrimSpace(lines[0])

	// Try to match conventional commit format
	matches := conventionalCommitRegex.FindStringSubmatch(subject)
	if matches != nil {
		result.Type = parseType(matches[1])
		result.Scope = matches[2]
		result.IsBreaking = matches[3] == "!"
		result.Description = matches[4]
	}

	// Check for BREAKING CHANGE in body
	if len(lines) > 1 {
		body := lines[1]
		if strings.Contains(body, "BREAKING CHANGE:") || strings.Contains(body, "BREAKING-CHANGE:") {
			result.IsBreaking = true
			// Extract breaking change description
			for _, line := range strings.Split(body, "\n") {
				if strings.HasPrefix(line, "BREAKING CHANGE:") {
					result.BreakingChange = strings.TrimPrefix(line, "BREAKING CHANGE:")
					result.BreakingChange = strings.TrimSpace(result.BreakingChange)
					break
				}
				if strings.HasPrefix(line, "BREAKING-CHANGE:") {
					result.BreakingChange = strings.TrimPrefix(line, "BREAKING-CHANGE:")
					result.BreakingChange = strings.TrimSpace(result.BreakingChange)
					break
				}
			}
		}
	}

	return result
}

// parseType converts a string to CommitType
func parseType(t string) CommitType {
	switch strings.ToLower(t) {
	case "feat", "feature":
		return TypeFeat
	case "fix", "bugfix":
		return TypeFix
	case "docs", "doc":
		return TypeDocs
	case "style":
		return TypeStyle
	case "refactor":
		return TypeRefactor
	case "perf", "performance":
		return TypePerf
	case "test", "tests":
		return TypeTest
	case "chore":
		return TypeChore
	case "build":
		return TypeBuild
	case "ci":
		return TypeCI
	default:
		return TypeUnknown
	}
}

// IsBumpType returns true if the commit type should trigger a version bump
func (p ParsedCommit) IsBumpType() bool {
	switch p.Type {
	case TypeFeat, TypeFix, TypeRefactor, TypePerf:
		return true
	default:
		return false
	}
}
