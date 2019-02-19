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

	var r, g, b uint32
	var img [3]*dsputils.Matrix

	in.Seek(0, 0)
	src, _, _ := image.Decode(in)
	img[0] = dsputils.MakeEmptyMatrix([]int{w, h})
	img[1] = dsputils.MakeEmptyMatrix([]int{w, h})
	img[2] = dsputils.MakeEmptyMatrix([]int{w, h})

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ = src.At(x, y).RGBA()
			img[0].SetValue(complex(float64(r), 0), []int{x, y})
			img[1].SetValue(complex(float64(g), 0), []int{x, y})
			img[2].SetValue(complex(float64(b), 0), []int{x, y})
		}
	}

	// FFT
	var fftorig [3]*dsputils.Matrix

	fftorig[0] = fft.FFTN(img[0])
	fftorig[1] = fft.FFTN(img[1])
	fftorig[2] = fft.FFTN(img[2])

	// mid-processing
	var maxmin [4][3]float64

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			for z := 0; z < 3; z++ {
				fftr := real(fftorig[z].Value([]int{x, y}))
				ffti := imag(fftorig[z].Value([]int{x, y}))
				if maxmin[0][z] < fftr {
					maxmin[0][z] = fftr
				}
				if maxmin[1][z] > fftr {
					maxmin[1][z] = fftr
				}
				if maxmin[2][z] < ffti {
					maxmin[2][z] = ffti
				}
				if maxmin[3][z] > ffti {
					maxmin[3][z] = ffti
				}
			}
		}
	}

	linear := func(val, a, b float64) uint8 {
		return uint8(val * 255 / (a - b))
	}

	fftimg := image.NewRGBA(image.Rect(0, 0, w*2, h*2))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var temp [2][3]uint8
			for z := 0; z < 3; z++ {
				temp[0][z] = linear(real(fftorig[z].Value([]int{x, y})), maxmin[0][z], maxmin[1][z])
				temp[1][z] = linear(imag(fftorig[z].Value([]int{x, y})), maxmin[2][z], maxmin[3][z])
			}
			fftimg.Set(x, y, src.At(x, y))
			fftimg.Set(x+w, y, color.RGBA{temp[0][0], temp[0][1], temp[0][2], 255})
			fftimg.Set(x, y+h, color.RGBA{temp[0][0], temp[0][1], temp[0][2], 255})
			fftimg.Set(x+w+w/2, y+h, fftimg.At(x/8+w, y/8))
			fftimg.Set(x+w, y+h+h/2, fftimg.At(x/8, y/8+h))
		}
	}

	outfile, _ := os.Create(os.Args[2])
	defer outfile.Close()
	png.Encode(outfile, fftimg)

	// iFFT
	//vol = fft.IFFTN(fftorig)
}
