package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonBranch() lesson.Lesson {
	return lesson.Lesson{
		Number:   6,
		Title:    "Creating Branches",
		Subtitle: "The Book of Automobiles: Experimenting with Structure",
		Description: `Your book has three committed chapters. Now you want to experiment with adding
a Honda chapter, but you're not sure about the format yet. You don't want
to mess up the stable main book while you're experimenting.

The solution: branches. A branch lets you work in parallel without
affecting the main line of your book.`,
		Steps: []lesson.Step{
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					if !s.IsGitRepo() {
						if err := s.InitRepo(); err != nil {
							return err
						}
					}
					if !s.HasCommits() {
						if err := s.WriteFile("toyota.md", toyotaContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Toyota chapter"); err != nil {
							return err
						}
						if err := s.WriteFile("ford.md", fordContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Ford chapter"); err != nil {
							return err
						}
						if err := s.WriteFile("bmw.md", bmwContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add BMW chapter"); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What is a branch?",
				Content: `A branch is a lightweight, movable pointer to a commit. Creating a branch
doesn't copy any files — it's just a named pointer that advances as you
make new commits on it.

  main branch:
    [Toyota] ← [Ford] ← [BMW] ← HEAD

  After creating a "honda" branch and switching to it:
    [Toyota] ← [Ford] ← [BMW] ← main
                                  ↑
                                honda (same point, different pointer)

  After committing on "honda":
    [Toyota] ← [Ford] ← [BMW] ← main
                                  ↑
                                 [Honda] ← honda

  The main branch is untouched. Your Honda work is isolated on its own branch.
  You can switch between branches at any time — Git swaps your files instantly.

  Branches are the feature that makes collaborative development possible.
  Each developer works on their own branch; changes are merged when ready.`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "The main branch",
				Content: `By convention, the primary branch is called "main" (or "master" in older repos).
It's not special in any technical sense — it's just a branch like any other.
But teams treat it as the "stable" version of the project.

To see all your current branches:
  git branch

The branch with an asterisk (*) is the one you're currently on:

  * main      ← you are here`,
				Expected: `* main`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Create a honda branch",
				Content: `Create a new branch named "honda". This creates the branch but doesn't
switch to it — you stay on main.`,
				Command: `git branch honda`,
				Expected: `(no output — the branch was created silently)`,
				Hint:    "Run: git branch honda",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.HasBranch("honda") {
						return false, "Branch 'honda' not found. Run: git branch honda"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Listing branches",
				Content: `Run git branch to see all local branches. The asterisk marks the current branch:`,
				Expected: `  honda
* main`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Create and switch in one step",
				Content: `git branch just creates a branch. To create AND switch to it simultaneously,
use git checkout -b (or the modern git switch -c).

Create a new branch called "mercedes" and switch to it immediately:`,
				Command:  `git checkout -b mercedes`,
				Expected: `Switched to a new branch 'mercedes'`,
				Hint:     "Run: git checkout -b mercedes",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.HasBranch("mercedes") {
						return false, "Branch 'mercedes' not found. Run: git checkout -b mercedes"
					}
					current, err := s.CurrentBranch()
					if err != nil || current != "mercedes" {
						return false, "You're not on 'mercedes'. Run: git checkout -b mercedes"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Branch naming conventions",
				Content: `Good branch names are descriptive and lowercase with hyphens:

  feature/add-honda-chapter     ← common prefix pattern
  fix/typo-in-toyota
  honda-chapter
  mercedes

  Avoid spaces (use hyphens), avoid special characters except - and /.

  Common branch naming patterns:
    feature/* — new features
    fix/*     — bug fixes
    docs/*    — documentation
    chore/*   — maintenance tasks

  Run git branch to see all your branches now:`,
				Expected: `  honda
* mercedes
  main`,
			},
		},
	}
}
