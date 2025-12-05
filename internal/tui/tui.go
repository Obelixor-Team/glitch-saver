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
		s.Fini() // Finalize screen if Init fails
		return nil, err
	}

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
	go func() {
		for {
			eventChan <- s.PollEvent()
		}
	}()

	// Create a ticker for animation updates based on fps flag
	ticker := time.NewTicker(time.Second / time.Duration(opts.FPS))
	defer ticker.Stop()

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
					s.Fini()      // Finalize the screen before returning
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
