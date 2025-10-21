package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kriscoleman-testifysec/bonsai/internal/config"
	"github.com/kriscoleman-testifysec/bonsai/internal/git"
	"github.com/kriscoleman-testifysec/bonsai/internal/ui"
	"github.com/spf13/cobra"
)

var (
	localBulk   bool
	localAge    string
	localDryRun bool
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Clean up local branches",
	Long:  `Identify and clean up stale local Git branches.`,
	RunE:  runLocalCleanup,
}

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().BoolVar(&localBulk, "bulk", false, "Delete all stale branches without interaction")
	localCmd.Flags().StringVar(&localAge, "age", "2w", "Age threshold for stale branches (e.g., 2w, 14d, 336h)")
	localCmd.Flags().BoolVar(&localDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
}

func runLocalCleanup(cmd *cobra.Command, args []string) error {
	// Parse age threshold
	ageThreshold, err := config.ParseDuration(localAge)
	if err != nil {
		return fmt.Errorf("invalid age format: %w", err)
	}

	// Initialize repository
	repo := git.NewRepository("")

	// Check if we're in a git repository
	if err := repo.IsGitRepository(); err != nil {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	// Get all local branches
	branches, err := repo.ListLocalBranches()
	if err != nil {
		return err
	}

	// Filter stale branches
	staleBranches := filterStaleBranches(branches, ageThreshold)

	if len(staleBranches) == 0 {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(successStyle.Render("✓ No stale local branches found!"))
		return nil
	}

	// Show summary
	printBranchSummary(staleBranches, "local", ageThreshold, localDryRun)

	if localDryRun {
		return nil
	}

	if localBulk {
		return runBulkDeletion(repo, staleBranches, false)
	}

	return ui.RunInteractiveSelection(repo, staleBranches, false)
}

func filterStaleBranches(branches []*git.Branch, threshold time.Duration) []*git.Branch {
	var stale []*git.Branch

	for _, branch := range branches {
		// Skip current branch
		if branch.IsCurrent {
			continue
		}

		// Skip protected branches
		if branch.IsProtected {
			continue
		}

		// Check if stale
		if branch.IsStale(threshold) {
			stale = append(stale, branch)
		}
	}

	return stale
}

func printBranchSummary(branches []*git.Branch, branchType string, threshold time.Duration, dryRun bool) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	title := fmt.Sprintf("Found %d stale %s branch(es)", len(branches), branchType)
	if dryRun {
		title += " [DRY RUN]"
	}

	fmt.Println(titleStyle.Render(title))
	fmt.Println(infoStyle.Render(fmt.Sprintf("Age threshold: %v", threshold)))
	fmt.Println()

	if !dryRun {
		fmt.Println(warningStyle.Render("⚠ These branches will be deleted"))
		fmt.Println()
	}
}

func runBulkDeletion(repo *git.Repository, branches []*git.Branch, isRemote bool) error {
	// Confirm bulk deletion
	if !confirmBulkDeletion(len(branches)) {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	successCount := 0
	errorCount := 0

	for _, branch := range branches {
		var err error
		if isRemote {
			err = repo.DeleteRemoteBranch(branch.RemoteName, branch.Name)
		} else {
			err = repo.DeleteLocalBranch(branch.Name, false)
		}

		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("✗ Failed to delete %s: %v", branch.FullName(), err)))
			errorCount++
		} else {
			fmt.Println(successStyle.Render(fmt.Sprintf("✓ Deleted %s", branch.FullName())))
			successCount++
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d deleted, %d failed\n", successCount, errorCount)

	return nil
}

func confirmBulkDeletion(count int) bool {
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)

	fmt.Println(warningStyle.Render(fmt.Sprintf("⚠ This will delete %d branch(es). Are you sure? (y/N)", count)))

	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == "yes"
}
