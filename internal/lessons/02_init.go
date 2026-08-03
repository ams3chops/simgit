package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonInit() lesson.Lesson {
	return lesson.Lesson{
		Number:   2,
		Title:    "Initializing a Repository",
		Subtitle: "Creating the workspace for your book",
		Description: `The Book of Automobiles needs a home. In Git, that home is a "repository" —
a folder that Git watches and tracks. You create one with a single command: git init.

We've set up a sandbox directory for you to work in. Think of it as an empty
folder on your desktop named "book-of-automobiles". Let's turn it into a
Git repository.`,
		Steps: []lesson.Step{
			{
				Kind:  lesson.KindExplain,
				Title: "What is a repository?",
				Content: `A Git repository is just a regular folder with a special hidden subdirectory
called .git. That .git folder is where Git stores all the history, configuration,
and metadata for your project.

You never need to touch .git directly — Git manages it for you. But it's good
to know it's there.

  book-of-automobiles/
  ├── .git/              ← Git's brain — don't edit this!
  │   ├── config         ← repo-specific settings
  │   ├── HEAD           ← pointer to current branch
  │   ├── objects/       ← all commits, files stored here
  │   └── refs/          ← branch and tag pointers
  ├── toyota.md          ← your actual files go here
  └── ford.md`,
			},
			{
				Kind:    lesson.KindSetup,
				Content: "Ensure the sandbox directory exists",
				Setup: func(s *sandbox.Sandbox) error {
					// Sandbox directory is already created by New()
					// Just make sure we're clean
					return nil
				},
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Initialize the repository",
				Content: `Navigate to the sandbox directory (shown below) and initialize a new Git repository.
This command creates the .git folder and turns your plain directory into a tracked repository.`,
				Command:  `git init`,
				Expected: `Initialized empty Git repository in /tmp/simgit-sandbox-.../book-of-automobiles/.git/`,
				Hint:     `Make sure you cd into the book-of-automobiles directory first, then run: git init`,
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.IsGitRepo() {
						return false, "No .git directory found. Run 'git init' in the sandbox directory."
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Checking repository status",
				Content: `After initializing, run git status to see the current state of your repository.
Since it's brand new with no files, you'll see:`,
				Expected: `On branch main

No commits yet

nothing to commit (create/copy files and use "git add" to track)`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What git status tells you",
				Content: `git status is your most-used command. It tells you:

  • Which branch you're on (more on branches later)
  • Whether you have any uncommitted changes
  • Which files are staged (ready to commit)
  • Which files are untracked (Git doesn't know about them yet)

Think of git status as "what's going on right now?" — you'll run it constantly.
Get in the habit of running it before and after every git operation.`,
			},
		},
	}
}
