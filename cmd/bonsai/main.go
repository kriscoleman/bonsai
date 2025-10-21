package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Show beautiful startup banner only when showing help or version
	showBanner := len(os.Args) == 1 ||
		(len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help"))

	if showBanner {
		printBanner()
	}

	if err := rootCmd.Execute(); err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)
		fmt.Fprintf(os.Stderr, "\n%s\n", errorStyle.Render("✗ Error: "+err.Error()))
		os.Exit(1)
	}
}

func printBanner() {
	bannerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7FB069")).
		Bold(true).
		Align(lipgloss.Center).
		MarginTop(1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89DDFF")).
		Italic(true).
		Align(lipgloss.Center)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8F8F8F")).
		Align(lipgloss.Center).
		MarginBottom(1)

	banner := `
    ╭───────────────────────────────────────╮
    │                                       │
    │           🌳  B O N S A I  🌳         │
    │                                       │
    ╰───────────────────────────────────────╯
`

	subtitle := "The Art of Branch Pruning"
	version := "v0.1.0 • Built with Charm 💜"

	fmt.Println(bannerStyle.Render(banner))
	fmt.Println(subtitleStyle.Render(subtitle))
	fmt.Println(versionStyle.Render(version))
}
