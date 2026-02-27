package ui

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Spinner provides a simple progress indicator
type Spinner struct {
	frames   []string
	interval time.Duration
	message  string
	writer   io.Writer
	done     chan struct{}
	running  bool
}

var defaultFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:   defaultFrames,
		interval: 80 * time.Millisecond,
		message:  message,
		writer:   os.Stderr,
		done:     make(chan struct{}),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if s.running {
		return
	}
	s.running = true
	s.done = make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				frame := Cyan.Sprint(s.frames[i%len(s.frames)])
				fmt.Fprintf(s.writer, "\r%s %s", frame, s.message)
				i++
				time.Sleep(s.interval)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	if !s.running {
		return
	}
	s.running = false
	close(s.done)
	// Clear the line
	fmt.Fprintf(s.writer, "\r\033[K")
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Fprintf(s.writer, "%s %s\n", Success.Sprint(SymbolSuccess), message)
}

// Fail stops the spinner and shows an error message
func (s *Spinner) Fail(message string) {
	s.Stop()
	fmt.Fprintf(s.writer, "%s %s\n", Error.Sprint(SymbolError), message)
}

// Update updates the spinner message
func (s *Spinner) Update(message string) {
	s.message = message
}
