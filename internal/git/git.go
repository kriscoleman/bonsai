package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var (
	// DefaultProtectedBranches are branches that should never be deleted
	DefaultProtectedBranches = []string{"main", "master", "develop"}
)

// Repository represents a Git repository
type Repository struct {
	Path string
}

// NewRepository creates a new Repository instance
func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

// IsGitRepository checks if the current directory is a Git repository
func (r *Repository) IsGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository")
	}
	return nil
}

// GetCurrentBranch returns the name of the currently checked out branch
func (r *Repository) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// ListLocalBranches returns a list of all local branches with their metadata
func (r *Repository) ListLocalBranches() ([]*Branch, error) {
	// Use git for-each-ref for efficient branch listing with all metadata
	// Format: refname|committerdate:iso8601|subject|authorname
	format := "%(refname:short)|%(committerdate:iso8601)|%(subject)|%(authorname)"
	cmd := exec.Command("git", "for-each-ref", "--format="+format, "refs/heads/")
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list local branches: %w", err)
	}

	currentBranch, err := r.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	return parseBranches(output, false, currentBranch)
}

// ListRemoteBranches returns a list of all remote branches with their metadata
func (r *Repository) ListRemoteBranches(remote string) ([]*Branch, error) {
	// Format: refname|committerdate:iso8601|subject|authorname
	format := "%(refname:short)|%(committerdate:iso8601)|%(subject)|%(authorname)"
	refPattern := fmt.Sprintf("refs/remotes/%s/", remote)
	cmd := exec.Command("git", "for-each-ref", "--format="+format, refPattern)
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list remote branches: %w", err)
	}

	branches, err := parseBranches(output, true, "")
	if err != nil {
		return nil, err
	}

	// Set remote name for all branches
	for _, branch := range branches {
		branch.RemoteName = remote
		// Remove remote prefix from branch name
		branch.Name = strings.TrimPrefix(branch.Name, remote+"/")
	}

	return branches, nil
}

// parseBranches parses the output from git for-each-ref
func parseBranches(output []byte, isRemote bool, currentBranch string) ([]*Branch, error) {
	var branches []*Branch

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		name := parts[0]
		commitDate := parts[1]
		commitMsg := parts[2]
		author := parts[3]

		// Parse commit date
		lastCommitAt, err := time.Parse("2006-01-02 15:04:05 -0700", commitDate)
		if err != nil {
			// Try alternative format
			lastCommitAt, err = time.Parse(time.RFC3339, commitDate)
			if err != nil {
				continue
			}
		}

		branch := &Branch{
			Name:          name,
			LastCommitAt:  lastCommitAt,
			LastCommitMsg: commitMsg,
			LastAuthor:    author,
			IsRemote:      isRemote,
			IsCurrent:     name == currentBranch,
			IsProtected:   isProtectedBranch(name),
		}

		branches = append(branches, branch)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse branch list: %w", err)
	}

	return branches, nil
}

// isProtectedBranch checks if a branch name is in the protected list
func isProtectedBranch(name string) bool {
	// Remove remote prefix if present
	branchName := name
	if idx := strings.LastIndex(name, "/"); idx != -1 {
		branchName = name[idx+1:]
	}

	for _, protected := range DefaultProtectedBranches {
		if branchName == protected {
			return true
		}
	}
	return false
}

// DeleteLocalBranch deletes a local branch
func (r *Repository) DeleteLocalBranch(branchName string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}

	cmd := exec.Command("git", "branch", flag, branchName)
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete local branch %s: %w", branchName, err)
	}

	return nil
}

// DeleteRemoteBranch deletes a remote branch
func (r *Repository) DeleteRemoteBranch(remote, branchName string) error {
	cmd := exec.Command("git", "push", remote, "--delete", branchName)
	if r.Path != "" {
		cmd.Dir = r.Path
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete remote branch %s/%s: %w", remote, branchName, err)
	}

	return nil
}
