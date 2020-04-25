package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	err = save(mosaic(img, 13), os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}

func save(img image.Image, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	out, err := os.Create(filename)

	switch ext {
	case ".jpg":
		err = jpeg.Encode(out, img, &jpeg.Options{})
		break
	case ".png":
	default:
		err = png.Encode(out, img)
		break
	}
	return err
}

func mosaic(img image.Image, dot int) image.Image {
	out := image.NewRGBA(img.Bounds())
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for y := img.Bounds().Min.Y; y <= height/dot; y++ {
		for x := img.Bounds().Min.X; x <= width/dot; x++ {
			px, py := x*dot, y*dot
			num := 0

			var r, g, b, a uint32
			for dy := 0; dy < dot; dy++ {
				for dx := 0; dx < dot; dx++ {
					if width <= px+dx || height <= py+dy {
						continue
					}

					pr, pg, pb, pa := img.At(px+dx, py+dy).RGBA()
					r += pr
					g += pg
					b += pb
					a += pa
					num++
				}
			}
			if num == 0 {
				continue
			}

			or := uint16(int(r) / num)
			og := uint16(int(g) / num)
			ob := uint16(int(b) / num)

			oa := uint16(int(a) / num)

			for dy := 0; dy < dot; dy++ {
				for dx := 0; dx < dot; dx++ {
					if width <= px+dx || height <= py+dy {
						continue
					}

					out.Set(px+dx, py+dy, color.RGBA64{R: or, G: og, B: ob, A: oa})
				}
			}
		}
	}

	return out
}
