package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bonsai",
	Short: "Git branch cleanup tool",
	Long: `Bonsai is a CLI tool for managing and cleaning up stale Git branches.
It helps you identify and remove both local and remote branches that haven't
been updated recently, keeping your repository neat and tidy.`,
	Version: "0.1.0",
}

func init() {
	// Global flags can be added here if needed
}
