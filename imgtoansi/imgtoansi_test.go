package imgtoansi

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *ImgToANSI
	}{
		{
			name: "success",
			want: &ImgToANSI{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImgToANSI_SetRGB(t *testing.T) {
	type fields struct {
		DefaultColor RGB
	}
	type args struct {
		rgb string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				rgb: "FFFFFF",
			},
			fields: fields{
				DefaultColor: RGB{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			wantErr: false,
		},
		{
			name: "success 2",
			args: args{
				rgb: "FFFF",
			},
			fields: fields{
				DefaultColor: RGB{
					R: 255,
					G: 255,
					B: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "success 3",
			args: args{
				rgb: "FF",
			},
			fields: fields{
				DefaultColor: RGB{
					R: 255,
					G: 0,
					B: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				rgb: "not hexa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImgToANSI{
				DefaultColor: tt.fields.DefaultColor,
			}
			if err := p.SetRGB(tt.args.rgb); (err != nil) != tt.wantErr {
				t.Errorf("ImgToANSI.SetRGB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImgToANSI_PrintFile(t *testing.T) {
	type fields struct {
		DefaultColor RGB
	}
	type args struct {
		fileName   string
		defaultRGB string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				fileName:   "../examples/debian.png",
				defaultRGB: "",
			},
			wantErr: false,
		},
		{
			name: "success 2",
			args: args{
				fileName:   "../examples/debian.png",
				defaultRGB: "FFFFFF",
			},
			wantErr: false,
		},
		{
			name: "success 3",
			args: args{
				fileName:   "../examples/test-01.png",
				defaultRGB: "FFFFFF",
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				fileName:   "../examples/test-01.png",
				defaultRGB: "not hexa",
			},
			wantErr: true,
		},
		{
			name: "error 2",
			args: args{
				fileName: "file error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImgToANSI{
				DefaultColor: tt.fields.DefaultColor,
			}
			if err := p.PrintFile(tt.args.fileName, tt.args.defaultRGB); (err != nil) != tt.wantErr {
				t.Errorf("ImgToANSI.PrintFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImgToANSI_FprintFile(t *testing.T) {
	type fields struct {
		DefaultColor RGB
	}
	type args struct {
		fileName   string
		defaultRGB string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				fileName: "../examples/test-02.png",
			},
			wantW: "\x1b[48;2;0;0;0m \x1b[38;2;255;255;255m▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀\n\x1b[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImgToANSI{
				DefaultColor: tt.fields.DefaultColor,
			}
			w := &bytes.Buffer{}
			if err := p.FprintFile(w, tt.args.fileName, tt.args.defaultRGB); (err != nil) != tt.wantErr {
				t.Errorf("ImgToANSI.FprintFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ImgToANSI.FprintFile() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}

func TestImgToANSI_Fprint(t *testing.T) {
	f, err := os.Open("../examples/test-02.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		DefaultColor RGB
	}
	type args struct {
		img image.Image
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantW  string
	}{
		{
			name: "success",
			args: args{
				img: img,
			},
			wantW: "\x1b[48;2;0;0;0m \x1b[38;2;255;255;255m▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀\n\x1b[0m\n",
		},
		{
			name: "success 2",
			args: args{
				img: img,
			},
			fields: fields{
				DefaultColor: RGB{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			wantW: "\x1b[48;2;0;0;0m \x1b[38;2;255;255;255m▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀\n\x1b[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImgToANSI{
				DefaultColor: tt.fields.DefaultColor,
			}
			w := &bytes.Buffer{}
			p.Fprint(w, tt.args.img)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ImgToANSI.Fprint() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
