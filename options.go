package main

import (
	"flag"
)

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
	JitterEnable        bool
	JitterProbability   float64
	MeltEnable          bool
	MeltProbability     float64
	BitRotEnable        bool
	BitRotProbability   float64
	SavePreset          string
	LoadPreset          string
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
	
	return &opts
}
