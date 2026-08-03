package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonConfig() lesson.Lesson {
	return lesson.Lesson{
		Number:   1,
		Title:    "Installing and Configuring Git",
		Subtitle: "Setting up your authoring environment",
		Description: `You've decided to write the definitive Book of Automobiles — a comprehensive reference
covering every major automobile manufacturer. Before you write a single word, you need
to set up Git: the version control system that will track every change you make.

Git is like a time machine for your files. It remembers every version of your work,
who changed what, and when. Professional developers, writers, and teams worldwide
rely on it. Let's get started.`,
		Steps: []lesson.Step{
			{
				Kind:  lesson.KindExplain,
				Title: "What is Git?",
				Content: `Git is a distributed version control system. Unlike saving a file (which overwrites
the previous version), Git keeps a complete history of every change ever made.

Key ideas:
  • A "repository" (repo) is a folder tracked by Git
  • A "commit" is a snapshot of your files at a point in time
  • Every commit has a unique ID, a message, an author, and a timestamp
  • You can go back to any previous commit at any time

Git works entirely on your local machine — no internet required for most operations.
You only need a network when sharing with others via services like GitHub or GitLab.`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Installing Git",
				Content: `Before configuring Git, you need to install it. Here's how on each platform:

  macOS:
    Open Terminal and run: xcode-select --install
    Or install via Homebrew: brew install git

  Linux (Ubuntu/Debian):
    sudo apt update && sudo apt install git

  Linux (Fedora/RHEL):
    sudo dnf install git

  Windows:
    Download from https://git-scm.com/download/win
    Or use: winget install Git.Git

  Verify your installation by running:
    git --version

  You should see something like: git version 2.43.0`,
				Expected: `git version 2.43.0`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Why configure Git?",
				Content: `Every Git commit is stamped with the author's name and email address. Without
this configuration, Git will either refuse to commit or use an auto-generated
identity that looks messy in history.

You set this once globally and it applies to every repository on your machine.
This is stored in your home directory at ~/.gitconfig.

The two essential settings are:
  • user.name  — your name (appears in every commit you make)
  • user.email — your email (links your commits to services like GitHub)`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Set your name",
				Content: `Set your Git identity — the name that will appear on every commit you make.
Replace "Your Name" with your actual name.`,
				Command:  `git config --global user.name "Your Name"`,
				Expected: `(no output — silence means success)`,
				Hint:     `Run: git config --global user.name "Your Name" (replace with your actual name)`,
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					val := sandbox.GlobalConfigGet("user.name")
					if val == "" {
						return false, "user.name is not set. Run: git config --global user.name \"Your Name\""
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Set your email",
				Content: `Now set your email address. Use the same email you plan to use with GitHub or GitLab
so your commits are linked to your account.`,
				Command:  `git config --global user.email "you@example.com"`,
				Expected: `(no output — silence means success)`,
				Hint:     `Run: git config --global user.email "you@example.com" (replace with your email)`,
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					val := sandbox.GlobalConfigGet("user.email")
					if val == "" {
						return false, "user.email is not set. Run: git config --global user.email \"you@example.com\""
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Viewing your configuration",
				Content: `To see all your Git settings, run:`,
				Expected: `$ git config --list

user.name=Your Name
user.email=you@example.com
core.editor=vim
init.defaultbranch=main
...`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Optional but recommended settings",
				Content: `A few more settings that make life easier:

  Set the default branch name to "main" (modern convention):
    git config --global init.defaultBranch main

  Set your preferred text editor for commit messages:
    git config --global core.editor nano      # beginner-friendly
    git config --global core.editor vim       # for vim users
    git config --global core.editor "code --wait"  # VS Code

  Enable colorized output (usually on by default):
    git config --global color.ui auto

  View your ~/.gitconfig file to see all settings:
    cat ~/.gitconfig`,
			},
		},
	}
}
