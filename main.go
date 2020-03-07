package main

import (
	"fmt"
	"image/png"
	"os"

	"pngtoansi/imgtoansi"

	"github.com/crgimenes/goconfig"
)

type config struct {
	FileName string `cfg:"f" cfgRequired:"true"`
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
	p.Print(img)
}
