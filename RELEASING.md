# Release Process

This project uses [release-please](https://github.com/googleapis/release-please) to automate releases.

## How it works

1. **Conventional Commits**: All commits to the `main` branch should follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

2. **Release PR**: When commits are merged to `main`, release-please will automatically create or update a release PR that:
   - Updates the version number
   - Updates CHANGELOG.md with the changes
   - Creates a GitHub release when the PR is merged

3. **Binary Building**: When a release is created, GoReleaser automatically builds binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)

## Commit Message Format

Use these prefixes to indicate the type of change:

- `feat:` - New features (triggers a minor version bump)
- `fix:` - Bug fixes (triggers a patch version bump)
- `feat!:` or `fix!:` - Breaking changes (triggers a major version bump)
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code refactoring
- `test:` - Test updates
- `ci:` - CI/CD changes
- `perf:` - Performance improvements

### Examples

```bash
# Feature (bumps 1.0.0 -> 1.1.0)
git commit -m "feat: add support for custom branch patterns"

# Bug fix (bumps 1.0.0 -> 1.0.1)
git commit -m "fix: handle empty git repositories gracefully"

# Breaking change (bumps 1.0.0 -> 2.0.0)
git commit -m "feat!: change configuration file format to TOML"

# With scope
git commit -m "feat(ui): add color-coded branch age indicators"

# With detailed description
git commit -m "feat: add support for GitLab repositories

- Add GitLab API client
- Support GitLab-specific branch protection rules
- Update documentation with GitLab examples"
```

## Release Workflow

1. **Make changes** and commit them using conventional commit messages
2. **Create a PR** to the `main` branch
3. **Merge the PR** - release-please will create or update a release PR
4. **Review the release PR** - check the version bump and changelog
5. **Merge the release PR** - this triggers:
   - Version tag creation
   - GitHub release creation
   - Binary builds via GoReleaser
   - Artifacts uploaded to the release

## Manual Release

If you need to manually trigger a release:

```bash
# Install release-please CLI
npm install -g release-please

# Create a release PR manually
release-please release-pr --token=$GITHUB_TOKEN --repo-url=kriscoleman/bonsai

# Create a GitHub release manually
release-please github-release --token=$GITHUB_TOKEN --repo-url=kriscoleman/bonsai
```

## Testing GoReleaser Locally

To test the GoReleaser configuration without publishing:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Run a snapshot build (doesn't publish)
goreleaser release --snapshot --clean

# Check the dist/ directory for built binaries
ls -la dist/
```

## Version Management

The current version is tracked in `.release-please-manifest.json`. Release-please automatically updates this file when creating releases.

Initial version: `0.1.0`

## Troubleshooting

### Release PR not created

- Ensure commits follow conventional commit format
- Check GitHub Actions logs for errors
- Verify GitHub token has proper permissions (contents: write, pull-requests: write)

### GoReleaser fails

- Check `.goreleaser.yml` syntax
- Verify Go version compatibility
- Review GitHub Actions logs for specific errors

## Resources

- [Release Please Documentation](https://github.com/googleapis/release-please)
- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
