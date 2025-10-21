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
	Short: "🌿 Prune stale local branches",
	Long: `🌿 Prune stale local branches

Carefully identify and remove local branches that are no longer actively growing.
Just as a bonsai master shapes their tree with intention, maintain your repository
with precision and care.`,
	RunE: runLocalCleanup,
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
		// Bonsai-themed success message
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7FB069")).
			Bold(true)

		successBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#89DDFF")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

		content := lipgloss.JoinVertical(lipgloss.Left,
			"🌳 Your repository is perfectly maintained!",
			"   No stale branches found - a true work of art.")

		fmt.Println(successBox.Render(successStyle.Render(content)))
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
	// Bonsai-themed colors
	leafGreen := lipgloss.Color("#7FB069")
	softCyan := lipgloss.Color("#89DDFF")
	mutedGray := lipgloss.Color("#8F8F8F")
	warningYellow := lipgloss.Color("#FFD43B")

	warningStyle := lipgloss.NewStyle().
		Foreground(warningYellow).
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(mutedGray).
		Italic(true)

	// Build the header content
	icon := "🌿"
	if dryRun {
		icon = "👁️"
	}

	title := fmt.Sprintf("%s Discovered %d stale %s branch(es)", icon, len(branches), branchType)
	if dryRun {
		title += " [PREVIEW MODE]"
	}

	info := fmt.Sprintf("Pruning threshold: %v", threshold)

	// Style each line
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(leafGreen)

	infoStyled := lipgloss.NewStyle().
		Foreground(mutedGray).
		Italic(true)

	// Render and join
	headerContent := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		infoStyled.Render(info))

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(softCyan).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	fmt.Println(headerBox.Render(headerContent))

	if !dryRun {
		warning := "⚠️  These branches are ready for careful pruning"
		fmt.Println(warningStyle.Render(warning))
		fmt.Println()
	} else {
		preview := "Preview mode: no changes will be made to your repository"
		fmt.Println(infoStyle.Render(preview))
		fmt.Println()
	}
}

func runBulkDeletion(repo *git.Repository, branches []*git.Branch, isRemote bool) error {
	// Confirm bulk deletion
	if !confirmBulkDeletion(len(branches)) {
		cancelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8F8F8F")).
			Italic(true)
		fmt.Println(cancelStyle.Render("🍃 Pruning cancelled. Your repository remains untouched."))
		return nil
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#51CF66"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))

	// Progress header
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7FB069")).
		Bold(true)
	fmt.Println(progressStyle.Render("🌀 Pruning in progress...\n"))

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
			fmt.Println(errorStyle.Render(fmt.Sprintf("  ✗ Failed to prune %s: %v", branch.FullName(), err)))
			errorCount++
		} else {
			fmt.Println(successStyle.Render(fmt.Sprintf("  ✓ Pruned %s", branch.FullName())))
			successCount++
		}
	}

	// Summary box
	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7FB069")).
		Bold(true)

	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89DDFF")).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	content := lipgloss.JoinVertical(lipgloss.Left,
		"🌳 Pruning complete!",
		fmt.Sprintf("   %d branches removed, %d failed", successCount, errorCount))

	fmt.Println(summaryBox.Render(summaryStyle.Render(content)))

	return nil
}

func confirmBulkDeletion(count int) bool {
	// Beautiful confirmation prompt with bonsai metaphor
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD43B")).
		Bold(true)

	warningBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFD43B")).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89DDFF")).
		Italic(true)

	content := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("⚠️  Ready to prune %d branch(es)", count),
		"   This action cannot be undone.")

	fmt.Println(warningBox.Render(warningStyle.Render(content)))
	fmt.Print(promptStyle.Render("Proceed with pruning? (y/N) "))

	var response string
	_, _ = fmt.Scanln(&response) // Ignore error - empty input is valid (defaults to No)

	return response == "y" || response == "Y" || response == "yes"
}
