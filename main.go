package main

import (
	"log"
	"os"
	"time"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

const glitchChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;':\",./<>?`~ "

// Using NewRGBColor for explicit color definitions
var glitchColors = []tcell.Color{
	tcell.NewRGBColor(0, 0, 0),       // Black
	tcell.NewRGBColor(255, 0, 0),     // Red
	tcell.NewRGBColor(0, 255, 0),     // Green
	tcell.NewRGBColor(255, 255, 0),   // Yellow
	tcell.NewRGBColor(0, 0, 255),     // Blue
	tcell.NewRGBColor(255, 0, 255),   // Magenta
	tcell.NewRGBColor(0, 255, 255),   // Cyan
	tcell.NewRGBColor(255, 255, 255), // White
}

// drawGlitch applies random character corruption to the screen
func drawGlitch(s tcell.Screen, width, height int) {
	numGlitch := rand.Intn(100) + 50 // Random number of glitches per frame

	for i := 0; i < numGlitch; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)

		// Get a random character
		r := rune(glitchChars[rand.Intn(len(glitchChars))])

		// Get a random color
		style := tcell.StyleDefault.Foreground(glitchColors[rand.Intn(len(glitchColors))])

		s.SetContent(x, y, r, nil, style)
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
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

	// Create a channel for events and a goroutine to listen for them
	eventChan := make(chan tcell.Event)
	go func() {
		for {
			eventChan <- s.PollEvent()
		}
	}()

	// Create a ticker for animation updates (e.g., 30 FPS)
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	// Main event loop
	for {
		select {
		case ev := <-eventChan: // Listen on our custom event channel
			switch ev := ev.(type) {
			case *tcell.EventResize:
				width, height = s.Size() // Update dimensions on resize
				s.Clear()                // Clear screen on resize to avoid artifacts
				s.Sync()                 // Sync screen after resize
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					return // Exit the application
				}
			}
		case <-ticker.C: // Handle animation tick
			drawGlitch(s, width, height) // Call the glitch drawing function
			s.Show()                     // Render the screen
		}
	}
}