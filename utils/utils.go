package utils

import (
	"bytes"
	"fmt"
	"image"
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
