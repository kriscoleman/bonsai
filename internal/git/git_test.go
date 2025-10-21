package git

import (
	"testing"
	"time"
)

func TestIsProtectedBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		want       bool
	}{
		{
			name:       "main is protected",
			branchName: "main",
			want:       true,
		},
		{
			name:       "master is protected",
			branchName: "master",
			want:       true,
		},
		{
			name:       "develop is protected",
			branchName: "develop",
			want:       true,
		},
		{
			name:       "feature branch is not protected",
			branchName: "feature/new-feature",
			want:       false,
		},
		{
			name:       "bugfix branch is not protected",
			branchName: "bugfix/issue-123",
			want:       false,
		},
		{
			name:       "remote main is protected",
			branchName: "origin/main",
			want:       true,
		},
		{
			name:       "remote master is protected",
			branchName: "origin/master",
			want:       true,
		},
		{
			name:       "remote develop is protected",
			branchName: "upstream/develop",
			want:       true,
		},
		{
			name:       "remote feature is not protected",
			branchName: "origin/feature/test",
			want:       false,
		},
		{
			name:       "empty string is not protected",
			branchName: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isProtectedBranch(tt.branchName); got != tt.want {
				t.Errorf("isProtectedBranch(%s) = %v, want %v", tt.branchName, got, tt.want)
			}
		})
	}
}

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name          string
		output        []byte
		isRemote      bool
		currentBranch string
		wantCount     int
		wantErr       bool
	}{
		{
			name: "single local branch",
			output: []byte(`feature/test|2024-01-15 10:30:00 -0800|Add new feature|John Doe
`),
			isRemote:      false,
			currentBranch: "main",
			wantCount:     1,
			wantErr:       false,
		},
		{
			name: "multiple local branches",
			output: []byte(`feature/test|2024-01-15 10:30:00 -0800|Add new feature|John Doe
bugfix/issue-123|2024-01-14 09:15:00 -0800|Fix critical bug|Jane Smith
`),
			isRemote:      false,
			currentBranch: "main",
			wantCount:     2,
			wantErr:       false,
		},
		{
			name: "current branch is identified",
			output: []byte(`main|2024-01-15 10:30:00 -0800|Update README|John Doe
feature/test|2024-01-14 09:15:00 -0800|Add feature|Jane Smith
`),
			isRemote:      false,
			currentBranch: "main",
			wantCount:     2,
			wantErr:       false,
		},
		{
			name:          "empty output",
			output:        []byte(``),
			isRemote:      false,
			currentBranch: "main",
			wantCount:     0,
			wantErr:       false,
		},
		{
			name: "malformed line - should skip",
			output: []byte(`feature/test|2024-01-15 10:30:00 -0800|Add new feature|John Doe
malformed-line
bugfix/issue-123|2024-01-14 09:15:00 -0800|Fix bug|Jane Smith
`),
			isRemote:      false,
			currentBranch: "main",
			wantCount:     2, // Should skip malformed line
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branches, err := parseBranches(tt.output, tt.isRemote, tt.currentBranch)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(branches) != tt.wantCount {
				t.Errorf("parseBranches() returned %d branches, want %d", len(branches), tt.wantCount)
			}

			// Check current branch flag
			if tt.currentBranch != "" {
				for _, b := range branches {
					if b.Name == tt.currentBranch && !b.IsCurrent {
						t.Errorf("Branch %s should be marked as current", tt.currentBranch)
					}
					if b.Name != tt.currentBranch && b.IsCurrent {
						t.Errorf("Branch %s should not be marked as current", b.Name)
					}
				}
			}

			// Check protected branch flag
			for _, b := range branches {
				expectedProtected := isProtectedBranch(b.Name)
				if b.IsProtected != expectedProtected {
					t.Errorf("Branch %s: IsProtected = %v, want %v", b.Name, b.IsProtected, expectedProtected)
				}
			}
		})
	}
}

func TestRepository_IsGitRepository(t *testing.T) {
	// This test requires an actual git repository
	// We'll test the basic structure
	repo := NewRepository("")

	// Test that the function exists and can be called
	_ = repo.IsGitRepository()
}

func TestBranchFiltering(t *testing.T) {
	// Create test branches
	now := time.Now()
	branches := []*Branch{
		{
			Name:         "feature/old",
			LastCommitAt: now.Add(-30 * 24 * time.Hour),
			IsCurrent:    false,
			IsProtected:  false,
		},
		{
			Name:         "feature/recent",
			LastCommitAt: now.Add(-5 * 24 * time.Hour),
			IsCurrent:    false,
			IsProtected:  false,
		},
		{
			Name:         "main",
			LastCommitAt: now.Add(-60 * 24 * time.Hour),
			IsCurrent:    false,
			IsProtected:  true,
		},
		{
			Name:         "feature/current",
			LastCommitAt: now.Add(-20 * 24 * time.Hour),
			IsCurrent:    true,
			IsProtected:  false,
		},
	}

	threshold := 14 * 24 * time.Hour

	// Filter logic that should be used in the CLI
	var stale []*Branch
	for _, b := range branches {
		if b.IsCurrent || b.IsProtected {
			continue
		}
		if b.IsStale(threshold) {
			stale = append(stale, b)
		}
	}

	// Should only have "feature/old"
	if len(stale) != 1 {
		t.Errorf("Expected 1 stale branch, got %d", len(stale))
	}

	if len(stale) > 0 && stale[0].Name != "feature/old" {
		t.Errorf("Expected stale branch to be 'feature/old', got '%s'", stale[0].Name)
	}

	// Verify filtering logic
	for _, b := range branches {
		shouldBeFiltered := !b.IsCurrent && !b.IsProtected && b.IsStale(threshold)
		isInStale := false
		for _, sb := range stale {
			if sb.Name == b.Name {
				isInStale = true
				break
			}
		}

		if shouldBeFiltered != isInStale {
			t.Errorf("Branch %s: shouldBeFiltered=%v, isInStale=%v", b.Name, shouldBeFiltered, isInStale)
		}
	}
}

func TestDefaultProtectedBranches(t *testing.T) {
	expected := []string{"main", "master", "develop"}

	if len(DefaultProtectedBranches) != len(expected) {
		t.Errorf("DefaultProtectedBranches length = %d, want %d", len(DefaultProtectedBranches), len(expected))
	}

	for i, branch := range expected {
		if DefaultProtectedBranches[i] != branch {
			t.Errorf("DefaultProtectedBranches[%d] = %s, want %s", i, DefaultProtectedBranches[i], branch)
		}
	}
}

func TestBranch_Metadata(t *testing.T) {
	// Test that Branch struct holds all required metadata
	branch := &Branch{
		Name:          "feature/test",
		LastCommitAt:  time.Now().Add(-10 * 24 * time.Hour),
		LastCommitMsg: "Add new feature",
		LastAuthor:    "John Doe",
		IsRemote:      false,
		RemoteName:    "",
		IsCurrent:     false,
		IsProtected:   false,
	}

	if branch.Name != "feature/test" {
		t.Errorf("Name = %s, want feature/test", branch.Name)
	}

	if branch.LastCommitMsg != "Add new feature" {
		t.Errorf("LastCommitMsg = %s, want 'Add new feature'", branch.LastCommitMsg)
	}

	if branch.LastAuthor != "John Doe" {
		t.Errorf("LastAuthor = %s, want 'John Doe'", branch.LastAuthor)
	}

	age := branch.Age()
	if age < 9*24*time.Hour || age > 11*24*time.Hour {
		t.Errorf("Age = %v, expected around 10 days", age)
	}
}
