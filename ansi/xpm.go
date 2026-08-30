// Copyright (c) 2020 Cesar Gimenes - MIT License

package ansi

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

// xpmCharset names palette entries in assignment order. It contains only
// characters that are safe inside a C string literal; space is reserved for
// the transparent entry, so pixel art stays readable in the file body.
const xpmCharset = ".+@#$%&*=-;:>,<" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

// xpmCode returns the palette code for entry i, cpp characters long, as a
// base-len(xpmCharset) number with the most significant character first.
func xpmCode(i, cpp int) string {
	code := make([]byte, cpp)
	for j := cpp - 1; j >= 0; j-- {
		code[j] = xpmCharset[i%len(xpmCharset)]
		i /= len(xpmCharset)
	}
	return string(code)
}

// FprintXPM writes img to w in XPM3 format, one palette entry per distinct
// color. Pixels follow the transparency rules of transparentPx (alpha below
// 50%, or an exact TransparentKey match) and become the None entry, written
// as spaces; opaque pixels are composited over DefaultColor. The name must be
// a valid C identifier, as XPM doubles as an includable C source file.
func (p *ImgToANSI) FprintXPM(w io.Writer, img image.Image, name string) error {
	bound := img.Bounds()

	// First pass: collect distinct colors in scan order, so the output is
	// deterministic.
	index := map[RGB]int{}
	var colors []RGB
	hasTransparent := false
	for y := bound.Min.Y; y < bound.Max.Y; y++ {
		for x := bound.Min.X; x < bound.Max.X; x++ {
			if p.transparentPx(x, y, img) {
				hasTransparent = true
				continue
			}
			c := p.pxColor(x, y, img)
			_, seen := index[c]
			if seen {
				continue
			}
			index[c] = len(colors)
			colors = append(colors, c)
		}
	}

	cpp := 1
	for n := len(xpmCharset); n < len(colors); n *= len(xpmCharset) {
		cpp++
	}

	entries := len(colors)
	if hasTransparent {
		entries++
	}

	bw := bufio.NewWriter(w)
	// Errors intentionally discarded: bufio.Writer records the first write
	// error and replays it from Flush below.
	_, _ = fmt.Fprintf(bw, "/* XPM */\nstatic const char *%s[] = {\n", name)
	_, _ = fmt.Fprintf(bw, "\"%d %d %d %d\",\n",
		bound.Dx(), bound.Dy(), entries, cpp)
	if hasTransparent {
		_, _ = fmt.Fprintf(bw, "\"%s\tc None\",\n", spaces(cpp))
	}
	for i, c := range colors {
		_, _ = fmt.Fprintf(bw, "\"%s\tc #%02X%02X%02X\",\n",
			xpmCode(i, cpp), c.R, c.G, c.B)
	}

	for y := bound.Min.Y; y < bound.Max.Y; y++ {
		_, _ = bw.WriteString("\"")
		for x := bound.Min.X; x < bound.Max.X; x++ {
			if p.transparentPx(x, y, img) {
				_, _ = bw.WriteString(spaces(cpp))
				continue
			}
			_, _ = bw.WriteString(xpmCode(index[p.pxColor(x, y, img)], cpp))
		}
		sep := ",\n"
		if y == bound.Max.Y-1 {
			sep = "};\n"
		}
		_, _ = bw.WriteString("\"" + sep)
	}

	return bw.Flush()
}

func spaces(n int) string {
	return strings.Repeat(" ", n)
}
