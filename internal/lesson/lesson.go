package lesson

import "simgit/internal/sandbox"

// StepKind distinguishes step types.
type StepKind int

const (
	KindExplain   StepKind = iota // display content, press Enter to continue
	KindChallenge                  // user must run a command; app verifies
	KindSetup                      // app auto-runs setup (not shown to user)
)

// VerifyFunc checks if the user completed a challenge.
// Returns (success, feedback message).
type VerifyFunc func(s *sandbox.Sandbox) (bool, string)

// SetupFunc sets up sandbox state before a step.
type SetupFunc func(s *sandbox.Sandbox) error

// Step is a single interactive unit within a lesson.
type Step struct {
	Kind     StepKind
	Title    string
	Content  string     // narrative prose
	Command  string     // git command to show (KindChallenge)
	Expected string     // expected terminal output to display
	Hint     string
	Verify   VerifyFunc // nil for KindExplain/KindSetup
	Setup    SetupFunc  // pre-step setup, run automatically
}

// Lesson groups related steps around a single git concept.
type Lesson struct {
	Number      int
	Title       string
	Subtitle    string
	Description string
	Steps       []Step
}
