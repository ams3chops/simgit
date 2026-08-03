package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonSwitch() lesson.Lesson {
	return lesson.Lesson{
		Number:   7,
		Title:    "Switching Branches",
		Subtitle: "The Book of Automobiles: Moving Between Timelines",
		Description: `You've created branches. Now learn to move between them. Switching branches
is how you jump between different lines of work — your files in the working
directory update instantly to reflect the state of the branch you switch to.`,
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
					}
					// Ensure we're on main
					s.Exec("checkout", "main") //nolint
					// Ensure honda branch exists
					if !s.HasBranch("honda") {
						s.Exec("branch", "honda") //nolint
					}
					// Ensure mercedes branch exists
					if !s.HasBranch("mercedes") {
						s.Exec("branch", "mercedes") //nolint
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "How branch switching works",
				Content: `When you switch branches, Git:
  1. Takes the snapshot from the target branch's latest commit
  2. Replaces your working directory files with that snapshot
  3. Updates HEAD to point to the new branch

  This happens instantly, even for large repositories. Git is efficient
  because it only copies files that differ between branches.

  Important: if you have uncommitted changes, Git may refuse to switch
  (to avoid losing your work). Always commit or stash before switching.

  Current state:
    main     — has Toyota, Ford, BMW chapters
    honda    — identical to main (just created, no new commits)
    mercedes — identical to main (just created, no new commits)`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Switch to the honda branch",
				Content: `Use git switch to move to the honda branch. "git switch" is the modern
command for switching branches (introduced in Git 2.23). It's clearer
than git checkout because it only switches branches — nothing else.`,
				Command:  `git switch honda`,
				Expected: `Switched to branch 'honda'`,
				Hint:     "Run: git switch honda",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					current, err := s.CurrentBranch()
					if err != nil {
						return false, "Could not determine current branch."
					}
					if current != "honda" {
						return false, "You're on '" + current + "', not 'honda'. Run: git switch honda"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Switch back to main",
				Content: `Now switch back to main using the classic git checkout syntax.
Both "git switch" and "git checkout" work for switching — you'll see
both in the wild, so it's good to know both.`,
				Command:  `git checkout main`,
				Expected: `Switched to branch 'main'`,
				Hint:     "Run: git checkout main",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					current, err := s.CurrentBranch()
					if err != nil {
						return false, "Could not determine current branch."
					}
					if current != "main" {
						return false, "You're on '" + current + "', not 'main'. Run: git checkout main"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "git switch vs git checkout",
				Content: `Git 2.23 (2019) introduced git switch and git restore to replace
the overloaded git checkout command:

  Old way (still works):
    git checkout <branch>           ← switch branch
    git checkout -b <new-branch>    ← create + switch
    git checkout -- <file>          ← discard file changes (dangerous!)

  New way (clearer):
    git switch <branch>             ← switch branch
    git switch -c <new-branch>      ← create + switch
    git restore <file>              ← discard file changes

  The old git checkout did too many things, which confused beginners.
  The new commands have one job each.

  Both work. You'll see git checkout in tutorials, Stack Overflow answers,
  and documentation written before 2019. Don't be surprised.

  Tip: If you switch and can't find a file, check git status — it may
  be on a different branch. Files only exist on their branch until merged.`,
			},
		},
	}
}
