package lessons

import (
	"fmt"
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonCheckout() lesson.Lesson {
	var firstSHA string

	return lesson.Lesson{
		Number:   9,
		Title:    "Navigating History with git checkout",
		Subtitle: "The Book of Automobiles: Time Travel",
		Description: `Your book has a rich history now. Git lets you travel back in time and see
exactly what the book looked like at any previous commit. This is one of
Git's most powerful features — and one of the reasons "losing work" is
nearly impossible once something is committed.`,
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
					if !s.FileExists("ford.md") {
						if err := s.WriteFile("ford.md", fordContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Ford chapter"); err != nil {
							return err
						}
					}
					if !s.FileExists("bmw.md") {
						if err := s.WriteFile("bmw.md", bmwContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add BMW chapter"); err != nil {
							return err
						}
					}
					if !s.FileExists("honda.md") {
						if err := s.WriteFile("honda.md", hondaContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Honda chapter"); err != nil {
							return err
						}
					}
					// Cache the first commit SHA
					sha, err := s.GetFirstCommitSHA()
					if err != nil {
						return err
					}
					firstSHA = sha
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Your book's history",
				Content: `Your repository now has four chapters committed. Run git log --oneline
to see the history:`,
				Expected: `d4f8a92 Add Honda chapter
c9a3f12 Add BMW chapter
b8e4d23 Add Ford chapter
a7f3c12 Add Toyota chapter`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What does git checkout do with a SHA?",
				Content: `So far, git checkout has been about branches. But you can also pass a
commit SHA to jump directly to that point in history:

  git checkout a7f3c12     ← go back to the Toyota-only state

When you do this, Git puts you in "detached HEAD" state. Normally,
HEAD points to a branch (like main), which points to the latest commit.
In detached HEAD, HEAD points directly to a commit — no branch.

  Normal:      HEAD → main → [latest commit]
  Detached:    HEAD → [some old commit]

In detached HEAD you can look around, run the code, read the files —
but any commits you make here are "orphaned" and will be lost when
you switch branches. Think of it as read-only time travel.

To do real work on an old commit, create a new branch from it:
  git checkout -b new-feature a7f3c12`,
			},
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					// Refresh firstSHA in case it wasn't set
					sha, err := s.GetFirstCommitSHA()
					if err != nil {
						return err
					}
					firstSHA = sha
					return nil
				},
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Travel back to the first commit",
				Content: fmt.Sprintf(`Travel back in time to the very first commit (Toyota only). Run git log --oneline
to find the SHA of the first commit (bottom of the list), then use it:

  git checkout <first-commit-sha>

The SHA in your sandbox is: %s
(Use your actual SHA from git log --oneline — the bottom entry)`, firstSHA),
				Command:  `git checkout <first-commit-sha>`,
				Expected: `Note: switching to 'a7f3c12'.

You are in 'detached HEAD' state. You can look around, make experimental
changes and commit them, but any commits you make in this state are not
part of any branch...

HEAD is now at a7f3c12 Add Toyota chapter`,
				Hint: "Run: git log --oneline to find the SHA of 'Add Toyota chapter' (bottom line), then: git checkout <that-sha>",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.IsDetachedHEAD() {
						return false, "HEAD is not detached. Try: git checkout <first-commit-sha> (find the SHA with git log --oneline)"
					}
					// Verify only toyota.md is present (ford.md should NOT be there)
					if s.FileExists("ford.md") {
						return false, "ford.md still exists — you may not have checked out the right commit. Look for the first (oldest) commit SHA."
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Look around in the past",
				Content: `You're now in the state where only toyota.md exists. Run ls to confirm:`,
				Expected: `toyota.md`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Return to the present",
				Content: `Time travel is fun, but you need to get back to the latest state.
Checkout the main branch to reattach HEAD and return to the present.`,
				Command:  `git checkout main`,
				Expected: `Previous HEAD position was a7f3c12 Add Toyota chapter
Switched to branch 'main'`,
				Hint: "Run: git checkout main",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					current, err := s.CurrentBranch()
					if err != nil {
						return false, "Could not determine current branch."
					}
					if current != "main" {
						return false, "You're still detached. Run: git checkout main"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Other uses of git checkout",
				Content: `git checkout has another important use: discarding changes to a file.

  Scenario: you edited toyota.md but made a mistake and want to undo it.

  git checkout -- toyota.md
  (or with modern git:)
  git restore toyota.md

  WARNING: This permanently discards your unsaved changes — they cannot
  be recovered. Only do this if you're sure you want to throw away edits.

  Summary of git checkout uses:
    git checkout main           ← switch to 'main' branch
    git checkout -b feature     ← create + switch to 'feature' branch
    git checkout a7f3c12        ← detached HEAD at that commit
    git checkout -- toyota.md   ← discard changes to toyota.md

  This overloaded behavior is why Git 2.23 introduced git switch and
  git restore as clearer alternatives.`,
			},
		},
	}
}
