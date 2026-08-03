package lessons

import (
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

const toyotaContent = `# Toyota

Toyota Motor Corporation was founded in 1937 by Kiichiro Toyoda.
Headquartered in Toyota City, Aichi, Japan, it is one of the world's
largest automobile manufacturers by production volume.

Founded initially as a division of Toyoda Automatic Loom Works, Toyota
became an independent company and grew into a global powerhouse known
for reliability, quality manufacturing, and pioneering hybrid technology
with the introduction of the Prius in 1997.
`

func lessonStaging() lesson.Lesson {
	return lesson.Lesson{
		Number:   3,
		Title:    "Staging Files",
		Subtitle: "The Book of Automobiles: Toyota Chapter",
		Description: `You've written the first chapter of your book: Toyota. The file is sitting in
your repository folder, but Git doesn't know about it yet.

This lesson introduces the staging area — one of Git's most important concepts,
and the one that confuses beginners most often.`,
		Steps: []lesson.Step{
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					// Ensure the repo is initialized
					if !s.IsGitRepo() {
						if err := s.InitRepo(); err != nil {
							return err
						}
					}
					// Write the Toyota chapter
					return s.WriteFile("toyota.md", toyotaContent)
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "The three areas of Git",
				Content: `Git has three distinct areas that files move through:

  1. Working Directory  — your actual files on disk. This is what you see in
                          your folder. You edit files here.

  2. Staging Area        — a "draft pile" of changes you're preparing to commit.
     (Index)               You choose exactly what goes into the next commit.

  3. Repository          — the permanent history. Once committed, a snapshot
     (.git)                is saved forever (or until you explicitly delete it).

  The flow looks like this:

     edit files       git add          git commit
  [Working Dir] ──────────────> [Staging Area] ──────────────> [Repository]

  This two-step process (add, then commit) gives you fine-grained control.
  You can stage only part of your changes — very useful when you've edited
  multiple files but only want to commit some of them.`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Seeing the untracked file",
				Content: `We've created the Toyota chapter file for you. Run git status to see it:`,
				Expected: `On branch main

No commits yet

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	toyota.md

nothing added to commit but untracked files present (use "git add" to track)`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Stage the Toyota chapter",
				Content: `The toyota.md file is "untracked" — Git sees it but isn't watching it.
Use git add to move it to the staging area. This tells Git: "I want this
file to be part of my next commit."`,
				Command:  `git add toyota.md`,
				Expected: `(no output — silence is success)`,
				Hint:     `Run: git add toyota.md — make sure you're in the book-of-automobiles directory`,
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.IsFileStaged("toyota.md") {
						return false, "toyota.md is not staged yet. Run: git add toyota.md"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Confirming the file is staged",
				Content: `Now run git status again. Notice the file moved from "Untracked" to
"Changes to be committed" — it's now in the staging area, shown in green:`,
				Expected: `On branch main

No commits yet

Changes to be committed:
  (use "git rm --cached <file>..." to unstage)
	new file:   toyota.md`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Staging tips",
				Content: `A few useful git add variations:

  git add toyota.md         ← stage one specific file
  git add *.md              ← stage all markdown files
  git add .                 ← stage everything in the current directory
  git add -p                ← interactive: stage specific lines within a file

  To unstage a file (remove from staging area but keep your edits):
    git restore --staged toyota.md
    (older git: git reset HEAD toyota.md)

  The staging area is your friend. Use it to craft clean, focused commits
  rather than dumping all your changes at once.`,
			},
		},
	}
}
