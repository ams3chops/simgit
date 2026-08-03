package lesson

import (
	"fmt"
	"os"
	"simgit/internal/sandbox"
	"simgit/internal/ui"
)

// Runner drives all lessons interactively.
type Runner struct {
	lessons      []Lesson
	sandbox      *sandbox.Sandbox
	startLesson  int
	startStep    int
	onProgress   func(lessonIdx, stepIdx int)
}

// NewRunner creates a Runner starting at the given lesson and step indices (0-based).
func NewRunner(lessons []Lesson, sb *sandbox.Sandbox, startLesson, startStep int, onProgress func(int, int)) *Runner {
	return &Runner{
		lessons:     lessons,
		sandbox:     sb,
		startLesson: startLesson,
		startStep:   startStep,
		onProgress:  onProgress,
	}
}

// Run drives all lessons from the start position.
func (r *Runner) Run() error {
	total := len(r.lessons)
	for i := r.startLesson; i < total; i++ {
		l := r.lessons[i]
		start := 0
		if i == r.startLesson {
			start = r.startStep
		}
		if err := r.runLesson(l, total, start); err != nil {
			return err
		}
		// Advance past this lesson
		r.startLesson = i + 1
		r.startStep = 0
	}
	ui.Congratulations()
	return nil
}

func (r *Runner) runLesson(l Lesson, totalLessons, startStep int) error {
	ui.Header(l.Number, totalLessons, l.Title)
	if l.Description != "" {
		ui.Body(l.Description)
		fmt.Println()
	}

	totalSteps := countDisplaySteps(l.Steps)
	displayIdx := 0

	for i := startStep; i < len(l.Steps); i++ {
		step := l.Steps[i]

		if step.Kind == KindSetup {
			if err := r.runSetupStep(step); err != nil {
				ui.Warn(fmt.Sprintf("Setup step failed: %v", err))
			}
			continue
		}

		displayIdx++
		if step.Title != "" {
			ui.StepHeader(displayIdx, totalSteps, step.Title)
		}

		var err error
		switch step.Kind {
		case KindExplain:
			err = r.runExplainStep(step)
		case KindChallenge:
			err = r.runChallengeStep(step)
		}
		if err != nil {
			return err
		}

		if r.onProgress != nil {
			r.onProgress(l.Number-1, i+1)
		}
	}

	ui.LessonComplete(l.Title)
	ui.WaitForEnter("Press Enter to continue to the next lesson...")
	return nil
}

func countDisplaySteps(steps []Step) int {
	n := 0
	for _, s := range steps {
		if s.Kind != KindSetup {
			n++
		}
	}
	return n
}

func (r *Runner) runExplainStep(s Step) error {
	if s.Content != "" {
		ui.Body(s.Content)
		fmt.Println()
	}
	if s.Expected != "" {
		ui.ExpectedOutput(s.Expected)
	}
	ui.WaitForEnter("")
	return nil
}

func (r *Runner) runChallengeStep(s Step) error {
	if s.Content != "" {
		ui.Body(s.Content)
		fmt.Println()
	}

	ui.Challenge()
	ui.CommandWithPrompt(s.Command)

	if s.Expected != "" {
		ui.ExpectedOutput(s.Expected)
	}

	ui.SandboxPath(r.sandbox.RepoDir)
	if r.sandbox.GetRemoteDir() != "" {
		fmt.Printf("  %s %s\n", ui.InfoText("Remote (origin) path:"), r.sandbox.GetRemoteDir())
		fmt.Println()
	}
	fmt.Printf("  %s\n", ui.InfoText("Open a terminal, cd to the directory above, and run the command."))

	if s.Verify == nil {
		// Honor-system: just wait for Enter
		ui.WaitForEnter("When done, press Enter to continue...")
		return nil
	}

	for {
		ui.WaitForEnter("When done, press Enter here to verify...")

		ok, feedback := s.Verify(r.sandbox)
		if ok {
			ui.Success("Verified! Great work.")
			fmt.Println()
			return nil
		}

		ui.Error("Not quite. " + feedback)
		if s.Hint != "" {
			ui.Hint(s.Hint)
		}

		choice := ui.PromptRetrySkipQuit()
		switch choice {
		case "S":
			ui.Warn("Skipping this step.")
			return nil
		case "Q":
			fmt.Println(ui.InfoText("Quitting. Your progress has been saved."))
			os.Exit(0)
		default:
			// retry
		}
	}
}

func (r *Runner) runSetupStep(s Step) error {
	if s.Setup != nil {
		return s.Setup(r.sandbox)
	}
	return nil
}
