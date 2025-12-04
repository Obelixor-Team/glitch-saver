package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"glitch-saver/internal/options"
	"glitch-saver/internal/tui"
)

func main() {
	opts := options.ParseOptions()

	if opts.LoadPreset != "" {
		data, err := ioutil.ReadFile(opts.LoadPreset)
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
		if err := ioutil.WriteFile(opts.SavePreset, data, 0644); err != nil {
			log.Fatalf("failed to write preset file: %v", err)
		}
	}

	tui.RunTUI(opts)
}
