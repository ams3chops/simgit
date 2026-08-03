package lessons

import (
	"strings"
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonCommit() lesson.Lesson {
	return lesson.Lesson{
		Number:   4,
		Title:    "Committing Changes",
		Subtitle: "The Book of Automobiles: Sealing the Toyota Chapter",
		Description: `The Toyota chapter is staged and ready. Now we create a "commit" — a permanent
snapshot that will be preserved in your book's history forever.

A commit is the fundamental unit of Git. Everything else (branches, merges,
history navigation) revolves around commits.`,
		Steps: []lesson.Step{
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					if !s.IsGitRepo() {
						if err := s.InitRepo(); err != nil {
							return err
						}
					}
					if !s.FileExists("toyota.md") {
						if err := s.WriteFile("toyota.md", toyotaContent); err != nil {
							return err
						}
					}
					if !s.IsFileStaged("toyota.md") {
						_, err := s.Exec("add", "toyota.md")
						return err
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What is a commit?",
				Content: `A commit is a snapshot of your files at a specific moment in time. Each commit has:

  • A unique SHA hash    — e.g., a3f2c1d7e9b4... (40 hex characters, usually shown as 7)
  • An author           — your name and email from git config
  • A timestamp         — when the commit was created
  • A parent reference  — the commit that came before this one
  • A commit message    — your description of what changed and why

  Commits form a chain:

    [root commit] ← [2nd commit] ← [3rd commit] ← HEAD

  HEAD is just a pointer to "where you are now" — usually the latest commit.

  A good commit message is the single most important habit in Git. Future you
  (and your teammates) will be grateful for clear messages like:
    "Add Toyota chapter"
  rather than:
    "stuff" or "fix" or "asdfgh"`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Commit the Toyota chapter",
				Content: `Create your first commit. The -m flag lets you write the commit message inline.
The message should be concise (under 72 characters) and describe what this commit adds.`,
				Command:  `git commit -m "Add Toyota chapter"`,
				Expected: `[main (root-commit) a3f2c1d] Add Toyota chapter
 1 file changed, 17 insertions(+)
 create mode 100644 toyota.md`,
				Hint: `Run: git commit -m "Add Toyota chapter" — make sure toyota.md is staged first (git add toyota.md)`,
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.HasCommits() {
						return false, "No commits found. Run: git commit -m \"Add Toyota chapter\""
					}
					msg, err := s.LastCommitMessage()
					if err != nil || !strings.Contains(strings.ToLower(msg), "toyota") {
						return false, "Commit found, but message should mention Toyota. Try: git commit -m \"Add Toyota chapter\""
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Understanding the commit output",
				Content: `After committing, Git shows you a summary:

  [main (root-commit) a3f2c1d] Add Toyota chapter
   1 file changed, 17 insertions(+)
   create mode 100644 toyota.md

  Breaking this down:
  • main              — the branch the commit was made on
  • (root-commit)     — this is the first commit (no parent)
  • a3f2c1d           — the first 7 characters of the commit's SHA hash
  • Add Toyota chapter — your commit message
  • 1 file changed    — number of files affected
  • 17 insertions     — lines added
  • create mode 100644 toyota.md — the file was newly created`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Writing good commit messages",
				Content: `The commit message is your gift to the future. A few guidelines:

  Subject line (first line):
    • Keep it under 72 characters
    • Use the imperative mood: "Add feature" not "Added feature"
    • Be specific: "Add Toyota chapter with founding history" not just "update"

  For longer messages, leave a blank line after the subject, then add details:

    Add Toyota chapter

    Covers Toyota's founding history, growth into a global manufacturer,
    and pioneering role in hybrid vehicle technology.

  The subject line shows in git log --oneline. The body is for context.

  Run git status after committing — it should now say:
    "nothing to commit, working tree clean"`,
			},
		},
	}
}
