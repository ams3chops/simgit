package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

const hondaContent = `# Honda

Honda Motor Co., Ltd. was founded in 1948 by Soichiro Honda in Hamamatsu, Japan.
Starting as a motorcycle manufacturer, Honda expanded into automobiles in 1963
and is now one of the world's largest automakers.

Honda is renowned for its engineering culture, fuel-efficient engines, and the
principle of "The Power of Dreams." It is the world's largest manufacturer of
internal combustion engines by volume and has been a pioneer in low-emission
vehicle technology.
`

func lessonMerge() lesson.Lesson {
	return lesson.Lesson{
		Number:   8,
		Title:    "Merging Branches",
		Subtitle: "The Book of Automobiles: Honda Joins the Main Book",
		Description: `The Honda chapter has been written on a separate branch. Now it's time to
bring it back into the main book. This is called "merging" — combining the
work from one branch into another.`,
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
					}
					// Ensure we're on main
					s.Exec("checkout", "main") //nolint
					// Create honda branch with Honda content
					if !s.HasBranch("honda") {
						s.Exec("checkout", "-b", "honda") //nolint
						s.WriteFile("honda.md", hondaContent)  //nolint
						s.QuickCommit("Add Honda chapter")      //nolint
						s.Exec("checkout", "main")             //nolint
					} else {
						// Check if honda has the file
						s.Exec("checkout", "honda") //nolint
						if !s.FileExists("honda.md") {
							s.WriteFile("honda.md", hondaContent) //nolint
							s.QuickCommit("Add Honda chapter")    //nolint
						}
						s.Exec("checkout", "main") //nolint
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What merging does",
				Content: `Merging takes the changes from one branch and incorporates them into another.
You run it from the branch you want to receive the changes (the "target").

  Before merge:
    main:  [Toyota] ← [Ford] ← HEAD
    honda: [Toyota] ← [Ford] ← [Honda] ← HEAD

  After: git merge honda (while on main):
    main:  [Toyota] ← [Ford] ← [Honda] ← HEAD

  The honda branch still exists unchanged. main now includes Honda's commit.`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Fast-forward vs three-way merge",
				Content: `There are two merge strategies Git uses automatically:

  Fast-forward merge (when main hasn't moved since the branch was created):
    Git simply moves the main pointer forward — no new commit needed.

    main:         A ← B
    feature:      A ← B ← C ← D

    After merge:  A ← B ← C ← D  ← main (fast-forwarded)

  Three-way merge (when both branches have diverged):
    Git creates a new "merge commit" that has two parents.

    main:          A ← B ← E
    feature:       A ← B ← C ← D

    After merge:   A ← B ← C ← D ← M  ← main
                         ↖         ↗
                           E ──────

  For our honda branch, it's a fast-forward since main hasn't advanced
  since we branched. Git will say "Fast-forward".`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Merge the honda branch into main",
				Content: `You should be on the main branch. Merge the honda branch into main.
This will bring the Honda chapter into the main book.`,
				Command:  `git merge honda`,
				Expected: `Updating b8e4d23..c9a3f12
Fast-forward
 honda.md | 18 ++++++++++++++++++
 1 file changed, 18 insertions(+)
 create mode 100644 honda.md`,
				Hint: "Make sure you're on main first (git checkout main), then run: git merge honda",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					current, _ := s.CurrentBranch()
					if current != "main" {
						return false, "You need to be on 'main' to merge into it. Run: git checkout main"
					}
					if !s.FileExists("honda.md") {
						return false, "honda.md is not in main yet. Run: git merge honda"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Verifying the merge",
				Content: `After merging, run git log --oneline to see the Honda commit is now on main:`,
				Expected: `c9a3f12 Add Honda chapter
b8e4d23 Add Ford chapter
a7f3c12 Add Toyota chapter`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Merge conflicts",
				Content: `Sometimes two branches modify the same lines in the same file. Git can't
automatically decide which version to keep — this is a merge conflict.

When a conflict occurs, Git marks the file with conflict markers:

  <<<<<<< HEAD
  Corolla – Japan's most reliable car
  =======
  Corolla – Budget-friendly family sedan
  >>>>>>> honda

  You must:
  1. Edit the file to resolve the conflict (choose one, combine both, or rewrite)
  2. Remove the <<<, ===, >>> markers
  3. git add <file>
  4. git commit (Git auto-generates the merge commit message)

  Tips to avoid conflicts:
  • Keep branches short-lived — merge often
  • Communicate with teammates about who's changing what
  • Use git pull --rebase to keep your branch up-to-date`,
			},
		},
	}
}
