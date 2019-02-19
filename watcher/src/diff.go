package main

import (
	//"fmt"
	"github.com/mjibson/go-dsp/dsputils"
	//	"github.com/mjibson/go-dsp/fft"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	w = 1920
	h = 1080
)

func main() {
	var in [2]*os.File
	var tui8 uint8

	vol := dsputils.MakeEmptyMatrix([]int{w, h})

	in[0], _ = os.Open(os.Args[1])
	in[1], _ = os.Open(os.Args[2])
	defer in[0].Close()
	defer in[1].Close()

	src1, _, _ := image.Decode(in[0])
	src2, _, _ := image.Decode(in[1])

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, _, _, _ := src1.At(x, y).RGBA()
			r2, _, _, _ := src2.At(x, y).RGBA()
			rc := complex(math.Abs(float64(int32(r1)-int32(r2))), 0)
			vol.SetValue(rc, []int{x, y})
		}
	}

	/*bias := real(vol.Value([]int{0, 0}))/5 - real(save.Value([]int{0, 0}))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ravg := real(vol.Value([]int{x, y})) / 5
			rorig := real(save.Value([]int{x, y}))
			vol.SetValue(complex(math.Abs(ravg-rorig-bias), 0), []int{x, y})
		}
	}*/
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			tui8 = uint8(real(vol.Value([]int{x, y})))
			gray.Set(x, y, color.Gray{tui8})
		}
	}

	outfile, _ := os.Create(os.Args[3])
	defer outfile.Close()
	png.Encode(outfile, gray)

	// FFT
	//fftvol := fft.FFTN(vol)

	// transform
	//threashold := 100
	//for x := 0; x < w; x++ {
	//	for y := 0; y < h; y++ {
	//		if (x > threashold && x <= w-threashold-1) && (y > threashold && y < h-1-threashold) {
	//			fftvol.SetValue(complex(0, 0), []int{0, 0})
	//		}
	//	}
	//}

	// iFFT
	//vol = fft.IFFTN(fftvol)

}
