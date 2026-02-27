package ui

import (
	"github.com/fatih/color"
)

var (
	// Primary colors
	Cyan    = color.New(color.FgCyan)
	Green   = color.New(color.FgGreen)
	Yellow  = color.New(color.FgYellow)
	Red     = color.New(color.FgRed)
	Blue    = color.New(color.FgBlue)
	Magenta = color.New(color.FgMagenta)
	White   = color.New(color.FgWhite)

	// Styled colors
	Bold      = color.New(color.Bold)
	Dim       = color.New(color.Faint)
	Italic    = color.New(color.Italic)
	Underline = color.New(color.Underline)

	// Combined styles
	CyanBold    = color.New(color.FgCyan, color.Bold)
	GreenBold   = color.New(color.FgGreen, color.Bold)
	YellowBold  = color.New(color.FgYellow, color.Bold)
	RedBold     = color.New(color.FgRed, color.Bold)
	BlueBold    = color.New(color.FgBlue, color.Bold)
	MagentaBold = color.New(color.FgMagenta, color.Bold)
	WhiteBold   = color.New(color.FgWhite, color.Bold)

	// Semantic colors
	Success = Green
	Error   = Red
	Warning = Yellow
	Info    = Cyan
	Hint    = Dim
	Version = GreenBold
	Command = CyanBold
	Path    = Blue
)

// Symbols for output
const (
	SymbolSuccess   = "✓"
	SymbolError     = "✗"
	SymbolWarning   = "⚠"
	SymbolInfo      = "ℹ"
	SymbolArrow     = "→"
	SymbolBullet    = "•"
	SymbolCheck     = "✓"
	SymbolCross     = "✗"
	SymbolStar      = "★"
	SymbolDot       = "·"
	SymbolPipe      = "│"
	SymbolCorner    = "└"
	SymbolTee       = "├"
	SymbolHorizontal = "─"
)

// DisableColors disables all color output
func DisableColors() {
	color.NoColor = true
}

// EnableColors enables color output
func EnableColors() {
	color.NoColor = false
}
