package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kriscoleman-testifysec/bonsai/internal/git"
)

var (
	// Bonsai color palette
	leafGreen    = lipgloss.Color("#7FB069") // Fresh leaf green
	accentPurple = lipgloss.Color("#C792EA") // Charm purple accent
	softCyan     = lipgloss.Color("#89DDFF") // Soft cyan highlight
	mutedGray    = lipgloss.Color("#8F8F8F") // Elegant gray
	warningRed   = lipgloss.Color("#FF6B6B") // Gentle warning red
	successGreen = lipgloss.Color("#51CF66") // Success green

	// Bonsai tree ASCII art for header
	bonsaiHeader = `
   ╭─────────────────────────────────────────────╮
   │         🌳  The Art of Pruning  🌳          │
   │                                             │
   │        "Cultivate with intention"           │
   ╰─────────────────────────────────────────────╯
`

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(leafGreen).
			MarginLeft(2).
			MarginTop(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(softCyan).
			Align(lipgloss.Center).
			Bold(true)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(lipgloss.Color("252"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(leafGreen).
				Bold(true)

	branchNameStyle = lipgloss.NewStyle().
			Foreground(softCyan).
			Bold(true)

	ageStyle = lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true)

	authorStyle = lipgloss.NewStyle().
			Foreground(accentPurple)

	commitMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4).
			Foreground(mutedGray)

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(mutedGray)

	deletingStyle = lipgloss.NewStyle().
			Margin(1, 0, 2, 4).
			Foreground(warningRed).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successGreen).
			Bold(true)
)

type branchItem struct {
	branch   *git.Branch
	selected bool
}

func (i branchItem) FilterValue() string {
	return i.branch.Name
}

func (i branchItem) Title() string {
	checkbox := "○"
	checkboxColor := mutedGray
	if i.selected {
		checkbox = "●"
		checkboxColor = leafGreen
	}

	checkboxStyle := lipgloss.NewStyle().Foreground(checkboxColor).Bold(true)
	age := formatAge(i.branch.Age())

	return fmt.Sprintf("%s %s %s",
		checkboxStyle.Render(checkbox),
		branchNameStyle.Render(i.branch.FullName()),
		ageStyle.Render("("+age+")"),
	)
}

func (i branchItem) Description() string {
	commitMsg := i.branch.LastCommitMsg
	if len(commitMsg) > 80 {
		commitMsg = commitMsg[:77] + "..."
	}

	// Use emoji for visual flair
	authorPrefix := "👤"
	commitPrefix := "💬"

	return fmt.Sprintf("  %s %s  %s %s",
		authorPrefix,
		authorStyle.Render(i.branch.LastAuthor),
		commitPrefix,
		commitMsgStyle.Render(commitMsg),
	)
}

type model struct {
	list     list.Model
	repo     *git.Repository
	branches []*git.Branch
	items    []branchItem
	isRemote bool
	quitting bool
	deleting bool
	message  string
}

type deleteCompleteMsg struct {
	success int
	failed  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.deleting {
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys(" ", "x"))):
			// Toggle selection
			if _, ok := m.list.SelectedItem().(branchItem); ok {
				idx := m.list.Index()
				m.items[idx].selected = !m.items[idx].selected
				m.list.SetItem(idx, m.items[idx])
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("a"))):
			// Select all
			for i := range m.items {
				m.items[i].selected = true
				m.list.SetItem(i, m.items[i])
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			// Select none
			for i := range m.items {
				m.items[i].selected = false
				m.list.SetItem(i, m.items[i])
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", "d"))):
			// Delete selected branches
			selected := m.getSelectedBranches()
			if len(selected) == 0 {
				return m, nil
			}

			m.deleting = true
			return m, m.deleteBranches(selected)
		}

	case deleteCompleteMsg:
		m.message = fmt.Sprintf("Deleted %d branch(es), %d failed", msg.success, msg.failed)
		m.quitting = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		if m.message != "" {
			// Create an elegant exit message with proper box alignment
			content := successStyle.Render("✓ " + m.message)

			box := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(softCyan).
				Padding(0, 1).
				MarginTop(1).
				MarginBottom(1)

			return "\n" + box.Render(content) + "\n"
		}

		cancelMsg := lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true).
			Render("Pruning session cancelled. Repository unchanged.")

		return "\n  " + cancelMsg + "\n"
	}

	if m.deleting {
		spinner := "🌀"
		deleteMsg := fmt.Sprintf("%s Carefully pruning selected branches...", spinner)

		return deletingStyle.Render("\n" + deleteMsg + "\n")
	}

	// Show the header with bonsai art
	header := headerStyle.Render(bonsaiHeader)

	// Get selected count
	selectedCount := 0
	for _, item := range m.items {
		if item.selected {
			selectedCount++
		}
	}

	statusBar := ""
	if selectedCount > 0 {
		statusBar = lipgloss.NewStyle().
			Foreground(leafGreen).
			Bold(true).
			MarginLeft(2).
			Render(fmt.Sprintf("🌿 %d branch(es) selected for pruning", selectedCount))
	} else {
		statusBar = lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true).
			MarginLeft(2).
			Render("Select branches to prune with space/x • a = all • n = none • enter/d = delete")
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", header, m.list.View(), statusBar)
}

func (m model) getSelectedBranches() []*git.Branch {
	var selected []*git.Branch
	for _, item := range m.items {
		if item.selected {
			selected = append(selected, item.branch)
		}
	}
	return selected
}

func (m model) deleteBranches(branches []*git.Branch) tea.Cmd {
	return func() tea.Msg {
		successCount := 0
		errorCount := 0

		for _, branch := range branches {
			var err error
			if m.isRemote {
				err = m.repo.DeleteRemoteBranch(branch.RemoteName, branch.Name)
			} else {
				err = m.repo.DeleteLocalBranch(branch.Name, false)
			}

			if err != nil {
				errorCount++
			} else {
				successCount++
			}
		}

		return deleteCompleteMsg{
			success: successCount,
			failed:  errorCount,
		}
	}
}

func formatAge(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	if days == 0 {
		return "today"
	} else if days == 1 {
		return "1 day ago"
	} else if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	} else if days < 30 {
		weeks := days / 7
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if days < 365 {
		months := days / 30
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := days / 365
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// RunInteractiveSelection starts the interactive branch selection UI
func RunInteractiveSelection(repo *git.Repository, branches []*git.Branch, isRemote bool) error {
	items := make([]list.Item, len(branches))
	branchItems := make([]branchItem, len(branches))

	for i, branch := range branches {
		branchItems[i] = branchItem{
			branch:   branch,
			selected: false,
		}
		items[i] = branchItems[i]
	}

	const defaultWidth = 80
	const listHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)

	// Elegant title with bonsai metaphor
	branchType := "local"
	if isRemote {
		branchType = "remote"
	}
	l.Title = fmt.Sprintf("🌿 Branches ready for pruning (%s)", branchType)

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Add custom key bindings to help
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("space", "x"),
				key.WithHelp("space/x", "toggle"),
			),
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "select all"),
			),
			key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "select none"),
			),
			key.NewBinding(
				key.WithKeys("enter", "d"),
				key.WithHelp("enter/d", "delete"),
			),
		}
	}

	m := model{
		list:     l,
		repo:     repo,
		branches: branches,
		items:    branchItems,
		isRemote: isRemote,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running interactive selection: %w", err)
	}

	return nil
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(branchItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s\n%s", i.Title(), i.Description())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
