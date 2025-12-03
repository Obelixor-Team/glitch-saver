# Glitch Saver

A terminal-based glitch art screensaver written in Go.

This application takes over your terminal and displays a chaotic, animated "glitch art" effect by rapidly drawing random characters, colors, and distorted blocks of text.

## Building

To build the application, you need to have Go installed.

```bash
go build -o glitch-saver main.go
```

## Running

Once built, you can run the screensaver with the following command:

```bash
./glitch-saver
```

To exit the screensaver, press `ESC` or `q`.

### Configuration

You can configure the speed and intensity of the glitch effect using command-line flags:

- `-fps`: Sets the frames per second for the animation. (Default: 30)
- `-intensity`: Sets the glitch intensity on a scale from 1 to 10. (Default: 5)

**Example:** Run at a slower 10 FPS with maximum intensity.
```bash
./glitch-saver -fps 10 -intensity 10
```
