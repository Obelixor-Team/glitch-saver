package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const glitchChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;':\",./<>?`~ "
const cp437Chars = "ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■ "

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
func shiftLineGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) {
	if height == 0 {
		return
	}
	y := rGen.Intn(height)
	offset := rGen.Intn(width/2) - (width / 4)

	line := make([]struct {
		r     rune
		style tcell.Style
	}, width)

	for x := 0; x < width; x++ {
		rawVal, style, _ := s.Get(x, y)
		var r rune
		if len(rawVal) > 0 {
			r = []rune(rawVal)[0]
		}
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
func blockDistortionGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) {
	if width == 0 || height == 0 {
		return
	}
	srcX, srcY := rGen.Intn(width), rGen.Intn(height)
	blockW := rGen.Intn(width/2) + 1
	blockH := rGen.Intn(height/2) + 1

	destX, destY := rGen.Intn(width), rGen.Intn(height)

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
				rawVal, style, _ := s.Get(srcX+x, srcY+y)
				var r rune
				if len(rawVal) > 0 {
					r = []rune(rawVal)[0]
				}
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
func drawGlitch(s tcell.Screen, width, height, intensity int, rGen *rand.Rand, useCP437 bool) {
	charSet := glitchChars
	if useCP437 {
		charSet = cp437Chars
	}

	numGlitch := rGen.Intn(100*intensity) + (50 * intensity)
	for i := 0; i < numGlitch; i++ {
		x := rGen.Intn(width)
		y := rGen.Intn(height)

		r := rune(charSet[rGen.Intn(len(charSet))])
		style := tcell.StyleDefault.Foreground(glitchColors[rGen.Intn(len(glitchColors))])

		s.SetContent(x, y, r, nil, style)
	}

	if rGen.Intn(10) < 2 {
		shiftLineGlitch(s, width, height, rGen)
	}

	if rGen.Intn(10) < 1 {
		blockDistortionGlitch(s, width, height, rGen)
	}
}

func main() {
	// Define command-line flags
	fps := flag.Int("fps", 30, "frames per second for the animation")
	intensity := flag.Int("intensity", 5, "glitch intensity (1-10)")
	useCP437 := flag.Bool("cp437", false, "use Code Page 437 characters for a retro effect")
	flag.Parse()

	// Clamp intensity
	if *intensity < 1 {
		*intensity = 1
	}
	if *intensity > 10 {
		*intensity = 10
	}

	// Create a local random number generator
	rGen := rand.New(rand.NewSource(time.Now().UnixNano()))

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
			drawGlitch(s, width, height, *intensity, rGen, *useCP437)
			s.Show()
		}
	}
}
