package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kriscoleman/bonsai/internal/config"
	"github.com/kriscoleman/bonsai/internal/git"
	"github.com/kriscoleman/bonsai/internal/ui"
	"github.com/spf13/cobra"
)

var (
	remoteBulk    bool
	remoteAge     string
	remoteDryRun  bool
	remoteName    string
	remoteVerbose bool
	remoteForce   bool
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "🌍 Prune stale remote branches",
	Long: `🌍 Prune stale remote branches

Identify and carefully remove remote branches that have stopped growing.
Maintain your shared repository with the same dedication and artistry
as tending to a cherished bonsai tree.`,
	RunE: runRemoteCleanup,
}

func init() {
	rootCmd.AddCommand(remoteCmd)

	remoteCmd.Flags().BoolVar(&remoteBulk, "bulk", false, "Delete all stale branches without interaction")
	remoteCmd.Flags().StringVar(&remoteAge, "age", "4w", "Age threshold for stale branches (e.g., 4w, 28d, 672h)")
	remoteCmd.Flags().BoolVar(&remoteDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	remoteCmd.Flags().StringVar(&remoteName, "remote", "origin", "Remote name to clean up")
	remoteCmd.Flags().BoolVarP(&remoteVerbose, "verbose", "v", false, "Show detailed error messages")
	remoteCmd.Flags().BoolVarP(&remoteForce, "force", "f", false, "Force delete branches even if not fully merged")
}

func runRemoteCleanup(cmd *cobra.Command, args []string) error {
	// Parse age threshold
	ageThreshold, err := config.ParseDuration(remoteAge)
	if err != nil {
		return fmt.Errorf("invalid age format: %w", err)
	}

	// Initialize repository
	repo := git.NewRepository("")

	// Check if we're in a git repository
	if err := repo.IsGitRepository(); err != nil {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	// Get all remote branches
	branches, err := repo.ListRemoteBranches(remoteName)
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
			fmt.Sprintf("🌳 Your %s remote is perfectly maintained!", remoteName),
			"   No stale branches found - a true work of art.")

		fmt.Println(successBox.Render(successStyle.Render(content)))
		return nil
	}

	// Show summary
	printBranchSummary(staleBranches, "remote", ageThreshold, remoteDryRun)

	if remoteDryRun {
		return nil
	}

	if remoteBulk {
		return runBulkDeletion(repo, staleBranches, true, remoteVerbose, remoteForce)
	}

	return ui.RunInteractiveSelection(repo, staleBranches, true, remoteVerbose, remoteForce)
}
