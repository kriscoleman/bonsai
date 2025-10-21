# Testing Guide

This document describes the testing infrastructure for Bonsai.

## Test Organization

Tests are organized into two categories:

### Unit Tests
Located in `*_test.go` files without build tags. These tests:
- Test individual functions and methods in isolation
- Use mocked or stubbed dependencies
- Run quickly without external dependencies
- Include tests for:
  - Branch age calculation
  - Duration parsing
  - Protected branch detection
  - Branch metadata handling
  - Configuration file loading

### Integration Tests
Located in `integration_test.go` with `// +build integration` tag. These tests:
- Create real Git repositories in temporary directories
- Test actual Git operations end-to-end
- Verify branch listing, filtering, and deletion workflows
- Test with various branch ages and configurations

## Running Tests

### Run Unit Tests Only
```bash
make test
```

### Run Integration Tests Only
```bash
make test-integration
```

### Run All Tests (Unit + Integration)
```bash
make test-all
```

### Generate Coverage Report (Unit Tests)
```bash
make coverage
# Opens coverage.html in browser
```

### Generate Coverage Report (All Tests)
```bash
make coverage-all
# Opens coverage.html with complete coverage
```

## Test Coverage

Current test coverage:
- **Config Package**: 75.9%
- **Git Package**: 67.8% (includes integration tests)
- **Overall Internal Packages**: ~70%

## Integration Test Examples

### Test Scenarios Covered

1. **Repository Validation**
   - Verifying Git repository detection
   - Testing non-repository error handling

2. **Branch Listing**
   - Listing local branches with metadata
   - Testing branches with different ages
   - Verifying current branch detection

3. **Branch Deletion**
   - Standard branch deletion
   - Force deletion of unmerged branches
   - Error handling for deletion failures

4. **Branch Filtering**
   - Age-based filtering
   - Protected branch exclusion
   - Current branch exclusion
   - Complete workflow testing

5. **Protected Branches**
   - main/master/develop protection
   - Protection flag verification

## Writing New Tests

### Unit Test Template
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "test case description",
            input:   testInput,
            want:    expectedOutput,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Test Template
```go
// +build integration

func TestIntegration_FeatureName(t *testing.T) {
    helper := NewTestHelper(t)
    helper.InitRepo()
    helper.CreateInitialCommit()
    
    // Test setup and execution
    
    // Assertions
}
```

## Test Utilities

### TestHelper Methods
- `InitRepo()` - Initialize a test Git repository
- `CreateInitialCommit()` - Create first commit
- `CreateBranch(name, checkout)` - Create a new branch
- `CreateBranchWithCommit(name, msg)` - Create branch with commit
- `CreateBranchWithAge(name, daysAgo)` - Create branch with old commit
- `CheckoutBranch(name)` - Switch branches
- `GetCurrentBranch()` - Get active branch
- `BranchExists(name)` - Check if branch exists
- `ListBranches()` - List all branches

## Continuous Integration

Tests are designed to run in CI environments:
- No external dependencies required
- Self-contained test repositories
- Automatic cleanup via `t.TempDir()`
- Platform-independent

## Troubleshooting

### Integration Tests Failing
- Ensure Git is installed and in PATH
- Check Git version (2.23+recommended)
- Verify temp directory permissions

### Coverage Not Generated
- Ensure all packages compile successfully
- Run `make build` first to verify
- Check for syntax errors in test files
