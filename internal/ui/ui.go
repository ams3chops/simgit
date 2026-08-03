package ui

import (
	"fmt"
	"os"
	"strings"
)

// ANSI escape codes
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Italic  = "\033[3m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue   = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"

	BgBlue  = "\033[44m"
	BgGreen = "\033[42m"
)

var colorEnabled = os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

func c(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + Reset
}

func Clear() {
	fmt.Print("\033[2J\033[H")
}

func Header(lessonNum, totalLessons int, title string) {
	width := 66
	line := strings.Repeat("═", width)
	padding := width - len([]rune(title)) - 4
	if padding < 0 {
		padding = 0
	}
	fmt.Println()
	fmt.Println(c(BrightCyan+Bold, "╔"+line+"╗"))
	fmt.Printf(c(BrightCyan+Bold, "║")+" %s %s"+c(BrightCyan+Bold, "║")+"\n",
		c(BrightWhite+Bold, title),
		strings.Repeat(" ", padding),
	)
	progress := fmt.Sprintf("Lesson %d / %d", lessonNum, totalLessons)
	progressPad := width - len(progress) - 2
	if progressPad < 0 {
		progressPad = 0
	}
	fmt.Printf(c(BrightCyan+Bold, "║")+" %s%s"+c(BrightCyan+Bold, "║")+"\n",
		c(Dim, progress),
		strings.Repeat(" ", progressPad),
	)
	fmt.Println(c(BrightCyan+Bold, "╚"+line+"╝"))
	fmt.Println()
}

func StepHeader(stepNum, totalSteps int, title string) {
	bar := fmt.Sprintf("[Step %d/%d]", stepNum, totalSteps)
	fmt.Printf("%s %s\n", c(Yellow+Bold, bar), c(BrightWhite+Bold, title))
	fmt.Println(c(Dim, strings.Repeat("─", 60)))
	fmt.Println()
}

func Section(title string) {
	fmt.Println()
	fmt.Printf("%s\n", c(BrightYellow+Bold, title))
	fmt.Println(c(Yellow, strings.Repeat("─", len(title)+4)))
}

func Body(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			fmt.Println()
			continue
		}
		// Detect indented lines (preserve their indentation, don't re-wrap)
		trimmed := strings.TrimLeft(line, " \t")
		indent := line[:len(line)-len(trimmed)]
		if len(indent) > 0 {
			// Preserve indented lines as-is (code, lists, ASCII art)
			fmt.Println("  " + line)
		} else {
			fmt.Println(WordWrap("  "+line, 80))
		}
	}
}

func Command(cmd string) {
	width := 50
	inner := cmd
	if len(inner) < width-4 {
		inner = inner + strings.Repeat(" ", width-4-len(inner))
	}
	fmt.Println()
	fmt.Println(c(BrightGreen, "  ┌"+strings.Repeat("─", width)+"┐"))
	fmt.Printf(c(BrightGreen, "  │")+"  %s  "+c(BrightGreen, "│")+"\n", c(BrightWhite+Bold, cmd))
	fmt.Println(c(BrightGreen, "  └"+strings.Repeat("─", width)+"┘"))
	fmt.Println()
	_ = inner
}

func CommandWithPrompt(cmd string) {
	width := 50
	inner := "$ " + cmd
	padding := ""
	if len(inner) < width-4 {
		padding = strings.Repeat(" ", width-4-len(inner))
	}
	fmt.Println()
	fmt.Println(c(BrightGreen, "  ┌"+strings.Repeat("─", width)+"┐"))
	fmt.Printf(c(BrightGreen, "  │")+"  %s%s  "+c(BrightGreen, "│")+"\n",
		c(BrightWhite+Bold, "$ ")+c(BrightCyan+Bold, cmd),
		padding,
	)
	fmt.Println(c(BrightGreen, "  └"+strings.Repeat("─", width)+"┘"))
	fmt.Println()
}

func ExpectedOutput(out string) {
	if out == "" {
		fmt.Println(c(Dim, "  (no output — silence is success)"))
		fmt.Println()
		return
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	fmt.Println(c(Dim, "  Expected output:"))
	fmt.Println(c(Dim, "  ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄"))
	for _, line := range lines {
		fmt.Println(c(Dim, "  "+line))
	}
	fmt.Println(c(Dim, "  ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄"))
	fmt.Println()
}

func Success(msg string) {
	fmt.Println(c(BrightGreen+Bold, "  ✓ "+msg))
}

func Error(msg string) {
	fmt.Println(c(BrightRed+Bold, "  ✗ "+msg))
}

func Info(msg string) {
	fmt.Println(c(BrightCyan, "  ℹ "+msg))
}

func Warn(msg string) {
	fmt.Println(c(BrightYellow, "  ⚠ "+msg))
}

func Hint(msg string) {
	fmt.Println(c(Dim+Italic, "  Hint: "+msg))
}

func Separator() {
	fmt.Println(c(Dim, strings.Repeat("─", 68)))
}

func SandboxPath(path string) {
	fmt.Println()
	fmt.Printf("  %s %s\n", c(Yellow+Bold, "Sandbox directory:"), c(BrightWhite+Bold, path))
	fmt.Println()
}

func Challenge() {
	fmt.Println(c(Magenta+Bold, "  YOUR CHALLENGE:"))
}

func LessonComplete(title string) {
	inner := "  ✓ " + title + " complete!"
	width := len([]rune(inner)) + 4
	if width < 45 {
		width = 45
	}
	border := strings.Repeat("═", width)
	padding := strings.Repeat(" ", width-len([]rune(inner)))
	fmt.Println()
	fmt.Println(c(BrightGreen+Bold, "  ╔"+border+"╗"))
	fmt.Println(c(BrightGreen+Bold, "  ║  "+inner+padding+"  ║"))
	fmt.Println(c(BrightGreen+Bold, "  ╚"+border+"╝"))
	fmt.Println()
}

func Congratulations() {
	fmt.Println()
	fmt.Println(c(BrightYellow+Bold, "╔══════════════════════════════════════════════════════════════════╗"))
	fmt.Println(c(BrightYellow+Bold, "║                                                                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║   🎉  Congratulations! You've completed SimGit!                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║                                                                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║   You now know the core of git:                                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Configure git identity                                     ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Initialize repositories                                    ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Stage and commit changes                                   ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ View history with git log                                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Create and switch branches                                 ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Merge branches together                                    ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Navigate history with checkout                             ║"))
	fmt.Println(c(BrightYellow+Bold, "║     ✓ Work with remote repositories                              ║"))
	fmt.Println(c(BrightYellow+Bold, "║                                                                  ║"))
	fmt.Println(c(BrightYellow+Bold, "║   The Book of Automobiles is complete!                           ║"))
	fmt.Println(c(BrightYellow+Bold, "╚══════════════════════════════════════════════════════════════════╝"))
	fmt.Println()
}

// InfoText returns styled info text without printing (for inline use).
func InfoText(msg string) string {
	return c(BrightCyan, msg)
}

func WordWrap(s string, width int) string {
	if len(s) <= width {
		return s
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	// detect leading whitespace
	indent := ""
	for _, ch := range s {
		if ch == ' ' || ch == '\t' {
			indent += string(ch)
		} else {
			break
		}
	}

	var sb strings.Builder
	line := ""
	for i, word := range words {
		if i == 0 {
			line = indent + word
		} else if len(line)+1+len(word) > width {
			sb.WriteString(line + "\n")
			line = indent + word
		} else {
			line += " " + word
		}
	}
	sb.WriteString(line)
	return sb.String()
}
