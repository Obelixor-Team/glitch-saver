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

	// Main event loop
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			width, height = s.Size() // Update dimensions on resize
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
				return
			}
		}
		// Add a small delay to prevent busy-waiting for events, though tcell.PollEvent() is blocking
		time.Sleep(time.Millisecond * 50)
	}
}