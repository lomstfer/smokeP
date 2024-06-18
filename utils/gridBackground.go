package utils

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type GridBackground struct {
	imageOp paint.ImageOp
	WindowSize image.Point
}

func NewGridBackground() *GridBackground {
	gb := &GridBackground{}
	gb.imageOp = paint.NewImageOp(GenerateGridImage(160, 90, color.NRGBA{200, 200, 200, 255}, color.NRGBA{100, 100, 100, 255}))
	gb.imageOp.Filter = paint.FilterNearest
	return gb
}

func (gb *GridBackground) Draw(ops *op.Ops, area image.Rectangle) {
	defer clip.Rect(area).Push(ops).Pop()
	gb.imageOp.Add(ops)
	scale := max(float32(gb.WindowSize.X)/float32(gb.imageOp.Size().X), float32(gb.WindowSize.Y)/float32(gb.imageOp.Size().Y))
	tStack := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(ops)
	paint.PaintOp{}.Add(ops)
	tStack.Pop()
}