package lessons

import (
	"fmt"
	"os"
	"path/filepath"
	"simgit/internal/lesson"
	"simgit/internal/sandbox"
)

func lessonRemote() lesson.Lesson {
	return lesson.Lesson{
		Number:   10,
		Title:    "Remote Repositories",
		Subtitle: "The Book of Automobiles: Publishing to the World",
		Description: `Your book is complete locally. Now it's time to share it with the world.
Remote repositories are copies of your repo hosted on servers — GitHub, GitLab,
Bitbucket, or any git server. They enable collaboration, backup, and deployment.

Since we can't connect to GitHub in this sandbox, we'll use a local bare repository
as a stand-in "remote" — it behaves identically to GitHub from Git's perspective.`,
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
					// Set up the bare remote repository
					if s.GetRemoteDir() == "" {
						if err := s.SetupRemote(); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "What is a remote?",
				Content: `A remote is a named reference to another copy of your repository hosted
somewhere else — usually a server. By convention, the primary remote
is named "origin".

  Your machine                    GitHub / GitLab / Server
  ─────────────                   ───────────────────────
  book-of-automobiles/    ←push→  origin/book-of-automobiles
  (local repo)           ←pull←  (remote repo)

  Common services:
  • GitHub (github.com)    — most popular, owned by Microsoft
  • GitLab (gitlab.com)    — popular for self-hosting, CI/CD built-in
  • Bitbucket              — popular in enterprise, integrates with Jira
  • Gitea                  — open-source self-hosted option

  Remotes let you:
  • Backup your work (push to remote = backup)
  • Collaborate (others pull your changes)
  • Deploy (servers pull from the remote)
  • Use pull requests / merge requests for code review`,
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Creating a repo on GitHub",
				Content: `In real life, before connecting a remote, you'd create the repository on GitHub:

  1. Go to github.com and sign in
  2. Click the "+" icon → "New repository"
  3. Name it "book-of-automobiles"
  4. Leave it empty (don't initialize with README)
  5. Click "Create repository"

  GitHub will show you commands to connect your existing repo:
    git remote add origin https://github.com/username/book-of-automobiles.git
    git branch -M main
    git push -u origin main

  For this lesson, we'll use a local directory as our "remote" instead of
  GitHub — it's identical from Git's perspective.`,
			},
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					// Ensure remote is set up
					if s.GetRemoteDir() == "" {
						return s.SetupRemote()
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Add a remote",
				Content: `Add the simulated remote as "origin". In real projects this would be a GitHub URL
like: https://github.com/username/book-of-automobiles.git

For this lesson, your simulated remote lives at:
  <sandbox-dir>/book-of-automobiles-remote.git

The exact path is your sandbox directory (shown below) + "/book-of-automobiles-remote.git".
Run the command substituting that full path.`,
				Command:  `git remote add origin /tmp/simgit-sandbox-.../book-of-automobiles-remote.git`,
				Expected: `(no output — success is silent)`,
				Hint:     "The remote path is: <your-sandbox-dir>/book-of-automobiles-remote.git",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.HasRemote("origin") {
						return false, fmt.Sprintf("Remote 'origin' not found. Run:\n  git remote add origin %s", s.GetRemoteDir())
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Viewing remotes",
				Content: `Run git remote -v to see your configured remotes. The -v flag shows the
fetch and push URLs:`,
				Expected: `origin  /tmp/simgit-sandbox-abc123/book-of-automobiles-remote.git (fetch)
origin  /tmp/simgit-sandbox-abc123/book-of-automobiles-remote.git (push)`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Push to origin",
				Content: `Send your local commits to the remote repository. The -u flag sets "upstream"
tracking — after this, plain git push and git pull work without specifying
the remote and branch name.`,
				Command:  `git push -u origin main`,
				Expected: `Enumerating objects: 6, done.
Counting objects: 100% (6/6), done.
Writing objects: 100% (6/6), 1.23 KiB | 1.23 MiB/s, done.
Total 6 (delta 0), reused 0 (delta 0), pack-reused 0
To /tmp/simgit-sandbox-.../book-of-automobiles-remote.git
 * [new branch]      main -> main
Branch 'main' set up to track remote branch 'main' from 'origin'.`,
				Hint: "Run: git push -u origin main",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					// Verify by checking if the remote has commits
					remoteDir := s.GetRemoteDir()
					if remoteDir == "" {
						return false, "Remote not configured in sandbox."
					}
					out, err := s.ExecIn(remoteDir, "log", "--oneline")
					if err != nil || out == "" {
						return false, "Remote appears empty. Run: git push -u origin main"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "Cloning a repository",
				Content: `git clone copies an entire repository from a remote to your machine.
It creates a new directory, initializes a git repo inside it, and
downloads all commits, branches, and files.

  git clone https://github.com/username/book-of-automobiles.git
  cd book-of-automobiles

  Or clone to a specific directory name:
  git clone https://github.com/username/repo.git my-folder

  After cloning:
  • The remote is automatically named "origin"
  • The default branch is checked out
  • You have a full copy of the history

  Clone is typically your starting point when joining an existing project.`,
			},
			{
				Kind: lesson.KindChallenge,
				Title: "Clone the repository",
				Content: `Clone the simulated remote into a new directory. This creates a fresh
copy of the book as if you were a new collaborator joining the project.

Run this command in the PARENT directory of your sandbox (one level up).`,
				Command:  `git clone <remote-path> book-of-automobiles-clone`,
				Expected: `Cloning into 'book-of-automobiles-clone'...
done.`,
				Hint: "cd to the parent directory first, then: git clone <remote-path> book-of-automobiles-clone",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					cloneDir := filepath.Join(s.Dir, "book-of-automobiles-clone")
					if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
						return false, fmt.Sprintf("Clone directory not found at %s. Make sure to cd to %s first, then clone.", cloneDir, s.Dir)
					}
					return true, ""
				},
			},
			{
				Kind: lesson.KindSetup,
				Setup: func(s *sandbox.Sandbox) error {
					// Add a new commit to the remote so fetch/pull has something to get
					cloneDir := filepath.Join(s.Dir, "book-of-automobiles-clone")
					if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
						// Create the clone automatically if user skipped
						s.ExecIn(s.Dir, "clone", s.GetRemoteDir(), "book-of-automobiles-clone") //nolint
					}
					// Push a new commit from the clone to simulate a teammate's work
					mercedesInClone := filepath.Join(cloneDir, "mercedes.md")
					if _, err := os.Stat(mercedesInClone); os.IsNotExist(err) {
						mercedesContent := `# Mercedes-Benz

Mercedes-Benz AG was founded in 1926 in Stuttgart, Germany, through the
merger of Benz & Cie. and Daimler-Motoren-Gesellschaft. The brand is
known for luxury vehicles, buses, coaches, and trucks.

Mercedes-Benz is one of the oldest automobile manufacturers in the world
and is credited with producing the first true automobile — the Benz
Patent-Motorwagen of 1886. Today it is a global leader in the premium
and luxury vehicle segment.
`
						os.WriteFile(mercedesInClone, []byte(mercedesContent), 0644)
						s.ExecIn(cloneDir, "config", "user.name", "Collaborator")    //nolint
						s.ExecIn(cloneDir, "config", "user.email", "co@example.com") //nolint
						s.ExecIn(cloneDir, "add", ".")                               //nolint
						s.ExecIn(cloneDir, "commit", "-m", "Add Mercedes chapter")  //nolint
						s.ExecIn(cloneDir, "push", "origin", "main")                 //nolint
					}
					return nil
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "fetch vs pull",
				Content: `Two commands download changes from a remote — and they're not the same:

  git fetch origin
    Downloads new commits from remote into your local repo.
    Does NOT merge them into your working branch.
    Safe to run anytime — it won't change your files.
    After fetching, you can review with: git log origin/main

  git pull
    = git fetch + git merge (in one step)
    Downloads and immediately merges the remote changes.
    Fastest way to update your local branch.

  When to use each:
  • git fetch — when you want to see what changed before merging
  • git pull  — when you trust the remote and want to update quickly
  • git pull --rebase — like pull, but keeps history linear (no merge commit)

  A collaborator has pushed a new commit (Mercedes chapter) to origin.
  Let's pull it down.`,
			},
			{
				Kind:  lesson.KindChallenge,
				Title: "Pull the latest changes",
				Content: `A collaborator has added the Mercedes chapter to the remote repository.
Pull their changes into your local copy.`,
				Command:  `git pull`,
				Expected: `remote: Enumerating objects: 4, done.
Updating b8e4d23..e2f1a34
Fast-forward
 mercedes.md | 18 ++++++++++++++++++
 1 file changed, 18 insertions(+)
 create mode 100644 mercedes.md`,
				Hint: "Run: git pull — make sure you're in the book-of-automobiles directory",
				Verify: func(s *sandbox.Sandbox) (bool, string) {
					if !s.FileExists("mercedes.md") {
						return false, "mercedes.md not found. Run: git pull in the book-of-automobiles directory"
					}
					return true, ""
				},
			},
			{
				Kind:  lesson.KindExplain,
				Title: "The full remote workflow",
				Content: `Here's the complete day-to-day workflow with remotes:

  Starting a new project:
    git init
    git add .
    git commit -m "Initial commit"
    git remote add origin <url>
    git push -u origin main

  Joining an existing project:
    git clone <url>
    cd project-name

  Daily workflow:
    git pull                    ← get latest changes from team
    # ... make your changes ...
    git add .
    git commit -m "Your message"
    git push                    ← send your changes to remote

  Branch-based workflow (recommended for teams):
    git checkout -b my-feature  ← create feature branch
    # ... make changes ...
    git push -u origin my-feature ← push branch to remote
    # Open a Pull Request on GitHub/GitLab
    # After review, merge into main
    git checkout main && git pull ← sync after merge

  Useful remote commands:
    git remote -v               ← list remotes with URLs
    git remote add upstream <url> ← add another remote (e.g., original fork)
    git fetch --all             ← fetch from all remotes
    git branch -r               ← list remote-tracking branches
    git push origin --delete my-branch ← delete a remote branch`,
			},
		},
	}
}
