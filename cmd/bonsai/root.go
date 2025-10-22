package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Color palette inspired by bonsai aesthetics
	leafGreen    = lipgloss.Color("#7FB069") // Fresh leaf green
	accentPurple = lipgloss.Color("#C792EA") // Charm purple accent
	softCyan     = lipgloss.Color("#89DDFF") // Soft cyan highlight
	mutedGray    = lipgloss.Color("#8F8F8F") // Elegant gray

	// Bonsai ASCII art
	bonsaiArt = `
           ,.,
      MMMM_    ,..,
        "_ "__"MMMMM          ,...,,
 ,..., __." --"    ,.,     _-"MMMMMMM
MMMMMM"___ "_._   MMM"_."" _ """"""
 """""    "" , \_.   "_. ."
        ,., _"__ \__./ ."
       MMMMM_"  "_    ./
        ''''      (    )
 ._______________.-'____"---._.
  \                          /
   \________________________/
   (_)                    (_)
`
)

var rootCmd = &cobra.Command{
	Use:     "bonsai",
	Short:   "🌳 The Art of Branch Pruning",
	Long:    renderLongDescription(),
	Version: "0.1.0",
}

func init() {
	// Set custom templates for help and usage
	cobra.AddTemplateFunc("StyleHeading", styleHeading)
	rootCmd.SetUsageTemplate(getUsageTemplate())
}

func renderLongDescription() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(leafGreen).
		MarginBottom(1)

	subtitleStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(mutedGray)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	quoteStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(accentPurple).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(softCyan).
		Padding(0, 1).
		MarginTop(1)

	artStyle := lipgloss.NewStyle().
		Foreground(leafGreen).
		Align(lipgloss.Center)

	return fmt.Sprintf(`%s

%s

%s

%s

%s`,
		artStyle.Render(bonsaiArt),
		titleStyle.Render("🌳 Bonsai - The Art of Branch Pruning"),
		subtitleStyle.Render("Keep your Git repository as elegant and intentional as a carefully cultivated bonsai tree."),
		descStyle.Render("\nJust as a bonsai master carefully prunes their tree to maintain its beauty and health,\nBonsai helps you maintain a clean, healthy repository by removing stale branches."),
		quoteStyle.Render("\"The best time to prune was 20 years ago. The second best time is now.\""),
	)
}

func styleHeading(s string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(leafGreen).
		Render(s)
}

func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
