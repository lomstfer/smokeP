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