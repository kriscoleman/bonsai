package git

import (
	"time"
)

// Branch represents a Git branch with metadata
type Branch struct {
	Name          string
	LastCommitAt  time.Time
	LastCommitMsg string
	LastAuthor    string
	IsRemote      bool
	RemoteName    string // e.g., "origin"
	IsCurrent     bool
	IsProtected   bool
}

// Age returns the duration since the last commit
func (b *Branch) Age() time.Duration {
	return time.Since(b.LastCommitAt)
}

// IsStale checks if the branch is older than the given threshold
func (b *Branch) IsStale(threshold time.Duration) bool {
	return b.Age() > threshold
}

// FullName returns the full branch name (with remote prefix if applicable)
func (b *Branch) FullName() string {
	if b.IsRemote {
		return b.RemoteName + "/" + b.Name
	}
	return b.Name
}
