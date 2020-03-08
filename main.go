package main

import (
	"fmt"
	"os"

	"github.com/crgimenes/goconfig"
	"github.com/crgimenes/pngtoansi/pngtoansi"
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

	p := pngtoansi.New()
	err = p.PrintFile(cfg.FileName, cfg.RGB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
