// Copyright (c) 2020 Cesar Gimenes - MIT License

// Package ansi converts a PNG image into the ANSI escape sequences needed to
// recreate it in a terminal using Unicode half-block characters.
//
// Every two rows of pixels in the source image become a single row of terminal
// text: each cell shows its two pixels using " ", "█", "▀" or "▄", with the
// glyph and colors chosen to reuse the terminal state and keep the output
// small. If the image has an odd number of rows, the missing bottom row is
// rendered with the configured default color.
package ansi
