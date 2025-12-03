package main

import (
	"flag" // Added
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

// shiftLineGlitch shifts a random line horizontally
func shiftLineGlitch(s tcell.Screen, width, height int) {
	if height == 0 { // Avoid panic on empty screen
		return
	}
	y := rand.Intn(height)
	offset := rand.Intn(width/2) - (width / 4) // Shift left or right

	// Buffer the line
	line := make([]struct {
		r     rune
		style tcell.Style
	}, width)

	for x := 0; x < width; x++ {
		r, _, style, _ := s.GetContent(x, y)
		line[x].r = r
		line[x].style = style
	}

	// Write the line back with an offset
	for x := 0; x < width; x++ {
		newX := x + offset
		if newX >= 0 && newX < width {
			s.SetContent(newX, y, line[x].r, nil, line[x].style)
		}
	}
}

// blockDistortionGlitch copies a random block of the screen to another random location
func blockDistortionGlitch(s tcell.Screen, width, height int) {
	if width == 0 || height == 0 {
		return
	}
	// Define the source block
	srcX, srcY := rand.Intn(width), rand.Intn(height)
	blockW := rand.Intn(width/2) + 1  // Max half screen width
	blockH := rand.Intn(height/2) + 1 // Max half screen height

	// Define the destination
	destX, destY := rand.Intn(width), rand.Intn(height)

	// Buffer the block
	block := make([][]struct {
		r     rune
		style tcell.Style
	}, blockH)

	for y := 0; y < blockH; y++ {
		block[y] = make([]struct {
			r     rune
			style tcell.Style
		}, blockW)
		for x := 0; x < blockW; x++ {
			if srcX+x < width && srcY+y < height {
				r, _, style, _ := s.GetContent(srcX+x, srcY+y)
				block[y][x].r = r
				block[y][x].style = style
			}
		}
	}

	// Write the block to the destination
	for y := 0; y < blockH; y++ {
		for x := 0; x < blockW; x++ {
			if destX+x < width && destY+y < height {
				s.SetContent(destX+x, destY+y, block[y][x].r, nil, block[y][x].style)
			}
		}
	}
}

// drawGlitch applies random character corruption and other effects to the screen
func drawGlitch(s tcell.Screen, width, height, intensity int) { // Added intensity
	// Character corruption
	numGlitch := rand.Intn(100 * intensity) + (50 * intensity) // Use intensity
	for i := 0; i < numGlitch; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)

		// Get a random character
		r := rune(glitchChars[rand.Intn(len(glitchChars))])

		// Get a random color
		style := tcell.StyleDefault.Foreground(glitchColors[rand.Intn(len(glitchColors))])

		s.SetContent(x, y, r, nil, style)
	}

	// Line shifts
	if rand.Intn(10) < 2 { // 20% chance to shift a line
		shiftLineGlitch(s, width, height)
	}

	// Block distortion
	if rand.Intn(10) < 1 { // 10% chance to distort a block
		blockDistortionGlitch(s, width, height)
	}
}

func main() {
	// Define command-line flags
	fps := flag.Int("fps", 30, "frames per second for the animation")
	intensity := flag.Int("intensity", 5, "glitch intensity (1-10)")
	flag.Parse()

	// Clamp intensity
	if *intensity < 1 { *intensity = 1 }
	if *intensity > 10 { *intensity = 10 }

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

	// Create a ticker for animation updates based on fps flag
	ticker := time.NewTicker(time.Second / time.Duration(*fps))
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
			drawGlitch(s, width, height, *intensity) // Pass intensity to drawGlitch
			s.Show()                     // Render the screen
		}
	}
}
