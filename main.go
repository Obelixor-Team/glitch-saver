package main

import (
	"flag" // Added
	"log"
	"math/rand"
	"os"
	"time"

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
func shiftLineGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) { // rGen added
	if height == 0 {
		return
	}
	y := rGen.Intn(height) // Use rGen
	offset := rGen.Intn(width/2) - (width / 4) // Use rGen

	line := make([]struct {
		r     rune
		style tcell.Style
	}, width)

	for x := 0; x < width; x++ {
		r, style, _ := s.Get(x, y) // Changed from GetContent
		line[x].r = r
		line[x].style = style
	}

	for x := 0; x < width; x++ {
		newX := x + offset
		if newX >= 0 && newX < width {
			s.SetContent(newX, y, line[x].r, nil, line[x].style)
		}
	}
}

// blockDistortionGlitch copies a random block of the screen to another random location
func blockDistortionGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) { // rGen added
	if width == 0 || height == 0 {
		return
	}
	srcX, srcY := rGen.Intn(width), rGen.Intn(height) // Use rGen
	blockW := rGen.Intn(width/2) + 1                   // Use rGen
	blockH := rGen.Intn(height/2) + 1                  // Use rGen

	destX, destY := rGen.Intn(width), rGen.Intn(height) // Use rGen

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
				r, style, _ := s.Get(srcX+x, srcY+y) // Changed from GetContent
				block[y][x].r = r
				block[y][x].style = style
			}
		}
	}

	for y := 0; y < blockH; y++ {
		for x := 0; x < blockW; x++ {
			if destX+x < width && destY+y < height {
				s.SetContent(destX+x, destY+y, block[y][x].r, nil, block[y][x].style)
			}
		}
	}
}

// drawGlitch applies random character corruption and other effects to the screen
func drawGlitch(s tcell.Screen, width, height, intensity int, rGen *rand.Rand) { // rGen added
	numGlitch := rGen.Intn(100*intensity) + (50 * intensity) // Use rGen
	for i := 0; i < numGlitch; i++ {
		x := rGen.Intn(width) // Use rGen
		y := rGen.Intn(height) // Use rGen

		r := rune(glitchChars[rGen.Intn(len(glitchChars))]) // Use rGen
		style := tcell.StyleDefault.Foreground(glitchColors[rGen.Intn(len(glitchColors))]) // Use rGen

		s.SetContent(x, y, r, nil, style)
	}

	if rGen.Intn(10) < 2 { // Use rGen
		shiftLineGlitch(s, width, height, rGen) // Pass rGen
	}

	if rGen.Intn(10) < 1 { // Use rGen
		blockDistortionGlitch(s, width, height, rGen) // Pass rGen
	}
}

func main() {
	// Define command-line flags
	fps := flag.Int("fps", 30, "frames per second for the animation")
	intensity := flag.Int("intensity", 5, "glitch intensity (1-10)")
	flag.Parse()

	// Clamp intensity
	if *intensity < 1 {
		*intensity = 1
	}
	if *intensity > 10 {
		*intensity = 10
	}

	// Create a local random number generator
	rGen := rand.New(rand.NewSource(time.Now().UnixNano())) // Changed: local generator

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
			drawGlitch(s, width, height, *intensity, rGen) // Pass rGen to drawGlitch
			s.Show()                                 // Render the screen
		}
	}
}