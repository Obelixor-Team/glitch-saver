# Glitch Saver

A terminal-based glitch art screensaver written in Go.

This application takes over your terminal and displays a chaotic, animated "glitch art" effect by rapidly drawing random characters, colors, and distorted blocks of text.

## Showcase

Here's a glimpse of the glitch effects in action:

![Glitch Saver Animation](animation.gif)

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

You can configure the speed and intensity of the glitch effect and the character set using command-line flags:

- `-fps`: Sets the frames per second for the animation. (Default: 30)
- `-intensity`: Sets the glitch intensity on a scale from 1 to 10. (Default: 5)
- `-cp437`: Use Code Page 437 characters for a retro, text-mode art effect. (Default: false)
- `-blocks`: Use only block characters (e.g., `░▒▓█`) for a heavy, block-based glitch effect. (Default: false)
- `-bg`: Enable random background coloring for an even more chaotic effect. (Default: false)

**Example:** Run at a slower 10 FPS with maximum intensity.
```bash
./glitch-saver -fps 10 -intensity 10
```

**Example:** Run with Code Page 437 characters.
```bash
./glitch-saver -cp437
```

**Example:** Run with block characters and background coloring.
```bash
./glitch-saver -blocks -bg
```

**Example:** Run with the default character set and background coloring.
```bash
./glitch-saver -bg
```
