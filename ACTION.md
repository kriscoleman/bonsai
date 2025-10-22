# 🌳 Bonsai Branch Cleanup Action

Automatically prune stale Git branches in your GitHub Actions workflows with the elegance and precision of bonsai cultivation.

Perfect for trunk-based development teams who want to keep their repository clean and tidy!

[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Bonsai%20Branch%20Cleanup-blue?logo=github)](https://github.com/marketplace/actions/bonsai-branch-cleanup)

## 🚀 Quick Setup (1 Minute!)

Copy [`.github/STARTER_WORKFLOW.yml`](.github/STARTER_WORKFLOW.yml) to your repository's `.github/workflows/` directory and you're done!

Or add this minimal workflow:

```yaml
name: Branch Cleanup
on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly

permissions:
  contents: write

jobs:
  cleanup:
    uses: kriscoleman/bonsai/.github/workflows/cleanup-branches.yml@v1
    with:
      age: '4w'
```

## Features

- 🌿 **Automated Cleanup**: Remove stale branches on schedule or after merges
- 🛡️ **Safe by Default**: Protected branches (main/master/develop) are never deleted
- ⚙️ **Highly Configurable**: Use inline parameters or `.bonsai.yaml` config files
- 🎯 **Flexible Time Units**: Support for years, months, weeks, days, hours
- 📊 **Detailed Reporting**: Outputs deleted/failed counts for workflow tracking
- 🔒 **Force Mode**: Optional force delete for unmerged branches

## Quick Start

### Option 1: Reusable Workflow (Easiest!)

The simplest way to use Bonsai is via our reusable workflow:

```yaml
name: Branch Cleanup
on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday

permissions:
  contents: write

jobs:
  cleanup:
    uses: kriscoleman/bonsai/.github/workflows/cleanup-branches.yml@v1
    with:
      mode: 'remote'
      age: '4w'
      dry-run: false
```

That's it! The reusable workflow handles checkout, fetch, and cleanup automatically.

### Option 2: Direct Action Usage

For more control, use the action directly:

#### Basic Usage - Clean Remote Branches

```yaml
name: Branch Cleanup
on:
  schedule:
    - cron: '0 0 * * 0'  # Run weekly on Sunday
  workflow_dispatch:      # Allow manual trigger

jobs:
  cleanup:
    runs-on: ubuntu-latest
    permissions:
      contents: write  # Required to delete branches
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Fetch all branches

      - name: Prune stale remote branches
        uses: kriscoleman/bonsai@v1
        with:
          mode: 'remote'
          age: '4w'
          remote-name: 'origin'
```

### Advanced Usage - With Config File

```yaml
- name: Prune branches using config
  uses: kriscoleman/bonsai@v1
  with:
    config-file: '.bonsai.yaml'
    mode: 'remote'
```

Create `.bonsai.yaml` in your repo:

```yaml
remote:
  age_threshold: "4w"
  remote_name: "origin"

protected_branches:
  - "production"
  - "staging"
```

## Reusable Workflow Inputs

When using the reusable workflow (`kriscoleman/bonsai/.github/workflows/cleanup-branches.yml@v1`):

| Input | Description | Type | Default |
|-------|-------------|------|---------|
| `mode` | Cleanup mode: `local` or `remote` | string | `remote` |
| `age` | Age threshold (e.g., `4w`, `30d`, `1y`) | string | `4w` |
| `remote-name` | Remote name to clean up | string | `origin` |
| `dry-run` | Preview mode (no deletions) | boolean | `false` |
| `force` | Force delete unmerged branches | boolean | `false` |
| `config-file` | Path to bonsai config file | string | `''` |

### Reusable Workflow Outputs

| Output | Description |
|--------|-------------|
| `deleted-count` | Number of branches successfully deleted |
| `failed-count` | Number of branches that failed to delete |

## Action Inputs

When using the action directly (`kriscoleman/bonsai@v1`):

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `mode` | Cleanup mode: `local` or `remote` | No | `remote` |
| `age` | Age threshold (e.g., `4w`, `30d`, `1y`, `12M`) | No | `4w` |
| `remote-name` | Remote name to clean up | No | `origin` |
| `dry-run` | Preview mode (no deletions) | No | `false` |
| `force` | Force delete unmerged branches | No | `false` |
| `config-file` | Path to bonsai config file | No | `''` |
| `github-token` | GitHub token for authentication | No | `${{ github.token }}` |

### Supported Time Units

| Unit | Description | Example |
|------|-------------|---------|
| `y` | Years (365 days) | `1y` |
| `M` | Months (30 days, 365 for 12M) | `12M` |
| `w` | Weeks | `4w` |
| `d` | Days | `30d` |
| `h` | Hours | `720h` |
| `m` | Minutes (lowercase) | `60m` |
| `s` | Seconds | `3600s` |

## Outputs

| Output | Description |
|--------|-------------|
| `deleted-count` | Number of branches successfully deleted |
| `failed-count` | Number of branches that failed to delete |
| `branch-list` | JSON array of deleted branch names |

## Usage Examples

### Trunk-Based Development - Weekly Cleanup

```yaml
name: Weekly Branch Cleanup
on:
  schedule:
    - cron: '0 2 * * 1'  # Monday at 2 AM

jobs:
  cleanup-old-branches:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Clean up stale remote branches
        uses: kriscoleman/bonsai@v1
        with:
          mode: 'remote'
          age: '4w'
          dry-run: 'false'
```

### After Pull Request Merge

```yaml
name: Cleanup on Merge
on:
  pull_request:
    types: [closed]
    branches: [main]

jobs:
  cleanup:
    if: github.event.pull_request.merged == true
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Prune stale branches after merge
        uses: kriscoleman/bonsai@v1
        with:
          mode: 'remote'
          age: '2w'
```

### Aggressive Cleanup with Force

```yaml
- name: Force cleanup very old branches
  uses: kriscoleman/bonsai@v1
  with:
    mode: 'remote'
    age: '1y'
    force: 'true'  # Delete even if not fully merged
```

### Dry Run for Testing

```yaml
- name: Preview what would be deleted
  uses: kriscoleman/bonsai@v1
  with:
    mode: 'remote'
    age: '4w'
    dry-run: 'true'
```

### Custom Protected Branches

```yaml
- name: Cleanup with custom protection
  uses: kriscoleman/bonsai@v1
  with:
    config-file: '.github/bonsai-config.yaml'
    mode: 'remote'
```

Create `.github/bonsai-config.yaml`:

```yaml
remote:
  age_threshold: "6w"

protected_branches:
  - "main"
  - "develop"
  - "production"
  - "staging"
  - "release/*"  # Protect all release branches
```

## Safety Features

### Built-in Protection

The following branches are **automatically protected** from deletion:
- Current branch (in local mode)
- `main`
- `master`
- `develop`
- Any branches specified in `protected_branches` config

### Permissions Required

The action needs these permissions in your workflow:

```yaml
permissions:
  contents: write  # Required to delete branches
```

## Outputs Usage

Use outputs to track cleanup metrics:

```yaml
- name: Cleanup branches
  id: bonsai
  uses: kriscoleman/bonsai@v1
  with:
    mode: 'remote'
    age: '4w'

- name: Report results
  run: |
    echo "Deleted: ${{ steps.bonsai.outputs.deleted-count }}"
    echo "Failed: ${{ steps.bonsai.outputs.failed-count }}"
```

## Troubleshooting

### Authentication Issues

If you encounter permission errors:

1. Ensure `contents: write` permission is set
2. Verify the repository settings allow Actions to delete branches
3. Check that protected branch rules aren't blocking deletion

### Branches Not Deleting

- **"Not fully merged" errors**: Use `force: 'true'` to force delete
- **Protected branches**: Add them to `.bonsai.yaml` protected list
- **Recent activity**: Adjust the `age` threshold

## Best Practices

### For Trunk-Based Development

```yaml
# Clean up feature branches weekly
on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly

with:
  mode: 'remote'
  age: '2w'  # Delete branches older than 2 weeks
```

### For Long-Lived Feature Branches

```yaml
with:
  mode: 'remote'
  age: '3M'  # Keep branches for 3 months
```

### Test First with Dry Run

```yaml
# Step 1: Test with dry-run
with:
  dry-run: 'true'
  age: '4w'

# Step 2: After verification, remove dry-run
```

## Philosophy

Just as a bonsai master carefully prunes their tree to maintain its beauty and health, this action helps you maintain a clean, healthy repository by removing branches that no longer serve you.

**"The best time to prune was 20 years ago. The second best time is now."**

---

## Related

- [Bonsai CLI Tool](https://github.com/kriscoleman/bonsai) - The standalone CLI version
- [Documentation](https://github.com/kriscoleman/bonsai#readme) - Full documentation
- [Report Issues](https://github.com/kriscoleman/bonsai/issues) - Bug reports and feature requests

## License

MIT License - See [LICENSE](LICENSE.md) for details

---

Built with 💜 using [Charm Bracelet](https://charm.sh/)
