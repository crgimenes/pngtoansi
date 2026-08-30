package ansi

/*
SGR escape sequences used by the renderer:

	\033[38;2;r;g;bm   24-bit foreground color
	\033[48;2;r;g;bm   24-bit background color
	\033[<fg>;<bg>m    4-bit VGA foreground/background color
	\033[m             reset all attributes

Cell glyphs: " ", "█", "▀", "▄" — all present in CP437, so the output can be
transcoded for classic terminals. The glyph for each cell is chosen to reuse
the colors already set, keeping the output small.
*/

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// resetln resets all SGR attributes and terminates the current terminal line.
var resetln = []byte("\033[m\r\n")

// RGB holds an 8-bit-per-channel color.
type RGB struct {
	R, G, B uint8
}

// ImgToANSI holds the conversion parameters.
type ImgToANSI struct {
	// DefaultColor is the opaque background that transparent and partially
	// transparent pixels are composited over.
	DefaultColor RGB
}

// New creates a new ImgToANSI instance.
func New() *ImgToANSI {
	return &ImgToANSI{}
}

// closer closes c and reports any failure on stderr so that it never corrupts
// the ANSI output written to stdout.
func closer(c io.Closer) {
	err := c.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error closing file:", err)
	}
}

// SetRGB sets DefaultColor from a hexadecimal RGB string such as "FFFFFF".
func (p *ImgToANSI) SetRGB(rgb string) error {
	x, err := strconv.ParseUint(rgb, 16, 24)
	if err != nil {
		return err
	}

	p.DefaultColor.R = uint8((x >> 16) & 0xff)
	p.DefaultColor.G = uint8((x >> 8) & 0xff)
	p.DefaultColor.B = uint8(x & 0xff)
	return nil
}

// PrintFile writes the ANSI rendering of the PNG file at fileName to stdout.
func (p *ImgToANSI) PrintFile(fileName, defaultRGB string) error {
	return p.FprintFile(os.Stdout, fileName, defaultRGB)
}

// FprintFile writes the ANSI rendering of the PNG file at fileName to w. The
// default color is validated before the file is opened so invalid input fails
// fast.
func (p *ImgToANSI) FprintFile(w io.Writer, fileName, defaultRGB string) error {
	if defaultRGB != "" {
		err := p.SetRGB(defaultRGB)
		if err != nil {
			return err
		}
	}

	// #nosec G304 -- fileName is the path the caller asked to render; reading it
	// is the sole purpose of this function.
	f, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return err
	}
	defer closer(f)

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	return p.Fprint(w, img)
}

// Print writes the ANSI rendering of img to stdout.
func (p *ImgToANSI) Print(img image.Image) error {
	return p.Fprint(os.Stdout, img)
}

// pxColor returns the 8-bit color of the pixel at (x, y) composited over the
// opaque DefaultColor using the Porter-Duff "over" operator. Fully opaque
// pixels keep their own color, fully transparent pixels (including coordinates
// outside the image bounds) become DefaultColor, and partially transparent
// pixels blend the two so anti-aliased edges follow the chosen background.
func (p *ImgToANSI) pxColor(x, y int, img image.Image) RGB {
	cr, cg, cb, ca := img.At(x, y).RGBA()
	return RGB{
		R: overChannel(cr, ca, p.DefaultColor.R),
		G: overChannel(cg, ca, p.DefaultColor.G),
		B: overChannel(cb, ca, p.DefaultColor.B),
	}
}

// overChannel composites one premultiplied 16-bit source channel over an
// opaque 8-bit background channel and returns the 8-bit result. RGBA channels
// are in [0, 0xffff] and src <= alpha, so the result never exceeds 0xff.
func overChannel(src, alpha uint32, bg uint8) uint8 {
	// bg scaled to 16 bits; the background is opaque, so it contributes the
	// fraction of the cell the source does not cover (0xffff - alpha).
	bg16 := uint32(bg) * 0x101
	out := src + bg16*(0xffff-alpha)/0xffff
	return uint8((out >> 8) & 0xff)
}

/*
4-bit VGA ANSI color palette.

	Name            fg   bg   RGB
	Black           30   40   0,0,0
	Red             31   41   170,0,0
	Green           32   42   0,170,0
	Yellow          33   43   170,85,0
	Blue            34   44   0,0,170
	Magenta         35   45   170,0,170
	Cyan            36   46   0,170,170
	White           37   47   170,170,170
	Bright Black    90   100  85,85,85
	Bright Red      91   101  255,85,85
	Bright Green    92   102  85,255,85
	Bright Yellow   93   103  255,255,85
	Bright Blue     94   104  85,85,255
	Bright Magenta  95   105  255,85,255
	Bright Cyan     96   106  85,255,255
	Bright White    97   107  255,255,255

Background codes are always the matching foreground code plus 10.
*/

// vgaColor maps an exact 8-bit RGB triple to its 4-bit VGA foreground SGR code
// (30-37 or 90-97). It reports ok=false when the color is not part of the
// 16-color VGA palette.
func vgaColor(r, g, b uint8) (code int, ok bool) {
	switch {
	case r == 0 && g == 0 && b == 0:
		return 30, true
	case r == 170 && g == 0 && b == 0:
		return 31, true
	case r == 0 && g == 170 && b == 0:
		return 32, true
	case r == 170 && g == 85 && b == 0:
		return 33, true
	case r == 0 && g == 0 && b == 170:
		return 34, true
	case r == 170 && g == 0 && b == 170:
		return 35, true
	case r == 0 && g == 170 && b == 170:
		return 36, true
	case r == 170 && g == 170 && b == 170:
		return 37, true
	case r == 85 && g == 85 && b == 85:
		return 90, true
	case r == 255 && g == 85 && b == 85:
		return 91, true
	case r == 85 && g == 255 && b == 85:
		return 92, true
	case r == 255 && g == 255 && b == 85:
		return 93, true
	case r == 85 && g == 85 && b == 255:
		return 94, true
	case r == 255 && g == 85 && b == 255:
		return 95, true
	case r == 85 && g == 255 && b == 255:
		return 96, true
	case r == 255 && g == 255 && b == 255:
		return 97, true
	default:
		return 0, false
	}
}

// RGB2VGAFg returns the 4-bit VGA foreground SGR code for an exact palette
// match, or ok=false when the color is not a VGA color.
func RGB2VGAFg(r, g, b uint8) (code int, ok bool) {
	return vgaColor(r, g, b)
}

// RGB2VGABg returns the 4-bit VGA background SGR code for an exact palette
// match, or ok=false when the color is not a VGA color.
func RGB2VGABg(r, g, b uint8) (code int, ok bool) {
	code, ok = vgaColor(r, g, b)
	if !ok {
		return 0, false
	}
	return code + 10, true
}

// sgrState tracks the colors the terminal currently has set, so that runs of
// identical colors cost no escape sequences. The zero value means "unknown",
// the state right after a reset.
type sgrState struct {
	fg, bg       RGB
	fgSet, bgSet bool
}

func (s *sgrState) hasFg(c RGB) bool {
	return s.fgSet && s.fg == c
}

func (s *sgrState) hasBg(c RGB) bool {
	return s.bgSet && s.bg == c
}

// appendSGR appends a single CSI sequence setting the foreground and/or
// background color. Each part uses the compact 4-bit code when the color is an
// exact VGA palette match and the 24-bit form otherwise.
func appendSGR(buf []byte, setFg bool, fg RGB, setBg bool, bg RGB) []byte {
	if !setFg && !setBg {
		return buf
	}

	buf = append(buf, 0x1b, '[')
	if setFg {
		code, ok := RGB2VGAFg(fg.R, fg.G, fg.B)
		if ok {
			buf = strconv.AppendInt(buf, int64(code), 10)
		} else {
			buf = append(buf, "38;2;"...)
			buf = appendRGB(buf, fg)
		}
	}
	if setBg {
		if setFg {
			buf = append(buf, ';')
		}
		code, ok := RGB2VGABg(bg.R, bg.G, bg.B)
		if ok {
			buf = strconv.AppendInt(buf, int64(code), 10)
		} else {
			buf = append(buf, "48;2;"...)
			buf = appendRGB(buf, bg)
		}
	}
	return append(buf, 'm')
}

func appendRGB(buf []byte, c RGB) []byte {
	buf = strconv.AppendInt(buf, int64(c.R), 10)
	buf = append(buf, ';')
	buf = strconv.AppendInt(buf, int64(c.G), 10)
	buf = append(buf, ';')
	return strconv.AppendInt(buf, int64(c.B), 10)
}

// changes counts the SGR emissions a cell drawn with the given foreground and
// background would need under the current state.
func (s *sgrState) changes(fg, bg RGB) int {
	n := 0
	if !s.hasFg(fg) {
		n++
	}
	if !s.hasBg(bg) {
		n++
	}
	return n
}

// Fprint writes the ANSI rendering of img to w.
//
// Each output cell encodes two vertically adjacent pixels. The glyph is chosen
// to minimize escape output: a cell whose two pixels share one color becomes a
// space (background only) or "█" (foreground only), a split cell becomes "▀"
// or "▄", whichever reuses more of the colors already set. Colors are emitted
// as one merged SGR sequence, only when they change, using compact 4-bit codes
// for exact VGA palette matches and 24-bit truecolor otherwise. Rows are
// consumed two pixels at a time, so for images with an odd height the missing
// bottom row falls back to DefaultColor.
func (p *ImgToANSI) Fprint(w io.Writer, img image.Image) error {
	bw := bufio.NewWriter(w)
	bound := img.Bounds()
	buf := make([]byte, 0, 64)

	for y := bound.Min.Y; y < bound.Max.Y; y += 2 {
		var st sgrState
		for x := bound.Min.X; x < bound.Max.X; x++ {
			top := p.pxColor(x, y, img)
			bot := p.pxColor(x, y+1, img)

			var glyph string
			var setFg, setBg bool
			var fg, bg RGB
			switch {
			case top == bot && st.hasBg(top):
				glyph = " "
			case top == bot && st.hasFg(top):
				glyph = "█"
			case top == bot:
				glyph = " "
				setBg = true
				bg = top
			case st.changes(bot, top) < st.changes(top, bot):
				// "▄" shows the background on top and the foreground
				// at the bottom, reusing more of the current state
				// than "▀" would.
				glyph = "▄"
				fg = bot
				bg = top
				setFg = !st.hasFg(fg)
				setBg = !st.hasBg(bg)
			default:
				glyph = "▀"
				fg = top
				bg = bot
				setFg = !st.hasFg(fg)
				setBg = !st.hasBg(bg)
			}

			buf = appendSGR(buf[:0], setFg, fg, setBg, bg)
			buf = append(buf, glyph...)
			// Error intentionally discarded: bufio.Writer records the
			// first write error and replays it from Flush below.
			_, _ = bw.Write(buf)

			if setFg {
				st.fg = fg
				st.fgSet = true
			}
			if setBg {
				st.bg = bg
				st.bgSet = true
			}
		}
		_, _ = bw.Write(resetln)
	}

	// A single Flush check covers every buffered write above.
	return bw.Flush()
}
