package main

import (
	"image"
	"image/color"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type PixelBoard struct {
	pixelImg *image.NRGBA
	position f32.Point
	size f32.Point
	scale float32
}

func newPixelBoard() *PixelBoard {
	pb := &PixelBoard{}

	pb.pixelImg = image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for i := range pb.pixelImg.Pix {
		pb.pixelImg.Pix[i] = uint8(rand.Intn(255))
	}

	pb.scale = 20

	return pb
}

func (pb *PixelBoard) Layout(gtx layout.Context) layout.Dimensions {
	r := image.Rect(pb.position.Round().X, pb.position.Round().Y, pb.position.Round().X + pb.size.Round().X, pb.position.Round().Y + pb.size.Round().Y)

	imgOp := paint.NewImageOp(pb.pixelImg)
	imgOp.Filter = paint.FilterNearest
	imgOp.Add(gtx.Ops)

	pb.size = f32.Pt(float32(imgOp.Size().X) * pb.scale, float32(imgOp.Size().Y) * pb.scale)
	
	op.Affine(
		f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(pb.scale, pb.scale)).Offset(pb.position),
	).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(r.Dx(), r.Dy())}
}

func (pb * PixelBoard) CheckAndDraw(mousePos f32.Point) {
	onBoard := mousePos.X > pb.position.X &&
				mousePos.X < pb.position.X+pb.size.X &&
				mousePos.Y > pb.position.Y &&
				mousePos.Y < pb.position.Y+pb.size.Y
	if onBoard {
		rel := mousePos.Sub(pb.position).Div(pb.scale)
		pixelCoord := image.Pt(int(rel.X), int(rel.Y))
		pb.pixelImg.SetNRGBA(pixelCoord.X, pixelCoord.Y, color.NRGBA{255,255,255,255})
	}
}