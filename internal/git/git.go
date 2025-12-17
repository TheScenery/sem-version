package git

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Commit represents a git commit
type Commit struct {
	Hash    string
	Message string
}

// GetLatestTag returns the latest semver tag in the repository
func GetLatestTag(repoPath string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match", "v*")
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// No tags found is not an error for our use case
		if strings.Contains(stderr.String(), "No names found") ||
			strings.Contains(stderr.String(), "No tags can describe") ||
			strings.Contains(stderr.String(), "fatal") {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

// GetCommitsSince returns all commits since the given tag
// If tag is empty, returns all commits
func GetCommitsSince(repoPath, tag string) ([]Commit, error) {
	var args []string
	if tag == "" {
		args = []string{"log", "--pretty=format:%H|%s", "--reverse"}
	} else {
		args = []string{"log", "--pretty=format:%H|%s", "--reverse", tag + "..HEAD"}
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	lines := strings.Split(output, "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			commits = append(commits, Commit{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}

	return commits, nil
}

// GetFullCommitMessage returns the full commit message including body
func GetFullCommitMessage(repoPath, hash string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%B", hash)
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}
