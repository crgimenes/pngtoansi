# pngtoansi

Convert PNG image to ANSI art using UTF-8 characters.

For best results it is necessary to use a font compatible with characters "█", "▀", "▄", I recommend the [source code pro](https://github.com/adobe-fonts/source-code-pro) font or even better use the [3270font](https://github.com/rbanffy/3270font).

## Install

### Install as a Utility

```console
go install github.com/crgimenes/pngtoansi@latest
```

### Install as a Golang package

```console
go get github.com/crgimenes/pngtoansi/ansi
```

## Examples

### Convert PNG to ANSI in the terminal

```console
pngtoansi -f ./examples/gopher.png
```

Adjusted the background color. It is possible to change the color used to replace the transparent background using the *-rgb* parameter.

```console
pngtoansi -f ./examples/test-01.png -rgb FFFFFF
```

Reading from stdin is supported with `-f -`:

```console
magick photo.jpg -resize 160x png:- | pngtoansi -f -
```

### Sprite mode

With `-sprite` the output is relocatable: transparent cells (alpha, or a color
given with `-transparent`) are skipped with cursor movement so whatever is
already on screen stays visible, and rows end with cursor repositioning
instead of a line break. The image can then be drawn at any position:

```console
pngtoansi -f sprite.png -sprite > sprite.ans
printf '\033[10;40H'; cat sprite.ans
```

```console
pngtoansi -f logo.png -sprite -transparent 000000
```

### Golang example

```golang
import "github.com/crgimenes/pngtoansi/ansi"

p := ansi.New()
err := p.PrintFile("./examples/gopher.png", "FFFFFF")
if err != nil {
	fmt.Println(err)
	return
}
```

## Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b fantastic-feature`.
- Then _commit_ your changes `git commit -m 'Implementation of new fantastic feature'`
- Make a _push_ to your _branch_ `git push origin fantastic-feature`.
- Submit a **Pull Request** so that we can review your changes

