package main

import "testing"

func TestParsePushBranch(t *testing.T) {
	tests := []struct {
		name   string
		stderr string
		want   string
	}{
		{
			name:   "typical push",
			stderr: "To github.com:user/repo.git\n   abc1234..def5678  main -> main\n",
			want:   "main",
		},
		{
			name:   "feature branch",
			stderr: "To github.com:user/repo.git\n   111..222  feat/login -> feat/login\n",
			want:   "feat/login",
		},
		{
			name:   "no arrow",
			stderr: "Everything up-to-date\n",
			want:   "",
		},
		{
			name:   "empty",
			stderr: "",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePushBranch(tt.stderr)
			if got != tt.want {
				t.Errorf("parsePushBranch() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParsePullOutput(t *testing.T) {
	tests := []struct {
		name   string
		stdout string
		want   string
	}{
		{
			name:   "up to date",
			stdout: "Already up to date.\n",
			want:   "ok ✓ up to date",
		},
		{
			name:   "files changed",
			stdout: "Updating abc..def\nFast-forward\n main.go | 5 +++++\n 1 file changed, 5 insertions(+)\n",
			want:   "ok ✓ 1 file changed, 5 insertions(+)",
		},
		{
			name:   "multiple files",
			stdout: "Updating abc..def\nFast-forward\n a.go | 2 ++\n b.go | 3 ---\n 2 files changed, 2 insertions(+), 3 deletions(-)\n",
			want:   "ok ✓ 2 files changed, 2 insertions(+), 3 deletions(-)",
		},
		{
			name:   "empty",
			stdout: "",
			want:   "ok ✓",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePullOutput(tt.stdout)
			if got != tt.want {
				t.Errorf("parsePullOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseCommitHash(t *testing.T) {
	tests := []struct {
		name   string
		stdout string
		want   string
	}{
		{
			name:   "typical commit",
			stdout: "[main abc1234] fix: resolve bug\n 1 file changed, 2 insertions(+)\n",
			want:   "abc1234",
		},
		{
			name:   "detached HEAD",
			stdout: "[detached HEAD def5678] wip\n",
			want:   "def5678",
		},
		{
			name:   "root commit",
			stdout: "[main (root-commit) aaa1111] initial\n",
			want:   "aaa1111",
		},
		{
			name:   "no brackets",
			stdout: "nothing to commit\n",
			want:   "",
		},
		{
			name:   "empty",
			stdout: "",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommitHash(tt.stdout)
			if got != tt.want {
				t.Errorf("parseCommitHash() = %q, want %q", got, tt.want)
			}
		})
	}
}
