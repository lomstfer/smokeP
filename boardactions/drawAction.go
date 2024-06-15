package boardactions

import (
	"image"
	"image/color"
)

type DrawAction struct {
	PixelPoints         []image.Point
	color               color.NRGBA
	PreviousPixelcolors map[image.Point]color.NRGBA
}

func NewDrawAction(pixelPoints []image.Point, col color.NRGBA) *DrawAction {
	da := &DrawAction{PixelPoints: pixelPoints, color: col}
	da.PreviousPixelcolors = make(map[image.Point]color.NRGBA)
	return da
}

func (da DrawAction) Do(img *image.NRGBA) {
	for p := range da.PreviousPixelcolors {
		delete(da.PreviousPixelcolors, p)
	}
	
	for _, p := range da.PixelPoints {
		da.PreviousPixelcolors[p] = img.NRGBAAt(p.X, p.Y)
		img.SetNRGBA(p.X, p.Y, da.color)
	}
}

func (da DrawAction) Undo(img *image.NRGBA) {
	for p, c := range da.PreviousPixelcolors {
		img.SetNRGBA(p.X, p.Y, c)
	}
	for p := range da.PreviousPixelcolors {
		delete(da.PreviousPixelcolors, p)
	}
}
