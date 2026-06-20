// Copyright (c) 2020 Cesar Gimenes - MIT License

// Package ansi converts a PNG image into the ANSI escape sequences needed to
// recreate it in a terminal using Unicode half-block characters.
//
// Every two rows of pixels in the source image become a single row of terminal
// text: the upper pixel sets the foreground color, the lower pixel sets the
// background color, and the "▀" (upper half block) glyph splits the cell. If
// the image has an odd number of rows, the missing bottom row is rendered with
// the configured default color.
package ansi
