package boardactions

import (
	"image"
	"image/color"
)

type DrawAction struct {
	Pixels         map[image.Point]color.NRGBA
	PreviousPixels map[image.Point]color.NRGBA
}

func NewDrawAction(pixelPoints map[image.Point]color.NRGBA) *DrawAction {
	da := &DrawAction{}

	if pixelPoints == nil {
		da.Pixels = make(map[image.Point]color.NRGBA)
	} else {
		da.Pixels = pixelPoints
	}

	da.PreviousPixels = make(map[image.Point]color.NRGBA)
	return da
}

func (da DrawAction) Do(img *image.NRGBA) {
	for p := range da.PreviousPixels {
		delete(da.PreviousPixels, p)
	}

	for p, c := range da.Pixels {
		da.PreviousPixels[p] = img.NRGBAAt(p.X, p.Y)
		img.SetNRGBA(p.X, p.Y, c)
	}
}

func (da DrawAction) Undo(img *image.NRGBA) {
	for p, c := range da.PreviousPixels {
		img.SetNRGBA(p.X, p.Y, c)
	}
	for p := range da.PreviousPixels {
		delete(da.PreviousPixels, p)
	}
}
