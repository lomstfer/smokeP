package main

import (
	"math/rand"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type PixelBoard struct {
	boardObj       *fyne.Container
	pixelImage     *canvas.Image
	background     *canvas.Image
	scaleFromZoom  float32
	moveOffsetTemp fyne.Position
	moveOffset     fyne.Position
	paintingColor color.Color
}

func newPixelBoard(backgroundImage *canvas.Image) *PixelBoard {
	pixelImage := image.NewRGBA(image.Rect(0, 0, defaultBoardPixelWidth, defaultBoardPixelHeight))
	for i := range pixelImage.Pix {
		pixelImage.Pix[i] = uint8(rand.Intn(255))
	}
	pixelImageObj := canvas.NewImageFromImage(pixelImage)
	pixelImageObj.ScaleMode = canvas.ImageScalePixels

	pb := &PixelBoard{
		boardObj:      container.New(layout.NewStackLayout()),
		pixelImage:    pixelImageObj,
		background:    backgroundImage,
		scaleFromZoom: 1,
	}
	pb.boardObj.Add(pb.background)
	pb.boardObj.Add(pb.pixelImage)

	pb.boardObj.Resize(fyne.NewSize(defaultBoardWidth, defaultBoardHeight))
	pb.boardObj.Move(pb.getCenterDiff())

	

	return pb
}

func (pb *PixelBoard) Paint(relativePosition fyne.Position) {
	pixelPosX := int(relativePosition.X / pb.boardObj.Size().Width * float32(pb.pixelImage.Image.Bounds().Dx()))
	pixelPosY := int(relativePosition.Y / pb.boardObj.Size().Height * float32(pb.pixelImage.Image.Bounds().Dy()))

	rgbaImg := getPixelData(pb.pixelImage.Image)
	rgbaImg.Set(pixelPosX, pixelPosY, pb.paintingColor)
	pb.pixelImage.Image = rgbaImg
	pb.pixelImage.Refresh()
}

func (pb *PixelBoard) updateMove(absMousePosition fyne.Position, mousePositionStart fyne.Position) {
	pb.moveOffsetTemp = absMousePosition.Subtract(mousePositionStart)

	newPos := pb.moveOffset.Add(pb.moveOffsetTemp).Add(pb.getCenterDiff())

	pb.boardObj.Move(newPos)
}

func (pb *PixelBoard) endMove() {
	pb.moveOffset = pb.moveOffset.Add(pb.moveOffsetTemp)
	pb.moveOffsetTemp = fyne.NewPos(0, 0)
}

func (pb *PixelBoard) Zoom(scrollY float32, mouseRelBoard fyne.Position) {
	dims := pb.boardObj.Size()
	size := math.Sqrt(float64(dims.Width*dims.Width) + float64(dims.Height*dims.Height))
	pb.scaleFromZoom += scrollY * zoomMultiplier * float32(size)

	mouseRelBoard = mouseRelBoard.Add(pb.getCenterDiff())

	sizeBefore := pb.boardObj.Size()
	pb.boardObj.Resize(fyne.NewSize(
		defaultBoardWidth*pb.scaleFromZoom,
		defaultBoardHeight*pb.scaleFromZoom,
	))
	sizeDiff := pb.boardObj.Size().Subtract(sizeBefore)

	ratioX := mouseRelBoard.X / float32(sizeBefore.Width)
	ratioY := mouseRelBoard.Y / float32(sizeBefore.Height)

	pb.moveOffset = pb.moveOffset.Add(fyne.NewPos(
		-ratioX*float32(sizeDiff.Width),
		-ratioY*float32(sizeDiff.Height),
	))
	newPos := pb.moveOffset.Add(pb.moveOffsetTemp).Add(pb.getCenterDiff())
	pb.boardObj.Move(newPos)
}

func (pb *PixelBoard) getCenterDiff() fyne.Position {
	return fyne.NewPos(
		-pb.boardObj.Size().Width/2+defaultBoardPixelWidth/2,
		-pb.boardObj.Size().Height/2+defaultBoardPixelHeight/2,
	)
}
