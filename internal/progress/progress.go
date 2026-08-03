package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Progress tracks where the user is in the SimGit curriculum.
type Progress struct {
	LessonIndex int    `json:"lesson_index"`
	StepIndex   int    `json:"step_index"`
	Completed   bool   `json:"completed"`
	SandboxDir  string `json:"sandbox_dir"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".simgit")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func progressPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "progress.json"), nil
}

// Load reads progress from disk. Returns a fresh Progress if file doesn't exist.
func Load() (*Progress, error) {
	path, err := progressPath()
	if err != nil {
		return &Progress{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Progress{}, nil
	}
	if err != nil {
		return &Progress{}, nil
	}
	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return &Progress{}, nil
	}
	return &p, nil
}

// Save writes progress to disk.
func (p *Progress) Save() error {
	path, err := progressPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Advance moves to the next step, or next lesson if done with steps.
func (p *Progress) Advance(stepIdx int) {
	p.StepIndex = stepIdx
}

// AdvanceLesson moves to the start of the next lesson.
func (p *Progress) AdvanceLesson(lessonIdx int) {
	p.LessonIndex = lessonIdx
	p.StepIndex = 0
}

// Reset clears all progress.
func (p *Progress) Reset() error {
	p.LessonIndex = 0
	p.StepIndex = 0
	p.Completed = false
	p.SandboxDir = ""
	return p.Save()
}
