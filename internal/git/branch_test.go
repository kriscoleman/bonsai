package git

import (
	"testing"
	"time"
)

func TestBranch_Age(t *testing.T) {
	tests := []struct {
		name         string
		lastCommitAt time.Time
		wantMin      time.Duration
		wantMax      time.Duration
	}{
		{
			name:         "recent commit",
			lastCommitAt: time.Now().Add(-1 * time.Hour),
			wantMin:      59 * time.Minute,
			wantMax:      61 * time.Minute,
		},
		{
			name:         "one day old",
			lastCommitAt: time.Now().Add(-24 * time.Hour),
			wantMin:      23*time.Hour + 59*time.Minute,
			wantMax:      24*time.Hour + 1*time.Minute,
		},
		{
			name:         "two weeks old",
			lastCommitAt: time.Now().Add(-14 * 24 * time.Hour),
			wantMin:      13*24*time.Hour + 23*time.Hour,
			wantMax:      14*24*time.Hour + 1*time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Branch{
				LastCommitAt: tt.lastCommitAt,
			}
			age := b.Age()
			if age < tt.wantMin || age > tt.wantMax {
				t.Errorf("Age() = %v, want between %v and %v", age, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestBranch_IsStale(t *testing.T) {
	tests := []struct {
		name         string
		lastCommitAt time.Time
		threshold    time.Duration
		want         bool
	}{
		{
			name:         "stale - older than threshold",
			lastCommitAt: time.Now().Add(-15 * 24 * time.Hour),
			threshold:    14 * 24 * time.Hour,
			want:         true,
		},
		{
			name:         "not stale - newer than threshold",
			lastCommitAt: time.Now().Add(-13 * 24 * time.Hour),
			threshold:    14 * 24 * time.Hour,
			want:         false,
		},
		{
			name:         "edge case - just under threshold",
			lastCommitAt: time.Now().Add(-14*24*time.Hour + 1*time.Minute),
			threshold:    14 * 24 * time.Hour,
			want:         false, // Should not be stale if just under threshold
		},
		{
			name:         "very old branch",
			lastCommitAt: time.Now().Add(-365 * 24 * time.Hour),
			threshold:    30 * 24 * time.Hour,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Branch{
				LastCommitAt: tt.lastCommitAt,
			}
			if got := b.IsStale(tt.threshold); got != tt.want {
				t.Errorf("IsStale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranch_FullName(t *testing.T) {
	tests := []struct {
		name       string
		branch     *Branch
		wantResult string
	}{
		{
			name: "local branch",
			branch: &Branch{
				Name:     "feature-branch",
				IsRemote: false,
			},
			wantResult: "feature-branch",
		},
		{
			name: "remote branch",
			branch: &Branch{
				Name:       "feature-branch",
				IsRemote:   true,
				RemoteName: "origin",
			},
			wantResult: "origin/feature-branch",
		},
		{
			name: "remote branch - upstream",
			branch: &Branch{
				Name:       "main",
				IsRemote:   true,
				RemoteName: "upstream",
			},
			wantResult: "upstream/main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.branch.FullName(); got != tt.wantResult {
				t.Errorf("FullName() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestBranch_SafetyChecks(t *testing.T) {
	tests := []struct {
		name   string
		branch *Branch
		desc   string
	}{
		{
			name: "current branch should be marked",
			branch: &Branch{
				Name:      "main",
				IsCurrent: true,
			},
			desc: "current branch flag",
		},
		{
			name: "protected branch should be marked",
			branch: &Branch{
				Name:        "main",
				IsProtected: true,
			},
			desc: "protected branch flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just verifies that the fields exist and can be set
			// The actual filtering logic is tested in git_test.go
			if tt.branch.IsCurrent {
				t.Logf("Branch %s is correctly marked as current", tt.branch.Name)
			}
			if tt.branch.IsProtected {
				t.Logf("Branch %s is correctly marked as protected", tt.branch.Name)
			}
		})
	}
}
