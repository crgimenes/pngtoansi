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
  -h    show this help

example:
  pngtoansi -f examples/gopher.png -rgb FFFFFF
`)
}

func main() {
	var help bool
	fileName := flag.String("f", "", "path to the PNG file to convert, or \"-\" for stdin (required)")
	rgb := flag.String("rgb", "", "hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels")
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
