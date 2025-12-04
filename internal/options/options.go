package options

import (
	"flag"
)

// GlitchOptions holds all configurable parameters for the glitch effects.
type GlitchOptions struct {
	FPS                     int
	Intensity               int
	UseCP437                bool
	UseBlocks               bool
	UseBG                   bool
	ShiftLineEnable         bool
	BlockDistortionEnable   bool
	CharCorruptionEnable    bool
	ScanlineEnable          bool
	ScanlineProbability     float64
	ScanlineIntensity       int
	ScanlineChar            string
	ColorCycleEnable        bool
	ColorCycleSpeed         int
	SmearEnable             bool
	SmearProbability        float64
	SmearLength             int
	StaticEnable            bool
	StaticProbability       float64
	StaticDuration          int
	StaticChar              string
	ScrollEnable            bool
	ScrollProbability       float64
	ScrollSpeed             int
	ScrollDirection         string
	JitterEnable            bool
	JitterProbability       float64
	MeltEnable              bool
	MeltProbability         float64
	BitRotEnable            bool
	BitRotProbability       float64
	VerticalLineEnable      bool
	VerticalLineProbability float64
	InvertColorsEnable      bool
	InvertColorsProbability float64
	CharScrambleEnable      bool
	CharScrambleProbability float64
	GhostingEnable          bool
	GhostingProbability     float64
	TunnelEnable            bool
	TunnelProbability       float64
	TunnelSpeed             int
	AllEffectsEnable        bool
	SavePreset              string
	LoadPreset              string
	// Add more options here later
}

func ParseOptions() *GlitchOptions {
	var opts GlitchOptions

	// Define command-line flags and populate opts
	flag.StringVar(&opts.SavePreset, "save-preset", "", "save the current options to a file")
	flag.StringVar(&opts.LoadPreset, "load-preset", "", "load options from a file")
	flag.IntVar(&opts.FPS, "fps", 30, "frames per second for the animation")
	flag.IntVar(&opts.Intensity, "intensity", 5, "glitch intensity (1-10)")
	flag.BoolVar(&opts.UseCP437, "cp437", false, "use Code Page 437 characters for a retro effect")
	flag.BoolVar(&opts.UseBlocks, "blocks", false, "use only block characters for a heavy glitch effect")
	flag.BoolVar(&opts.UseBG, "bg", false, "enable random background coloring")
	flag.BoolVar(&opts.ShiftLineEnable, "shift-line", false, "enable horizontal line shift glitch effect")
	flag.BoolVar(&opts.BlockDistortionEnable, "block-distort", false, "enable block distortion glitch effect")
	flag.BoolVar(&opts.CharCorruptionEnable, "char-corrupt", false, "enable character corruption glitch effect")
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
	flag.BoolVar(&opts.JitterEnable, "jitter", false, "enable jitter effect")
	flag.Float64Var(&opts.JitterProbability, "jitter-prob", 0.1, "probability (0.0-1.0) of a character jittering")
	flag.BoolVar(&opts.MeltEnable, "melt", false, "enable melt effect")
	flag.Float64Var(&opts.MeltProbability, "melt-prob", 0.1, "probability (0.0-1.0) of a character melting")
	flag.BoolVar(&opts.BitRotEnable, "bitrot", false, "enable bit-rot effect")
	flag.Float64Var(&opts.BitRotProbability, "bitrot-prob", 0.1, "probability (0.0-1.0) of a character bit-rotting")
	flag.BoolVar(&opts.VerticalLineEnable, "vert-line", false, "enable vertical line glitch effect")
	flag.Float64Var(&opts.VerticalLineProbability, "vert-line-prob", 0.1, "probability (0.0-1.0) of a vertical line appearing each frame")
	flag.BoolVar(&opts.InvertColorsEnable, "invert-colors", false, "enable invert colors glitch effect")
	flag.Float64Var(&opts.InvertColorsProbability, "invert-colors-prob", 0.1, "probability (0.0-1.0) of a color inversion appearing each frame")
	flag.BoolVar(&opts.CharScrambleEnable, "char-scramble", false, "enable character scramble glitch effect")
	flag.Float64Var(&opts.CharScrambleProbability, "char-scramble-prob", 0.1, "probability (0.0-1.0) of a character scramble appearing each frame")
	flag.BoolVar(&opts.GhostingEnable, "ghosting", false, "enable ghosting trail effect")
	flag.Float64Var(&opts.GhostingProbability, "ghosting-prob", 0.1, "probability (0.0-1.0) of a character starting to ghost")
	flag.BoolVar(&opts.TunnelEnable, "tunnel", false, "enable tunnel/zoom effect")
	flag.Float64Var(&opts.TunnelProbability, "tunnel-prob", 0.1, "probability (0.0-1.0) of a tunnel/zoom effect appearing each frame")
	flag.IntVar(&opts.TunnelSpeed, "tunnel-speed", 1, "speed of the tunnel/zoom effect")
	flag.BoolVar(&opts.AllEffectsEnable, "all-effects", false, "enable all glitch effects")
	flag.Parse()

	if opts.AllEffectsEnable {
		opts.UseCP437 = true
		opts.UseBlocks = true
		opts.UseBG = true
		opts.ShiftLineEnable = true
		opts.BlockDistortionEnable = true
		opts.CharCorruptionEnable = true
		opts.ScanlineEnable = true
		opts.ScanlineProbability = 1.0
		opts.ScanlineIntensity = 10
		opts.ColorCycleEnable = true
		opts.ColorCycleSpeed = 10
		opts.SmearEnable = true
		opts.SmearProbability = 1.0
		opts.SmearLength = 10
		opts.StaticEnable = true
		opts.StaticProbability = 0.5 // High probability, but not 1.0 to allow other effects
		opts.StaticDuration = 5
		opts.ScrollEnable = true
		opts.ScrollProbability = 0.5 // High probability, but not 1.0 to allow other effects
		opts.ScrollSpeed = 5
		opts.JitterEnable = true
		opts.JitterProbability = 1.0
		opts.MeltEnable = true
		opts.MeltProbability = 1.0
		opts.BitRotEnable = true
		opts.BitRotProbability = 1.0
		opts.VerticalLineEnable = true
		opts.VerticalLineProbability = 1.0
		opts.InvertColorsEnable = true
		opts.InvertColorsProbability = 1.0
		opts.CharScrambleEnable = true
		opts.CharScrambleProbability = 1.0
		opts.GhostingEnable = true
		opts.GhostingProbability = 1.0
		opts.TunnelEnable = true
		opts.TunnelProbability = 0.5 // High probability, but not 1.0
		opts.TunnelSpeed = 5
		opts.Intensity = 10
	}

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
	// Clamp jitter probability
	if opts.JitterProbability < 0.0 {
		opts.JitterProbability = 0.0
	}
	if opts.JitterProbability > 1.0 {
		opts.JitterProbability = 1.0
	}
	// Clamp melt probability
	if opts.MeltProbability < 0.0 {
		opts.MeltProbability = 0.0
	}
	if opts.MeltProbability > 1.0 {
		opts.MeltProbability = 1.0
	}
	// Clamp bit-rot probability
	if opts.BitRotProbability < 0.0 {
		opts.BitRotProbability = 0.0
	}
	if opts.BitRotProbability > 1.0 {
		opts.BitRotProbability = 1.0
	}
	// Clamp vertical line probability
	if opts.VerticalLineProbability < 0.0 {
		opts.VerticalLineProbability = 0.0
	}
	if opts.VerticalLineProbability > 1.0 {
		opts.VerticalLineProbability = 1.0
	}
	// Clamp invert colors probability
	if opts.InvertColorsProbability < 0.0 {
		opts.InvertColorsProbability = 0.0
	}
	if opts.InvertColorsProbability > 1.0 {
		opts.InvertColorsProbability = 1.0
	}
	// Clamp char scramble probability
	if opts.CharScrambleProbability < 0.0 {
		opts.CharScrambleProbability = 0.0
	}
	if opts.CharScrambleProbability > 1.0 {
		opts.CharScrambleProbability = 1.0
	}
	// Clamp ghosting probability
	if opts.GhostingProbability < 0.0 {
		opts.GhostingProbability = 0.0
	}
	if opts.GhostingProbability > 1.0 {
		opts.GhostingProbability = 1.0
	}
	// Clamp tunnel probability
	if opts.TunnelProbability < 0.0 {
		opts.TunnelProbability = 0.0
	}
	if opts.TunnelProbability > 1.0 {
		opts.TunnelProbability = 1.0
	}
	// Clamp tunnel speed
	if opts.TunnelSpeed < 1 {
		opts.TunnelSpeed = 1
	}
	if opts.TunnelSpeed > 10 {
		opts.TunnelSpeed = 10
	}

	return &opts
}
