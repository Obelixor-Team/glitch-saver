package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialize tcell screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err = s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default style and clear screen
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	// Hide cursor
	s.HideCursor()

	// Event loop for handling input and drawing
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	defer quit() // Ensure screen is finalized on exit

	// Get initial screen dimensions
	width, height := s.Size()
	_ = width // Suppress unused warning for now
	_ = height // Suppress unused warning for now

	// Create a ticker for animation updates (e.g., 30 FPS)
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	// Main event loop
	for {
		select {
		case ev := <-s.Events(): // Non-blocking read from event channel
			switch ev := ev.(type) {
			case *tcell.EventResize:
				width, height = s.Size() // Update dimensions on resize
				s.Sync()                 // Sync screen after resize
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					return // Exit the application
				}
			}
		case <-ticker.C: // Handle animation tick
			// Placeholder for update logic (e.g., update glitch state)
			// Placeholder for drawing logic
			s.Show() // Render the screen
		}
	}
}