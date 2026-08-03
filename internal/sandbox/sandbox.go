package sandbox

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sandbox manages a temporary git repository for exercises.
type Sandbox struct {
	Dir      string
	RepoDir  string // subdirectory where the book repo lives
	RemoteDir string // bare repo acting as "origin"
}

// New creates a new sandbox in a temporary directory.
func New() (*Sandbox, error) {
	dir, err := os.MkdirTemp("", "simgit-sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("create sandbox dir: %w", err)
	}
	repoDir := filepath.Join(dir, "book-of-automobiles")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return nil, fmt.Errorf("create repo dir: %w", err)
	}
	return &Sandbox{Dir: dir, RepoDir: repoDir}, nil
}

// Resume reconnects to an existing sandbox directory.
func Resume(dir string) (*Sandbox, error) {
	repoDir := filepath.Join(dir, "book-of-automobiles")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("sandbox directory no longer exists: %s", dir)
	}
	s := &Sandbox{Dir: dir, RepoDir: repoDir}
	remoteDir := filepath.Join(dir, "book-of-automobiles-remote.git")
	if _, err := os.Stat(remoteDir); err == nil {
		s.RemoteDir = remoteDir
	}
	return s, nil
}

// Destroy removes the entire sandbox directory.
func (s *Sandbox) Destroy() error {
	return os.RemoveAll(s.Dir)
}

// Exec runs a git command inside the repo directory. Returns combined output.
func (s *Sandbox) Exec(args ...string) (string, error) {
	return s.execIn(s.RepoDir, args...)
}

// ExecIn runs a git command in a specific directory.
func (s *Sandbox) ExecIn(dir string, args ...string) (string, error) {
	return s.execIn(dir, args...)
}

func (s *Sandbox) execIn(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	// Isolate from user's global git config for setup steps
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=SimGit User",
		"GIT_AUTHOR_EMAIL=simgit@example.com",
		"GIT_COMMITTER_NAME=SimGit User",
		"GIT_COMMITTER_EMAIL=simgit@example.com",
		"GIT_TERMINAL_PROMPT=0",
	)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// ExecShell runs an arbitrary shell command in the repo directory.
func (s *Sandbox) ExecShell(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = s.RepoDir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// WriteFile creates or overwrites a file in the repo directory.
func (s *Sandbox) WriteFile(name, content string) error {
	path := filepath.Join(s.RepoDir, name)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ReadFile reads a file from the repo directory.
func (s *Sandbox) ReadFile(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(s.RepoDir, name))
	return string(data), err
}

// FileExists checks if a path exists in the repo directory.
func (s *Sandbox) FileExists(name string) bool {
	_, err := os.Stat(filepath.Join(s.RepoDir, name))
	return err == nil
}

// SetupRemote creates a bare repository to act as the "origin" remote.
func (s *Sandbox) SetupRemote() error {
	remoteDir := filepath.Join(s.Dir, "book-of-automobiles-remote.git")
	out, err := s.execIn(s.Dir, "init", "--bare", "book-of-automobiles-remote.git")
	if err != nil {
		return fmt.Errorf("init bare repo: %w\n%s", err, out)
	}
	s.RemoteDir = remoteDir
	return nil
}

// GetRemoteDir returns the path to the simulated remote.
func (s *Sandbox) GetRemoteDir() string {
	return s.RemoteDir
}

// ConfigureIdentity sets git user.name and user.email in the repo's local config.
func (s *Sandbox) ConfigureIdentity(name, email string) error {
	if _, err := s.Exec("config", "user.name", name); err != nil {
		return err
	}
	if _, err := s.Exec("config", "user.email", email); err != nil {
		return err
	}
	return nil
}

// InitRepo runs git init in the repo directory.
func (s *Sandbox) InitRepo() error {
	out, err := s.Exec("init", "-b", "main")
	if err != nil {
		// older git may not support -b flag
		out, err = s.Exec("init")
		if err != nil {
			return fmt.Errorf("git init: %w\n%s", err, out)
		}
		// rename master to main if needed
		s.Exec("checkout", "-b", "main") //nolint
	}
	_ = out
	// Configure identity in this repo
	s.Exec("config", "user.name", "SimGit User")  //nolint
	s.Exec("config", "user.email", "simgit@example.com") //nolint
	return nil
}

// QuickCommit adds all files and commits with the given message.
func (s *Sandbox) QuickCommit(message string) error {
	if _, err := s.Exec("add", "."); err != nil {
		return err
	}
	out, err := s.Exec("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("commit: %w\n%s", err, out)
	}
	return nil
}

// GetFirstCommitSHA returns the SHA of the first (root) commit.
func (s *Sandbox) GetFirstCommitSHA() (string, error) {
	out, err := s.Exec("rev-list", "--max-parents=0", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetLog returns git log --oneline output.
func (s *Sandbox) GetLog() (string, error) {
	return s.Exec("log", "--oneline")
}
