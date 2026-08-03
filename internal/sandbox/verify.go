package sandbox

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitRepo returns true if .git exists in the repo directory.
func (s *Sandbox) IsGitRepo() bool {
	_, err := os.Stat(filepath.Join(s.RepoDir, ".git"))
	return err == nil
}

// HasCommits returns true if there is at least one commit.
func (s *Sandbox) HasCommits() bool {
	out, err := s.Exec("log", "--oneline", "-1")
	return err == nil && strings.TrimSpace(out) != ""
}

// CommitCount returns number of commits reachable from HEAD.
func (s *Sandbox) CommitCount() int {
	out, err := s.Exec("rev-list", "--count", "HEAD")
	if err != nil {
		return 0
	}
	count := 0
	for _, ch := range strings.TrimSpace(out) {
		if ch >= '0' && ch <= '9' {
			count = count*10 + int(ch-'0')
		}
	}
	return count
}

// HasBranch checks if a local branch name exists.
func (s *Sandbox) HasBranch(name string) bool {
	out, err := s.Exec("branch", "--list", name)
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// CurrentBranch returns the current branch name. Returns "HEAD" if detached.
func (s *Sandbox) CurrentBranch() (string, error) {
	out, err := s.Exec("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// IsDetachedHEAD returns true if HEAD is detached.
func (s *Sandbox) IsDetachedHEAD() bool {
	branch, err := s.CurrentBranch()
	if err != nil {
		return false
	}
	return branch == "HEAD"
}

// IsFileStaged checks if a file has been added to the index (staging area).
func (s *Sandbox) IsFileStaged(filename string) bool {
	out, err := s.Exec("diff", "--cached", "--name-only")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == filename {
			return true
		}
	}
	return false
}

// HasRemote checks if a remote with the given name exists.
func (s *Sandbox) HasRemote(name string) bool {
	out, err := s.Exec("remote")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == name {
			return true
		}
	}
	return false
}

// LastCommitMessage returns the message of the most recent commit.
func (s *Sandbox) LastCommitMessage() (string, error) {
	out, err := s.Exec("log", "-1", "--pretty=%s")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// BranchCommitCount returns commit count on a named branch.
func (s *Sandbox) BranchCommitCount(branch string) int {
	out, err := s.Exec("rev-list", "--count", branch)
	if err != nil {
		return 0
	}
	count := 0
	for _, ch := range strings.TrimSpace(out) {
		if ch >= '0' && ch <= '9' {
			count = count*10 + int(ch-'0')
		}
	}
	return count
}

// GlobalConfigGet returns a global git config value by key.
func GlobalConfigGet(key string) string {
	cmd := exec.Command("git", "config", "--global", key)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	_ = cmd.Run()
	return strings.TrimSpace(buf.String())
}
