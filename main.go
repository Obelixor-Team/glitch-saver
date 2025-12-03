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
const cp437Chars = "ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜⛛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■ "
const blockChars = "░▒▓█"
const staticChars = " .*"

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

var staticColors = []tcell.Color{
	tcell.NewRGBColor(0, 0, 0),       // Black
	tcell.NewRGBColor(128, 128, 128), // Grey
	tcell.NewRGBColor(255, 255, 255), // White
}

// Point represents a coordinate on the screen.
type Point struct {
	X, Y int
}

// cyclingCells holds the state of cells that are cycling colors.
var cyclingCells = make(map[Point]int)

// SmearCell represents a cell with a trail life.
type SmearCell struct {
	r     rune
	style tcell.Style
	lifetime  int
}

var smearBuffer [][]SmearCell

// staticFrames tracks the remaining duration of a static burst.
var staticFrames int

// ScrollingBlock represents a block of the screen that is scrolling.
type ScrollingBlock struct {
	srcX, srcY, destX, destY, w, h, dx, dy, life int
	cells                                       [][]SmearCell
}

var scrollingBlocks []*ScrollingBlock

// GlitchOptions holds all configurable parameters for the glitch effects.
type GlitchOptions struct {
	FPS                 int
	Intensity           int
	UseCP437            bool
	UseBlocks           bool
	UseBG               bool
	ScanlineEnable      bool
	ScanlineProbability float64
	ScanlineIntensity   int
	ScanlineChar        string
	ColorCycleEnable    bool
	ColorCycleSpeed     int
	SmearEnable         bool
	SmearProbability    float64
	SmearLength         int
	StaticEnable        bool
	StaticProbability   float64
	StaticDuration      int
	StaticChar          string
	ScrollEnable        bool
	ScrollProbability   float64
	ScrollSpeed         int
	ScrollDirection     string
	// Add more options here later
}

// shiftLineGlitch shifts a random line horizontally
func shiftLineGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) { // opts added
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
			if line[x].r != 0 { // Only draw if the buffered rune is not a zero value
				s.SetContent(newX, y, line[x].r, nil, line[x].style)
			}
		}
	}
}

// blockDistortionGlitch copies a random block of the screen to another random location
func blockDistortionGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) { // opts added
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
				if block[y][x].r != 0 { // Only draw if the buffered rune is not a zero value
					s.SetContent(destX+x, destY+y, block[y][x].r, nil, block[y][x].style)
				}
			}
		}
	}
}

// applyCharCorruption draws random characters with glitch effects to the screen.
func applyCharCorruption(s tcell.Screen, width, height int, rGen *rand.Rand, charSet []rune, fgColors []tcell.Color, opts *GlitchOptions, bgColors []tcell.Color) {
	numGlitch := rGen.Intn(100*opts.Intensity) + (50 * opts.Intensity)
	for i := 0; i < numGlitch; i++ {
		x := rGen.Intn(width)
		y := rGen.Intn(height)

		r := charSet[rGen.Intn(len(charSet))]
		fg := fgColors[rGen.Intn(len(fgColors))]

		style := tcell.StyleDefault.Foreground(fg)

		if opts.UseBG {
			bg := bgColors[rGen.Intn(len(bgColors))]
			style = style.Background(bg)
		}

		s.SetContent(x, y, r, nil, style)

		// Add to color cycling
		if opts.ColorCycleEnable {
			if rGen.Float64() < 0.1 { // 10% chance to add to cycling
				cyclingCells[Point{x, y}] = rGen.Intn(len(glitchColors))
			}
		}

		// Add to smear buffer
		if opts.SmearEnable {
			if rGen.Float64() < opts.SmearProbability {
				smearBuffer[y][x] = SmearCell{r, style, opts.SmearLength}
			}
		}
	}
}

// applyScanlineEffect draws a horizontal scanline with glitch effects.
func applyScanlineEffect(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if height == 0 || !opts.ScanlineEnable {
		return
	}
	if rGen.Float64() > opts.ScanlineProbability { // Check probability
		return
	}

	y := rGen.Intn(height) // Random row

	scanlineRunes := []rune(glitchChars)
	if opts.ScanlineChar != "" {
		scanlineRunes = []rune(opts.ScanlineChar)
	} else if opts.UseBlocks {
		scanlineRunes = []rune(blockChars)
	} else if opts.UseCP437 {
		scanlineRunes = []rune(cp437Chars)
	}

	numScanlineChars := width / 2 // Default density
	if opts.ScanlineIntensity > 0 {
		numScanlineChars = rGen.Intn(width/2) + (width/4 * opts.ScanlineIntensity/10) // Scale with intensity
	}
	if numScanlineChars > width {
		numScanlineChars = width
	}


	for i := 0; i < numScanlineChars; i++ {
		x := rGen.Intn(width) // Random position within the row
		
		r := scanlineRunes[rGen.Intn(len(scanlineRunes))]
		fg := glitchColors[rGen.Intn(len(glitchColors))]
		
		style := tcell.StyleDefault.Foreground(fg)
		if opts.UseBG {
			bg := glitchColors[rGen.Intn(len(glitchColors))]
			style = style.Background(bg)
		}

		s.SetContent(x, y, r, nil, style)
	}
}

// applyColorCycle updates the colors of cycling cells.
func applyColorCycle(s tcell.Screen, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.ColorCycleEnable {
		return
	}

	for p, colorIndex := range cyclingCells {
		rawVal, style, _ := s.Get(p.X, p.Y)
		var r rune
		if len(rawVal) > 0 {
			r = []rune(rawVal)[0]
		}
		if r == 0 {
			delete(cyclingCells, p)
			continue
		}

		// Update color index
		colorIndex = (colorIndex + opts.ColorCycleSpeed) % len(glitchColors)
		cyclingCells[p] = colorIndex
		
		newStyle := style.Foreground(glitchColors[colorIndex])
		
		if opts.UseBG {
			bg := glitchColors[(colorIndex+len(glitchColors)/2)%len(glitchColors)] // Offset background color
			newStyle = newStyle.Background(bg)
		}

		s.SetContent(p.X, p.Y, r, nil, newStyle)
	}
}

// applySmear draws and fades smeared characters.
func applySmear(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.SmearEnable {
		return
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if smearBuffer[y][x].lifetime > 0 {
				smearBuffer[y][x].lifetime--
				s.SetContent(x, y, smearBuffer[y][x].r, nil, smearBuffer[y][x].style.Dim(true))
				if smearBuffer[y][x].lifetime == 0 {
					s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
				}
			}
		}
	}
}

// applyStaticBurst fills the screen with static noise.
func applyStaticBurst(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	staticRunes := []rune(staticChars)
	if opts.StaticChar != "" {
		staticRunes = []rune(opts.StaticChar)
	}
	
	numStaticChars := (width * height) / 4 // Cover a quarter of the screen with static
	for i := 0; i < numStaticChars; i++ {
		x := rGen.Intn(width)
		y := rGen.Intn(height)
		
		r := staticRunes[rGen.Intn(len(staticRunes))]
		fg := staticColors[rGen.Intn(len(staticColors))]
		bg := staticColors[rGen.Intn(len(staticColors))]
		
		style := tcell.StyleDefault.Foreground(fg).Background(bg)
		s.SetContent(x, y, r, nil, style)
	}
}

// applyScrollingBlocks scrolls blocks of the screen.
func applyScrollingBlocks(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.ScrollEnable {
		return
	}

	// Remove dead blocks
	newScrollingBlocks := scrollingBlocks[:0]
	for _, b := range scrollingBlocks {
		if b.life > 0 {
			newScrollingBlocks = append(newScrollingBlocks, b)
		}
	}
	scrollingBlocks = newScrollingBlocks

	// Update and draw existing blocks
	for _, b := range scrollingBlocks {
		b.life--
		b.destX += b.dx
		b.destY += b.dy

		for y := 0; y < b.h; y++ {
			for x := 0; x < b.w; x++ {
				if b.destX+x < width && b.destY+y < height && b.destX+x >= 0 && b.destY+y >= 0 {
					s.SetContent(b.destX+x, b.destY+y, b.cells[y][x].r, nil, b.cells[y][x].style)
				}
			}
		}
	}

	// Trigger new blocks
	if rGen.Float64() < opts.ScrollProbability {
		srcX, srcY := rGen.Intn(width), rGen.Intn(height)
		blockW := rGen.Intn(width/4) + 5
		blockH := rGen.Intn(height/4) + 5

		if srcX+blockW > width {
			blockW = width - srcX
		}
		if srcY+blockH > height {
			blockH = height - srcY
		}

		cells := make([][]SmearCell, blockH)
		for y := 0; y < blockH; y++ {
			cells[y] = make([]SmearCell, blockW)
			for x := 0; x < blockW; x++ {
				rawVal, style, _ := s.Get(srcX+x, srcY+y)
				var r rune
				if len(rawVal) > 0 {
					r = []rune(rawVal)[0]
				}
				cells[y][x] = SmearCell{r, style, 1}
			}
		}

		var dx, dy int
		switch opts.ScrollDirection {
		case "horizontal":
			dx = opts.ScrollSpeed
			if rGen.Intn(2) == 0 {
				dx = -dx
			}
		case "vertical":
			dy = opts.ScrollSpeed
			if rGen.Intn(2) == 0 {
				dy = -dy
			}
		default: // random
			dx = rGen.Intn(opts.ScrollSpeed*2+1) - opts.ScrollSpeed
			dy = rGen.Intn(opts.ScrollSpeed*2+1) - opts.ScrollSpeed
		}
		if dx == 0 && dy == 0 {
			dx = 1 // Ensure movement
		}

		scrollingBlocks = append(scrollingBlocks, &ScrollingBlock{
			srcX:  srcX,
			srcY:  srcY,
			destX: srcX,
			destY: srcY,
			w:     blockW,
			h:     blockH,
			dx:    dx,
			dy:    dy,
			life:  (width + height) / (abs(dx) + abs(dy) + 1),
			cells: cells,
		})
	}
}


// drawGlitch orchestrates various glitch effects on the screen.
func drawGlitch(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) { // opts replaces many args
	if staticFrames > 0 {
		applyStaticBurst(s, width, height, rGen, opts)
		staticFrames--
		return
	}
	if opts.StaticEnable && rGen.Float64() < opts.StaticProbability {
		staticFrames = opts.StaticDuration
		return
	}

	var charSet []rune
	if opts.UseBlocks {
		charSet = []rune(blockChars)
	} else if opts.UseCP437 {
		charSet = []rune(cp437Chars)
	} else {
		charSet = []rune(glitchChars)
	}

	applyCharCorruption(s, width, height, rGen, charSet, glitchColors, opts, glitchColors)

	if rGen.Intn(10) < 2 {
		shiftLineGlitch(s, width, height, rGen)
	}

	if rGen.Intn(10) < 1 {
		blockDistortionGlitch(s, width, height, rGen)
	}

	applyScanlineEffect(s, width, height, rGen, opts) // Call new scanline effect
	applyColorCycle(s, rGen, opts) // Call new color cycle effect
	applySmear(s, width, height, rGen, opts) // Call new smear effect
	applyScrollingBlocks(s, width, height, rGen, opts) // Call new scrolling blocks effect
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	var opts GlitchOptions

	// Define command-line flags and populate opts
	flag.IntVar(&opts.FPS, "fps", 30, "frames per second for the animation")
	flag.IntVar(&opts.Intensity, "intensity", 5, "glitch intensity (1-10)")
	flag.BoolVar(&opts.UseCP437, "cp437", false, "use Code Page 437 characters for a retro effect")
	flag.BoolVar(&opts.UseBlocks, "blocks", false, "use only block characters for a heavy glitch effect")
	flag.BoolVar(&opts.UseBG, "bg", false, "enable random background coloring")
	flag.BoolVar(&opts.ScanlineEnable, "scanline", false, "enable scanline glitch effect")
	flag.Float64Var(&opts.ScanlineProbability, "scanline-prob", 0.1, "probability (0.0-1.0) of a scanline appearing each frame")
	flag.IntVar(&opts.ScanlineIntensity, "scanline-intensity", 5, "intensity (1-10) of scanlines")
	flag.StringVar(&opts.ScanlineChar, "scanline-char", "", "character to use for scanlines (default: random from current charSet)")
	flag.BoolVar(&opts.ColorCycleEnable, "color-cycle", false, "enable color cycling effect")
	flag.IntVar(&opts.ColorCycleSpeed, "color-cycle-speed", 5, "speed (1-10) of color cycling")
	flag.BoolVar(&opts.SmearEnable, "smear", false, "enable character smearing/trails effect")
	flag.Float64Var(&opts.SmearProbability, "smear-prob", 0.1, "probability (0.0-1.0) of a character starting to smear")
	flag.IntVar(&opts.SmearLength, "smear-length", 5, "length of the smear trail (in frames)")
	flag.BoolVar(&opts.StaticEnable, "static", false, "enable static burst effect")
	flag.Float64Var(&opts.StaticProbability, "static-prob", 0.01, "probability (0.0-1.0) of a static burst occurring each frame")
	flag.IntVar(&opts.StaticDuration, "static-duration", 3, "duration of a static burst (in frames)")
	flag.StringVar(&opts.StaticChar, "static-char", "", "character to use for static bursts (default: random from '. *')")
	flag.BoolVar(&opts.ScrollEnable, "scroll", false, "enable scrolling blocks effect")
	flag.Float64Var(&opts.ScrollProbability, "scroll-prob", 0.05, "probability (0.0-1.0) of a new scrolling block appearing each frame")
	flag.IntVar(&opts.ScrollSpeed, "scroll-speed", 1, "speed of scrolling blocks")
	flag.StringVar(&opts.ScrollDirection, "scroll-direction", "random", "direction of scrolling blocks (horizontal, vertical, random)")
	flag.Parse()

	// Clamp intensity
	if opts.Intensity < 1 {
		opts.Intensity = 1
	}
	if opts.Intensity > 10 {
		opts.Intensity = 10
	}
	// Clamp scanline probability
	if opts.ScanlineProbability < 0.0 {
		opts.ScanlineProbability = 0.0
	}
	if opts.ScanlineProbability > 1.0 {
		opts.ScanlineProbability = 1.0
	}
	// Clamp scanline intensity
	if opts.ScanlineIntensity < 1 {
		opts.ScanlineIntensity = 1
	}
	if opts.ScanlineIntensity > 10 {
		opts.ScanlineIntensity = 10
	}
	// Clamp color cycle speed
	if opts.ColorCycleSpeed < 1 {
		opts.ColorCycleSpeed = 1
	}
	if opts.ColorCycleSpeed > 10 {
		opts.ColorCycleSpeed = 10
	}
	// Clamp smear probability
	if opts.SmearProbability < 0.0 {
		opts.SmearProbability = 0.0
	}
	if opts.SmearProbability > 1.0 {
		opts.SmearProbability = 1.0
	}
	// Clamp smear length
	if opts.SmearLength < 1 {
		opts.SmearLength = 1
	}
	// Clamp static probability
	if opts.StaticProbability < 0.0 {
		opts.StaticProbability = 0.0
	}
	if opts.StaticProbability > 1.0 {
		opts.StaticProbability = 1.0
	}
	// Clamp static duration
	if opts.StaticDuration < 1 {
		opts.StaticDuration = 1
	}
	// Clamp scroll probability
	if opts.ScrollProbability < 0.0 {
		opts.ScrollProbability = 0.0
	}
	if opts.ScrollProbability > 1.0 {
		opts.ScrollProbability = 1.0
	}
	// Clamp scroll speed
	if opts.ScrollSpeed < 1 {
		opts.ScrollSpeed = 1
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

	// Initialize smearBuffer
	smearBuffer = make([][]SmearCell, height)
	for i := range smearBuffer {
		smearBuffer[i] = make([]SmearCell, width)
	}

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
				// Re-initialize smearBuffer on resize
				smearBuffer = make([][]SmearCell, height)
				for i := range smearBuffer {
					smearBuffer[i] = make([]SmearCell, width)
				}
				s.Clear()                // Clear screen on resize to avoid artifacts
				s.Sync()                 // Sync screen after resize
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					return // Exit the application
				}
			}
		case <-ticker.C: // Handle animation tick
			drawGlitch(s, width, height, rGen, &opts) // Pass opts struct
			s.Show()
		}
	}
}
