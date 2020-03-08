package imgtoansi

import (
	"bytes"
	"image"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *ImgToANSI
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
				t.Errorf("ImgToANSI.FprintFile() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestImgToANSI_Print(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImgToANSI{
				DefaultColor: tt.fields.DefaultColor,
			}
			p.Print(tt.args.img)
		})
	}
}

func TestImgToANSI_Fprint(t *testing.T) {
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
		// TODO: Add test cases.
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
