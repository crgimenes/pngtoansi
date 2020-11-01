package pngtoansi

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
	fgColor   = "\033[38;2"
	bgColor   = "\033[48;2"
	reset     = "\033[0m"
	ansiColor = "%v;%d;%d;%dm"
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

func closer(c io.Closer) {
	err := c.Close()
	if err != nil {
		fmt.Println("error closing file", err)
	}
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
	defer closer(f)

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

func (p *ImgToANSI) pxColor(x, y int, img image.Image) (r, g, b uint32) {
	px := img.At(x, y)
	r, g, b, a := px.RGBA()
	if a == 0 {
		r = p.DefaultColor.R
		g = p.DefaultColor.G
		b = p.DefaultColor.B
	}
	return r, g, b
}

// Fprint prints write a image to a writer using ANSI codes
func (p *ImgToANSI) Fprint(w io.Writer, img image.Image) error {
	bound := img.Bounds()

	var (
		fgCode     string
		bgCode     string
		lastBgCode string
		lastFgCode string
		err        error
	)

	for y := bound.Min.Y; y < bound.Max.Y; y += 2 {
		for x := bound.Min.X; x < bound.Max.X; x++ {

			fr, fg, fb := p.pxColor(x, y, img)
			fgCode = fmt.Sprintf(ansiColor,
				fgColor,
				uint8(fr), uint8(fg), uint8(fb))

			br, bg, bb := p.pxColor(x, y+1, img)
			bgCode = fmt.Sprintf(ansiColor,
				bgColor,
				uint8(br), uint8(bg), uint8(bb))

			//-=-=-=-=-=-=-=-=-=-=-=-=

			if lastBgCode != bgCode {
				_, err = fmt.Fprint(w, bgCode)
				if err != nil {
					return err
				}
				lastBgCode = bgCode
			}

			//-=-=-=-=-=-=-=-=-=-=-=-=-=-=
			if fr == br &&
				fg == bg &&
				fb == bb {
				_, err = fmt.Fprint(w, " ")
				if err != nil {
					return err
				}
				continue
			}
			//-=-=-=-=-=-=-=-=-=-=-=-=-=-=

			if lastFgCode != fgCode {
				_, err = fmt.Fprint(w, fgCode)
				if err != nil {
					return err
				}
				lastFgCode = fgCode
			}

			_, err = fmt.Fprint(w, "▀")
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(w, reset)
		if err != nil {
			return err
		}
		lastFgCode = ""
		lastBgCode = ""
	}
	return nil
}
