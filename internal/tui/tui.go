package tui

import (
	"log"
	"math/rand"
	"time"

	"glitch-saver/internal/effects"
	"glitch-saver/internal/options"

	"github.com/gdamore/tcell/v2"
)

func RunTUI(opts *options.GlitchOptions) (tcell.Screen, error) {
	log.Println("Starting RunTUI")

	// Create a local random number generator
	rGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize tcell screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Printf("ERROR: tcell.NewScreen() failed: %+v", err)
		return nil, err // Return to allow deferred quit() to run
	}
	if err = s.Init(); err != nil {
		log.Printf("ERROR: s.Init() failed: %+v", err)
		s.Fini() // Finalize screen if Init fails
		return nil, err
	}
	log.Println("Screen initialized successfully")

	// Set default style and clear screen
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	// Hide cursor
	s.HideCursor()

	// Event loop for handling input and drawing
	// quit := func() {
	// 	s.Fini()
	// 	os.Exit(0)
	// }
	// defer quit() // Ensure screen is finalized on exit

	// Get initial screen dimensions
	width, height := s.Size()

	effects.InitializeEffects(width, height)

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
				width, height = s.Size() // Update dimensions on resize
				effects.InitializeEffects(width, height)
				s.Clear() // Clear screen on resize to avoid artifacts
				s.Sync()  // Sync screen after resize
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					return s, nil // Exit the application, returning the screen
				}
			}
		case <-ticker.C: // Handle animation tick
			effects.DrawGlitch(s, width, height, rGen, opts) // Pass opts struct
			s.Show()
		}
	}
	// If the loop exits for some reason, return the screen.
	return s, nil
}
