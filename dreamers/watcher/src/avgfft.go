package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
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
	var in [5]*os.File
	var tc128 complex128
	var tui8 uint8

	vol := dsputils.MakeEmptyMatrix([]int{w, h})
	save := dsputils.MakeEmptyMatrix([]int{w, h})

	// Get a volume of 5 FullHD images
	for i := 0; i < 5; i++ {
		in[i], _ = os.Open(fmt.Sprintf("%d", i+1) + ".png")
		defer in[i].Close()

		src, _, _ := image.Decode(in[i])

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, _, _, _ := src.At(x, y).RGBA()
				rc := complex(real(vol.Value([]int{x, y}))+float64(r), 0)
				vol.SetValue(rc, []int{x, y})
				if i == 0 {
					save.SetValue(tc128, []int{x, y})
				}
			}
		}
	}

	bias := real(vol.Value([]int{0, 0}))/5 - real(save.Value([]int{0, 0}))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ravg := real(vol.Value([]int{x, y})) / 5
			rorig := real(save.Value([]int{x, y}))
			vol.SetValue(complex(math.Abs(ravg-rorig-bias), 0), []int{x, y})
		}
	}
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			tui8 = uint8(real(vol.Value([]int{x, y})))
			gray.Set(x, y, color.Gray{tui8})
		}
	}

	// Output the volume file
	outfile, _ := os.Create(os.Args[1])
	defer outfile.Close()
	png.Encode(outfile, gray)

	// FFT
	fft.SetWorkerPoolSize(3)
	fftvol := fft.FFTN(vol)

	// transform
	threashold := 100
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if (x > threashold && x <= w-threashold-1) && (y > threashold && y < h-1-threashold) {
				fftvol.SetValue(complex(0, 0), []int{0, 0})
			}
		}
	}

	// iFFT
	vol = fft.IFFTN(fftvol)

	gray2 := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			tui8 = uint8(real(vol.Value([]int{x, y})))
			gray2.Set(x, y, color.Gray{tui8})
		}
	}

	// Output the volume file
	outfile2, _ := os.Create(os.Args[2])
	defer outfile2.Close()
	png.Encode(outfile2, gray2)
}
