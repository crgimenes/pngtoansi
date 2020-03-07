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
	"io"
	"os"
)

const (
	fgColor string = "\033[38;2"
	bgColor string = "\033[48;2"
	reset   string = "\033[0m"
)

type RGB struct {
	R, G, B uint32
}

type ImgToANSI struct {
	DefaultColor RGB
}

func New() *ImgToANSI {
	return &ImgToANSI{}
}

func (p *ImgToANSI) Print(img image.Image) {
	p.Fprint(os.Stdout, img)
}

func (p *ImgToANSI) Fprint(w io.Writer, img image.Image) {
	var (
		fr, fg, fb, br, bg, bb       uint32
		lfr, lfg, lfb, lbr, lbg, lbb uint32
		fgCode, bgCode               string
		lastFgCode, lastBgCode       string
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
					fmt.Fprint(w, bgCode)
					fmt.Fprint(w, " ")
					continue
				}
				if lastBgCode == bgCode {
					fmt.Fprint(w, " ")
					continue
				}
				if lastFgCode == fgCode {
					fmt.Fprint(w, "█")
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
				fmt.Fprint(w, "▄")
				continue
			}
			//-=-=-=-=-=-=-=-=-=-
			if lastFgCode != fgCode {
				lastFgCode = fgCode
				lfr, lfg, lfb = fr, fg, fb
				fmt.Fprint(w, fgCode)
			}
			if lastBgCode != bgCode {
				lastBgCode = bgCode
				lbr, lbg, lbb = br, bg, bb
				fmt.Fprint(w, bgCode)
			}
			fmt.Fprint(w, "▀")
		}
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, reset)
}
