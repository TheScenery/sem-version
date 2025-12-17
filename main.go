package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"sem-version/internal/config"
	"sem-version/internal/git"
	"sem-version/internal/version"
)

func main() {
	// Parse command line flags
	prefix := flag.String("prefix", "v", "Version prefix (default: v)")
	repoPath := flag.String("path", ".", "Path to git repository (default: current directory)")
	configPath := flag.String("config", "", "Path to config file (default: auto-detect .sem-version.yaml)")
	noPrefix := flag.Bool("no-prefix", false, "Output version without prefix")
	verbose := flag.Bool("verbose", false, "Show verbose output")
	initConfig := flag.Bool("init", false, "Generate default config file")
	flag.Parse()

	// Resolve absolute path
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Handle --init command
	if *initConfig {
		if err := generateDefaultConfig(absPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated .sem-version.yaml")
		return
	}

	// Load configuration
	var cfg *config.Config
	if *configPath != "" {
		cfg, err = config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config %s: %v\n", *configPath, err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Using config: %s\n", *configPath)
		}
	} else {
		cfg, err = config.LoadDefault(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	}

	// Get the latest tag
	latestTag, err := git.GetLatestTag(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting latest tag: %v\n", err)
		os.Exit(1)
	}

	var currentVersion version.Version
	if latestTag == "" {
		// No tags found, use initial version v0.0.0
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

	// Analyze commits using config
	bumpType := version.BumpNone
	for _, commit := range commits {
		// Get full commit message for better matching
		fullMessage, err := git.GetFullCommitMessage(absPath, commit.Hash)
		if err != nil {
			fullMessage = commit.Message
		}

		// Check patterns in order of priority: major > minor > patch
		if cfg.MatchMajor(fullMessage) {
			bumpType = version.BumpMajorType
			if *verbose {
				fmt.Fprintf(os.Stderr, "  - [MAJOR] %s\n", commit.Message)
			}
			break // Major is highest priority
		}

		if cfg.MatchMinor(fullMessage) && bumpType < version.BumpMinorType {
			bumpType = version.BumpMinorType
			if *verbose {
				fmt.Fprintf(os.Stderr, "  - [MINOR] %s\n", commit.Message)
			}
		} else if cfg.MatchPatch(fullMessage) && bumpType < version.BumpPatchType {
			bumpType = version.BumpPatchType
			if *verbose {
				fmt.Fprintf(os.Stderr, "  - [PATCH] %s\n", commit.Message)
			}
		} else if *verbose {
			fmt.Fprintf(os.Stderr, "  - [SKIP] %s\n", commit.Message)
		}
	}

	// Calculate next version
	var nextVersion version.Version
	if latestTag == "" {
		// No existing tag - determine initial version based on bump type
		nextVersion = calculateInitialVersion(bumpType)
	} else {
		nextVersion = applyBump(currentVersion, bumpType)
	}

	outputVersion(nextVersion, *prefix, *noPrefix)
}

// calculateInitialVersion determines the initial version based on bump type
func calculateInitialVersion(bumpType version.BumpType) version.Version {
	switch bumpType {
	case version.BumpMajorType:
		return version.Version{Major: 1, Minor: 0, Patch: 0}
	case version.BumpMinorType:
		return version.Version{Major: 0, Minor: 1, Patch: 0}
	case version.BumpPatchType:
		return version.Version{Major: 0, Minor: 0, Patch: 1}
	default:
		return version.Version{Major: 0, Minor: 1, Patch: 0}
	}
}

// applyBump applies the bump type to the current version
func applyBump(current version.Version, bumpType version.BumpType) version.Version {
	switch bumpType {
	case version.BumpMajorType:
		return current.BumpMajor()
	case version.BumpMinorType:
		return current.BumpMinor()
	case version.BumpPatchType:
		return current.BumpPatch()
	default:
		return current
	}
}

func outputVersion(v version.Version, prefix string, noPrefix bool) {
	if noPrefix {
		fmt.Printf("%d.%d.%d\n", v.Major, v.Minor, v.Patch)
	} else {
		fmt.Printf("%s%d.%d.%d\n", prefix, v.Major, v.Minor, v.Patch)
	}
}

func generateDefaultConfig(dir string) error {
	configPath := filepath.Join(dir, ".sem-version.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	return os.WriteFile(configPath, []byte(config.DefaultConfigYAML()), 0644)
}
