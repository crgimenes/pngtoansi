//Copyright (c) 2020 Cesar Gimenes - MIT License
//
// pngtogo receives an image as a parameter and generates the ANSI
// code to recreate the same image on the terminal using UTF-8 characters.
//
// Every two lines of pixels in the original image generate a line in the
// image on the terminal. If the original image has an odd number of lines,
// the last line of the final image will use the default color to generate
// an extra line.

package main
