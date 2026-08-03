package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

const fordContent = `# Ford

Ford Motor Company was founded in 1903 by Henry Ford in Dearborn, Michigan.
Ford pioneered the moving assembly line, revolutionizing mass production
and making cars affordable for ordinary Americans.

Henry Ford's introduction of the Model T in 1908 and the Highland Park
assembly line in 1913 transformed the automobile from a luxury item into
an everyday necessity. Ford remains one of the largest automakers in the
world and the second-largest in the United States.
`

const bmwContent = `# BMW

Bayerische Motoren Werke AG (BMW) was founded in 1916 in Munich, Germany.
Originally an aircraft engine manufacturer, BMW pivoted to motorcycles
and then automobiles, becoming synonymous with performance and luxury.

BMW is headquartered in Munich and operates production facilities across
the globe. The company is known for its "Ultimate Driving Machine" brand
philosophy and owns the Rolls-Royce and MINI marques in addition to the
core BMW brand.
`

func lessonLog() lesson.Lesson {
	return lesson.Lesson{
		Number:   5,
		Title:    "Viewing History with git log",
		Subtitle: "The Book of Automobiles: Three Chapters In",
		Description: `Your book now has the Toyota chapter committed. To make this lesson more
interesting, we've added the Ford and BMW chapters behind the scenes — so you'll
have a proper multi-commit history to explore.

git log is your window into that history. You'll use it constantly.`,
		Steps: []lesson.Step{
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					if !s.IsGitRepo() {
						if err := s.InitRepo(); err != nil {
							return err
						}
					}
					// Ensure Toyota is committed
					if !s.HasCommits() {
						if err := s.WriteFile("toyota.md", toyotaContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Toyota chapter"); err != nil {
							return err
						}
					}
					// Add Ford chapter if not present
					if !s.FileExists("ford.md") {
						if err := s.WriteFile("ford.md", fordContent); err != nil {
							return err
						}
						if err := s.QuickCommit("Add Ford chapter"); err != nil {
							return err
						}
					}
					// Add BMW chapter if not present
					if !s.FileExists("bmw.md") {
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
				Title: "The history so far",
				Content: `Your repository now has three commits — one for each manufacturer chapter.
git log shows these commits from newest to oldest, with full details about
each one (author, date, hash, message).

The basic format looks like this:`,
				Expected: `commit a7f3c12d8e9b4f5a6c7d8e9f0a1b2c3d4e5f6789
Author: Your Name <you@example.com>
Date:   Sun Jan 12 14:30:00 2025 +0000

    Add BMW chapter

commit b8e4d23c9f0a5b6c7d8e9f0a1b2c3d4e5f67890
Author: Your Name <you@example.com>
Date:   Sun Jan 12 14:29:00 2025 +0000

    Add Ford chapter

commit c9f5e34d0a1b6c7d8e9f0a1b2c3d4e5f678901
Author: Your Name <you@example.com>
Date:   Sun Jan 12 14:28:00 2025 +0000

    Add Toyota chapter`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "View the full commit log",
				Content: `Run git log to see the complete history of your repository. Press 'q' to exit
the pager when you're done reading.`,
				Command:  `git log`,
				Expected: `(see multi-line commit history above)`,
				Hint:     "Run: git log — press 'q' to quit the pager if it pauses",
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "View the compact log",
				Content: `The full log is verbose. For a quick overview, --oneline shows one commit per line:
each line is the short SHA (first 7 characters) plus the commit message.`,
				Command:  `git log --oneline`,
				Expected: `a7f3c12 Add BMW chapter
b8e4d23 Add Ford chapter
c9f5e34 Add Toyota chapter`,
				Hint: "Run: git log --oneline",
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Useful git log options",
				Content: `git log is highly customizable. Here are the most useful variants:

  git log --oneline
    ↳ One line per commit: SHA + message

  git log --oneline --graph
    ↳ ASCII art showing branch topology — essential for complex repos

  git log --oneline --graph --all
    ↳ Same, but shows ALL branches (not just current one)

  git log -5
    ↳ Show only the last 5 commits

  git log --author="Ford"
    ↳ Show only commits by that author

  git log --since="2 weeks ago"
    ↳ Show commits from the last two weeks

  git log -- toyota.md
    ↳ Show only commits that touched toyota.md

  git log -p
    ↳ Show the actual diff (changes) for each commit — very detailed

  git log --stat
    ↳ Show which files changed and how many lines

  Try a few of these now if you're curious — the sandbox directory is still there.`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Understanding the SHA hash",
				Content: `Every commit has a SHA-1 hash — a 40-character hexadecimal fingerprint.

  Full SHA:  a7f3c12d8e9b4f5a6c7d8e9f0a1b2c3d4e5f6789
  Short SHA: a7f3c12  (first 7 characters — unique enough for most repos)

  The SHA is deterministic: it's computed from the commit's content, author,
  timestamp, and parent SHA. Any change to any of these produces a completely
  different hash — that's how Git detects corruption and guarantees integrity.

  You can reference any commit by its SHA (or a unique prefix) in any git command:
    git show a7f3c12           ← show what changed in that commit
    git checkout a7f3c12       ← go back to that point in history
    git revert a7f3c12         ← create a new commit that undoes those changes`,
			},
		},
	}
}
