package main

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type PixelBoard struct {
	img               *canvas.Image
	scaleFromZoom     float32
	moveOffsetTemp    fyne.Position
	moveOffset        fyne.Position
}

func newPixelBoard(img *canvas.Image) *PixelBoard {
	c := &PixelBoard{
		img:           img,
		scaleFromZoom: 1,
	}

	return c
}

func (pb *PixelBoard) Paint(relativePosition fyne.Position) {
	pixelPosX := int(relativePosition.X / pb.img.Size().Width * float32(pb.img.Image.Bounds().Dx()))
	pixelPosY := int(relativePosition.Y / pb.img.Size().Height * float32(pb.img.Image.Bounds().Dy()))

	rgbaImg := getPixelData(pb.img.Image)
	rgbaImg.Set(pixelPosX, pixelPosY, color.RGBA{255, 255, 255, 255})
	pb.img.Image = rgbaImg
	pb.img.Refresh()
}

func (pb *PixelBoard) updateMove(absMousePosition fyne.Position, mousePositionStart fyne.Position) {
	pb.moveOffsetTemp = absMousePosition.Subtract(mousePositionStart)

	newPos := pb.moveOffset.Add(
		fyne.NewPos(
			-pb.img.Size().Width/2+defaultBoardPixelWidth/2,
			-pb.img.Size().Height/2+defaultBoardPixelHeight/2,
		)).Add(pb.moveOffsetTemp)

	pb.img.Move(newPos)
}

func (pb *PixelBoard) endMove() {
	pb.moveOffset = pb.moveOffset.Add(pb.moveOffsetTemp)
	pb.moveOffsetTemp = fyne.NewPos(0, 0)
}

func (pb *PixelBoard) Zoom(scrollY float32, mouseRelBoard fyne.Position) {
	dims := pb.img.Size()
	size := math.Sqrt(float64(dims.Width*dims.Width) + float64(dims.Height*dims.Height))
	pb.scaleFromZoom += scrollY * zoomMultiplier * float32(size)

	mouseRelBoard = mouseRelBoard.Add(pb.getCenterDiff())
	scrollDir := float32(math.Abs(float64(scrollY)) / float64(scrollY))
	// mouseRelBoard = fyne.NewPos(mouseRelBoard.X * scrollDir, mouseRelBoard.Y * scrollDir)

	// sizeBefore := pb.img.Size()

	pb.img.Resize(fyne.NewSize(
		defaultBoardWidth*pb.scaleFromZoom,
		defaultBoardHeight*pb.scaleFromZoom,
	))

	// sizeDiff := pb.img.Size().Subtract(sizeBefore)

	pb.moveOffset = pb.moveOffset.AddXY(mouseRelBoard.X / 10 * -scrollDir, mouseRelBoard.Y / 10 * -scrollDir)

	newPos := pb.moveOffset.Add(pb.moveOffsetTemp).Add(pb.getCenterDiff())

	pb.img.Move(newPos)
}

func (pb *PixelBoard) getCenterDiff() (fyne.Position) {
	return fyne.NewPos(
		-pb.img.Size().Width/2+defaultBoardPixelWidth/2,
		-pb.img.Size().Height/2+defaultBoardPixelHeight/2,
	)
}