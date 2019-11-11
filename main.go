package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"github.com/macroblock/htimage/graph"
)

func main() {
	in := []string{
		"150x100.png",
		"300x200.png",
		"400x268.jpg",
		"450x675.jpg",
	}
	reader, err := os.Open(in[0])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	// bounds := img.Bounds()
	step := 5.0
	radius := step * 0.68 // for rgb - 0.58

	img, err = graph.TestFilter(step, radius, img)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
