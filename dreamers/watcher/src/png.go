package main

import (
	"image"
	"image/color"
	"image/png" // register the PNG format with the image package
	"os"
)

func main() {
	infile, _ := os.Open(os.Args[1])
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, _ := image.Decode(infile)

	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			_, _, _, a := oldColor.RGBA()
			gray.Set(x, y, color.Gray{uint8(a)})
		}
	}
	// Encode the grayscale image to the output file
	outfile, _ := os.Create(os.Args[2])
	defer outfile.Close()
	png.Encode(outfile, gray)
}
