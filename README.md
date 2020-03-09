# pngtoansi

Convert PNG image to ANSI art using UTF-8 characters.


For best results it is necessary to use a font compatible with characters "█", "▀", "▄", I recommend the [source code pro](https://github.com/adobe-fonts/source-code-pro) fount.


## Install

### Install as a Utility

```console
go install github.com/crgimenes/pngtoansi/cmd/pngtoansi 
```

### Install as a Golang package

```console
go get github.com/crgimenes/pngtoansi
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

### Golang example

```golang
...
p := imgtoansi.New()
err = p.PrintFile("./examples/gopher.png", "FFFFFF")
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

