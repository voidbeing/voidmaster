package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	"image"
	"image/color"
	"image/png"
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

	vol := dsputils.MakeEmptyMatrix([]int{5, w, h})

	// Get a volume of 5 FullHD images
	for i := 0; i < 5; i++ {
		in[i], _ = os.Open(fmt.Sprintf("%d", i+1) + ".png")
		defer in[i].Close()

		src, _, _ := image.Decode(in[i])

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, _, _, _ := src.At(x, y).RGBA()
				tc128 = complex(float64(r), 0)
				vol.SetValue(tc128, []int{i, x, y})
			}
		}
	}

	// FFT
	fft.SetWorkerPoolSize(3)
	fftvol := fft.FFTN(vol)

	// transform
	threashold := 100
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			for z := 0; z < 5; z++ {
				if (x > threashold && x <= w-threashold-1) && (y > threashold && y < h-1-threashold) && (z != 0 && z != 4) {
					fftvol.SetValue(complex(0, 0), []int{0, 0, 0})
					//	fftvol.SetValue(complex(0, 0), []int{0, 0, h - 1})
					//	fftvol.SetValue(complex(0, 0), []int{0, w - 1, 0})
					//	fftvol.SetValue(complex(0, 0), []int{0, w - 1, h - 1})
					//	fftvol.SetValue(complex(0, 0), []int{4, 0, 0})
					//	fftvol.SetValue(complex(0, 0), []int{4, 0, h - 1})
					//	fftvol.SetValue(complex(0, 0), []int{4, w - 1, 0})
					//	fftvol.SetValue(complex(0, 0), []int{4, w - 1, h - 1})
				}
			}
		}
	}

	// iFFT
	vol = fft.IFFTN(fftvol)

	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			tui8 = uint8(real(vol.Value([]int{0, x, y})))
			gray.Set(x, y, color.Gray{tui8})
		}
	}

	// Output the volume file
	outfile, _ := os.Create(os.Args[1])
	defer outfile.Close()
	png.Encode(outfile, gray)
}
