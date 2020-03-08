package imgtoansi

/*
\x1b[38;2;r;g;bm // fg
\x1b[48;2;r;g;bm // bg
\x1b[0m // reset

chars: "█", "▀", "▄", " "
*/

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
)

const (
	fgColor string = "\033[38;2"
	bgColor string = "\033[48;2"
	reset   string = "\033[0m"
)

// RGB color
type RGB struct {
	R, G, B uint32
}

// ImgToANSI holds module parameters
type ImgToANSI struct {
	DefaultColor RGB
}

// New create a new instance of ImgToANSI
func New() *ImgToANSI {
	return &ImgToANSI{}
}

// SetRGB update RGB values in current instance of ImgToANSI
func (p *ImgToANSI) SetRGB(rgb string) error {
	x, err := strconv.ParseUint(rgb, 16, 64)
	if err != nil {
		return err
	}

	r := uint8(x >> 16)
	g := uint8(x >> 8)
	b := uint8(x)

	p.DefaultColor.R = uint32(r)
	p.DefaultColor.G = uint32(g)
	p.DefaultColor.B = uint32(b)
	return nil
}

// PrintFile print a png file to the stdout using ANSI codes
func (p *ImgToANSI) PrintFile(fileName string, defaultRGB string) error {
	return p.FprintFile(os.Stdout, fileName, defaultRGB)
}

// FprintFile write a file to the Stdout using ANSI codes
func (p *ImgToANSI) FprintFile(w io.Writer, fileName string, defaultRGB string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	if defaultRGB != "" {
		err = p.SetRGB(defaultRGB)
		if err != nil {
			return err
		}
	}
	return p.Fprint(w, img)
}

// Print prints a image in the stdout using ANSI codes
func (p *ImgToANSI) Print(img image.Image) error {
	return p.Fprint(os.Stdout, img)
}

// Fprint prints write a image to a writer using ANSI codes
func (p *ImgToANSI) Fprint(w io.Writer, img image.Image) error {
	var (
		fr, fg, fb, br, bg, bb       uint32
		lfr, lfg, lfb, lbr, lbg, lbb uint32
		fgCode, bgCode               string
		lastFgCode, lastBgCode       string
		err                          error
	)

	bound := img.Bounds()

	fgCode = fmt.Sprintf("%v;%d;%d;%dm",
		fgColor,
		uint8(fr), uint8(fg), uint8(fb))
	bgCode = fmt.Sprintf("%v;%d;%d;%dm",
		bgColor,
		uint8(br), uint8(bg), uint8(bb))

	for y := bound.Min.Y; y < bound.Max.Y; y += 2 {
		for x := bound.Min.X; x < bound.Max.X; x++ {
			px := img.At(x, y)
			r, g, b, a := px.RGBA()
			if a == 0 {
				r = p.DefaultColor.R
				g = p.DefaultColor.G
				b = p.DefaultColor.B
			}
			if fr != r ||
				fg != g ||
				fb != b {
				fr, fg, fb = r, g, b
				fgCode = fmt.Sprintf("%v;%d;%d;%dm",
					fgColor,
					uint8(r), uint8(g), uint8(b))
			}

			px = img.At(x, y+1)
			r, g, b, a = px.RGBA()
			if a == 0 {
				r = p.DefaultColor.R
				g = p.DefaultColor.G
				b = p.DefaultColor.B
			}
			if br != r ||
				bg != g ||
				bb != b {
				br, bg, bb = r, g, b
				bgCode = fmt.Sprintf("%v;%d;%d;%dm",
					bgColor,
					uint8(r), uint8(g), uint8(b))
			}

			//-=-=-=-=-=-=-=-=-=-
			if fr == br &&
				fg == bg &&
				fb == bb {
				if lastBgCode != bgCode &&
					lastFgCode != fgCode {
					lastBgCode = bgCode
					lbr, lbg, lbb = br, bg, bb
					_, err = fmt.Fprint(w, bgCode, " ")
					if err != nil {
						return err
					}
					continue
				}
				if lastBgCode == bgCode {
					_, err = fmt.Fprint(w, " ")
					if err != nil {
						return err
					}
					continue
				}
				if lastFgCode == fgCode {
					_, err = fmt.Fprint(w, "█")
					if err != nil {
						return err
					}
					continue
				}
			}
			//-=-=-=-=-=-=-=-=-=-
			if lbr == fr &&
				lbg == fg &&
				lbb == fg &&
				lfr == br &&
				lfg == bg &&
				lfb == bg &&
				lastFgCode != "" &&
				lastBgCode != "" {
				_, err = fmt.Fprint(w, "▄")
				if err != nil {
					return err
				}
				continue
			}
			//-=-=-=-=-=-=-=-=-=-
			if lastFgCode != fgCode {
				lastFgCode = fgCode
				lfr, lfg, lfb = fr, fg, fb
				_, err = fmt.Fprint(w, fgCode)
				if err != nil {
					return err
				}
			}
			if lastBgCode != bgCode {
				lastBgCode = bgCode
				lbr, lbg, lbb = br, bg, bb
				_, err = fmt.Fprint(w, bgCode)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprint(w, "▀")
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(w, "")
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(w, reset)
	return err
}
