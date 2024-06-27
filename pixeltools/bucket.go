package pixeltools

import (
	"image"
	"image/color"
	"smokep/boardactions"
)

func BucketConstrainedOnClick(img *image.NRGBA, pos image.Point, fillColor color.NRGBA) *boardactions.DrawAction {
	colorAtClick := img.NRGBAAt(pos.X, pos.Y)
	if colorAtClick == fillColor {
		return nil
	}

	imgCopy := image.NewNRGBA(image.Rect(0, 0, img.Rect.Dx(), img.Rect.Dy()))
	copy(imgCopy.Pix, img.Pix)
	img = imgCopy

	drawActionPoints := make(map[image.Point]color.NRGBA, 0)

	queue := []image.Point{pos}
	img.SetNRGBA(pos.X, pos.Y, fillColor)
	drawActionPoints[pos] = fillColor

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		neighbours := make([]image.Point, 0)
		if pos.X > 0 {
			neighbours = append(neighbours, pos.Add(image.Pt(-1, 0)))
		}
		if pos.X < img.Rect.Dx()-1 {
			neighbours = append(neighbours, pos.Add(image.Pt(1, 0)))
		}
		if pos.Y > 0 {
			neighbours = append(neighbours, pos.Add(image.Pt(0, -1)))
		}
		if pos.Y < img.Rect.Dy()-1 {
			neighbours = append(neighbours, pos.Add(image.Pt(0, 1)))
		}

		for _, n := range neighbours {
			if img.NRGBAAt(n.X, n.Y) == colorAtClick {
				queue = append(queue, n)
				img.SetNRGBA(n.X, n.Y, fillColor)
				drawActionPoints[n] = fillColor
			}
		}
	}

	action := boardactions.NewDrawAction(drawActionPoints)
	return action
}

func BucketAllOnClick(img *image.NRGBA, pos image.Point, fillColor color.NRGBA) *boardactions.DrawAction {
	colorAtClick := img.NRGBAAt(pos.X, pos.Y)
	if colorAtClick == fillColor {
		return nil
	}

	drawActionPoints := make(map[image.Point]color.NRGBA, 0)

	for x := range img.Rect.Dx() {
		for y := range img.Rect.Dy() {
			if img.NRGBAAt(x, y) == colorAtClick {
				drawActionPoints[image.Pt(x, y)] = fillColor
			}
		}
	}

	action := boardactions.NewDrawAction(drawActionPoints)
	return action
}
