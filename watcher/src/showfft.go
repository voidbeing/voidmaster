package main

import (
	//"fmt"
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	"image"
	"image/color"
	"image/png"
	//"math"
	"os"
)

func main() {
	var in *os.File

	in, _ = os.Open(os.Args[1])
	defer in.Close()

	conf, _ := png.DecodeConfig(in)
	w := conf.Width
	h := conf.Height

	in.Seek(0, 0)
	src, _, _ := image.Decode(in)
	vol := dsputils.MakeEmptyMatrix([]int{w, h})

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, _ := src.At(x, y).RGBA()
			rc := complex(float64(r), 0)
			vol.SetValue(rc, []int{x, y})
		}
	}

	// FFT
	fftorig := fft.FFTN(vol)

	// transform
	//threashold := 100
	//for x := 0; x < w; x++ {
	//	for y := 0; y < h; y++ {
	//		if (x > threashold && x <= w-threashold-1) && (y > threashold && y < h-1-threashold) {
	//			fftorig.SetValue(complex(0, 0), []int{0, 0})
	//		}
	//	}
	//}

	maxr := 0.0
	maxi := 0.0
	minr := 0.0
	mini := 0.0

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fftr := real(fftorig.Value([]int{x, y}))
			ffti := imag(fftorig.Value([]int{x, y}))
			if maxr < fftr {
				maxr = fftr
			}
			if minr > fftr {
				minr = fftr
			}
			if maxi < ffti {
				maxi = ffti
			}
			if mini > ffti {
				mini = ffti
			}
		}
	}

	linear := func(val, a, b float64) uint8 {
		return uint8(val * 255 / (a - b))
	}

	fftimg := image.NewRGBA(image.Rect(0, 0, w*2, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fftr := linear(real(fftorig.Value([]int{x, y})), maxr, minr)
			ffti := linear(imag(fftorig.Value([]int{x, y})), maxi, mini)
			fftimg.Set(x, y, color.RGBA{fftr, ffti, ffti, 255})
			fftimg.Set(x+w, y, src.At(x, y))
		}
	}

	outfile, _ := os.Create(os.Args[2])
	defer outfile.Close()
	png.Encode(outfile, fftimg)

	// iFFT
	//vol = fft.IFFTN(fftorig)
}
