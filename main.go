package main

import (
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/crgimenes/pngtoansi/ansi"
)

func usage(w io.Writer) {
	_, _ = fmt.Fprint(w, `pngtoansi converts a PNG image to ANSI art on stdout.

usage: pngtoansi -f <file.png> [options]

  -f string
        path to the PNG file to convert, or "-" for stdin (required)
  -rgb string
        hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels
  -sprite
        relocatable output: transparent cells are skipped with cursor
        movement and rows end with cursor repositioning instead of a
        line break, so the image can be drawn at any cursor position
  -transparent string
        hexadecimal RGB color treated as transparent (requires -sprite)
  -h    show this help

examples:
  pngtoansi -f examples/gopher.png -rgb FFFFFF
  printf '\033[10;40H'; pngtoansi -f sprite.png -sprite
`)
}

func main() {
	var help bool
	fileName := flag.String("f", "", "path to the PNG file to convert, or \"-\" for stdin (required)")
	rgb := flag.String("rgb", "", "hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels")
	sprite := flag.Bool("sprite", false, "relocatable output with transparent cells skipped")
	transparent := flag.String("transparent", "", "hexadecimal RGB color treated as transparent (requires -sprite)")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.BoolVar(&help, "help", false, "show this help")
	flag.Usage = func() { usage(os.Stderr) }
	flag.Parse()

	if help {
		usage(os.Stdout)
		return
	}

	if *fileName == "" {
		fmt.Fprintln(os.Stderr, "error: the -f flag is required")
		usage(os.Stderr)
		os.Exit(2)
	}

	p := ansi.New()
	p.Sprite = *sprite

	if *transparent != "" {
		if !*sprite {
			fmt.Fprintln(os.Stderr, "error: -transparent requires -sprite")
			usage(os.Stderr)
			os.Exit(2)
		}
		key, err := ansi.ParseRGB(*transparent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid -transparent color: %v\n", err)
			os.Exit(2)
		}
		p.TransparentKey = &key
	}

	err := run(p, *fileName, *rgb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(p *ansi.ImgToANSI, fileName, rgb string) error {
	if fileName != "-" {
		return p.PrintFile(fileName, rgb)
	}

	if rgb != "" {
		err := p.SetRGB(rgb)
		if err != nil {
			return err
		}
	}
	img, err := png.Decode(os.Stdin)
	if err != nil {
		return err
	}
	return p.Print(img)
}
