# Bonsai

A beautiful CLI tool for managing and pruning stale Git branches, just like trimming a bonsai tree.

Built with [Charm Bracelet](https://charm.sh/) libraries for an elegant terminal experience.

## Features

- **Local & Remote Branch Cleanup**: Manage both local and remote branches
- **Interactive Mode**: Beautiful TUI for selecting branches to delete
- **Bulk Mode**: Quickly delete all stale branches at once
- **Configurable Age Thresholds**: Customize what "stale" means for your workflow
- **Dry-Run Mode**: Preview what would be deleted without making changes
- **Safety First**: Never deletes current branch or protected branches (main/master/develop)
- **Rich Terminal UI**: Built with Bubble Tea, Lip Gloss, and Bubbles

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/kriscoleman-testifysec/bonsai.git
cd bonsai

# Install locally
make install

# Make sure ~/.local/bin is in your PATH
export PATH="$HOME/.local/bin:$PATH"
```

### Build Only

```bash
make build
# Binary will be in ./build/bonsai
```

## Usage

### Basic Commands

```bash
# Clean up local branches (interactive mode)
bonsai local

# Clean up remote branches (interactive mode)
bonsai remote

# Show what would be deleted without deleting (dry-run)
bonsai local --dry-run

# Delete all stale local branches at once (bulk mode)
bonsai local --bulk

# Delete all stale remote branches at once (bulk mode)
bonsai remote --bulk
```

### Customizing Age Thresholds

```bash
# Local branches older than 1 week
bonsai local --age 1w

# Local branches older than 7 days (equivalent)
bonsai local --age 7d

# Remote branches older than 30 days
bonsai remote --age 30d

# Remote branches older than 720 hours (equivalent)
bonsai remote --age 720h
```

Supported time formats:
- `w` - weeks (e.g., `2w` = 2 weeks)
- `d` - days (e.g., `14d` = 14 days)
- `h` - hours (e.g., `336h` = 336 hours)
- `m` - minutes (e.g., `20160m`)
- `s` - seconds (e.g., `1209600s`)

### Remote Branch Options

```bash
# Specify a different remote (default is "origin")
bonsai remote --remote upstream

# Combine with other options
bonsai remote --remote upstream --age 4w --dry-run
```

## Configuration File

Bonsai supports configuration files to set default values. Configuration files are searched in the following order:

1. `.bonsai.yaml` or `.bonsai.yml` in the current directory
2. `~/.bonsai.yaml` or `~/.bonsai.yml` in your home directory
3. `$XDG_CONFIG_HOME/bonsai/config.yaml`

### Example Configuration

Create a `.bonsai.yaml` file (see `.bonsai.example.yaml` for a full example):

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

**Note**: Command-line flags always take precedence over configuration file settings.

## Interactive Mode

When running in interactive mode (default), you can:

- **Navigate**: Use arrow keys or `j`/`k` to move up and down
- **Toggle Selection**: Press `space` or `x` to select/deselect a branch
- **Select All**: Press `a` to select all branches
- **Select None**: Press `n` to deselect all branches
- **Delete Selected**: Press `enter` or `d` to delete selected branches
- **Search/Filter**: Start typing to filter branches by name
- **Quit**: Press `q`, `esc`, or `ctrl+c` to cancel

Each branch displays:
- Branch name
- Age (e.g., "2 weeks ago")
- Last commit message
- Last commit author

## Examples

### Clean up local branches older than 2 weeks (default)

```bash
bonsai local
```

### Preview what would be deleted

```bash
bonsai local --dry-run
```

Output:
```
Found 5 stale local branch(es)
Age threshold: 336h0m0s

[Lists branches that would be deleted]
```

### Bulk delete all stale remote branches older than 4 weeks

```bash
bonsai remote --bulk
```

You'll be prompted to confirm:
```
⚠ This will delete 3 branch(es). Are you sure? (y/N)
```

### Clean up feature branches after 1 week

```bash
bonsai local --age 1w --bulk
```

## How It Works

### Stale Branch Detection

Bonsai identifies stale branches based on the **last commit date**, not the branch creation date. This ensures accuracy.

**Default Thresholds:**
- Local branches: 2 weeks (14 days)
- Remote branches: 4 weeks (28 days)

### Protected Branches

The following branches are **never** deleted:
- Current branch (the one you're on)
- `main`
- `master`
- `develop`

### Git Operations

Bonsai uses `git for-each-ref` for efficient branch listing with metadata, ensuring good performance even in repositories with hundreds of branches.

## Development

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

### Building and Testing

```bash
# Build the project
make build

# Run tests
make test

# Generate test coverage
make coverage

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean

# Install locally
make install

# Uninstall
make uninstall
```

### Running from Source

```bash
go run ./cmd/bonsai local --dry-run
```

## Technology Stack

- **Language**: Go
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles)

## Requirements

- Go 1.21 or higher
- Git

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated releases. See [RELEASING.md](RELEASING.md) for details on the release process.

## License

See [LICENSE.md](LICENSE.md) for details.

## Credits

Built with love using [Charm Bracelet](https://charm.sh/) libraries.

---

**Keep your repository clean and tidy, just like a well-maintained bonsai tree!** 🌳
