package fakescreenshots

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func MakeFakeSetpArtifacts(timestamp string) {
	makeFakeScreenshot(timestamp)
	makeFakeHTML(timestamp)
	makeFakeLink(timestamp)
}

func makeFakeScreenshot(timestamp string) {
	var w, h int = 580, 540
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := float64(w / 4)
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 60}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 60}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cg.Brightness(float64(x), float64(y)),
				cb.Brightness(float64(x), float64(y)),
				255,
			}
			m.Set(x, y, c)
		}
	}

	addLabel(m, 20, 10, timestamp)

	f, err := os.OpenFile(`out`+string(os.PathSeparator)+timestamp+`.png`, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	png.Encode(f, m)

}

func makeFakeHTML(timestamp string) {
	str := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" dir="ltr" lang="en" style="background-color: rgb(0, 114, 198);">
<head>
    <title>` + timestamp + `</title>
</head>
<body>    
` + timestamp + `
</body>
</html>`

	file, err := os.Create(`out` + string(os.PathSeparator) + timestamp + `.html`)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(str)
}

func makeFakeLink(timestamp string) {
	str := `https://google.com/` + timestamp

	file, err := os.Create(`out` + string(os.PathSeparator) + timestamp + `.txt`)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(str)
}
