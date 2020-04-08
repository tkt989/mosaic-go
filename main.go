package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

func main() {
	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create("./out.png")
	_ = png.Encode(out, mosaic(img, 13))
}

func mosaic(img image.Image, dot int) image.Image {
	out := image.NewRGBA(img.Bounds())

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y/dot; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X/dot; x++ {
			px, py := x*dot, y*dot

			var r, g, b, a uint32
			for dy := 0; dy < dot; dy++ {
				for dx := 0; dx < dot; dx++ {
					pr, pg, pb, pa := img.At(px+dx, py+dy).RGBA()
					r += pr
					g += pg
					b += pb
					a += pa
				}
			}
			or := uint8(int(r) / (dot * dot))
			og := uint8(int(g) / (dot * dot))
			ob := uint8(int(b) / (dot * dot))

			oa := uint8(int(a) / (dot * dot))

			for dy := 0; dy < dot; dy++ {
				for dx := 0; dx < dot; dx++ {
					out.Set(px+dx, py+dy, color.RGBA{R: or, G: og, B: ob, A: oa})
				}
			}
		}
	}

	return out
}
