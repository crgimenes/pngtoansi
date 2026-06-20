package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/crgimenes/pngtoansi"
)

func main() {
	fileName := flag.String("f", "", "path to the PNG file to convert (required)")
	rgb := flag.String("rgb", "", "hexadecimal RGB color (e.g. FFFFFF) used for transparent pixels")
	flag.Parse()

	if *fileName == "" {
		fmt.Fprintln(os.Stderr, "error: the -f flag is required")
		flag.Usage()
		os.Exit(2)
	}

	p := pngtoansi.New()
	err := p.PrintFile(*fileName, *rgb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
