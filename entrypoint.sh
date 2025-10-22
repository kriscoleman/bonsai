#!/bin/sh
set -e

# GitHub Action entrypoint for Bonsai
# This script runs inside the Docker container

# Fix git ownership issue in GitHub Actions
git config --global --add safe.directory /github/workspace
git config --global --add safe.directory '*'

MODE="${1:-remote}"
AGE="${2:-4w}"
REMOTE_NAME="${3:-origin}"
DRY_RUN="${4:-false}"
FORCE="${5:-false}"
CONFIG_FILE="${6:-}"

echo "🌳 Bonsai Branch Cleanup Action"
echo "================================"
echo "Mode: $MODE"
echo "Age threshold: $AGE"
echo "Remote: $REMOTE_NAME"
echo "Dry run: $DRY_RUN"
echo "Force: $FORCE"
echo ""

# Build the bonsai command
CMD="bonsai $MODE --age $AGE"

# Add flags based on inputs
if [ "$DRY_RUN" = "true" ]; then
  CMD="$CMD --dry-run"
fi

if [ "$FORCE" = "true" ]; then
  CMD="$CMD --force"
fi

if [ "$MODE" = "remote" ]; then
  CMD="$CMD --remote $REMOTE_NAME"
fi

# Use config file if provided
if [ -n "$CONFIG_FILE" ] && [ -f "$CONFIG_FILE" ]; then
  echo "📄 Using config file: $CONFIG_FILE"
  # Copy config to expected location
  cp "$CONFIG_FILE" .bonsai.yaml
fi

# For remote cleanup, we need to fetch first
if [ "$MODE" = "remote" ]; then
  echo "🔄 Fetching remote branches..."
  git fetch "$REMOTE_NAME" --prune
  echo ""
fi

# Run in bulk mode for CI/CD (non-interactive)
# --no-prompt skips confirmation prompts
CMD="$CMD --bulk --verbose --no-prompt"

echo "Running: $CMD"
echo ""

# Run bonsai and capture output
if OUTPUT=$($CMD 2>&1); then
  echo "$OUTPUT"

  # Try to extract counts from output (for GitHub Action outputs)
  # This is a best-effort parse of the summary
  DELETED_COUNT=$(echo "$OUTPUT" | grep -oE "[0-9]+ branches removed" | grep -oE "[0-9]+" || echo "0")
  FAILED_COUNT=$(echo "$OUTPUT" | grep -oE "[0-9]+ failed" | grep -oE "[0-9]+" || echo "0")

  # Set GitHub Action outputs
  echo "deleted-count=$DELETED_COUNT" >> "$GITHUB_OUTPUT"
  echo "failed-count=$FAILED_COUNT" >> "$GITHUB_OUTPUT"

  # If dry-run, show what would be deleted
  if [ "$DRY_RUN" = "true" ]; then
    echo ""
    echo "✅ Dry run complete - no changes made"
    exit 0
  fi

  # Check if any branches failed
  if [ "$FAILED_COUNT" != "0" ]; then
    echo ""
    echo "⚠️  Some branches failed to delete. Check the detailed error report above."
    exit 1
  fi

  echo ""
  echo "✅ Branch cleanup complete!"
  exit 0
else
  echo "$OUTPUT"
  echo ""
  echo "❌ Bonsai failed to run"
  exit 1
fi
