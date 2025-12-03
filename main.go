package main

import (
	"log"
	"os"
	"time"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

const glitchChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;':\",./<>?`~ "

var glitchColors = []tcell.Color{
	tcell.ColorRed,
	tcell.ColorGreen,
	tcell.ColorBlue,
	tcell.ColorMagenta,
	tcell.ColorYellow,
	tcell.ColorCyan,
	tcell.ColorWhite,
	tcell.ColorLightGreen,
	tcell.ColorLightCyan,
	tcell.ColorLightSkyBlue,
	tcell.ColorLightGoldenrodYellow,
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
	// No longer suppressing unused warnings, as width/height are now used

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
