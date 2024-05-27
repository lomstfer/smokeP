package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"smokep/utils"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

const (
	zoomMultiplier     = 0.001
	defaultBoardWidth  = 8
	defaultBoardHeight = 8
)

type PixelBoard struct {
	pixelImg     *image.NRGBA
	position     f32.Point
	scale        float32
	drawingColor color.NRGBA
	bgImage      paint.ImageOp
}

func newPixelBoard() *PixelBoard {
	pb := &PixelBoard{}

	pb.pixelImg = image.NewNRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight))
	for i := range pb.pixelImg.Pix {
		pb.pixelImg.Pix[i] = uint8(rand.Intn(255))
	}

	pb.scale = 20

	pb.bgImage = paint.NewImageOp(utils.LoadImage("transp.jpg"))

	return pb
}

func (pb *PixelBoard) Size() f32.Point {
	return f32.Pt(defaultBoardWidth*pb.scale, defaultBoardHeight*pb.scale)
}

func (pb *PixelBoard) Draw(ops *op.Ops) {
	imgOp := paint.NewImageOp(pb.pixelImg)
	imgOp.Filter = paint.FilterNearest
	imgOp.Add(ops)

	tStack := op.Affine(
		f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(pb.scale, pb.scale)).Offset(pb.position),
	).Push(ops)
	paint.PaintOp{}.Add(ops)
	tStack.Pop()
}

func (pb *PixelBoard) CheckIfOnBoardAndDraw(mousePos f32.Point) {
	size := pb.Size()
	onBoard := mousePos.X > pb.position.X &&
		mousePos.X < pb.position.X+size.X &&
		mousePos.Y > pb.position.Y &&
		mousePos.Y < pb.position.Y+size.Y
	if onBoard {
		rel := mousePos.Sub(pb.position).Div(pb.scale)
		pixelCoord := image.Pt(int(rel.X), int(rel.Y))
		pb.pixelImg.SetNRGBA(pixelCoord.X, pixelCoord.Y, pb.drawingColor)
	}
}

func (pb *PixelBoard) Zoom(scrollY float32, mousePos f32.Point) {
	size := pb.Size()
	sizeNum := math.Sqrt(float64(size.X*size.X) + float64(size.Y*size.Y))
	scaleChange := -scrollY * zoomMultiplier * float32(sizeNum)
	pb.scale += scaleChange

	mouseRelBoard := mousePos.Sub(pb.position)

	ratioX := mouseRelBoard.X / size.X
	ratioY := mouseRelBoard.Y / size.Y
	pb.position = pb.position.Sub(f32.Pt(
		ratioX*scaleChange*defaultBoardWidth,
		ratioY*scaleChange*defaultBoardHeight,
	))
}
