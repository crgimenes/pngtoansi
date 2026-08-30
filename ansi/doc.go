// Copyright (c) 2020 Cesar Gimenes - MIT License

// Package ansi converts a PNG image into the ANSI escape sequences needed to
// recreate it in a terminal using Unicode half-block characters.
//
// Every two rows of pixels in the source image become a single row of terminal
// text: each cell shows its two pixels using " ", "█", "▀" or "▄", with the
// glyph and colors chosen to reuse the terminal state and keep the output
// small. If the image has an odd number of rows, the missing bottom row is
// rendered with the configured default color.
//
// In Sprite mode the output is relocatable: transparent cells (alpha below
// 50%, or an exact TransparentKey match) are skipped with cursor movement so
// the screen content behind them stays visible, and rows end with cursor
// repositioning instead of a line break, so the image can be drawn at any
// cursor position.
package ansi
