package tui

import (
	"math/rand"
	"sync"
	"time"

	"glitch-saver/internal/effects"
	"glitch-saver/internal/options"

	"github.com/gdamore/tcell/v2"
)

func RunTUI(opts *options.GlitchOptions) (tcell.Screen, error) {
	// Create a local random number generator
	rGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize tcell screen
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err = s.Init(); err != nil {
		return nil, err
	}
	// We'll call s.Fini() in main.go after TUI returns,
	// so we don't need a defer here that would call it prematurely

	// Set default style and clear screen
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	// Hide cursor
	s.HideCursor()

	// Get initial screen dimensions
	width, height := s.Size()

	effects.InitializeEffects(width, height)

	// Use mutex to protect width and height variables that are accessed by both the main loop and event handler
	var mu sync.Mutex

	// Create a channel for events and a goroutine to listen for them
	eventChan := make(chan tcell.Event)
	done := make(chan bool, 1) // Channel to signal when to stop the polling goroutine
	go func() {
		for {
			select {
			case eventChan <- s.PollEvent():
			case <-done:
				return // Exit the goroutine when done is signaled
			}
		}
	}()

	// Create a ticker for animation updates based on fps flag
	ticker := time.NewTicker(time.Second / time.Duration(opts.FPS))
	defer func() {
		ticker.Stop()
		done <- true // Signal the goroutine to stop
	}()

	// Main event loop
	for {
		select {
		case ev := <-eventChan: // Listen on our custom event channel
			switch ev := ev.(type) {
			case *tcell.EventResize:
				mu.Lock()
				width, height = s.Size() // Update dimensions on resize
				mu.Unlock()
				effects.InitializeEffects(width, height)
				s.Clear() // Clear screen on resize to avoid artifacts
				s.Sync()  // Sync screen after resize
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					return s, nil // Exit the application, returning the screen
				}
			}
		case <-ticker.C: // Handle animation tick
			mu.Lock()
			currentWidth := width
			currentHeight := height
			mu.Unlock()
			effects.DrawGlitch(s, currentWidth, currentHeight, rGen, opts) // Pass opts struct
			s.Show()
		}
	}
}
