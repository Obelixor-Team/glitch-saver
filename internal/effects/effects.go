package effects

import (
	"glitch-saver/internal/options"
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
	r        rune
	style    tcell.Style
	lifetime int
}

var smearBuffer [][]SmearCell
var ghostBuffer [][]SmearCell

// staticFrames tracks the remaining duration of a static burst.
var staticFrames int

// ScrollingBlock represents a block of the screen that is scrolling.
type ScrollingBlock struct {
	srcX, srcY, destX, destY, w, h, dx, dy, life int
	cells                                        [][]SmearCell
}

var scrollingBlocks []*ScrollingBlock

func InitializeEffects(width, height int) {
	smearBuffer = make([][]SmearCell, height)
	ghostBuffer = make([][]SmearCell, height)
	for i := range smearBuffer {
		smearBuffer[i] = make([]SmearCell, width)
		ghostBuffer[i] = make([]SmearCell, width)
	}
	scrollingBlocks = nil
	cyclingCells = make(map[Point]int)
}

// shiftLineGlitch shifts a random line horizontally
func shiftLineGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) { // opts added
	if height == 0 || width == 0 {
		return
	}
	y := rGen.Intn(height)
	offset := rGen.Intn(width/2) - (width / 4)

	line := make([]struct {
		r     rune
		style tcell.Style
	}, width)

	for x := 0; x < width; x++ {
		// Bounds checking to prevent access beyond screen dimensions
		if y >= 0 && y < height && x >= 0 && x < width {
			mainc, style, _ := s.Get(x, y)
			line[x].r = rune(mainc[0])
			line[x].style = style
		}
	}

	for x := 0; x < width; x++ {
		newX := x + offset
		if newX >= 0 && newX < width && x >= 0 && x < width {
			if line[x].r != 0 { // Only draw if the buffered rune is not a zero value
				s.SetContent(newX, y, line[x].r, nil, line[x].style)
			}
		}
	}
}

// applyVerticalLineGlitch shifts a random column vertically
func applyVerticalLineGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) {
	if width == 0 || height == 0 {
		return
	}
	x := rGen.Intn(width)
	offset := rGen.Intn(height/2) - (height / 4)

	column := make([]struct {
		r     rune
		style tcell.Style
	}, height)

	for y := 0; y < height; y++ {
		// Bounds checking to prevent access beyond screen dimensions
		if y >= 0 && y < height && x >= 0 && x < width {
			mainc, style, _ := s.Get(x, y)
			column[y].r = rune(mainc[0])
			column[y].style = style
		}
	}

	for y := 0; y < height; y++ {
		newY := y + offset
		if newY >= 0 && newY < height && x >= 0 && x < width {
			if column[y].r != 0 {
				s.SetContent(x, newY, column[y].r, nil, column[y].style)
			}
		}
	}
}

// applyInvertColorsGlitch inverts the colors of a random block of the screen
func applyInvertColorsGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) {
	if width == 0 || height == 0 {
		return
	}
	blockX := rGen.Intn(width)
	blockY := rGen.Intn(height)
	blockW := rGen.Intn(width/2) + 1
	blockH := rGen.Intn(height/2) + 1

	// Ensure block dimensions don't exceed screen boundaries
	if blockX+blockW >= width {
		blockW = width - blockX
	}
	if blockY+blockH >= height {
		blockH = height - blockY
	}

	for y := blockY; y < blockY+blockH && y < height; y++ {
		for x := blockX; x < blockX+blockW && x < width; x++ {
			mainc, style, _ := s.Get(x, y)
			fg, bg, _ := style.Decompose()
			newStyle := style.Foreground(bg).Background(fg)
			s.SetContent(x, y, rune(mainc[0]), nil, newStyle)
		}
	}
}

// applyCharScrambleGlitch scrambles the characters in a random block of the screen
func applyCharScrambleGlitch(s tcell.Screen, width, height int, rGen *rand.Rand) {
	if width == 0 || height == 0 {
		return
	}
	blockX := rGen.Intn(width)
	blockY := rGen.Intn(height)
	blockW := rGen.Intn(width/4) + 2
	blockH := rGen.Intn(height/4) + 2

	// Ensure block dimensions don't exceed screen boundaries
	if blockX+blockW >= width {
		blockW = width - blockX
	}
	if blockY+blockH >= height {
		blockH = height - blockY
	}

	// Read the block's content
	cells := make([][]struct {
		r     rune
		style tcell.Style
	}, blockH)
	for y := 0; y < blockH; y++ {
		cells[y] = make([]struct {
			r     rune
			style tcell.Style
		}, blockW)
		for x := 0; x < blockW; x++ {
			if blockX+x < width && blockY+y < height {
				mainc, style, _ := s.Get(blockX+x, blockY+y)
				cells[y][x].r = rune(mainc[0])
				cells[y][x].style = style
			}
		}
	}

	// Flatten, shuffle, and rewrite the characters
	runes := make([]rune, 0, blockW*blockH)
	for _, row := range cells {
		for _, cell := range row {
			runes = append(runes, cell.r)
		}
	}
	rGen.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	// Rewrite the scrambled characters
	i := 0
	for y := 0; y < blockH; y++ {
		for x := 0; x < blockW; x++ {
			if blockX+x < width && blockY+y < height {
				if i < len(runes) {
					s.SetContent(blockX+x, blockY+y, runes[i], nil, cells[y][x].style)
					i++
				}
			}
		}
	}
}

// applyTunnelEffect creates a zoom/tunnel effect by shifting characters
func applyTunnelEffect(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.TunnelEnable {
		return
	}
	if rGen.Float64() > opts.TunnelProbability {
		return
	}

	centerX := width / 2
	centerY := height / 2

	// Create a temporary buffer to hold the original screen state for this effect
	originalScreen := make([][]struct {
		r     rune
		style tcell.Style
	}, height)
	for i := range originalScreen {
		originalScreen[i] = make([]struct {
			r     rune
			style tcell.Style
		}, width)
	}

	// Read the current screen content first
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mainc, style, _ := s.Get(x, y)
			originalScreen[y][x].r = rune(mainc[0])
			originalScreen[y][x].style = style
		}
	}

	// Create a temporary buffer to hold the new screen state
	newScreen := make([][]struct {
		r     rune
		style tcell.Style
	}, height)
	for i := range newScreen {
		newScreen[i] = make([]struct {
			r     rune
			style tcell.Style
		}, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := x - centerX
			dy := y - centerY

			// Simple zoom in/out logic
			srcX := x + (dx*opts.TunnelSpeed)/10
			srcY := y + (dy*opts.TunnelSpeed)/10

			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				newScreen[y][x] = originalScreen[srcY][srcX]
			} else {
				newScreen[y][x].r = ' '
				newScreen[y][x].style = tcell.StyleDefault
			}
		}
	}

	// Apply the new screen state
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			s.SetContent(x, y, newScreen[y][x].r, nil, newScreen[y][x].style)
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

	// Ensure block dimensions don't exceed screen boundaries
	if srcX+blockW >= width {
		blockW = width - srcX
	}
	if srcY+blockH >= height {
		blockH = height - srcY
	}

	destX, destY := rGen.Intn(width), rGen.Intn(height)

	// Ensure destination block doesn't exceed screen boundaries
	if destX+blockW >= width {
		blockW = width - destX
	}
	if destY+blockH >= height {
		blockH = height - destY
	}

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
				mainc, style, _ := s.Get(srcX+x, srcY+y)
				block[y][x].r = rune(mainc[0])
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
func applyCharCorruption(s tcell.Screen, width, height int, rGen *rand.Rand, charSet []rune, fgColors []tcell.Color, opts *options.GlitchOptions, bgColors []tcell.Color) {
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

		// Add to ghost buffer
		if opts.GhostingEnable {
			if rGen.Float64() < opts.GhostingProbability {
				ghostBuffer[y][x] = SmearCell{r, style, 10} // 10 frames lifetime for ghost
			}
		}
	}
}

// applyScanlineEffect draws a horizontal scanline with glitch effects.
func applyScanlineEffect(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
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
		numScanlineChars = rGen.Intn(width/2) + (width / 4 * opts.ScanlineIntensity / 10) // Scale with intensity
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
func applyColorCycle(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.ColorCycleEnable {
		return
	}

	// Clean up out-of-bounds positions first
	for p := range cyclingCells {
		if p.X < 0 || p.X >= width || p.Y < 0 || p.Y >= height {
			delete(cyclingCells, p)
		}
	}

	for p, colorIndex := range cyclingCells {
		// Double-check bounds after potential cleanup
		if p.X < 0 || p.X >= width || p.Y < 0 || p.Y >= height {
			delete(cyclingCells, p)
			continue
		}

		mainc, style, _ := s.Get(p.X, p.Y)
		if rune(mainc[0]) == ' ' {
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

		s.SetContent(p.X, p.Y, rune(mainc[0]), nil, newStyle)
	}
}

// applySmear draws and fades smeared characters.
func applySmear(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.SmearEnable {
		return
	}

	// Ensure buffer dimensions match screen dimensions to prevent out-of-bounds access
	if len(smearBuffer) != height {
		// Reinitialize smearBuffer if dimensions don't match
		InitializeEffects(width, height)
	}

	for y := 0; y < height && y < len(smearBuffer); y++ {
		for x := 0; x < width && x < len(smearBuffer[y]); x++ {
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

// applyGhostingEffect draws and fades ghosted characters.
func applyGhostingEffect(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.GhostingEnable {
		return
	}

	// Ensure buffer dimensions match screen dimensions to prevent out-of-bounds access
	if len(ghostBuffer) != height {
		// Reinitialize ghostBuffer if dimensions don't match
		InitializeEffects(width, height)
	}

	for y := 0; y < height && y < len(ghostBuffer); y++ {
		for x := 0; x < width && x < len(ghostBuffer[y]); x++ {
			if ghostBuffer[y][x].lifetime > 0 {
				ghostBuffer[y][x].lifetime--
				// Draw the ghost with a dimmer style
				fg, bg, _ := ghostBuffer[y][x].style.Decompose()
				ghostStyle := tcell.StyleDefault.Foreground(fg).Background(bg).Dim(true)
				s.SetContent(x, y, ghostBuffer[y][x].r, nil, ghostStyle)

				if ghostBuffer[y][x].lifetime == 0 {
					s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
				}
			}
		}
	}
}

// applyStaticBurst fills the screen with static noise.
func applyStaticBurst(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
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
func applyScrollingBlocks(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.ScrollEnable {
		return
	}

	// Remove dead blocks and blocks that would be out of bounds after resize
	newScrollingBlocks := scrollingBlocks[:0]
	for _, b := range scrollingBlocks {
		if b.life > 0 {
			// Remove blocks that would be outside the screen after resize
			if b.destX < 0 || b.destY < 0 || b.destX+b.w > width || b.destY+b.h > height {
				// Check if the block is completely outside the new screen dimensions
				if b.destX > width || b.destY > height || b.destX+b.w < 0 || b.destY+b.h < 0 {
					continue // Skip this block
				}
			}
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
				mainc, style, _ := s.Get(srcX+x, srcY+y)
				cells[y][x] = SmearCell{rune(mainc[0]), style, 1}
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

// DrawGlitch orchestrates various glitch effects on the screen.
func DrawGlitch(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) { // opts replaces many args
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

	// If both UseBlocks and UseCP437 are enabled, combine the character sets
	if opts.UseBlocks && opts.UseCP437 {
		charSet = append([]rune(blockChars), []rune(cp437Chars)...)
	}

	if opts.CharCorruptionEnable {
		applyCharCorruption(s, width, height, rGen, charSet, glitchColors, opts, glitchColors)
	}

	if opts.ShiftLineEnable && rGen.Intn(10) < 2 {
		shiftLineGlitch(s, width, height, rGen)
	}

	if opts.VerticalLineEnable && rGen.Float64() < opts.VerticalLineProbability {
		applyVerticalLineGlitch(s, width, height, rGen)
	}

	if opts.InvertColorsEnable && rGen.Float64() < opts.InvertColorsProbability {
		applyInvertColorsGlitch(s, width, height, rGen)
	}

	if opts.CharScrambleEnable && rGen.Float64() < opts.CharScrambleProbability {
		applyCharScrambleGlitch(s, width, height, rGen)
	}

	if opts.TunnelEnable && rGen.Float64() < opts.TunnelProbability {
		applyTunnelEffect(s, width, height, rGen, opts)
	}

	if opts.BlockDistortionEnable && rGen.Intn(10) < 1 {
		blockDistortionGlitch(s, width, height, rGen)
	}

	applyScanlineEffect(s, width, height, rGen, opts)  // Call new scanline effect
	applyColorCycle(s, width, height, rGen, opts)      // Call new color cycle effect
	applySmear(s, width, height, rGen, opts)           // Call new smear effect
	applyGhostingEffect(s, width, height, rGen, opts)  // Call new ghosting effect
	applyScrollingBlocks(s, width, height, rGen, opts) // Call new scrolling blocks effect
	applyBitRot(s, width, height, rGen, opts)
	applyMelt(s, width, height, rGen, opts)
	applyJitter(s, width, height, rGen, opts)
}

func applyBitRot(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.BitRotEnable {
		return
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rGen.Float64() < opts.BitRotProbability {
				_, style, _ := s.Get(x, y)
				r := []rune(cp437Chars)[rGen.Intn(len(cp437Chars))] // Use the pre-converted slice

				s.SetContent(x, y, r, nil, style)
			}
		}
	}
}

func applyMelt(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
	if !opts.MeltEnable {
		return
	}

	for y := height - 2; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if rGen.Float64() < opts.MeltProbability {
				c, style, _ := s.Get(x, y)
				if y+1 < height { // Bounds check
					below, _, _ := s.Get(x, y+1)

					if rune(below[0]) == ' ' {
						s.SetContent(x, y+1, rune(c[0]), nil, style)
						s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
					}
				}
			}
		}
	}
}

func applyJitter(s tcell.Screen, width, height int, rGen *rand.Rand, opts *options.GlitchOptions) {
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

				// Only proceed if the calculated coordinates are within bounds
				if x >= 0 && x < width && y >= 0 && y < height && nx >= 0 && nx < width && ny >= 0 && ny < height {
					// Swap cells
					c1, style1, _ := s.Get(x, y)
					c2, style2, _ := s.Get(nx, ny)
					s.SetContent(x, y, rune(c2[0]), nil, style2)
					s.SetContent(nx, ny, rune(c1[0]), nil, style1)
				}
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
