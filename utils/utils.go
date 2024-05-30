package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
)

func LoadImage(name string) image.Image {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	return img
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
