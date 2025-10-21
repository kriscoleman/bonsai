package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kriscoleman-testifysec/bonsai/internal/config"
	"github.com/kriscoleman-testifysec/bonsai/internal/git"
	"github.com/kriscoleman-testifysec/bonsai/internal/ui"
	"github.com/spf13/cobra"
)

var (
	remoteBulk   bool
	remoteAge    string
	remoteDryRun bool
	remoteName   string
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Clean up remote branches",
	Long:  `Identify and clean up stale remote Git branches.`,
	RunE:  runRemoteCleanup,
}

func init() {
	rootCmd.AddCommand(remoteCmd)

	remoteCmd.Flags().BoolVar(&remoteBulk, "bulk", false, "Delete all stale branches without interaction")
	remoteCmd.Flags().StringVar(&remoteAge, "age", "4w", "Age threshold for stale branches (e.g., 4w, 28d, 672h)")
	remoteCmd.Flags().BoolVar(&remoteDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	remoteCmd.Flags().StringVar(&remoteName, "remote", "origin", "Remote name to clean up")
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
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(successStyle.Render("✓ No stale remote branches found!"))
		return nil
	}

	// Show summary
	printBranchSummary(staleBranches, "remote", ageThreshold, remoteDryRun)

	if remoteDryRun {
		return nil
	}

	if remoteBulk {
		return runBulkDeletion(repo, staleBranches, true)
	}

	return ui.RunInteractiveSelection(repo, staleBranches, true)
}
