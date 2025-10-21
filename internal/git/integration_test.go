// +build integration

package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestHelper provides utilities for integration tests
type TestHelper struct {
	t       *testing.T
	TempDir string
	RepoDir string
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "test-repo")

	return &TestHelper{
		t:       t,
		TempDir: tmpDir,
		RepoDir: repoDir,
	}
}

// InitRepo initializes a new Git repository
func (h *TestHelper) InitRepo() {
	h.runGitCommand("init", h.RepoDir)
	h.runGitCommand("-C", h.RepoDir, "config", "user.name", "Test User")
	h.runGitCommand("-C", h.RepoDir, "config", "user.email", "test@example.com")
	h.runGitCommand("-C", h.RepoDir, "config", "commit.gpgsign", "false")
}

// CreateInitialCommit creates the first commit on main branch
func (h *TestHelper) CreateInitialCommit() {
	readmePath := filepath.Join(h.RepoDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Repo\n"), 0644); err != nil {
		h.t.Fatalf("Failed to create README: %v", err)
	}

	h.runGitCommand("-C", h.RepoDir, "add", "README.md")
	h.runGitCommand("-C", h.RepoDir, "commit", "-m", "Initial commit")
}

// CreateBranch creates a new branch and optionally checks it out
func (h *TestHelper) CreateBranch(name string, checkout bool) {
	if checkout {
		h.runGitCommand("-C", h.RepoDir, "checkout", "-b", name)
	} else {
		h.runGitCommand("-C", h.RepoDir, "branch", name)
	}
}

// CreateBranchWithCommit creates a branch with a commit
func (h *TestHelper) CreateBranchWithCommit(name, message string) {
	currentBranch := h.GetCurrentBranch()
	h.CreateBranch(name, true)

	// Create a file and commit it
	filePath := filepath.Join(h.RepoDir, fmt.Sprintf("%s.txt", name))
	if err := os.WriteFile(filePath, []byte(fmt.Sprintf("Content for %s\n", name)), 0644); err != nil {
		h.t.Fatalf("Failed to create file: %v", err)
	}

	h.runGitCommand("-C", h.RepoDir, "add", ".")
	h.runGitCommand("-C", h.RepoDir, "commit", "-m", message)

	// Go back to original branch
	h.runGitCommand("-C", h.RepoDir, "checkout", currentBranch)
}

// CreateBranchWithAge creates a branch with a commit at a specific time in the past
func (h *TestHelper) CreateBranchWithAge(name string, daysAgo int) {
	currentBranch := h.GetCurrentBranch()
	h.CreateBranch(name, true)

	filePath := filepath.Join(h.RepoDir, fmt.Sprintf("%s.txt", name))
	if err := os.WriteFile(filePath, []byte(fmt.Sprintf("Content for %s\n", name)), 0644); err != nil {
		h.t.Fatalf("Failed to create file: %v", err)
	}

	h.runGitCommand("-C", h.RepoDir, "add", ".")

	// Set the commit date to the past
	commitDate := time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour)
	dateStr := commitDate.Format("Mon Jan 2 15:04:05 2006 -0700")

	cmd := exec.Command("git", "-C", h.RepoDir, "commit", "-m", fmt.Sprintf("Commit for %s", name), "--date", dateStr)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))

	if output, err := cmd.CombinedOutput(); err != nil {
		h.t.Fatalf("Failed to create commit with custom date: %v\nOutput: %s", err, output)
	}

	h.runGitCommand("-C", h.RepoDir, "checkout", currentBranch)
}

// CheckoutBranch checks out a branch
func (h *TestHelper) CheckoutBranch(name string) {
	h.runGitCommand("-C", h.RepoDir, "checkout", name)
}

// GetCurrentBranch returns the current branch name
func (h *TestHelper) GetCurrentBranch() string {
	cmd := exec.Command("git", "-C", h.RepoDir, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		h.t.Fatalf("Failed to get current branch: %v", err)
	}
	return strings.TrimSpace(string(output))
}

// BranchExists checks if a branch exists
func (h *TestHelper) BranchExists(name string) bool {
	cmd := exec.Command("git", "-C", h.RepoDir, "rev-parse", "--verify", name)
	err := cmd.Run()
	return err == nil
}

// ListBranches lists all local branches
func (h *TestHelper) ListBranches() []string {
	cmd := exec.Command("git", "-C", h.RepoDir, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		h.t.Fatalf("Failed to list branches: %v", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, b := range branches {
		if b != "" {
			result = append(result, b)
		}
	}
	return result
}

// runGitCommand runs a git command and fails the test on error
func (h *TestHelper) runGitCommand(args ...string) {
	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		h.t.Fatalf("Git command failed: git %v\nError: %v\nOutput: %s", args, err, output)
	}
}

// Integration Tests

func TestIntegration_Repository_IsGitRepository(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()

	repo := NewRepository(helper.RepoDir)
	err := repo.IsGitRepository()

	if err != nil {
		t.Errorf("IsGitRepository() failed for valid repo: %v", err)
	}

	// Test non-git directory
	nonGitDir := filepath.Join(helper.TempDir, "not-a-repo")
	if err := os.Mkdir(nonGitDir, 0755); err != nil {
		t.Fatalf("Failed to create non-git directory: %v", err)
	}

	repo2 := NewRepository(nonGitDir)
	err2 := repo2.IsGitRepository()

	if err2 == nil {
		t.Error("IsGitRepository() should fail for non-git directory")
	}
}

func TestIntegration_Repository_GetCurrentBranch(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	repo := NewRepository(helper.RepoDir)

	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch() error = %v", err)
	}

	// Should be on main or master (depending on git config)
	if currentBranch != "main" && currentBranch != "master" {
		t.Errorf("GetCurrentBranch() = %s, want 'main' or 'master'", currentBranch)
	}

	// Create and switch to a new branch
	helper.CreateBranch("feature-test", true)

	currentBranch2, err := repo.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch() error = %v", err)
	}

	if currentBranch2 != "feature-test" {
		t.Errorf("GetCurrentBranch() = %s, want 'feature-test'", currentBranch2)
	}
}

func TestIntegration_Repository_ListLocalBranches(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	defaultBranch := helper.GetCurrentBranch()

	// Create several branches
	helper.CreateBranchWithCommit("feature-1", "Add feature 1")
	helper.CreateBranchWithCommit("feature-2", "Add feature 2")
	helper.CreateBranchWithCommit("bugfix-123", "Fix bug 123")

	repo := NewRepository(helper.RepoDir)
	branches, err := repo.ListLocalBranches()

	if err != nil {
		t.Fatalf("ListLocalBranches() error = %v", err)
	}

	// Should have default branch + 3 feature branches = 4 total
	expectedCount := 4
	if len(branches) != expectedCount {
		t.Errorf("ListLocalBranches() returned %d branches, want %d", len(branches), expectedCount)
	}

	// Check that all expected branches are present
	branchNames := make(map[string]bool)
	for _, b := range branches {
		branchNames[b.Name] = true

		// Verify metadata is populated
		if b.LastCommitMsg == "" {
			t.Errorf("Branch %s has empty commit message", b.Name)
		}
		if b.LastAuthor == "" {
			t.Errorf("Branch %s has empty author", b.Name)
		}
		if b.LastCommitAt.IsZero() {
			t.Errorf("Branch %s has zero commit time", b.Name)
		}
	}

	expectedBranches := []string{defaultBranch, "feature-1", "feature-2", "bugfix-123"}
	for _, expected := range expectedBranches {
		if !branchNames[expected] {
			t.Errorf("Expected branch %s not found in results", expected)
		}
	}

	// Check current branch flag
	currentSet := false
	for _, b := range branches {
		if b.IsCurrent {
			currentSet = true
			if b.Name != defaultBranch {
				t.Errorf("IsCurrent flag set on %s, but should be on %s", b.Name, defaultBranch)
			}
		}
	}
	if !currentSet {
		t.Error("No branch marked as current")
	}
}

func TestIntegration_Repository_ListLocalBranches_WithAge(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	// Create branches with different ages
	helper.CreateBranchWithAge("old-branch", 30)    // 30 days old
	helper.CreateBranchWithAge("recent-branch", 5)  // 5 days old
	helper.CreateBranchWithAge("very-old", 100)     // 100 days old

	repo := NewRepository(helper.RepoDir)
	branches, err := repo.ListLocalBranches()

	if err != nil {
		t.Fatalf("ListLocalBranches() error = %v", err)
	}

	// Check ages
	for _, b := range branches {
		age := b.Age()
		t.Logf("Branch %s is %v old", b.Name, age)

		switch b.Name {
		case "old-branch":
			minAge := 29 * 24 * time.Hour
			maxAge := 31 * 24 * time.Hour
			if age < minAge || age > maxAge {
				t.Errorf("Branch %s age = %v, expected around 30 days", b.Name, age)
			}
		case "recent-branch":
			minAge := 4 * 24 * time.Hour
			maxAge := 6 * 24 * time.Hour
			if age < minAge || age > maxAge {
				t.Errorf("Branch %s age = %v, expected around 5 days", b.Name, age)
			}
		case "very-old":
			minAge := 99 * 24 * time.Hour
			if age < minAge {
				t.Errorf("Branch %s age = %v, expected > 99 days", b.Name, age)
			}
		}
	}
}

func TestIntegration_Repository_DeleteLocalBranch(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	// Create a branch to delete
	helper.CreateBranchWithCommit("to-delete", "Branch to be deleted")

	// Verify it exists
	if !helper.BranchExists("to-delete") {
		t.Fatal("Test branch 'to-delete' was not created")
	}

	// Delete it with force (since it has unmerged commits)
	repo := NewRepository(helper.RepoDir)
	err := repo.DeleteLocalBranch("to-delete", true)

	if err != nil {
		t.Fatalf("DeleteLocalBranch() error = %v", err)
	}

	// Verify it's gone
	if helper.BranchExists("to-delete") {
		t.Error("Branch 'to-delete' still exists after deletion")
	}
}

func TestIntegration_Repository_DeleteLocalBranch_ForceDelete(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	defaultBranch := helper.GetCurrentBranch()

	// Create a branch with unmerged changes
	helper.CreateBranch("unmerged", true)

	// Create multiple commits on the branch
	for i := 0; i < 3; i++ {
		filePath := filepath.Join(helper.RepoDir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(filePath, []byte(fmt.Sprintf("Content %d\n", i)), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		helper.runGitCommand("-C", helper.RepoDir, "add", ".")
		helper.runGitCommand("-C", helper.RepoDir, "commit", "-m", fmt.Sprintf("Commit %d", i))
	}

	helper.CheckoutBranch(defaultBranch)

	// Try to delete without force (should fail)
	repo := NewRepository(helper.RepoDir)
	err := repo.DeleteLocalBranch("unmerged", false)

	if err == nil {
		t.Error("DeleteLocalBranch() should fail for unmerged branch without force flag")
	}

	// Delete with force
	err = repo.DeleteLocalBranch("unmerged", true)

	if err != nil {
		t.Fatalf("DeleteLocalBranch() with force flag error = %v", err)
	}

	// Verify it's gone
	if helper.BranchExists("unmerged") {
		t.Error("Branch 'unmerged' still exists after force deletion")
	}
}

func TestIntegration_BranchFiltering_Workflow(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	// Create branches with different characteristics
	helper.CreateBranchWithAge("stale-feature", 20)    // Older than 2 weeks
	helper.CreateBranchWithAge("recent-feature", 5)    // Newer than 2 weeks
	helper.CreateBranchWithAge("very-stale", 40)       // Very old
	helper.CreateBranchWithCommit("develop", "Develop branch")  // Protected
	helper.CreateBranchWithCommit("current-work", "Current work")

	// Checkout current-work to make it current
	helper.CheckoutBranch("current-work")

	repo := NewRepository(helper.RepoDir)
	branches, err := repo.ListLocalBranches()

	if err != nil {
		t.Fatalf("ListLocalBranches() error = %v", err)
	}

	// Filter stale branches (older than 14 days)
	threshold := 14 * 24 * time.Hour
	var staleBranches []*Branch

	for _, b := range branches {
		// Skip current branch
		if b.IsCurrent {
			continue
		}

		// Skip protected branches
		if b.IsProtected {
			continue
		}

		// Check if stale
		if b.IsStale(threshold) {
			staleBranches = append(staleBranches, b)
		}
	}

	// Should have 2 stale branches: stale-feature and very-stale
	if len(staleBranches) != 2 {
		t.Errorf("Found %d stale branches, want 2", len(staleBranches))
		for _, b := range staleBranches {
			t.Logf("Stale branch: %s (age: %v)", b.Name, b.Age())
		}
	}

	// Verify the correct branches were filtered
	staleNames := make(map[string]bool)
	for _, b := range staleBranches {
		staleNames[b.Name] = true
	}

	if !staleNames["stale-feature"] {
		t.Error("stale-feature should be in stale branches")
	}
	if !staleNames["very-stale"] {
		t.Error("very-stale should be in stale branches")
	}
	if staleNames["recent-feature"] {
		t.Error("recent-feature should NOT be in stale branches")
	}
	if staleNames["develop"] {
		t.Error("develop (protected) should NOT be in stale branches")
	}
	if staleNames["current-work"] {
		t.Error("current-work (current branch) should NOT be in stale branches")
	}
}

func TestIntegration_ProtectedBranches(t *testing.T) {
	helper := NewTestHelper(t)
	helper.InitRepo()
	helper.CreateInitialCommit()

	defaultBranch := helper.GetCurrentBranch()

	// Create protected branches (only create if they don't exist)
	if defaultBranch != "master" {
		helper.CreateBranchWithCommit("master", "Master branch")
	}
	if defaultBranch != "main" {
		helper.CreateBranchWithCommit("main", "Main branch")
	}
	helper.CreateBranchWithCommit("develop", "Develop branch")

	repo := NewRepository(helper.RepoDir)
	branches, err := repo.ListLocalBranches()

	if err != nil {
		t.Fatalf("ListLocalBranches() error = %v", err)
	}

	// Check that protected branches are marked correctly
	for _, b := range branches {
		switch b.Name {
		case "main", "master", "develop":
			if !b.IsProtected {
				t.Errorf("Branch %s should be marked as protected", b.Name)
			}
		default:
			if b.IsProtected {
				t.Errorf("Branch %s should NOT be marked as protected", b.Name)
			}
		}
	}
}
