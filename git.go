package main

import (
	"fmt"
	"os"
	"strings"
)

func gitHandler(args []string) int {
	if len(args) == 0 {
		return passthrough("git", args)
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "status":
		return gitStatus(rest)
	case "log":
		return gitLog(rest)
	case "diff":
		return gitDiff(rest)
	case "add":
		return gitAdd(rest)
	case "fetch":
		return gitFetch(rest)
	case "push":
		return gitPush(rest)
	case "pull":
		return gitPull(rest)
	case "commit":
		return gitCommit(rest)
	default:
		return passthrough("git", args)
	}
}

// gitStatus runs git status -sb for compact branch + changed files.
func gitStatus(args []string) int {
	a := append([]string{"status", "-sb"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	fmt.Print(r.Stdout)
	return 0
}

// gitLog runs git log --oneline -10, compact by default.
func gitLog(args []string) int {
	// If the user passed their own flags, passthrough instead.
	if len(args) > 0 {
		return passthrough("git", append([]string{"log"}, args...))
	}
	r := run("git", "log", "--oneline", "-10")
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	fmt.Print(r.Stdout)
	return 0
}

// gitDiff runs git diff --stat plus any user args.
func gitDiff(args []string) int {
	a := append([]string{"diff", "--stat"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	fmt.Print(r.Stdout)
	return 0
}

// gitAdd runs git add and prints a compact confirmation.
func gitAdd(args []string) int {
	a := append([]string{"add"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	fmt.Println("ok ✓")
	return 0
}

// gitFetch runs git fetch --quiet and prints a compact confirmation.
func gitFetch(args []string) int {
	a := append([]string{"fetch", "--quiet"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	fmt.Println("ok ✓")
	return 0
}

// gitPush runs git push and parses output for branch info.
func gitPush(args []string) int {
	a := append([]string{"push"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	branch := parsePushBranch(r.Stderr) // git push outputs to stderr
	if branch != "" {
		fmt.Printf("ok ✓ %s\n", branch)
	} else {
		fmt.Println("ok ✓")
	}
	return 0
}

// parsePushBranch extracts the branch name from git push stderr.
// Looks for lines like "   abc1234..def5678  main -> main"
func parsePushBranch(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, " -> "); idx != -1 {
			branch := strings.TrimSpace(line[idx+4:])
			return branch
		}
	}
	return ""
}

// gitPull runs git pull and prints a compact summary.
func gitPull(args []string) int {
	a := append([]string{"pull"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		return r.ExitCode
	}
	summary := parsePullOutput(r.Stdout)
	fmt.Println(summary)
	return 0
}

// parsePullOutput extracts a compact summary from git pull stdout.
func parsePullOutput(stdout string) string {
	if strings.Contains(stdout, "Already up to date") {
		return "ok ✓ up to date"
	}
	// Count files changed from the summary line like "3 files changed, 10 insertions(+)"
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "files changed") || strings.Contains(line, "file changed") {
			return "ok ✓ " + line
		}
	}
	return "ok ✓"
}

// gitCommit runs git commit and prints a compact result with short hash.
func gitCommit(args []string) int {
	a := append([]string{"commit"}, args...)
	r := run("git", a...)
	if r.ExitCode != 0 {
		fmt.Fprint(os.Stderr, r.Stderr)
		// Also print stdout since git commit may output there on error
		if r.Stdout != "" {
			fmt.Print(r.Stdout)
		}
		return r.ExitCode
	}
	hash := parseCommitHash(r.Stdout)
	if hash != "" {
		fmt.Printf("ok ✓ %s\n", hash)
	} else {
		fmt.Println("ok ✓")
	}
	return 0
}

// parseCommitHash extracts the short hash from git commit output.
// First line looks like: "[main abc1234] commit message"
func parseCommitHash(stdout string) string {
	line := strings.SplitN(stdout, "\n", 2)[0]
	// Find content inside brackets: [branch hash]
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	inner := line[start+1 : end]
	parts := strings.Fields(inner)
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}
