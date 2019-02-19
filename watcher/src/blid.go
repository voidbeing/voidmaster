package main

import (
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	//"github.com/anthonynsimon/bild/segment"
	"os"
)

func main() {
	img, err := imgio.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	//seg := segment.Threshold(img, 128)
	//med := effect.Median(seg, 4.0)

	med := effect.Outline(img)

	if err := imgio.Save(os.Args[2], med, imgio.PNGEncoder()); err != nil {
		panic(err)
	}
}
