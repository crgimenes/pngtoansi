package main

import (
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/crgimenes/pngtoansi/ansi"
)

func usage(w io.Writer) {
	_, _ = fmt.Fprint(w, `pngtoansi converts a PNG image to ANSI art on stdout.

usage: pngtoansi -f <file.png> [options]

  -f string
        path to the PNG file to convert, or "-" for stdin (required)
  -o string
        output format: "ansi" or "xpm" (default "ansi")
  -rgb string
        hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels
  -sprite
        relocatable output: transparent cells are skipped with cursor
        movement and rows end with cursor repositioning instead of a
        line break, so the image can be drawn at any cursor position
  -transparent string
        hexadecimal RGB color treated as transparent
        (requires -sprite or -o xpm)
  -h    show this help

examples:
  pngtoansi -f examples/gopher.png -rgb FFFFFF
  printf '\033[10;40H'; pngtoansi -f sprite.png -sprite
  pngtoansi -f sprite.png -o xpm > sprite.xpm
`)
}

func main() {
	var help bool
	fileName := flag.String("f", "", "path to the PNG file to convert, or \"-\" for stdin (required)")
	output := flag.String("o", "ansi", "output format: \"ansi\" or \"xpm\"")
	rgb := flag.String("rgb", "", "hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels")
	sprite := flag.Bool("sprite", false, "relocatable output with transparent cells skipped")
	transparent := flag.String("transparent", "", "hexadecimal RGB color treated as transparent (requires -sprite or -o xpm)")
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

	if *output != "ansi" && *output != "xpm" {
		fmt.Fprintf(os.Stderr, "error: unknown output format %q\n", *output)
		usage(os.Stderr)
		os.Exit(2)
	}

	if *output == "xpm" && *sprite {
		fmt.Fprintln(os.Stderr, "error: -sprite only applies to -o ansi")
		os.Exit(2)
	}

	p := ansi.New()
	p.Sprite = *sprite

	if *transparent != "" {
		if !*sprite && *output != "xpm" {
			fmt.Fprintln(os.Stderr, "error: -transparent requires -sprite or -o xpm")
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

	err := run(p, *fileName, *rgb, *output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// cName derives a C identifier for the XPM array from the input file name.
func cName(fileName string) string {
	if fileName == "-" {
		return "image"
	}

	base := filepath.Base(fileName)
	ext := filepath.Ext(base)
	base = strings.TrimSuffix(base, ext)
	var b strings.Builder
	for _, r := range base {
		isWord := r == '_' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')
		if isWord {
			b.WriteRune(r)
			continue
		}
		b.WriteByte('_')
	}

	name := b.String()
	if name == "" || (name[0] >= '0' && name[0] <= '9') {
		name = "img_" + name
	}
	return name
}

func run(p *ansi.ImgToANSI, fileName, rgb, output string) error {
	if rgb != "" {
		err := p.SetRGB(rgb)
		if err != nil {
			return err
		}
	}

	var r io.Reader = os.Stdin
	if fileName != "-" {
		// #nosec G304 -- fileName is the path the user asked to convert.
		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error closing file:", err)
			}
		}()
		r = f
	}

	img, err := png.Decode(r)
	if err != nil {
		return err
	}

	if output == "xpm" {
		return p.FprintXPM(os.Stdout, img, cName(fileName))
	}
	return p.Print(img)
}
