package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"os"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
)

func LoadImage(name string) *image.NRGBA {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var nrgba *image.NRGBA

	switch img := img.(type) {
	case *image.RGBA:
		nrgba = image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)
	case *image.NRGBA:
		nrgba = img
	default:
		fmt.Println("Unsupported image type")
	}

	return nrgba
}

func Clamp(x float64, min float64, max float64) float64 {
	return math.Max(math.Min(x, max), min)
}

func ClampInt(x int, min int, max int) int {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}

func GenerateGridImage(width, height int, color1 color.NRGBA, color2 color.NRGBA) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	doColor1 := false
	for i := 0; i < len(img.Pix); i += 4 {
		if (i / 4 % width == 0) {
			doColor1 = !doColor1
		}

		if doColor1 {
			img.Pix[i] = color1.R
			img.Pix[i + 1] = color1.G
			img.Pix[i + 2] = color1.B
			img.Pix[i + 3] = color1.A
		} else {
			img.Pix[i] = color2.R
			img.Pix[i + 1] = color2.G
			img.Pix[i + 2] = color2.B
			img.Pix[i + 3] = color2.A
		}

		doColor1 = !doColor1
	}

	return img
}

// func getPixelData(img image.Image) *image.NRGBA {
// 	rgba := image.NewNRGBA(img.Bounds())

// 	for y := 0; y < img.Bounds().Dx(); y++ {
// 		for x := 0; x < img.Bounds().Dy(); x++ {
// 			c := img.At(x, y)
// 			rgba.Set(x, y, c)
// 		}
// 	}

// 	return rgba
// }

func IsLight(c color.NRGBA) bool {
    r := float64(c.R) / 255.0
    g := float64(c.G) / 255.0
    b := float64(c.B) / 255.0

    luminance := 0.299*r + 0.587*g + 0.114*b

    return luminance > 0.5
}

func SaveImageToFile(img *image.NRGBA, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    err = png.Encode(file, img)
	if err != nil {
        return err
    }

    return nil
}

func GetLineBetweenPoints(p1 image.Point, p2 image.Point) []image.Point {
	var points []image.Point

	d := p2.Sub(p1)
	if int(math.Abs(float64(d.X))) > int(math.Abs(float64(d.Y))) {
		k := float32(d.Y) / float32(d.X)
		for i := 1; i <= int(math.Abs(float64(d.X))); i++ {
			x := i * int(float32(d.X)/float32(math.Abs(float64(d.X))))
			y := int(math.Round(float64(float32(x) * k)))
			pos := p1.Add(image.Pt(x, y))
			points = append(points, pos)
		}
	} else {
		k := float32(d.X) / float32(d.Y)
		for i := 1; i <= int(math.Abs(float64(d.Y))); i++ {
			y := i * int(float32(d.Y)/float32(math.Abs(float64(d.Y))))
			x := int(math.Round(float64(float32(y) * k)))
			pos := p1.Add(image.Pt(x, y))
			points = append(points, pos)
		}
	}

	return points
}

func FocusSelfOnClick(selfTag interface{}, gtx layout.Context) {
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:       selfTag,
			Kinds:        pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if !e.Buttons.Contain(pointer.ButtonPrimary) {
			continue
		}

		gtx.Execute(key.FocusCmd{Tag: selfTag})
	}
}