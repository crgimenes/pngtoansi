package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
)

/*
\x1b[38;2;r;g;bm // fg
\x1b[48;2;r;g;bm // bg
\x1b[0m // reset

chars: "█", "▀", "▄", " "
*/

const (
	fgColor string = "\033[38;2"
	bgColor string = "\033[48;2"
	reset   string = "\033[0m"
)

var (
	fr, fg, fb, br, bg, bb       uint32
	lfr, lfg, lfb, lbr, lbg, lbb uint32
	fgCode, bgCode               string
	lastFgCode, lastBgCode       string
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("pngtoansi usage: png2ansi <pngfile>")
		os.Exit(0)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
				r, g, b = 0xff, 0xff, 0xff
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
				r, g, b = 0xff, 0xff, 0xff
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
					fmt.Print(bgCode)
				}
				if lastBgCode == bgCode {
					fmt.Print(" ")
					continue
				}
				if lastFgCode == fgCode {
					fmt.Print("█")
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
				fmt.Print("▄")
				continue
			}
			//-=-=-=-=-=-=-=-=-=-
			if lastFgCode != fgCode {
				lastFgCode = fgCode
				lfr, lfg, lfb = fr, fg, fb
				fmt.Print(fgCode)
			}
			if lastBgCode != bgCode {
				lastBgCode = bgCode
				lbr, lbg, lbb = br, bg, bb
				fmt.Print(bgCode)
			}
			fmt.Print("▀")
		}
		fmt.Println("")
	}
	fmt.Println(reset)
}
