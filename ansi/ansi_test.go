package ansi

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"testing"
)

// resetSeq is the SGR reset plus line terminator emitted at the end of every
// rendered row.
const resetSeq = "\033[m\r\n"

func opaque(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}

func TestNew(t *testing.T) {
	got := New()
	if *got != (ImgToANSI{}) {
		t.Errorf("New() = %+v, want zero value", *got)
	}
}

func TestImgToANSI_SetRGB(t *testing.T) {
	tests := []struct {
		name    string
		rgb     string
		want    RGB
		wantErr bool
	}{
		{name: "six digits", rgb: "FFFFFF", want: RGB{R: 255, G: 255, B: 255}},
		{name: "lowercase", rgb: "0a141e", want: RGB{R: 10, G: 20, B: 30}},
		{name: "black", rgb: "000000", want: RGB{R: 0, G: 0, B: 0}},
		{name: "four digits zero padded", rgb: "FFFF", want: RGB{R: 0, G: 255, B: 255}},
		{name: "two digits zero padded", rgb: "FF", want: RGB{R: 0, G: 0, B: 255}},
		{name: "not hexadecimal", rgb: "not hexa", wantErr: true},
		{name: "out of range", rgb: "FFFFFFF", wantErr: true},
		{name: "empty", rgb: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			err := p.SetRGB(tt.rgb)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetRGB(%q) error = %v, wantErr %v", tt.rgb, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if p.DefaultColor != tt.want {
				t.Errorf("SetRGB(%q) DefaultColor = %+v, want %+v", tt.rgb, p.DefaultColor, tt.want)
			}
		})
	}
}

func TestRGB2VGA(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		wantFg  int
		wantBg  int
		wantOk  bool
	}{
		{name: "black", r: 0, g: 0, b: 0, wantFg: 30, wantBg: 40, wantOk: true},
		{name: "vga red", r: 170, g: 0, b: 0, wantFg: 31, wantBg: 41, wantOk: true},
		{name: "light gray", r: 170, g: 170, b: 170, wantFg: 37, wantBg: 47, wantOk: true},
		{name: "bright black", r: 85, g: 85, b: 85, wantFg: 90, wantBg: 100, wantOk: true},
		{name: "bright white", r: 255, g: 255, b: 255, wantFg: 97, wantBg: 107, wantOk: true},
		{name: "not a palette color", r: 1, g: 2, b: 3, wantOk: false},
		{name: "near miss", r: 255, g: 0, b: 0, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg, okFg := RGB2VGAFg(tt.r, tt.g, tt.b)
			if okFg != tt.wantOk || (tt.wantOk && fg != tt.wantFg) {
				t.Errorf("RGB2VGAFg(%d,%d,%d) = %d,%v want %d,%v",
					tt.r, tt.g, tt.b, fg, okFg, tt.wantFg, tt.wantOk)
			}
			bg, okBg := RGB2VGABg(tt.r, tt.g, tt.b)
			if okBg != tt.wantOk || (tt.wantOk && bg != tt.wantBg) {
				t.Errorf("RGB2VGABg(%d,%d,%d) = %d,%v want %d,%v",
					tt.r, tt.g, tt.b, bg, okBg, tt.wantBg, tt.wantOk)
			}
		})
	}
}

// TestFprint pins the exact ANSI output for small, deterministic images. It
// covers the truecolor path, the 4-bit VGA path, the transparent default-color
// fallback, and the odd-height fallback.
func TestFprint(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (image.Image, *ImgToANSI)
		want  string
	}{
		{
			name: "truecolor cell",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 2))
				m.Set(0, 0, opaque(10, 20, 30))
				m.Set(0, 1, opaque(40, 50, 60))
				return m, New()
			},
			want: "\033[38;2;10;20;30;48;2;40;50;60m▀" + resetSeq,
		},
		{
			name: "vga palette cell",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 2))
				m.Set(0, 0, opaque(0, 0, 0))       // black -> fg 30
				m.Set(0, 1, opaque(255, 255, 255)) // bright white -> bg 107
				return m, New()
			},
			want: "\033[30;107m▀" + resetSeq,
		},
		{
			name: "transparent bottom uses default color",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 2))
				m.Set(0, 0, opaque(255, 0, 0)) // not a VGA color
				// bottom pixel left transparent (alpha 0)
				p := New()
				err := p.SetRGB("0000FF")
				if err != nil {
					t.Fatal(err)
				}
				return m, p
			},
			want: "\033[38;2;255;0;0;48;2;0;0;255m▀" + resetSeq,
		},
		{
			// Both halves end up bright white, so the whole cell is a
			// space over a 4-bit white background.
			name: "transparent default matches vga palette",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 2))
				m.Set(0, 0, opaque(255, 255, 255))
				// bottom transparent -> default white
				p := New()
				err := p.SetRGB("FFFFFF")
				if err != nil {
					t.Fatal(err)
				}
				return m, p
			},
			want: "\033[107m " + resetSeq,
		},
		{
			name: "semi-transparent pixel blends over default",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 2))
				// 50% black over white default -> mid gray (127,127,127).
				m.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 128})
				m.Set(0, 1, opaque(0, 0, 0)) // VGA black, but fg is not VGA
				p := New()
				err := p.SetRGB("FFFFFF")
				if err != nil {
					t.Fatal(err)
				}
				return m, p
			},
			want: "\033[38;2;127;127;127;40m▀" + resetSeq,
		},
		{
			name: "odd height falls back to default for bottom row",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 1, 1))
				m.Set(0, 0, opaque(10, 20, 30))
				return m, New() // default color is black
			},
			want: "\033[38;2;10;20;30;40m▀" + resetSeq,
		},
		{
			name: "two columns render in order",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 2, 2))
				m.Set(0, 0, opaque(1, 2, 3))
				m.Set(1, 0, opaque(4, 5, 6))
				m.Set(0, 1, opaque(7, 8, 9))
				m.Set(1, 1, opaque(10, 11, 12))
				return m, New()
			},
			want: "\033[38;2;1;2;3;48;2;7;8;9m▀" +
				"\033[38;2;4;5;6;48;2;10;11;12m▀" + resetSeq,
		},
		{
			// A run of identical cells pays for its colors once.
			name: "repeated colors emit a single escape",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 3, 2))
				for x := range 3 {
					m.Set(x, 0, opaque(10, 20, 30))
					m.Set(x, 1, opaque(40, 50, 60))
				}
				return m, New()
			},
			want: "\033[38;2;10;20;30;48;2;40;50;60m▀▀▀" + resetSeq,
		},
		{
			// Second cell has the same colors vertically flipped: "▄"
			// reuses both and needs no escape at all.
			name: "flipped cell reuses state via lower half block",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 2, 2))
				m.Set(0, 0, opaque(10, 20, 30))
				m.Set(0, 1, opaque(40, 50, 60))
				m.Set(1, 0, opaque(40, 50, 60))
				m.Set(1, 1, opaque(10, 20, 30))
				return m, New()
			},
			want: "\033[38;2;10;20;30;48;2;40;50;60m▀▄" + resetSeq,
		},
		{
			// A uniform cell whose color is already the foreground is
			// drawn as "█" without touching the background.
			name: "uniform cell reuses foreground via full block",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 2, 2))
				m.Set(0, 0, opaque(10, 20, 30))
				m.Set(0, 1, opaque(40, 50, 60))
				m.Set(1, 0, opaque(10, 20, 30))
				m.Set(1, 1, opaque(10, 20, 30))
				return m, New()
			},
			want: "\033[38;2;10;20;30;48;2;40;50;60m▀█" + resetSeq,
		},
		{
			// Uniform runs become a background color and spaces.
			name: "uniform area becomes spaces",
			setup: func() (image.Image, *ImgToANSI) {
				m := image.NewRGBA(image.Rect(0, 0, 3, 2))
				for x := range 3 {
					m.Set(x, 0, opaque(10, 20, 30))
					m.Set(x, 1, opaque(10, 20, 30))
				}
				return m, New()
			},
			want: "\033[48;2;10;20;30m   " + resetSeq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, p := tt.setup()
			var buf bytes.Buffer
			err := p.Fprint(&buf, img)
			if err != nil {
				t.Fatalf("Fprint() error = %v", err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("Fprint() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestImgToANSI_FprintFile(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		defaultRGB string
		wantErr    bool
	}{
		{name: "success", fileName: "../examples/debian.png"},
		{name: "success with rgb", fileName: "../examples/debian.png", defaultRGB: "FFFFFF"},
		{name: "success with transparency", fileName: "../examples/test-01.png", defaultRGB: "FFFFFF"},
		{name: "invalid rgb", fileName: "../examples/test-01.png", defaultRGB: "not hexa", wantErr: true},
		{name: "missing file", fileName: "does-not-exist.png", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := New()
			err := p.FprintFile(&buf, tt.fileName, tt.defaultRGB)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FprintFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && buf.Len() == 0 {
				t.Errorf("FprintFile() produced no output")
			}
		})
	}
}

func ExampleImgToANSI_PrintFile() {
	p := New()
	err := p.PrintFile("../examples/gopher.png", "FFFFFF")
	if err != nil {
		fmt.Println(err)
	}
}
