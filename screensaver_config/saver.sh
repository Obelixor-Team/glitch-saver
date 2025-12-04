#!/bin/bash

# --- Configuration ---
# Timeout in seconds (300 seconds = 5 minutes)
export TMOUT=300

# The full path to your glitch-saver executable
# IMPORTANT: You MUST change this path to the actual location of your glitch-saver
SAVER_APP="/path/to/your/glitch-saver"

# --- Screensaver Function ---
run_saver() {
    # Clear the screen and run the screensaver
    clear
    $SAVER_APP
}

# --- Trap Setup ---
# Catch the signal when TMOUT timer expires and call run_saver
trap 'run_saver' ALRM

echo "Terminal screensaver is active. It will start after 5 minutes of inactivity."
