package boardactions

import (
	"image"
	"math"

	"gioui.org/f32"
)

type ResizeAction struct {
	Size image.Point
	resizeOrigin f32.Point // (0, 0) to (1, 1)
	PreviousSize image.Point
	PreviousSizePixels []uint8
}

func NewResizeAction(size image.Point, resizeOrigin f32.Point) *ResizeAction {
	da := &ResizeAction{Size: size, resizeOrigin: resizeOrigin}
	return da
}

func (da *ResizeAction) Do(img *image.NRGBA) {
	da.PreviousSize = img.Bounds().Size()
	da.PreviousSizePixels = img.Pix

	difference := da.Size.Sub(da.PreviousSize)
	offsetX := int(math.Round(float64(float32(difference.X) * da.resizeOrigin.X)))
	offsetY := int(math.Round(float64(float32(difference.Y) * da.resizeOrigin.Y)))

	newImage := image.NewNRGBA(image.Rect(0, 0, da.Size.X, da.Size.Y))
	for x := 0; x < da.Size.X; x++ {
		for y := 0; y < da.Size.Y; y++ {
			newImage.SetNRGBA(x, y, img.NRGBAAt(x - offsetX, y - offsetY))
		}
	}

	*img = *newImage
}

func (da *ResizeAction) Undo(img *image.NRGBA) {
	difference := da.Size.Sub(da.PreviousSize)
	offsetX := int(math.Round(float64(float32(difference.X) * da.resizeOrigin.X)))
	offsetY := int(math.Round(float64(float32(difference.Y) * da.resizeOrigin.Y)))

	newImage := image.NewNRGBA(image.Rect(0, 0, da.PreviousSize.X, da.PreviousSize.Y))
	for x := 0; x < da.PreviousSize.X; x++ {
		for y := 0; y < da.PreviousSize.Y; y++ {
			newImage.SetNRGBA(x, y, img.NRGBAAt(x + offsetX, y + offsetY))
		}
	}

	newImage.Pix = da.PreviousSizePixels

	*img = *newImage
}
