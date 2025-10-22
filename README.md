<div align="center">

# 🌳 Bonsai

### *The Art of Branch Pruning*

**Keep your Git repository as elegant and intentional as a carefully cultivated bonsai tree.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE.md)
[![Built with Charm](https://img.shields.io/badge/Built%20with-Charm%20%F0%9F%92%9C-ff69b4)](https://charm.sh/)
[![GitHub Release](https://img.shields.io/github/v/release/kriscoleman/bonsai)](https://github.com/kriscoleman/bonsai/releases)
[![Tests](https://github.com/kriscoleman/bonsai/actions/workflows/test.yml/badge.svg)](https://github.com/kriscoleman/bonsai/actions/workflows/test.yml)

---

**Interactive** • **Intelligent** • **Beautiful**

</div>

## Why Bonsai?

Just as a bonsai master carefully prunes their tree to maintain its beauty and health, Bonsai helps you maintain a clean, healthy repository by removing stale branches that no longer serve you.

**The Problem:** Over time, repositories accumulate forgotten feature branches, old experiments, and merged PRs that clutter your workspace and slow down your workflow.

**The Solution:** Bonsai makes branch cleanup *delightful* with an elegant terminal interface that puts you in control.

Built with the exceptional [Charm Bracelet](https://charm.sh/) toolkit for a terminal experience that feels modern, responsive, and *joyful*.

## ✨ Features That Spark Joy

<table>
<tr>
<td width="33%" valign="top">

### 🎨 **Beautiful Interface**

Experience branch cleanup through an elegant, keyboard-driven TUI that makes maintenance feel effortless and even *enjoyable*.

</td>
<td width="33%" valign="top">

### 🛡️ **Safety First**

Smart protection prevents deletion of your current branch, main/master/develop, and any custom protected branches you specify.

</td>
<td width="33%" valign="top">

### 🎯 **Precision Control**

Interactive selection, bulk operations, dry-run previews—work exactly the way *you* want to work.

</td>
</tr>
<tr>
<td width="33%" valign="top">

### 🌍 **Local & Remote**

Manage both local and remote branches from a single, unified interface. No more scattered git commands.

</td>
<td width="33%" valign="top">

### ⚙️ **Highly Configurable**

Define your own age thresholds, protect custom branches, and save preferences in configuration files.

</td>
<td width="33%" valign="top">

### ⚡ **Blazing Fast**

Optimized Git operations handle repositories with hundreds of branches without breaking a sweat.

</td>
</tr>
</table>

---

## 🚀 Quick Start

### Installation

**Pre-built Binaries (Recommended):**

Visit the [latest release](https://github.com/kriscoleman/bonsai/releases/latest) or download for your platform:

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/kriscoleman/bonsai/releases/latest/download/bonsai_Darwin_arm64.tar.gz
tar -xzf bonsai_Darwin_arm64.tar.gz
sudo mv bonsai /usr/local/bin/

# macOS (Intel)
curl -LO https://github.com/kriscoleman/bonsai/releases/latest/download/bonsai_Darwin_x86_64.tar.gz
tar -xzf bonsai_Darwin_x86_64.tar.gz
sudo mv bonsai /usr/local/bin/

# Linux (x86_64)
curl -LO https://github.com/kriscoleman/bonsai/releases/latest/download/bonsai_Linux_x86_64.tar.gz
tar -xzf bonsai_Linux_x86_64.tar.gz
sudo mv bonsai /usr/local/bin/

# Linux (ARM64)
curl -LO https://github.com/kriscoleman/bonsai/releases/latest/download/bonsai_Linux_arm64.tar.gz
tar -xzf bonsai_Linux_arm64.tar.gz
sudo mv bonsai /usr/local/bin/
```

**From Source:**

```bash
git clone https://github.com/kriscoleman/bonsai.git
cd bonsai
make install

# Ensure ~/.local/bin is in your PATH
export PATH="$HOME/.local/bin:$PATH"
```

**Build Only:**

```bash
make build
# Binary will be available at ./build/bonsai
```

### Your First Pruning Session

```bash
# Preview stale local branches (safe, no changes made)
bonsai local --dry-run

# Interactive cleanup - pick exactly what to remove
bonsai local

# Quick bulk cleanup of remote branches
bonsai remote --bulk --age 4w
```

> **💡 Pro Tip:** Start with `--dry-run` to see what Bonsai would delete before making any changes.

---

## 📖 Usage Guide

### Core Commands

| Command | Description |
|---------|-------------|
| `bonsai local` | Clean up local branches (interactive mode) |
| `bonsai remote` | Clean up remote branches (interactive mode) |
| `bonsai local --dry-run` | Preview what would be deleted (safe!) |
| `bonsai local --bulk` | Delete all stale local branches at once |
| `bonsai remote --bulk` | Delete all stale remote branches at once |
| `bonsai local --bulk -v` | Show detailed error messages for failed deletions |

### Fine-Tune Your Pruning

**Age Thresholds** - Define "stale" on your terms:

```bash
bonsai local --age 1y      # 1 year old
bonsai local --age 12M     # 12 months old
bonsai local --age 1w      # 1 week old
bonsai local --age 7d      # 7 days old (equivalent)
bonsai remote --age 30d    # 30 days old
bonsai remote --age 720h   # 720 hours old (equivalent)
```

**Supported Time Units:**
`y` (years) • `M` (months, uppercase) • `w` (weeks) • `d` (days) • `h` (hours) • `m` (minutes, lowercase) • `s` (seconds)

**Remote Options** - Work with any remote:

```bash
# Target a specific remote (default: "origin")
bonsai remote --remote upstream

# Chain multiple options together
bonsai remote --remote upstream --age 4w --dry-run
```

**Debugging & Force Deletion**:

```bash
# Show detailed error messages when deletions fail
bonsai local --bulk --verbose
bonsai local --bulk -v  # Short form

# Verbose mode shows:
# - Full git error messages for each failed deletion
# - Detailed error report at the end
# - Smart suggestions (e.g., use --force for unmerged branches)

# Force delete unmerged branches (equivalent to git branch -D)
bonsai local --force
bonsai local -f  # Short form

# Combine flags for maximum control
bonsai local --bulk --force --verbose  # Force delete all, show details
bonsai local -bfv --age 1y             # Short form: bulk + force + verbose
```

---

## ⚙️ Configuration

Save your preferences in a configuration file and let Bonsai remember how you like to work.

**Config File Locations** (searched in order):

1. `.bonsai.yaml` / `.bonsai.yml` — Current directory
2. `~/.bonsai.yaml` / `~/.bonsai.yml` — Home directory
3. `$XDG_CONFIG_HOME/bonsai/config.yaml` — XDG config

### Example Configuration

Create `.bonsai.yaml` in your repository or home directory:

```yaml
# Local branch settings
local:
  age_threshold: "2w"  # 2 weeks

# Remote branch settings
remote:
  age_threshold: "4w"  # 4 weeks
  remote_name: "origin"

# Additional protected branches (beyond main/master/develop)
protected_branches:
  - "production"
  - "staging"
```

> **Note:** Command-line flags always override configuration file settings.

---

## 🎹 Interactive Mode

When you run Bonsai in interactive mode (the default), you get a beautiful terminal UI with full keyboard control:

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` `↓` or `j` `k` | Navigate up and down |
| `space` or `x` | Toggle branch selection |
| `a` | Select all branches |
| `n` | Deselect all branches |
| `enter` or `d` | Delete selected branches |
| **Start typing** | Filter/search branches by name |
| `q` `esc` or `ctrl+c` | Quit without changes |

### Branch Information Display

Each branch shows you everything you need to make informed decisions:
- 📌 Branch name
- ⏰ Age (e.g., "2 weeks ago")
- 💬 Last commit message
- 👤 Last commit author

---

## 💡 Examples & Workflows

**Scenario 1: Weekly Maintenance**
```bash
# Your weekly ritual - clean up what's no longer needed
bonsai local
```

**Scenario 2: Safety First Approach**
```bash
# Preview before you prune (always a good idea!)
bonsai local --dry-run
```
Output example:
```
Found 5 stale local branch(es)
Age threshold: 336h0m0s

[Detailed list of branches that would be deleted]
```

**Scenario 3: Aggressive Cleanup**
```bash
# Quickly remove all stale remote branches
bonsai remote --bulk
```
You'll be asked to confirm:
```
⚠ This will delete 3 branch(es). Are you sure? (y/N)
```

**Scenario 4: Fast-Moving Project**
```bash
# Clean up feature branches after just 1 week
bonsai local --age 1w --bulk
```

---

## 🔧 How It Works

### Smart Detection

Bonsai identifies stale branches based on the **last commit date**, not the branch creation date—ensuring accuracy and fairness.

**Default Age Thresholds:**
- 🏠 Local branches: **2 weeks** (14 days)
- 🌐 Remote branches: **4 weeks** (28 days)

### Built-in Safety

The following branches are **automatically protected** from deletion:
- ✓ Your current branch (the one you're on)
- ✓ `main` / `master` / `develop`
- ✓ Any additional branches you specify in config

### Performance

Under the hood, Bonsai uses `git for-each-ref` for efficient branch listing with full metadata. This means it stays fast even in repositories with hundreds of branches.

---

## 👩‍💻 For Developers

### Project Structure

```
bonsai/
├── cmd/bonsai/         # CLI commands and entry point
│   ├── main.go
│   ├── root.go
│   ├── local.go
│   └── remote.go
├── internal/
│   ├── git/            # Git operations and branch management
│   │   ├── git.go
│   │   └── branch.go
│   ├── ui/             # Terminal UI components
│   │   └── interactive.go
│   └── config/         # Configuration and parsing
│       └── config.go
├── Makefile
├── go.mod
└── README.md
```

### Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile the project |
| `make test` | Run the test suite |
| `make coverage` | Generate test coverage report |
| `make fmt` | Format code |
| `make lint` | Run linter |
| `make clean` | Remove build artifacts |
| `make install` | Install to `~/.local/bin` |
| `make uninstall` | Remove from system |

### Running from Source

```bash
go run ./cmd/bonsai local --dry-run
```

### Technology Stack

Built with modern, battle-tested Go libraries:

| Component | Library |
|-----------|---------|
| **Language** | [Go 1.21+](https://go.dev) |
| **CLI Framework** | [Cobra](https://github.com/spf13/cobra) |
| **TUI Framework** | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| **Styling** | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| **UI Components** | [Bubbles](https://github.com/charmbracelet/bubbles) |

---

## 🤝 Contributing

We welcome contributions! Whether it's:
- 🐛 Bug reports
- 💡 Feature requests
- 📖 Documentation improvements
- 🔧 Code contributions

Please feel free to submit a Pull Request or open an issue.

**Commit Convention:** This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated releases. See [RELEASING.md](RELEASING.md) for details.

---

## 📄 License

Released under the MIT License. See [LICENSE.md](LICENSE.md) for details.

---

<div align="center">

## 🌸 The Bonsai Philosophy

*"The best time to prune a bonsai tree was 20 years ago. The second best time is now."*

Just like the ancient art of bonsai cultivation, keeping a clean repository is about **mindful maintenance**, **intentional growth**, and **respect for your craft**.

---

**Keep your repository clean and tidy, just like a well-maintained bonsai tree!** 🌳

Built with care using [Charm Bracelet](https://charm.sh/) 💜

</div>
