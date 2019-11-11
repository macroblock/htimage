package graph

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
)

var angleSlice = []float64{15, 75, 0, 45}
var layer []*gg.Context
var radian = math.Pi / 180

// TestFilter -
func TestFilter(step, radius float64, src image.Image) (image.Image, error) {
	b := src.Bounds()

	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	src = m

	w := src.Bounds().Dx()
	h := src.Bounds().Dy()

	// dst := image.NewRGBA(image.Rect(0, 0, int(float64(w)*step+step), int(float64(h)*step+step)))

	dst := gg.NewContext(int(float64(w)*step+step), int(float64(h)*step+step))
	dst.SetRGBA255(0, 0, 0, 255)
	dst.Clear()

	for range angleSlice {
		x := gg.NewContext(int(float64(w)*step+step), int(float64(h)*step+step))
		x.SetRGBA(0, 0, 0, 0)
		x.Clear()
		layer = append(layer, x)
	}

	for i, angle := range angleSlice {
		cosAng := math.Cos(angle * radian)
		sinAng := math.Sin(angle * radian)

		w1 := float64(w) * math.Cos(angle*radian)
		w2 := float64(w) * math.Sin(angle*radian)

		h1 := float64(h) * math.Sin((angle)*radian)
		h2 := float64(h) * math.Cos((angle)*radian)

		w0 := w1 + w2
		h0 := h1 + h2

		wOffs := h1 * math.Cos(angle*radian)
		hOffs := float64(h) - h1*math.Sin(angle*radian)
		_, _ = wOffs, hOffs

		// step = 10.0

		fmt.Printf("woffs, hoffs: %v, %v\n", wOffs, hOffs)
		fmt.Printf("woffs, hoffs: %v, %v\n", wOffs, hOffs)
		y := 0.0
		for y < h0 {
			x := 0.0
			for x < w0 {
				x++

				xImg := x*cosAng + y*sinAng - wOffs
				yImg := x*sinAng - y*cosAng + hOffs
				xImgInt := int(math.Round(xImg))
				yImgInt := int(math.Round(yImg))

				r, g, b, _ := src.At(xImgInt, yImgInt).RGBA()
				rr := float64(r>>8) / 255
				gg := float64(g>>8) / 255
				bb := float64(b>>8) / 255
				kk := math.Min(1-math.Max(math.Max(rr, gg), bb), 0.9999999999)
				cc := (1 - rr - kk) / (1 - kk)
				mm := (1 - gg - kk) / (1 - kk)
				yy := (1 - bb - kk) / (1 - kk)

				// cc := (1 - rr)
				// mm := (1 - gg)
				// yy := (1 - bb)
				// kk := 0
				lum := 255.0
				switch i {
				case 0:
					lum = cc
					rr = 1
				case 1:
					lum = mm
					rr = 1
				case 2:
					lum = yy
					rr = 1
				case 3:
					lum = kk * 1
					rr = 1
				}
				// lum = (float64(r>>8) + float64(g>>8) + float64(b>>8)) / 3
				s0 := (radius * radius) * math.Pi
				rad := math.Sqrt(s0 * lum / math.Pi)
				// rad = (radius * lum)
				// fmt.Println(rr, gg, bb)
				// fmt.Println(cc, mm, yy, kk, lum, " --- ")
				layer[i].SetRGBA(rr, rr, rr, 1)
				layer[i].DrawCircle(xImg*step+step/2, yImg*step+step/2, rad)
				layer[i].Fill()
			}
			y++
		}
	}

	for y := 0; y < dst.Image().Bounds().Dy(); y++ {
		for x := 0; x < dst.Image().Bounds().Dx(); x++ {
			c1, _, _, _ := layer[0].Image().At(x, y).RGBA()
			_, m1, _, _ := layer[1].Image().At(x, y).RGBA()
			_, _, y1, _ := layer[2].Image().At(x, y).RGBA()
			_, _, _, k1 := layer[3].Image().At(x, y).RGBA()

			cc := float64(c1>>8) / 255
			mm := float64(m1>>8) / 255
			yy := float64(y1>>8) / 255
			kk := float64(k1>>8) / 255

			zz := 1.0
			r := (1 - cc) * (1 - kk*zz)
			g := (1 - mm) * (1 - kk*zz)
			b := (1 - yy) * (1 - kk*zz)
			// fmt.Print(cc, mm, yy, kk, r, g, b, "-")

			// c := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(255)}
			dst.SetRGBA(r, g, b, 1)
			dst.SetPixel(x, y)
		}
	}

	// for y := 0; y < h; y++ {
	// 	for x := 0; x < w; x++ {
	// 		// Circle(dst, float64(x)*step, float64(y)*step, radius, src.At(x, y))
	// 		r, g, b, _ := src.At(x, y).RGBA()
	// 		red = r >> 8
	// 		lum = (float64(r>>8) + float64(g>>8) + float64(b>>8)) / 3
	// 		rad = (radius / 255 * lum)
	// 		dst.SetRGBA255(0, 0, 0, 255)
	// 		dst.DrawCircle(float64(x)*step, float64(y)*step, rad)
	// 		dst.Fill()
	// 	}
	// 	fmt.Printf("%v/%v r: %v radius:%v luma:%v               \x0d", y*w, w*h, red, rad, lum)
	// }

	return dst.Image(), nil
}

// Circle -
func Circle(img *gg.Context, x, y float64, radius float64, col color.Color) {
	// img.Set(int(math.Round(x)), int(math.Round(y)), col)
	// x0, y0 := iround(x), iround(y)
	// for yy := iround(y - radius); yy < iround(y+radius); yy++ {
	// 	for xx := iround(x - radius); xx < iround(x+radius); xx++ {
	// 		r, g, b, a := img.At(xx, yy).RGBA()
	// 		if (xx-x0)*(xx-x0)+(yy-y0)*(yy-y0) < iround(radius*radius) {
	// 			r, g, b, a = col.RGBA()
	// 		}

	// 		img.Set(xx, yy, color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	// 	}
	// }
}

func iround(v float64) int {
	return int(math.Round(v))
}
