package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"sem-version/internal/parser"
)

// Version represents a semantic version
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Metadata   string
}

// semverRegex matches semantic version format
var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9.-]+))?(?:\+([a-zA-Z0-9.-]+))?$`)

// Parse parses a version string into a Version struct
func Parse(v string) (Version, error) {
	matches := semverRegex.FindStringSubmatch(v)
	if matches == nil {
		return Version{}, fmt.Errorf("invalid version format: %s", v)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: matches[4],
		Metadata:   matches[5],
	}, nil
}

// String returns the version as a string with v prefix
func (v Version) String() string {
	base := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		base += "-" + v.Prerelease
	}
	if v.Metadata != "" {
		base += "+" + v.Metadata
	}
	return base
}

// StringWithoutPrefix returns the version as a string without v prefix
func (v Version) StringWithoutPrefix() string {
	return strings.TrimPrefix(v.String(), "v")
}

// BumpMajor increments the major version and resets minor and patch
func (v Version) BumpMajor() Version {
	return Version{
		Major: v.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

// BumpMinor increments the minor version and resets patch
func (v Version) BumpMinor() Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor + 1,
		Patch: 0,
	}
}

// BumpPatch increments the patch version
func (v Version) BumpPatch() Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch + 1,
	}
}

// BumpType represents the type of version bump
type BumpType int

const (
	BumpNone BumpType = iota
	BumpPatchType
	BumpMinorType
	BumpMajorType
)

// CalculateNextVersion determines the next version based on parsed commits
func CalculateNextVersion(current Version, commits []parser.ParsedCommit) Version {
	bumpType := BumpNone

	for _, commit := range commits {
		// Breaking change always bumps major
		if commit.IsBreaking {
			bumpType = BumpMajorType
			break // No need to check further
		}

		// feat bumps minor
		if commit.Type == parser.TypeFeat && bumpType < BumpMinorType {
			bumpType = BumpMinorType
		}

		// fix, refactor, perf bump patch
		if (commit.Type == parser.TypeFix ||
			commit.Type == parser.TypeRefactor ||
			commit.Type == parser.TypePerf) && bumpType < BumpPatchType {
			bumpType = BumpPatchType
		}
	}

	switch bumpType {
	case BumpMajorType:
		return current.BumpMajor()
	case BumpMinorType:
		return current.BumpMinor()
	case BumpPatchType:
		return current.BumpPatch()
	default:
		return current
	}
}

// DefaultInitialVersion returns the default initial version (v0.1.0)
func DefaultInitialVersion() Version {
	return Version{
		Major: 0,
		Minor: 1,
		Patch: 0,
	}
}
