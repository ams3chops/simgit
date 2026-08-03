package main

import (
	"flag"
	"fmt"
	"os"
	"simgit/internal/lesson"
	"simgit/internal/lessons"
	"simgit/internal/progress"
	"simgit/internal/sandbox"
	"simgit/internal/ui"
	"strconv"
)

const version = "1.0.0"

func main() {
	var (
		resetFlag     = flag.Bool("reset", false, "Reset all progress and start over")
		lessonFlag    = flag.Int("lesson", 0, "Jump to a specific lesson number (1-10)")
		listFlag      = flag.Bool("list", false, "List all lessons and exit")
		keepSandbox   = flag.Bool("keep-sandbox", false, "Don't delete sandbox directory on exit")
		versionFlag   = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()

	if *versionFlag {
		fmt.Printf("SimGit version %s\n", version)
		os.Exit(0)
	}

	allLessons := lessons.All()

	if *listFlag {
		fmt.Println()
		fmt.Println(ui.InfoText("SimGit Lessons:"))
		fmt.Println()
		for _, l := range allLessons {
			fmt.Printf("  %s  %s\n",
				ui.InfoText(fmt.Sprintf("[%2d]", l.Number)),
				l.Title,
			)
			if l.Description != "" {
				// Print first line of description
				first := firstLine(l.Description)
				fmt.Printf("       %s\n", dimText(first))
			}
		}
		fmt.Println()
		os.Exit(0)
	}

	// Load progress
	prog, err := progress.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading progress: %v\n", err)
		os.Exit(1)
	}

	if *resetFlag {
		if err := prog.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, "Error resetting progress: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(ui.InfoText("Progress reset. Starting from the beginning."))
	}

	if *lessonFlag > 0 {
		if *lessonFlag > len(allLessons) {
			fmt.Fprintf(os.Stderr, "Lesson %d doesn't exist. Valid range: 1-%d\n", *lessonFlag, len(allLessons))
			os.Exit(1)
		}
		prog.LessonIndex = *lessonFlag - 1
		prog.StepIndex = 0
	}

	// Set up sandbox
	var sb *sandbox.Sandbox
	if prog.SandboxDir != "" && !*resetFlag {
		sb, err = sandbox.Resume(prog.SandboxDir)
		if err != nil {
			// Sandbox no longer exists, create a new one
			sb, err = sandbox.New()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating sandbox: %v\n", err)
				os.Exit(1)
			}
			prog.SandboxDir = sb.Dir
		}
	} else {
		sb, err = sandbox.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating sandbox: %v\n", err)
			os.Exit(1)
		}
		prog.SandboxDir = sb.Dir
		if err := prog.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving progress: %v\n", err)
		}
	}

	if !*keepSandbox {
		defer func() {
			if err := sb.Destroy(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not clean up sandbox: %v\n", err)
			}
		}()
	} else {
		fmt.Printf("\n  %s %s\n\n", ui.InfoText("Sandbox will be kept at:"), sb.Dir)
	}

	// Print welcome banner
	printWelcome(prog, allLessons)

	// Create runner
	runner := lesson.NewRunner(
		allLessons,
		sb,
		prog.LessonIndex,
		prog.StepIndex,
		func(lessonIdx, stepIdx int) {
			prog.LessonIndex = lessonIdx
			prog.StepIndex = stepIdx
			prog.Save() //nolint
		},
	)

	if err := runner.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		os.Exit(1)
	}

	prog.Completed = true
	prog.Save() //nolint
}

func printWelcome(prog *progress.Progress, allLessons []lesson.Lesson) {
	fmt.Println()
	fmt.Println(ui.InfoText("╔══════════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.InfoText("║                                                              ║"))
	fmt.Println(ui.InfoText("║   SimGit — Learn Git with the Book of Automobiles            ║"))
	fmt.Println(ui.InfoText("║                                                              ║"))
	fmt.Println(ui.InfoText("╚══════════════════════════════════════════════════════════════╝"))
	fmt.Println()

	if prog.LessonIndex > 0 || prog.StepIndex > 0 {
		total := len(allLessons)
		fmt.Printf("  %s Resuming at Lesson %s/%s\n",
			ui.InfoText("→"),
			ui.InfoText(strconv.Itoa(prog.LessonIndex+1)),
			ui.InfoText(strconv.Itoa(total)),
		)
		fmt.Println()
	} else {
		fmt.Println("  Welcome to SimGit!")
		fmt.Println()
		fmt.Println("  You will write a Book of Automobiles while learning Git.")
		fmt.Println("  Each lesson teaches a core Git concept using real commands.")
		fmt.Println()
		fmt.Println("  When challenged to run a command:")
		fmt.Println("    1. Note the sandbox directory shown on screen")
		fmt.Println("    2. Open a terminal and cd to that directory")
		fmt.Println("    3. Run the command")
		fmt.Println("    4. Press Enter back here to verify")
		fmt.Println()
		fmt.Println("  Flags:")
		fmt.Println("    --reset         Start over from lesson 1")
		fmt.Println("    --lesson N      Jump to lesson N")
		fmt.Println("    --list          Show all lessons")
		fmt.Println("    --keep-sandbox  Keep the sandbox directory after exit")
		fmt.Println()
	}
	ui.WaitForEnter("Press Enter to begin...")
}

func firstLine(s string) string {
	for i, ch := range s {
		if ch == '\n' {
			return s[:i]
		}
	}
	return s
}

func dimText(s string) string {
	return "\033[2m" + s + "\033[0m"
}
