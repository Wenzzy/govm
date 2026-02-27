package ui

import (
	"fmt"
	"os"
	"strings"
)

// Print prints a message to stdout
func Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println prints a message to stdout with newline
func Println(args ...interface{}) {
	fmt.Println(args...)
}

// Printf prints a formatted message to stdout
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Success.Sprint(SymbolSuccess), msg)
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s %s\n", Error.Sprint(SymbolError), msg)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Warning.Sprint(SymbolWarning), msg)
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Info.Sprint(SymbolInfo), msg)
}

// PrintHint prints a hint message (dimmed)
func PrintHint(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s %s\n", Hint.Sprint(SymbolArrow), Hint.Sprint(msg))
}

// PrintVersion prints a version with styling
func PrintVersion(ver string, current bool, installed bool) {
	prefix := "  "
	suffix := ""

	if current {
		prefix = Green.Sprint(SymbolArrow) + " "
		ver = GreenBold.Sprint(ver)
		suffix = Dim.Sprint(" (current)")
	} else if installed {
		prefix = "  "
		ver = Green.Sprint(ver)
		suffix = Dim.Sprint(" (installed)")
	} else {
		ver = White.Sprint(ver)
	}

	fmt.Printf("%s%s%s\n", prefix, ver, suffix)
}

// PrintHeader prints a section header
func PrintHeader(title string) {
	fmt.Printf("\n%s\n", Bold.Sprint(title))
	fmt.Println(Dim.Sprint(strings.Repeat(SymbolHorizontal, len(title)+2)))
}

// PrintKeyValue prints a key-value pair
func PrintKeyValue(key, value string) {
	fmt.Printf("  %s: %s\n", Dim.Sprint(key), value)
}

// PrintBullet prints a bullet point
func PrintBullet(text string) {
	fmt.Printf("  %s %s\n", Cyan.Sprint(SymbolBullet), text)
}

// PrintCommand prints a command example
func PrintCommand(cmd string) {
	fmt.Printf("  %s %s\n", Dim.Sprint("$"), Command.Sprint(cmd))
}

// Confirm asks for user confirmation
func Confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// PrintLogo prints the govm logo
func PrintLogo() {
	logo := `
   __ _  _____   ___ __ ___
  / _` + "`" + ` |/ _ \ \ / / '_ ` + "`" + ` _ \
 | (_| | (_) \ V /| | | | | |
  \__, |\___/ \_/ |_| |_| |_|
   __/ |
  |___/   Go Version Manager
`
	fmt.Println(Cyan.Sprint(logo))
}

// PrintVersionInfo prints version information
func PrintVersionInfo(version, buildTime string) {
	fmt.Printf("%s %s\n", Dim.Sprint("Version:"), Version.Sprint(version))
	if buildTime != "unknown" {
		fmt.Printf("%s %s\n", Dim.Sprint("Built:"), buildTime)
	}
}

// ClearLine clears the current line
func ClearLine() {
	fmt.Print("\r\033[K")
}

// MoveCursorUp moves cursor up n lines
func MoveCursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}
