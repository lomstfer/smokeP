package pixeltools

import (
	"image"
	"image/color"
	"smokep/boardactions"
	"smokep/utils"
)

type Pencil struct {
	previousDrawPoint *image.Point
	currentDrawAction *boardactions.DrawAction
}

func (pen *Pencil) OnEnd() *boardactions.DrawAction {
	pen.previousDrawPoint = nil
	action := pen.currentDrawAction
	pen.currentDrawAction = nil
	return action
}

func (pen *Pencil) OnDraw(img *image.NRGBA, pos image.Point, col color.NRGBA) {
	if img.NRGBAAt(pos.X, pos.Y) == col {
		pen.previousDrawPoint = &pos
		return
	}

	if pen.currentDrawAction == nil {
		pen.currentDrawAction = boardactions.NewDrawAction(nil)
	}

	if pen.previousDrawPoint != nil {
		var points []image.Point
		{
			betweenPoints := utils.GetLineBetweenPoints(*pen.previousDrawPoint, pos)
			for _, p := range betweenPoints {
				if img.NRGBAAt(p.X, p.Y) != col {
					points = append(points, p)
				}
			}
		}

		for _, p := range points {
			pen.currentDrawAction.PreviousPixels[p] = img.NRGBAAt(p.X, p.Y)
			img.SetNRGBA(p.X, p.Y, col)
			pen.currentDrawAction.Pixels[p] = col
		}
	} else {
		pen.currentDrawAction.PreviousPixels[pos] = img.NRGBAAt(pos.X, pos.Y)
		img.SetNRGBA(pos.X, pos.Y, col)
		pen.currentDrawAction.Pixels[pos] = col
	}

	pen.previousDrawPoint = &pos
}
