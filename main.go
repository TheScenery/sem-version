package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"sem-version/internal/git"
	"sem-version/internal/parser"
	"sem-version/internal/version"
)

func main() {
	// Parse command line flags
	prefix := flag.String("prefix", "v", "Version prefix (default: v)")
	repoPath := flag.String("path", ".", "Path to git repository (default: current directory)")
	noPrefix := flag.Bool("no-prefix", false, "Output version without prefix")
	verbose := flag.Bool("verbose", false, "Show verbose output")
	flag.Parse()

	// Resolve absolute path
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Get the latest tag
	latestTag, err := git.GetLatestTag(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting latest tag: %v\n", err)
		os.Exit(1)
	}

	var currentVersion version.Version
	if latestTag == "" {
		// No tags found, use initial version v0.0.0 (will bump to v0.1.0 on first feat)
		currentVersion = version.Version{Major: 0, Minor: 0, Patch: 0}
		if *verbose {
			fmt.Fprintln(os.Stderr, "No existing tags found, starting from v0.0.0")
		}
	} else {
		currentVersion, err = version.Parse(latestTag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing version %s: %v\n", latestTag, err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Current version: %s\n", latestTag)
		}
	}

	// Get commits since last tag
	commits, err := git.GetCommitsSince(absPath, latestTag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting commits: %v\n", err)
		os.Exit(1)
	}

	if len(commits) == 0 {
		if *verbose {
			fmt.Fprintln(os.Stderr, "No new commits since last tag")
		}
		// If no commits and no tag, output initial version
		if latestTag == "" {
			outputVersion(version.DefaultInitialVersion(), *prefix, *noPrefix)
		} else {
			outputVersion(currentVersion, *prefix, *noPrefix)
		}
		return
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Found %d commits since last tag\n", len(commits))
	}

	// Parse all commits
	parsedCommits := make([]parser.ParsedCommit, 0, len(commits))
	for _, commit := range commits {
		// Get full commit message for BREAKING CHANGE detection
		fullMessage, err := git.GetFullCommitMessage(absPath, commit.Hash)
		if err != nil {
			fullMessage = commit.Message
		}
		parsed := parser.ParseCommit(fullMessage)
		parsedCommits = append(parsedCommits, parsed)

		if *verbose {
			breakingStr := ""
			if parsed.IsBreaking {
				breakingStr = " [BREAKING]"
			}
			fmt.Fprintf(os.Stderr, "  - %s: %s%s\n", parsed.Type, parsed.Description, breakingStr)
		}
	}

	// Calculate next version
	var nextVersion version.Version
	if latestTag == "" {
		// No existing tag - determine initial version based on commits
		nextVersion = calculateInitialVersion(parsedCommits)
	} else {
		nextVersion = version.CalculateNextVersion(currentVersion, parsedCommits)
	}

	outputVersion(nextVersion, *prefix, *noPrefix)
}

// calculateInitialVersion determines the initial version based on commits
// If there's a breaking change: v1.0.0
// If there's a feat: v0.1.0
// Otherwise: v0.0.1
func calculateInitialVersion(commits []parser.ParsedCommit) version.Version {
	hasBreaking := false
	hasFeat := false
	hasFix := false

	for _, commit := range commits {
		if commit.IsBreaking {
			hasBreaking = true
			break
		}
		if commit.Type == parser.TypeFeat {
			hasFeat = true
		}
		if commit.Type == parser.TypeFix || commit.Type == parser.TypeRefactor || commit.Type == parser.TypePerf {
			hasFix = true
		}
	}

	if hasBreaking {
		return version.Version{Major: 1, Minor: 0, Patch: 0}
	}
	if hasFeat {
		return version.Version{Major: 0, Minor: 1, Patch: 0}
	}
	if hasFix {
		return version.Version{Major: 0, Minor: 0, Patch: 1}
	}
	// Default to v0.1.0 for initial version
	return version.Version{Major: 0, Minor: 1, Patch: 0}
}

func outputVersion(v version.Version, prefix string, noPrefix bool) {
	if noPrefix {
		fmt.Printf("%d.%d.%d\n", v.Major, v.Minor, v.Patch)
	} else {
		fmt.Printf("%s%d.%d.%d\n", prefix, v.Major, v.Minor, v.Patch)
	}
}
