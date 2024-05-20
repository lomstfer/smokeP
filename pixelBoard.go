package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

const (
	zoomMultiplier     = 0.001
	defaultBoardWidth  = 8
	defaultBoardHeight = 8
)

type PixelBoard struct {
	pixelImg *image.NRGBA
	position f32.Point
	size     f32.Point
	scale    float32
}

func newPixelBoard() *PixelBoard {
	pb := &PixelBoard{}

	pb.pixelImg = image.NewNRGBA(image.Rect(0, 0, defaultBoardWidth, defaultBoardHeight))
	for i := range pb.pixelImg.Pix {
		pb.pixelImg.Pix[i] = uint8(rand.Intn(255))
	}

	pb.scale = 20

	return pb
}

func (pb *PixelBoard) Layout(gtx layout.Context) layout.Dimensions {
	r := image.Rect(pb.position.Round().X, pb.position.Round().Y, pb.position.Round().X+pb.size.Round().X, pb.position.Round().Y+pb.size.Round().Y)

	imgOp := paint.NewImageOp(pb.pixelImg)
	imgOp.Filter = paint.FilterNearest
	imgOp.Add(gtx.Ops)

	pb.size = f32.Pt(defaultBoardWidth*pb.scale, defaultBoardHeight*pb.scale)

	op.Affine(
		f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(pb.scale, pb.scale)).Offset(pb.position),
	).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(r.Dx(), r.Dy())}
}

func (pb *PixelBoard) CheckAndDraw(mousePos f32.Point) {
	onBoard := mousePos.X > pb.position.X &&
		mousePos.X < pb.position.X+pb.size.X &&
		mousePos.Y > pb.position.Y &&
		mousePos.Y < pb.position.Y+pb.size.Y
	if onBoard {
		rel := mousePos.Sub(pb.position).Div(pb.scale)
		pixelCoord := image.Pt(int(rel.X), int(rel.Y))
		pb.pixelImg.SetNRGBA(pixelCoord.X, pixelCoord.Y, color.NRGBA{255, 255, 255, 255})
	}
}

func (pb *PixelBoard) Zoom(scrollY float32, mousePos f32.Point) {
	dims := pb.size
	size := math.Sqrt(float64(dims.X*dims.X) + float64(dims.Y*dims.Y))
	scaleChange := -scrollY * zoomMultiplier * float32(size)
	pb.scale += scaleChange

	mouseRelBoard := mousePos.Sub(pb.position)

	ratioX := mouseRelBoard.X / pb.size.X
	ratioY := mouseRelBoard.Y / pb.size.Y
	pb.position = pb.position.Sub(f32.Pt(
		ratioX*scaleChange*defaultBoardWidth,
		ratioY*scaleChange*defaultBoardHeight,
	))

	pb.size = f32.Pt(
		defaultBoardWidth*pb.scale,
		defaultBoardHeight*pb.scale,
	)
}
