package main

import (
	"fmt"
	"image/png"
	"os"
	"strconv"

	"pngtoansi/imgtoansi"

	"github.com/crgimenes/goconfig"
)

type config struct {
	FileName string `cfg:"f" cfgRequired:"true"`
	RGB      string `cfg:"rgb"`
}

func main() {
	cfg := config{}
	goconfig.PrefixEnv = "PNGTOANSI"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(cfg.FileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := imgtoansi.New()

	if cfg.RGB != "" {
		rgb, err := strconv.ParseUint(cfg.RGB, 16, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		r := uint8(rgb >> 16)
		g := uint8(rgb >> 8)
		b := uint8(rgb)

		p.DefaultColor.R = uint32(r)
		p.DefaultColor.G = uint32(g)
		p.DefaultColor.B = uint32(b)
	}
	p.Print(img, os.Stdout)
}
