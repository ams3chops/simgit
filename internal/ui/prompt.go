package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func WaitForEnter(msg string) {
	if msg == "" {
		msg = "Press Enter to continue..."
	}
	fmt.Printf("\n  %s", c(Dim, msg))
	_, _ = reader.ReadString('\n')
}

func AskYesNo(question string) bool {
	for {
		fmt.Printf("  %s %s ", c(BrightWhite, question), c(Dim, "[y/n]:"))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}

func Menu(prompt string, choices []string) int {
	fmt.Printf("\n  %s\n", c(BrightWhite+Bold, prompt))
	for i, ch := range choices {
		fmt.Printf("    %s  %s\n", c(BrightCyan, fmt.Sprintf("[%d]", i+1)), ch)
	}
	for {
		fmt.Printf("\n  %s ", c(Dim, "Enter choice:"))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		for i := range choices {
			if input == fmt.Sprintf("%d", i+1) {
				return i
			}
		}
		fmt.Println(c(BrightRed, "  Invalid choice. Try again."))
	}
}

func ReadLine(prompt string) string {
	fmt.Printf("  %s ", c(BrightWhite, prompt))
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PromptRetrySkipQuit() string {
	fmt.Println()
	fmt.Printf("  %s  %s  %s\n",
		c(BrightYellow, "[R] Retry"),
		c(Dim, "[S] Skip this step"),
		c(BrightRed, "[Q] Quit"),
	)
	fmt.Printf("  %s ", c(Dim, "Choice:"))
	input, _ := reader.ReadString('\n')
	return strings.ToUpper(strings.TrimSpace(input))
}
