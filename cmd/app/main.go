package main

import (
	"encoding/json"
	"log"
	"os"

	"glitch-saver/internal/options"
	"glitch-saver/internal/tui"
)

func main() {
	opts := options.ParseOptions()

	if opts.LoadPreset != "" {
		data, err := os.ReadFile(opts.LoadPreset)
		if err != nil {
			log.Fatalf("failed to read preset file: %v", err)
		}
		if err := json.Unmarshal(data, opts); err != nil {
			log.Fatalf("failed to unmarshal preset file: %v", err)
		}
	}

	if opts.SavePreset != "" {
		data, err := json.MarshalIndent(opts, "", "  ")
		if err != nil {
			log.Fatalf("failed to marshal preset: %v", err)
		}
		if err := os.WriteFile(opts.SavePreset, data, 0644); err != nil {
			log.Fatalf("failed to write preset file: %v", err)
		}
	}

	log.Println("Calling RunTUI")
	s, err := tui.RunTUI(opts)
	if err != nil {
		log.Fatalf("TUI application failed: %v", err)
	}
	log.Println("RunTUI returned successfully, finalizing screen.")
	// Ensure the screen is finalized
	if s != nil {
		s.Fini()
	}
	log.Println("Application exited normally.")
}
