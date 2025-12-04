package main

import (
	"math/rand"

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

func InitializeEffects(width, height int) {
	smearBuffer = make([][]SmearCell, height)
	for i := range smearBuffer {
		smearBuffer[i] = make([]SmearCell, width)
	}
	scrollingBlocks = nil
	cyclingCells = make(map[Point]int)
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
		mainc, _, style, _ := s.GetContent(x, y)
		line[x].r = mainc
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
				mainc, _, style, _ := s.GetContent(srcX+x, srcY+y)
				block[y][x].r = mainc
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
		mainc, _, style, _ := s.GetContent(p.X, p.Y)
		if mainc == ' ' {
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

		s.SetContent(p.X, p.Y, mainc, nil, newStyle)
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
				mainc, _, style, _ := s.GetContent(srcX+x, srcY+y)
				cells[y][x] = SmearCell{mainc, style, 1}
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
	applyBitRot(s, width, height, rGen, opts)
	applyMelt(s, width, height, rGen, opts)
	applyJitter(s, width, height, rGen, opts)
}

func applyBitRot(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.BitRotEnable {
		return
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rGen.Float64() < opts.BitRotProbability {
				_, _, style, _ := s.GetContent(x, y)
				r := []rune(cp437Chars)[rGen.Intn(len(cp437Chars))]
				s.SetContent(x, y, r, nil, style)
			}
		}
	}
}

func applyMelt(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.MeltEnable {
		return
	}

	for y := height - 2; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if rGen.Float64() < opts.MeltProbability {
				c, _, style, _ := s.GetContent(x, y)
				below, _, _, _ := s.GetContent(x, y+1)

				if below == ' ' {
					s.SetContent(x, y+1, c, nil, style)
					s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
				}
			}
		}
	}
}

func applyJitter(s tcell.Screen, width, height int, rGen *rand.Rand, opts *GlitchOptions) {
	if !opts.JitterEnable {
		return
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rGen.Float64() < opts.JitterProbability {
				// Pick a random neighbor
				nx, ny := x+rGen.Intn(3)-1, y+rGen.Intn(3)-1

				// Clamp to screen bounds
				if nx < 0 {
					nx = 0
				}
				if nx >= width {
					nx = width - 1
				}
				if ny < 0 {
					ny = 0
				}
				if ny >= height {
					ny = height - 1
				}
				
				// Swap cells
				c1, _, style1, _ := s.GetContent(x, y)
				c2, _, style2, _ := s.GetContent(nx, ny)
				s.SetContent(x, y, c2, nil, style2)
				s.SetContent(nx, ny, c1, nil, style1)
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
